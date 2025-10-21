package mongodb

import (
	"context"
	"kenja2/documents"
	"kenja2/engine"
)

type Atlas struct {
	*atlasClient
}

func Connect(
	uri string,
	maxConn uint64,
	logger engine.Logger,
) (*Atlas, error) {
	atlasClient, err := connect(uri, maxConn, logger)
	if err != nil {
		return nil, err
	}

	return &Atlas{atlasClient}, nil
}

func ConnectAws(
	uri string,
	maxConn uint64,
	logger engine.Logger,
) (*Atlas, error) {
	atlasClient, err := connectAws(uri, maxConn, logger)
	if err != nil {
		return nil, err
	}

	return &Atlas{atlasClient}, nil
}

func (a *Atlas) Ping(ctx context.Context) error {
	return a.ping(ctx)
}

func (a *Atlas) Random(
	ctx context.Context, q documents.Query,
) ([]documents.FlatDocument, error) {
	stream, err := a.random(ctx, q)
	if err != nil {
		return nil, err
	}

	var docs []documents.FlatDocument
	if err := stream.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (a *Atlas) TextSearch(
	ctx context.Context, q documents.Query,
) ([]documents.FlatDocument, error) {
	stream, err := a.textSearch(ctx, q)
	if err != nil {
		return nil, err
	}

	var docs []documents.FlatDocument
	if err := stream.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (a *Atlas) VectorSeach(
	ctx context.Context, q documents.Query,
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

func (a *Atlas) TextSearchStream(
	ctx context.Context, q documents.Query,
) (engine.Stream[documents.FlatDocument], error) {
	stream, err := a.textSearch(ctx, q)
	if err != nil {
		return nil, err
	}

	return &Stream{stream}, nil
}

func (a *Atlas) VectorSearchStream(
	ctx context.Context, q documents.Query,
) (engine.Stream[documents.FlatDocument], error) {
	stream, err := a.vectorSearch(ctx, q)
	if err != nil {
		return nil, err
	}

	return &Stream{stream}, nil
}
