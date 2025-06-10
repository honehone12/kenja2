package mongodb

import (
	"context"
	"errors"
	"kenja2/documents"
	"kenja2/endec"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const ATLAS_SEARCH_LIMIT = 100
const ATLAS_KW_LIMIT = 100

type Atlas[E endec.Encoder, D endec.Decoder] struct {
	mongoClient
	encoder E
	decoder D
}

func Connect[E endec.Encoder, D endec.Decoder](
	uri string,
	encoder E,
	decoder D,
) (*Atlas[E, D], error) {
	mongoClient, err := connect(uri)
	if err != nil {
		return nil, err
	}

	return &Atlas[E, D]{
		mongoClient,
		encoder,
		decoder,
	}, nil
}

func (a *Atlas[E, D]) RequestContentType() string {
	return a.decoder.ContentType()
}

func (a *Atlas[E, D]) ResponseContentType() string {
	return a.encoder.ContentType()
}

func (a *Atlas[E, D]) TextSearch(ctx context.Context, input []byte) ([]byte, error) {
	q := documents.TextQuery{}
	if err := a.decoder.Unmarshal(input, &q); err != nil {
		return nil, err
	}
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
	)

	op := options.Aggregate().SetMaxAwaitTime(time.Second)
	stream, err := a.collection.Aggregate(ctx, p, op)
	if err != nil {
		return nil, err
	}

	var candidates []documents.Candidate
	if err := stream.All(ctx, &candidates); err != nil {
		return nil, err
	}

	b, err := a.encoder.Marshal(documents.QueryResult{
		Candidates: candidates,
	})
	if err != nil {
		return nil, err
	}

	return b, err
}

func (a *Atlas[E, D]) VectorSeach(ctx context.Context, input []byte) ([]byte, error) {
	q := documents.VectorQuery{}
	if err := a.decoder.Unmarshal(input, &q); err != nil {
		return nil, err
	}
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

		srcVec, err = v.BinaryField(q.SourceField)
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
				"numCandidates": 20 * ATLAS_SEARCH_LIMIT,
				"path":          targetField,
				"queryVector":   srcVec,
			}},
		},
		{
			{Key: "$project", Value: bson.M{
				"text_vector":  0,
				"image_vector": 0,
				"item_type":    0,
				"rating":       0,
				"description":  0,
			}},
		},
	}
	op := options.Aggregate().SetMaxAwaitTime(time.Second)
	stream, err := a.collection.Aggregate(ctx, p, op)
	if err != nil {
		return nil, err
	}

	var candidates []documents.Candidate
	if err := stream.All(ctx, &candidates); err != nil {
		return nil, err
	}

	b, err := a.encoder.Marshal(documents.QueryResult{
		Candidates: candidates,
	})
	if err != nil {
		return nil, err
	}

	return b, err
}

func (a *Atlas[E, D]) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
