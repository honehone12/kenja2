package mongodb

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"kenja2/documents"
	"kenja2/engine"
	"kenja2/keywords"
	"net/url"
	"os"
	"strings"

	"github.com/chewxy/math32"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorgonia.org/vecf32"
)

const ATLAS_KW_LIMIT = 300 //at least js accecpts 120 -> 120 * 2bytes
const ATLAS_SEARCH_LIMIT = 120
const ATLAS_SEARCH_SIZE = 50
const FUZZY_MAX_EDIT_TO_MATCH = 1 // 1 or 2
const FUZZY_MAX_TRY_PER_KW = 50
const FUZZY_NO_CHNGE_PREFIX = 1
const FUZZY_MAX_SINGLE_TKN = 3
const FUZZY_MAX_KW_LIMIT = 65
const VECTOR_LEN = 768
const VECTOR_ENCODED_LEN = 1024

type atlasClient struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	logger     engine.Logger

	replacer *strings.Replacer
}

func newAtlasClient(
	m *mongo.Client, l engine.Logger,
) (*atlasClient, error) {
	db := os.Getenv("MONGO_DATABASE")
	if len(db) == 0 {
		return nil, errors.New("env for db name is not set")
	}
	cl := os.Getenv("MONGO_COLLECTION")
	if len(cl) == 0 {
		return nil, errors.New("env for collection is not set")
	}

	d := m.Database(db)
	c := d.Collection(cl)

	return &atlasClient{
		client:     m,
		database:   d,
		collection: c,
		logger:     l,
		replacer:   __mongoReplacer,
	}, nil
}

func connect(
	uri string,
	maxConn uint64,
	l engine.Logger,
) (*atlasClient, error) {
	ops := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(maxConn).
		SetMaxConnecting(1)

	m, err := mongo.Connect(ops)
	if err != nil {
		return nil, err
	}

	return newAtlasClient(m, l)
}

func connectAws(
	format string,
	maxConn uint64,
	l engine.Logger,
) (*atlasClient, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if len(accessKey) == 0 {
		return nil, errors.New("env for aws access key id is not set")
	}
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if len(secretKey) == 0 {
		return nil, errors.New("env for aws secret access key is not set")
	}
	tkn := os.Getenv("AWS_SESSION_TOKEN")
	if len(tkn) == 0 {
		return nil, errors.New("env for aws session token is not set")
	}

	uri := fmt.Sprintf(
		format,
		url.QueryEscape(accessKey),
		url.QueryEscape(secretKey),
		url.QueryEscape(tkn),
	)
	return connect(uri, maxConn, l)
}

func fieldProjectionFlags() bson.M {
	return bson.M{
		"updated_at":   0,
		"image_vector": 0,
	}
}

func allITypeFilter() bson.M {
	return bson.M{
		"$in": []documents.ItemType{
			documents.ITEM_TYPE_ANIME,
			documents.ITEM_TYPE_CHARACTER,
			documents.ITEM_TYPE_YVIDEO,
		},
	}
}

func (a *atlasClient) ping(ctx context.Context) error {
	a.logger.Info("[ping]")

	return a.client.Ping(ctx, nil)
}

func (a *atlasClient) random(
	ctx context.Context, q documents.Query,
) (*mongo.Cursor, error) {
	it := q.RandomQuery()

	a.logger.Infof("[random] item_type=%d", it)

	if err := it.Ok(); err != nil {
		return nil, err
	}

	p := mongo.Pipeline{}
	if it != documents.ITEM_TYPE_UNSPECIFIED {
		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": it,
			}},
		})
	}

	p = append(p, bson.D{
		{Key: "$sample", Value: bson.M{
			"size": ATLAS_SEARCH_LIMIT,
		}},
	})

	op := options.Aggregate()
	return a.collection.Aggregate(ctx, p, op)
}

func (a *atlasClient) unsafeTextSearch(
	ctx context.Context,
	keywords keywords.Keywords,
	itemType documents.ItemType,
	fuzzy documents.FuzzyKind,
) (*mongo.Cursor, error) {
	text := bson.M{
		"query": keywords,
		"path": []string{
			"name",
			"name_english",
			"name_japanese",
			"parent.name",
			"parent.name_japanese",
		},
		"matchCriteria": "all",
	}
	if fuzzy == documents.FUZZY_USE {
		text["fuzzy"] = bson.M{
			"maxEdits":      FUZZY_MAX_EDIT_TO_MATCH,
			"maxExpansions": FUZZY_MAX_TRY_PER_KW,
			"prefixLength":  FUZZY_NO_CHNGE_PREFIX,
		}
	}

	p := mongo.Pipeline{
		{
			{Key: "$search", Value: bson.M{
				"index":              "text",
				"text":               text,
				"returnStoredSource": true,
			}},
		},
	}

	if itemType == documents.ITEM_TYPE_UNSPECIFIED {
		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": allITypeFilter(),
			}},
		})
	} else {
		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": itemType,
			}},
		})
	}

	p = append(p,
		bson.D{{Key: "$limit", Value: ATLAS_SEARCH_LIMIT}},
		bson.D{{Key: "$project", Value: fieldProjectionFlags()}},
	)

	op := options.Aggregate()
	return a.collection.Aggregate(ctx, p, op)
}

func (a *atlasClient) textSearch(
	ctx context.Context, q documents.Query,
) (*mongo.Cursor, error) {
	kw, it, fz := q.TextSearchQuery()

	a.logger.Infof("[text] keywords=%s item_type=%d fuzzy_kind=%d", kw, it, fz)

	if len(kw) > ATLAS_KW_LIMIT {
		return nil, errors.New("unexpected keywords length")
	}
	if err := it.Ok(); err != nil {
		return nil, err
	}
	if err := fz.Ok(); err != nil {
		return nil, err
	}

	cleanKw, singleTkns := keywords.CleanKeywords(a.replacer, kw)
	lenKw := len(cleanKw)
	if lenKw == 0 {
		return nil, errors.New("invalid keywords")
	}
	if fz == documents.FUZZY_USE {
		// turn fuzzy off to prevent error
		if singleTkns > FUZZY_MAX_SINGLE_TKN {
			fz = documents.FUZZY_UNSPECIFIED
		}
		if lenKw > FUZZY_MAX_KW_LIMIT {
			fz = documents.FUZZY_UNSPECIFIED
		}
	}

	a.logger.Infof("[text] (actual) keywords=%s", cleanKw)

	return a.unsafeTextSearch(
		ctx,
		cleanKw,
		it,
		fz,
	)
}

func (a *atlasClient) unsafeVectorSearch(
	ctx context.Context,
	itemType documents.ItemType,
	srcVec bson.Binary,
) (*mongo.Cursor, error) {
	var filter bson.M
	switch itemType {
	case documents.ITEM_TYPE_ANIME:
		filter = bson.M{
			"item_type": documents.ITEM_TYPE_ANIME,
		}
	case documents.ITEM_TYPE_CHARACTER:
		filter = bson.M{
			"item_type": documents.ITEM_TYPE_CHARACTER,
		}
	default:
		return nil, fmt.Errorf("unexpected item type: %d", itemType)
	}

	p := mongo.Pipeline{
		{
			{Key: "$vectorSearch", Value: bson.M{
				"exact":         false,
				"filter":        filter,
				"index":         "vector",
				"limit":         ATLAS_SEARCH_LIMIT,
				"numCandidates": ATLAS_SEARCH_LIMIT * ATLAS_SEARCH_SIZE, // max 10000
				"path":          "image_vector",
				"queryVector":   srcVec,
			}},
		},
		{
			{Key: "$project", Value: fieldProjectionFlags()},
		},
	}

	op := options.Aggregate()
	return a.collection.Aggregate(ctx, p, op)
}

func (a *atlasClient) queryVector(
	ctx context.Context, id bson.ObjectID,
) (bson.Binary, error) {
	f := bson.D{{Key: "_id", Value: id}}
	op := options.FindOne().SetProjection(bson.M{
		"image_vector": 1,
	})
	res := a.collection.FindOne(ctx, f, op)
	if res.Err() != nil {
		return bson.Binary{Data: nil}, res.Err()
	}

	v := documents.ImageVector{}
	err := res.Decode(&v)
	if err != nil {
		return bson.Binary{Data: nil}, err
	}
	if v.ImageVector.IsZero() {
		return bson.Binary{Data: nil}, errors.New("image vector is empty")
	}

	return v.ImageVector, nil
}

func (a *atlasClient) queryMultiVector(
	ctx context.Context, ids []bson.ObjectID,
) ([]bson.Binary, error) {
	f := bson.D{{
		Key:   "_id",
		Value: bson.M{"$in": ids},
	}}
	op := options.Find().SetProjection(bson.M{
		"image_vector": 1,
	})
	stream, err := a.collection.Find(ctx, f, op)
	if err != nil {
		return nil, err
	}
	defer stream.Close(ctx)

	vs := []bson.Binary{}
	for stream.Next(ctx) {
		v := documents.ImageVector{}
		if err := stream.Decode(&v); err != nil {
			return nil, err
		}

		if v.ImageVector.IsZero() {
			return nil, errors.New("image vector is empty")
		}

		vs = append(vs, v.ImageVector)
	}
	if err := stream.Err(); err != nil {
		return nil, err
	}

	return vs, nil
}

func (a *atlasClient) mixVectors(vectors []bson.Binary) (bson.Binary, error) {
	l := len(vectors)
	if l <= 1 {
		return bson.Binary{Data: nil}, errors.New("invalid shape")
	}

	fvectors := make([][]float32, l)
	for i := range l {
		if len(vectors[i].Data) != VECTOR_LEN+2 {
			return bson.Binary{Data: nil}, errors.New("invalid vector length")
		}

		f := make([]float32, VECTOR_LEN)
		for j := range VECTOR_LEN {
			f[j] = float32(int8(vectors[i].Data[j+2]))
		}

		fvectors[i] = f
	}

	// just get average
	for i := 1; i < l; i++ {
		vecf32.Add(fvectors[0], fvectors[i])
	}
	vecf32.ScaleInv(fvectors[0], float32(l))

	// to int8
	b := make([]byte, VECTOR_LEN+2)
	b[0] = bson.Int8Vector
	b[1] = 0
	for i := range VECTOR_LEN {
		b[i+2] = byte(int8(math32.Round(fvectors[0][i])))
	}

	v := bson.Binary{
		Subtype: bson.TypeBinaryVector,
		Data:    b,
	}

	return v, nil
}

func (a *atlasClient) vectorSearch(
	ctx context.Context, q documents.Query,
) (*mongo.Cursor, error) {
	id, it := q.VectorSearchQuery()

	a.logger.Infof("[vector] id=%s item_type=%d", id, it)

	if err := it.Ok(); err != nil {
		return nil, err
	}
	if it == documents.ITEM_TYPE_UNSPECIFIED {
		return nil, errors.New("vector search requires item type")
	}

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	v, err := a.queryVector(ctx, oid)
	if err != nil {
		return nil, err
	}

	return a.unsafeVectorSearch(ctx, it, v)
}

func (a *atlasClient) directVectorSearch(
	ctx context.Context, q documents.Query,
) (*mongo.Cursor, error) {
	enc, it := q.DirectVectorSearchQuery()

	a.logger.Infof("[direct vector] vector=%s item_type=%d", enc, it)

	if err := it.Ok(); err != nil {
		return nil, err
	}
	if it == documents.ITEM_TYPE_UNSPECIFIED {
		return nil, errors.New("vector search requires item type")
	}

	if len(enc) != VECTOR_ENCODED_LEN {
		return nil, errors.New("unexpected encoded vecter length")
	}

	dec, err := base64.URLEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}
	if len(dec) != VECTOR_LEN {
		return nil, errors.New("unexpected vector length")
	}

	b := make([]byte, VECTOR_LEN+2)
	b[0] = bson.Int8Vector
	b[1] = 0
	copy(b[2:], dec)

	v := bson.Binary{
		Subtype: bson.TypeBinaryVector,
		Data:    b,
	}

	return a.unsafeVectorSearch(ctx, it, v)
}

func (a *atlasClient) multiVectorSearch(
	ctx context.Context,
	q documents.Query,
) (*mongo.Cursor, error) {
	ids, it := q.MultiVectorSearchQuery()

	a.logger.Infof("[multi vector] ids=%v item_type=%d", ids, it)

	if err := it.Ok(); err != nil {
		return nil, err
	}
	if it == documents.ITEM_TYPE_UNSPECIFIED {
		return nil, errors.New("vector search requires item type")
	}

	l := len(ids)
	if l <= 1 {
		return nil, errors.New("invalid ids length")
	}

	oids := make([]bson.ObjectID, l)
	for i := range l {
		oid, err := bson.ObjectIDFromHex(ids[i])
		if err != nil {
			return nil, err
		}

		oids[i] = oid
	}

	vs, err := a.queryMultiVector(ctx, oids)
	if err != nil {
		return nil, err
	}

	v, err := a.mixVectors(vs)
	if err != nil {
		return nil, err
	}

	return a.unsafeVectorSearch(ctx, it, v)
}

func (a *atlasClient) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
