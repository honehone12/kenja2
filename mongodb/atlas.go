package mongodb

import (
	"context"
	"errors"
	"kenja2/documents"
	"kenja2/ed"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Atlas[E ed.Encoder, D ed.Decoder] struct {
	mongoClient
	collections collections
	encoder     E
	decoder     D
}

func Connet[E ed.Encoder, D ed.Decoder](
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

func (a *Atlas[E, D]) TextSearch(ctx context.Context, input []byte) ([]byte, error) {
	q := documents.TextQuery{}
	if err := a.decoder.Unmarshal(input, &q); err != nil {
		return nil, err
	}

	var c *mongo.Collection
	switch q.Rating {
	case documents.RATING_ALL_AGES:
		c = a.collections.allAges
	case documents.RATING_HENTAI:
		return nil, errors.New("the search service is not available yet")
	default:
		return nil, errors.New("unexpected rating")
	}

	k := CleanKeywords(q.Keywords)
	if len(k) == 0 {
		return nil, errors.New("invalid keywords")
	}

	p := bson.D{
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
		}},
		{Key: "$limit", Value: 1000},
		{Key: "$project", Value: bson.M{
			"text_vector":  0,
			"image_vector": 0,
		}},
	}
	op := options.Aggregate().
		SetMaxAwaitTime(time.Duration(time.Second))
	stream, err := c.Aggregate(ctx, p, op)
	if err != nil {
		return nil, err
	}

	var candidates []documents.Candidate
	if err = stream.All(ctx, &candidates); err != nil {
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
}

func (a *Atlas[E, D]) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
