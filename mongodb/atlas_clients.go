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
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const ATLAS_SEARCH_SIZE = 100
const ATLAS_SEARCH_LIMIT = 100
const ATLAS_KW_LIMIT = 120
const ATLAS_CURSOR_MAX_AWAIT = time.Second
const FUZZY_MAX_EDIT_TO_MATCH = 1 // 1 or 2
const FUZZY_MAX_TRY_PER_KW = 100
const MAX_SINGLE_TKN_FUZZY = 2

var ErrFuzzySingleTknLimit = errors.New("keywords contain too many single tokens")

type AtlasClient struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	logger     engine.Logger

	replacer *strings.Replacer
}

func newAtlasClient(m *mongo.Client, l engine.Logger) (*AtlasClient, error) {
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

func DirectConnect(uri string, maxConn uint64, l engine.Logger) (*AtlasClient, error) {
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

func connectAws(format string, maxConn uint64, l engine.Logger) (*AtlasClient, error) {
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
	return DirectConnect(uri, maxConn, l)
}

func fieldFlags() bson.M {
	return bson.M{
		"item_type":     1,
		"rating":        1,
		"url":           1,
		"img":           1,
		"parent":        1,
		"name":          1,
		"name_japanese": 1,
		"name_english":  1,
	}
}

func animeITypeFilter() bson.M {
	return bson.M{
		"$in": []documents.ItemType{
			documents.ITEM_TYPE_ANIME,
			documents.ITEM_TYPE_ANIME_IMG_ONLY,
		},
	}
}

func characterITypeFilter() bson.M {
	return bson.M{
		"$in": []documents.ItemType{
			documents.ITEM_TYPE_CHARACTER,
			documents.ITEM_TYPE_CHARACTER_IMG_ONLY,
		},
	}
}

func animeAndCharacterITypeFilter() bson.M {
	return bson.M{
		"$in": []documents.ItemType{
			documents.ITEM_TYPE_ANIME,
			documents.ITEM_TYPE_ANIME_IMG_ONLY,
			documents.ITEM_TYPE_CHARACTER,
			documents.ITEM_TYPE_CHARACTER_IMG_ONLY,
		},
	}
}

func (a *AtlasClient) ping(ctx context.Context) error {
	a.logger.Info("ping")
	return a.client.Ping(ctx, nil)
}

func (a *AtlasClient) document(
	ctx context.Context,
	q *documents.SearchQuery,
) (*mongo.SingleResult, error) {
	id, err := bson.ObjectIDFromHex(q.Id)
	if err != nil {
		return nil, err
	}

	a.logger.Infof("[documet] q=%#v", q)

	f := bson.M{"_id": id}
	op := options.FindOne().SetProjection(fieldFlags())

	res := a.collection.FindOne(ctx, f, op)
	if res.Err() != nil {
		return nil, res.Err()
	}

	return res, nil
}

func (a *AtlasClient) UnsafeTextSearch(
	ctx context.Context,
	keywords keywords.Keywords,
	itemType documents.ItemType,
	rating documents.Rating,
	fuzzy documents.FuzzyKind,
) (*mongo.Cursor, error) {
	if rating == documents.RATING_UNSPECIFIED {
		rating = documents.RATING_ALL_AGES
	}

	text := bson.M{
		"query": keywords,
		"path": []string{
			"name",
			"name_english",
			"aliases",
			"description",
			"parent.name",
		},
		"matchCriteria": "all",
	}
	if fuzzy == documents.FUZZY_USE {
		text["fuzzy"] = bson.M{
			"maxEdits":      FUZZY_MAX_EDIT_TO_MATCH,
			"maxExpansions": FUZZY_MAX_TRY_PER_KW,
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

	if itemType.IsAnime() {
		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": animeITypeFilter(),
			}},
		})
	} else if itemType.IsCharacter() {
		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": characterITypeFilter(),
			}},
		})
	} else {
		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": animeAndCharacterITypeFilter(),
			}},
		})
	}

	p = append(p, bson.D{
		{Key: "$match", Value: bson.M{
			"rating": rating,
		}},
	})

	p = append(p,
		bson.D{{Key: "$limit", Value: ATLAS_SEARCH_LIMIT}},
		bson.D{{Key: "$project", Value: fieldFlags()}},
	)

	op := options.Aggregate().SetMaxAwaitTime(ATLAS_CURSOR_MAX_AWAIT)
	return a.collection.Aggregate(ctx, p, op)
}

func (a *AtlasClient) textSearch(
	ctx context.Context,
	q *documents.SearchQuery,
) (*mongo.Cursor, error) {
	a.logger.Infof("[text] q=%#v", q)

	if len(q.Keywords) > ATLAS_KW_LIMIT {
		return nil, errors.New("unexpected keywords length")
	}
	if err := q.ItemType.Ok(); err != nil {
		return nil, err
	}
	if err := q.Rating.Ok(); err != nil {
		return nil, err
	}
	if err := q.FuzzyKind.Ok(); err != nil {
		return nil, err
	}

	k, singleTkns := keywords.CleanKeywords(a.replacer, q.Keywords)
	if q.FuzzyKind == documents.FUZZY_USE && singleTkns > MAX_SINGLE_TKN_FUZZY {
		return nil, ErrFuzzySingleTknLimit
	}
	if len(k) == 0 {
		return nil, errors.New("invalid keywords")
	}

	return a.UnsafeTextSearch(
		ctx,
		k,
		q.ItemType,
		q.Rating,
		q.FuzzyKind,
	)
}

func (a *AtlasClient) UnsafeVectorSearch(
	ctx context.Context,
	itemType documents.ItemType,
	rating documents.Rating,
	vFieldName documents.VFieldName,
	srcVec bson.Binary,
) (*mongo.Cursor, error) {
	if rating == documents.RATING_UNSPECIFIED {
		rating = documents.RATING_ALL_AGES
	}

	var itemMap bson.M
	if itemType.IsAnime() {
		itemMap = animeITypeFilter()
	} else if itemType.IsCharacter() {
		itemMap = characterITypeFilter()
	} else {
		return nil, fmt.Errorf("unexpected item type: %d", itemType)
	}

	p := mongo.Pipeline{
		{
			{Key: "$vectorSearch", Value: bson.M{
				"exact": false,
				"filter": bson.M{
					"item_type": itemMap,
					"rating":    rating,
				},
				"index":         "vector",
				"limit":         ATLAS_SEARCH_LIMIT,
				"numCandidates": ATLAS_SEARCH_SIZE * ATLAS_SEARCH_LIMIT,
				"path":          vFieldName,
				"queryVector":   srcVec,
			}},
		},
		{
			{Key: "$project", Value: fieldFlags()},
		},
	}

	op := options.Aggregate().SetMaxAwaitTime(ATLAS_CURSOR_MAX_AWAIT)
	return a.collection.Aggregate(ctx, p, op)
}

func (a *AtlasClient) queryVector(
	ctx context.Context,
	id bson.ObjectID,
	vField documents.VectorField,
	vFieldName documents.VFieldName,
) (bson.Binary, error) {
	f := bson.D{{Key: "_id", Value: id}}
	op := options.FindOne().SetProjection(bson.M{
		string(vFieldName): 1,
	})
	res := a.collection.FindOne(ctx, f, op)
	if res.Err() != nil {
		return bson.Binary{}, res.Err()
	}
	v := documents.Vector{}
	if err := res.Decode(&v); err != nil {
		return bson.Binary{}, err
	}

	b, err := v.Binary(vField)
	if err != nil {
		return bson.Binary{}, err
	}

	return b, nil
}

func (a *AtlasClient) vectorSearch(
	ctx context.Context,
	q *documents.SearchQuery,
) (*mongo.Cursor, error) {
	a.logger.Infof("[vector] q=%#v", q)

	if q.ItemType == documents.ITEM_TYPE_UNSPECIFIED {
		return nil, errors.New("vector search requires item type")
	}
	if err := q.ItemType.Ok(); err != nil {
		return nil, err
	}
	if err := q.Rating.Ok(); err != nil {
		return nil, err
	}

	id, err := bson.ObjectIDFromHex(q.Id)
	if err != nil {
		return nil, err
	}
	vFieldName, err := q.VectorField.Name()
	if err != nil {
		return nil, err
	}

	v, err := a.queryVector(
		ctx,
		id,
		q.VectorField,
		vFieldName,
	)
	if err != nil {
		return nil, err
	}

	return a.UnsafeVectorSearch(ctx, q.ItemType, q.Rating, vFieldName, v)
}

func (a *AtlasClient) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
