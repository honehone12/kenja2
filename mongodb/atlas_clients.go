package mongodb

import (
	"context"
	"errors"
	"kenja2/documents"
	"kenja2/engine"
	"os"
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
}

func connect(uri string, logger engine.Logger) (atlasClient, error) {
	a := atlasClient{}

	db := os.Getenv("MONGO_DB_NAME")
	if len(db) == 0 {
		return a, errors.New("env for db name is not set")
	}

	ops := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ops)
	if err != nil {
		return a, err
	}

	database := client.Database(db)

	coll := os.Getenv("MONGO_COLLECTION")
	if len(coll) == 0 {
		return a, errors.New("env for collection is not set")
	}

	collection := database.Collection(coll)

	a.client = client
	a.database = database
	a.collection = collection
	a.logger = logger
	return a, nil
}

func (a *atlasClient) document(
	ctx context.Context,
	q documents.DocumentQuery,
) (*mongo.SingleResult, error) {
	id, err := bson.ObjectIDFromHex(q.Id)
	if err != nil {
		return nil, err
	}

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
	q documents.TextQuery,
) (*mongo.Cursor, error) {
	if len(q.Keywords) > ATLAS_KW_LIMIT {
		return nil, errors.New("unexpected keywords length")
	}

	k := CleanKeywords(q.Keywords)
	if len(k) == 0 {
		return nil, errors.New("invalid keywords")
	}

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
		i, err := q.ItemType.I32()
		if err != nil {
			return nil, err
		}

		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": i,
			}},
		})
	}

	if q.Rating != documents.RATING_UNSPECIFIED {
		r, err := q.Rating.I32()
		if err != nil {
			return nil, err
		}

		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"rating": r,
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
	q documents.VectorQuery,
) (*mongo.Cursor, error) {
	if q.ItemType == documents.ITEM_TYPE_UNSPECIFIED {
		return nil, errors.New("vector search requires item type")
	}
	if q.Rating == documents.RATING_UNSPECIFIED {
		return nil, errors.New("vector search requires rating")
	}

	i, err := q.ItemType.I32()
	if err != nil {
		return nil, err
	}

	r, err := q.Rating.I32()
	if err != nil {
		return nil, err
	}

	var srcVec bson.Binary
	{
		id, err := bson.ObjectIDFromHex(q.Id)
		if err != nil {
			return nil, err
		}

		src, err := q.SourceField.String()
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

		srcVec, err = v.Binary(q.SourceField)
		if err != nil {
			return nil, err
		}
	}

	targetField, err := q.TargetField.String()
	if err != nil {
		return nil, err
	}

	p := mongo.Pipeline{
		{
			{Key: "$vectorSearch", Value: bson.M{
				"exact": false,
				"filter": bson.M{
					"item_type": i,
					"rating":    r,
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
