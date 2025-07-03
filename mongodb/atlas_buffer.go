package mongodb

import (
	"context"
	"kenja2/documents"
	"kenja2/endec"
	"kenja2/engine"
)

type AtlasBuffer[E endec.Encoder, D endec.Decoder] struct {
	*atlasClient
	encoder E
	decoder D
}

func Connect[E endec.Encoder, D endec.Decoder](
	uri string,
	encoder E,
	decoder D,
	logger engine.Logger,
) (*AtlasBuffer[E, D], error) {
	mongoClient, err := connect(uri, logger)
	if err != nil {
		return nil, err
	}

	return &AtlasBuffer[E, D]{
		mongoClient,
		encoder,
		decoder,
	}, nil
}

func ConnectAws[E endec.Encoder, D endec.Decoder](
	uri string,
	encoder E,
	decoder D,
	logger engine.Logger,
) (*AtlasBuffer[E, D], error) {
	mongoClient, err := connectAws(uri, logger)
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

func (a *AtlasBuffer[E, D]) Document(
	ctx context.Context,
	q *documents.DocumentQuery,
) ([]byte, error) {
	res, err := a.document(ctx, q)
	if err != nil {
		return nil, err
	}

	d := documents.SingleDocument{}
	if err := res.Decode(&d); err != nil {
		return nil, err
	}

	d.Format()

	return a.encoder.Marshal(d)
}

func (a *AtlasBuffer[E, D]) TextSearch(
	ctx context.Context,
	q *documents.TextQuery,
) ([]byte, error) {
	stream, err := a.textSearch(ctx, q)
	if err != nil {
		return nil, err
	}

	var candidates []documents.Candidate
	if err := stream.All(ctx, &candidates); err != nil {
		return nil, err
	}

	l := len(candidates)
	a.logger.Infof("buffered %d candidates", l)

	for i := range l {
		candidates[i].Format()
	}

	return a.encoder.Marshal(documents.SearchResult{
		Candidates: candidates,
	})
}

func (a *AtlasBuffer[E, D]) VectorSeach(
	ctx context.Context,
	q *documents.VectorQuery,
) ([]byte, error) {
	stream, err := a.vectorSearch(ctx, q)
	if err != nil {
		return nil, err
	}

	var candidates []documents.Candidate
	if err := stream.All(ctx, &candidates); err != nil {
		return nil, err
	}

	l := len(candidates)
	a.logger.Infof("buffered %d candidates", l)

	for i := range l {
		candidates[i].Format()
	}

	return a.encoder.Marshal(documents.SearchResult{
		Candidates: candidates,
	})
}

func (a *AtlasBuffer[E, D]) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
