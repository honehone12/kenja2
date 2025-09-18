package mongodb

import (
	"context"
	"errors"
	"fmt"
	"kenja2/documents"
	"kenja2/engine"
	"kenja2/keywords"
	"net/url"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const ATLAS_KW_LIMIT = 240 //js accecpts 120 -> 120 * 2bytes
const ATLAS_SEARCH_LIMIT = 120
const ATLAS_SEARCH_SIZE = 50
const FUZZY_MAX_EDIT_TO_MATCH = 1 // 1 or 2
const FUZZY_MAX_TRY_PER_KW = 50
const FUZZY_NO_CHNGE_PREFIX = 1
const MAX_SINGLE_TKN_FUZZY = 3

var ErrFuzzySingleTknLimit = errors.New("keywords contain too many single tokens")

type AtlasClient struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	logger     engine.Logger

	replacer *strings.Replacer
}

func newAtlasClient(
	m *mongo.Client, l engine.Logger,
) (*AtlasClient, error) {
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

	return &AtlasClient{
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
) (*AtlasClient, error) {
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
) (*AtlasClient, error) {
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

func (a *AtlasClient) ping(ctx context.Context) error {
	a.logger.Info("[ping]")

	return a.client.Ping(ctx, nil)
}

func (a *AtlasClient) random(
	ctx context.Context, q documents.Query,
) (*mongo.Cursor, error) {
	it := q.RandomQuery()

	a.logger.Infof("[random] item_type=%d", it)

	if err := it.Ok(); err != nil {
		return nil, err
	}

	// add item type query here

	p := mongo.Pipeline{
		{
			{Key: "$sample", Value: bson.M{
				"size": ATLAS_SEARCH_LIMIT,
			}},
		},
	}

	op := options.Aggregate()
	return a.collection.Aggregate(ctx, p, op)
}

func (a *AtlasClient) unsafeTextSearch(
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

func (a *AtlasClient) textSearch(
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

	k, singleTkns := keywords.CleanKeywords(a.replacer, kw)
	if fz == documents.FUZZY_USE && singleTkns > MAX_SINGLE_TKN_FUZZY {
		return nil, ErrFuzzySingleTknLimit
	}
	if len(k) == 0 {
		return nil, errors.New("invalid keywords")
	}

	return a.unsafeTextSearch(
		ctx,
		k,
		it,
		fz,
	)
}

func (a *AtlasClient) unsafeVectorSearch(
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

func (a *AtlasClient) queryVector(
	ctx context.Context, id bson.ObjectID,
) (bson.Binary, error) {
	var b bson.Binary

	f := bson.D{{Key: "_id", Value: id}}
	op := options.FindOne().SetProjection(bson.M{
		"image_vector": 1,
	})
	res := a.collection.FindOne(ctx, f, op)
	if res.Err() != nil {
		return b, res.Err()
	}

	v := documents.ImageVector{}
	err := res.Decode(&v)
	if err != nil {
		return b, err
	}
	b = v.ImageVector.Binary()
	if b.IsZero() {
		return b, errors.New("image vector is empty")
	}

	return b, nil
}

func (a *AtlasClient) vectorSearch(
	ctx context.Context, q documents.Query,
) (*mongo.Cursor, error) {
	id, it := q.VectorSearchQuery()

	a.logger.Infof("[vector] id=%s item_type=%d", id, it)

	if err := it.Ok(); err != nil {
		return nil, err
	}

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if it == documents.ITEM_TYPE_UNSPECIFIED {
		return nil, errors.New("vector search requires item type")
	}

	v, err := a.queryVector(ctx, oid)
	if err != nil {
		return nil, err
	}

	return a.unsafeVectorSearch(ctx, it, v)
}

func (a *AtlasClient) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
