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

type Atlas[E endec.Encoder, D endec.Decoder] struct {
	mongoClient
	collections collections
	encoder     E
	decoder     D
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

	collections, err := newCollections(mongoClient.database)
	if err != nil {
		return nil, err
	}

	return &Atlas[E, D]{
		mongoClient,
		collections,
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

func (a *Atlas[E, D]) collection(rating documents.Rating) (*mongo.Collection, error) {
	switch rating {
	case documents.RATING_ALL_AGES:
		return a.collections.allAges, nil
	case documents.RATING_HENTAI:
		return nil, errors.New("the search service is not available yet")
	default:
		return nil, errors.New("unexpected rating")
	}
}

func (a *Atlas[E, D]) TextSearch(ctx context.Context, input []byte) ([]byte, error) {
	q := documents.TextQuery{}
	if err := a.decoder.Unmarshal(input, &q); err != nil {
		return nil, err
	}

	c, err := a.collection(q.Rating)
	if err != nil {
		return nil, err
	}

	i, err := q.ItemType.To32()
	if err != nil {
		return nil, err
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

	if i != documents.ITEM_TYPE_UNSPECIFIED {
		p = append(p, bson.D{
			{Key: "$match", Value: bson.M{
				"item_type": i,
			}},
		})
	}

	p = append(p,
		bson.D{{Key: "$limit", Value: ATLAS_SEARCH_LIMIT}},
		bson.D{{Key: "$project", Value: bson.M{
			"text_vector":  0,
			"image_vector": 0,
		}}},
	)

	op := options.Aggregate().SetMaxAwaitTime(time.Second)
	stream, err := c.Aggregate(ctx, p, op)
	if err != nil {
		return nil, err
	}

	var candidates []documents.Candidate
	if err := stream.All(ctx, &candidates); err != nil {
		return nil, err
	}

	r := documents.QueryResult{
		Candidates: candidates,
	}
	b, err := a.encoder.Marshal(r)
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

	c, err := a.collection(q.Rating)
	if err != nil {
		return nil, err
	}

	i, err := q.ItemType.To32()
	if err != nil {
		return nil, err
	}
	if i == documents.ITEM_TYPE_UNSPECIFIED {
		return nil, errors.New("vector search requires item type")
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
		res := c.FindOne(ctx, f, op)
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
				},
				"index":         "vector",
				"limit":         ATLAS_SEARCH_LIMIT,
				"numCandidates": 20 * ATLAS_SEARCH_LIMIT,
				"path":          targetField,
				"queryVector":   srcVec,
			}},
			{Key: "$project", Value: bson.M{
				"text_vector":  0,
				"image_vector": 0,
			}},
		},
	}
	op := options.Aggregate().SetMaxAwaitTime(time.Second)
	stream, err := c.Aggregate(ctx, p, op)
	if err != nil {
		return nil, err
	}

	var candidates []documents.Candidate
	if err := stream.All(ctx, &candidates); err != nil {
		return nil, err
	}

	r := documents.QueryResult{
		Candidates: candidates,
	}
	b, err := a.encoder.Marshal(r)
	if err != nil {
		return nil, err
	}

	return b, err
}

func (a *Atlas[E, D]) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
