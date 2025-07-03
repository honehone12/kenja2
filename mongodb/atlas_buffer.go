package mongodb

import (
	"context"
	"kenja2/documents"
	"kenja2/engine"
)

type AtlasBuffer struct {
	*atlasClient
}

func Connect(uri string, logger engine.Logger) (*AtlasBuffer, error) {
	mongoClient, err := connect(uri, logger)
	if err != nil {
		return nil, err
	}

	return &AtlasBuffer{atlasClient: mongoClient}, nil
}

func ConnectAws(uri string, logger engine.Logger) (*AtlasBuffer, error) {
	mongoClient, err := connectAws(uri, logger)
	if err != nil {
		return nil, err
	}

	return &AtlasBuffer{atlasClient: mongoClient}, nil
}

func (a *AtlasBuffer) Document(
	ctx context.Context,
	q *documents.DocumentQuery,
) (*documents.SingleDocument, error) {
	res, err := a.document(ctx, q)
	if err != nil {
		return nil, err
	}

	d := &documents.SingleDocument{}
	if err := res.Decode(d); err != nil {
		return nil, err
	}

	d.Format()
	return d, nil
}

func (a *AtlasBuffer) TextSearch(
	ctx context.Context,
	q *documents.TextQuery,
) (*documents.SearchResult, error) {
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

	return &documents.SearchResult{
		Candidates: candidates,
	}, nil
}

func (a *AtlasBuffer) VectorSeach(
	ctx context.Context,
	q *documents.VectorQuery,
) (*documents.SearchResult, error) {
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

	return &documents.SearchResult{
		Candidates: candidates,
	}, nil
}

func (a *AtlasBuffer) Close(ctx context.Context) error {
	return a.client.Disconnect(ctx)
}
