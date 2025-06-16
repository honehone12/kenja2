package mongodb

import (
	"context"
	"kenja2/documents"
	"kenja2/endec"
)

type AtlasBuffer[E endec.Encoder, D endec.Decoder] struct {
	atlasClient
	encoder E
	decoder D
}

func Connect[E endec.Encoder, D endec.Decoder](
	uri string,
	encoder E,
	decoder D,
) (*AtlasBuffer[E, D], error) {
	mongoClient, err := connect(uri)
	if err != nil {
		return nil, err
	}

	return &AtlasBuffer[E, D]{
		mongoClient,
		encoder,
		decoder,
	}, nil
}

func (a *AtlasBuffer[E, D]) Encoder() endec.Encoder {
	return a.encoder
}

func (a *AtlasBuffer[E, D]) Decoder() endec.Decoder {
	return a.decoder
}

func (a *AtlasBuffer[E, D]) TextSearch(ctx context.Context, input []byte) ([]byte, error) {
	q := documents.TextQuery{}
	if err := a.decoder.Unmarshal(input, &q); err != nil {
		return nil, err
	}

	stream, err := a.textSearch(ctx, q)
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

func (a *AtlasBuffer[E, D]) VectorSeach(ctx context.Context, input []byte) ([]byte, error) {
	q := documents.VectorQuery{}
	if err := a.decoder.Unmarshal(input, &q); err != nil {
		return nil, err
	}

	stream, err := a.vectorSearch(ctx, q)
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

func (a *AtlasBuffer[E, D]) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
