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
		"staff_vector": 0,
		"item_type":    0,
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

func (a *atlasClient) textSearch(
	ctx context.Context,
	q *documents.SearchQuery,
) (*mongo.Cursor, error) {
	if len(q.Keywords) > ATLAS_KW_LIMIT {
		return nil, errors.New("unexpected keywords length")
	}

	k := keywords.CleanKeywords(a.replacer, q.Keywords)
	if len(k) == 0 {
		return nil, errors.New("invalid keywords")
	}

	a.logger.Infof("[text] q=%#v", q)

	p := mongo.Pipeline{
		{
			{Key: "$search", Value: bson.M{
				"index": "text",
				"text": bson.M{
					"query": k,
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

	if q.ItemType != documents.ITEM_TYPE_UNSPECIFIED {
		if err := q.ItemType.Ok(); err != nil {
			return nil, err
		}

		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": q.ItemType,
			}},
		})
	}

	if q.Rating != documents.RATING_UNSPECIFIED {
		if err := q.Rating.Ok(); err != nil {
			return nil, err
		}

		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"rating": q.Rating,
			}},
		})
	}

	p = append(p,
		bson.D{{Key: "$limit", Value: ATLAS_SEARCH_LIMIT}},
		bson.D{{Key: "$project", Value: bson.M{
			"parent":        1,
			"name":          1,
			"name_japanese": 1,
		}}},
	)

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

	a.logger.Infof("[vector] q=%#v", q)

	if q.Rating == documents.RATING_UNSPECIFIED {
		q.Rating = documents.RATING_ALL_AGES
	}

	var srcVec bson.Binary
	{
		id, err := bson.ObjectIDFromHex(q.Id)
		if err != nil {
			return nil, err
		}

		src, err := q.VectorField.String()
		if err != nil {
			return nil, err
		}

		f := bson.D{{Key: "_id", Value: id}}
		op := options.FindOne().SetProjection(bson.D{{Key: src, Value: 1}})
		res := a.collection.FindOne(ctx, f, op)
		if res.Err() != nil {
			return nil, res.Err()
		}
		v := documents.Vector{}
		if err := res.Decode(&v); err != nil {
			return nil, err
		}

		srcVec, err = v.Binary(q.VectorField)
		if err != nil {
			return nil, err
		}
	}

	targetField, err := q.VectorField.String()
	if err != nil {
		return nil, err
	}

	p := mongo.Pipeline{
		{
			{Key: "$vectorSearch", Value: bson.M{
				"exact": false,
				"filter": bson.M{
					"item_type": q.ItemType,
					"rating":    q.Rating,
				},
				"index":         "vector",
				"limit":         ATLAS_SEARCH_LIMIT,
				"numCandidates": ATLAS_SEARCH_SIZE * ATLAS_SEARCH_LIMIT,
				"path":          targetField,
				"queryVector":   srcVec,
			}},
		},
		{
			{Key: "$project", Value: bson.M{
				"parent":        1,
				"name":          1,
				"name_japanese": 1,
			}},
		},
	}

	op := options.Aggregate().SetMaxAwaitTime(ATLAS_CURSOR_MAX_AWAIT)
	return a.collection.Aggregate(ctx, p, op)
}
