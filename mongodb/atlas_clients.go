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
const ATLAS_KW_LIMIT = 100
const ATLAS_CURSOR_MAX_AWAIT = time.Second

type atlasClient struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	logger     engine.Logger

	replacer *strings.Replacer
}

func newAtlasClient(m *mongo.Client, l engine.Logger) (*atlasClient, error) {
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
		replacer:   __MONGO_REPLACER,
	}, nil
}

func connect(uri string, maxConn uint64, l engine.Logger) (*atlasClient, error) {
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

func connectAws(format string, maxConn uint64, l engine.Logger) (*atlasClient, error) {
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

func (a *atlasClient) ping(ctx context.Context) error {
	a.logger.Info("ping")
	return a.client.Ping(ctx, nil)
}

func (a *atlasClient) document(
	ctx context.Context,
	q *documents.SearchQuery,
) (*mongo.SingleResult, error) {
	id, err := bson.ObjectIDFromHex(q.Id)
	if err != nil {
		return nil, err
	}

	a.logger.Infof("[documet] q=%#v", q)

	f := bson.M{"_id": id}
	op := options.FindOne().SetProjection(bson.M{
		"text_vector":  0,
		"image_vector": 0,
		"rating":       0,
		"staff":        0,
		"description":  0,
	})

	res := a.collection.FindOne(ctx, f, op)
	if res.Err() != nil {
		return nil, res.Err()
	}

	return res, nil
}

func (a *atlasClient) _textSearch(
	ctx context.Context,
	keywords keywords.Keywords,
	itemType documents.ItemType,
	rating documents.Rating,
) (*mongo.Cursor, error) {
	p := mongo.Pipeline{
		{
			{Key: "$search", Value: bson.M{
				"index": "text",
				"text": bson.M{
					"query": keywords,
					"path": []string{
						"name",
						"name_english",
						"aliases",
						"studios",
						"staff",
						"description",
						"parent.name",
					},
					"matchCriteria": "all",
				},
				"returnStoredSource": true,
			}},
		},
	}

	if itemType != documents.ITEM_TYPE_UNSPECIFIED {
		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": itemType,
			}},
		})
	}

	if rating != documents.RATING_UNSPECIFIED {
		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"rating": rating,
			}},
		})
	}

	p = append(p,
		bson.D{{Key: "$limit", Value: ATLAS_SEARCH_LIMIT}},
		bson.D{{Key: "$project", Value: bson.M{
			"item_type":     1,
			"url":           1,
			"parent":        1,
			"name":          1,
			"name_japanese": 1,
			"name_english":  1,
		}}},
	)

	op := options.Aggregate().SetMaxAwaitTime(ATLAS_CURSOR_MAX_AWAIT)
	return a.collection.Aggregate(ctx, p, op)
}

func (a *atlasClient) textSearch(
	ctx context.Context,
	q *documents.SearchQuery,
) (*mongo.Cursor, error) {
	if len(q.Keywords) > ATLAS_KW_LIMIT {
		return nil, errors.New("unexpected keywords length")
	}
	if err := q.ItemType.Ok(); err != nil {
		return nil, err
	}
	if err := q.Rating.Ok(); err != nil {
		return nil, err
	}

	k := keywords.CleanKeywords(a.replacer, q.Keywords)
	if len(k) == 0 {
		return nil, errors.New("invalid keywords")
	}

	a.logger.Infof("[text] q=%#v", q)

	return a._textSearch(
		ctx,
		k,
		q.ItemType,
		q.Rating,
	)
}

func (a *atlasClient) _vectorSearch(
	ctx context.Context,
	itemType documents.ItemType,
	rating documents.Rating,
	id bson.ObjectID,
	vField documents.VectorField,
	vFieldName documents.VFieldName,
) (*mongo.Cursor, error) {
	if rating == documents.RATING_UNSPECIFIED {
		rating = documents.RATING_ALL_AGES
	}

	var srcVec bson.Binary
	{
		f := bson.D{{Key: "_id", Value: id}}
		op := options.FindOne().SetProjection(bson.M{
			string(vFieldName): 1,
		})
		res := a.collection.FindOne(ctx, f, op)
		if res.Err() != nil {
			return nil, res.Err()
		}
		v := documents.Vector{}
		if err := res.Decode(&v); err != nil {
			return nil, err
		}

		src, err := v.Binary(vField)
		if err != nil {
			return nil, err
		}
		srcVec = src
	}

	p := mongo.Pipeline{
		{
			{Key: "$vectorSearch", Value: bson.M{
				"exact": false,
				"filter": bson.M{
					"item_type": itemType,
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
			{Key: "$project", Value: bson.M{
				"item_type":     1,
				"url":           1,
				"parent":        1,
				"name":          1,
				"name_japanese": 1,
				"name_english":  1,
			}},
		},
	}

	op := options.Aggregate().SetMaxAwaitTime(ATLAS_CURSOR_MAX_AWAIT)
	return a.collection.Aggregate(ctx, p, op)
}

func (a *atlasClient) vectorSearch(
	ctx context.Context,
	q *documents.SearchQuery,
) (*mongo.Cursor, error) {
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

	a.logger.Infof("[vector] q=%#v", q)

	return a._vectorSearch(
		ctx,
		q.ItemType,
		q.Rating,
		id,
		q.VectorField,
		vFieldName,
	)
}
