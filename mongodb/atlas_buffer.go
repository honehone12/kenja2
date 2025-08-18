package mongodb

import (
	"context"
	"errors"
	"kenja2/documents"
	"kenja2/engine"
)

type AtlasBuffer struct {
	*atlasClient
}

func Connect(uri string, maxConn uint64, logger engine.Logger) (*AtlasBuffer, error) {
	atlasClient, err := DirectConnect(uri, maxConn, logger)
	if err != nil {
		return nil, err
	}

	return &AtlasBuffer{atlasClient}, nil
}

func ConnectAws(uri string, maxConn uint64, logger engine.Logger) (*AtlasBuffer, error) {
	atlasClient, err := connectAws(uri, maxConn, logger)
	if err != nil {
		return nil, err
	}

	return &AtlasBuffer{atlasClient}, nil
}

func (a *AtlasBuffer) Ping(ctx context.Context) error {
	return a.ping(ctx)
}

func (a *AtlasBuffer) Document(
	ctx context.Context,
	q *documents.SearchQuery,
) ([]documents.FlatDocument, error) {
	res, err := a.document(ctx, q)
	if err != nil {
		return nil, err
	}

	d := documents.FlatDocument{}
	if err := res.Decode(&d); err != nil {
		return nil, err
	}

	return []documents.FlatDocument{d}, nil
}

func (a *AtlasBuffer) TextSearch(
	ctx context.Context,
	q *documents.SearchQuery,
) ([]documents.FlatDocument, error) {
	stream, err := a.textSearch(ctx, q)
	if errors.Is(err, ErrFuzzySingleTknLimit) {
		a.logger.Error(err)
		return []documents.FlatDocument{}, nil
	} else if err != nil {
		return nil, err
	}

	var docs []documents.FlatDocument
	if err := stream.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (a *AtlasBuffer) VectorSeach(
	ctx context.Context,
	q *documents.SearchQuery,
) ([]documents.FlatDocument, error) {
	stream, err := a.vectorSearch(ctx, q)
	if err != nil {
		return nil, err
	}

	var docs []documents.FlatDocument
	if err := stream.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (a *AtlasBuffer) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
