package engine

import (
	"context"
	"kenja2/documents"
)

type BufferingEngine interface {
	Ping(ctx context.Context) error
	TextSearch(
		ctx context.Context,
		q *documents.SearchQuery,
	) ([]documents.FlatDocument, error)
	VectorSeach(
		ctx context.Context,
		q *documents.SearchQuery,
	) ([]documents.FlatDocument, error)
	NewItems(ctx context.Context) ([]documents.FlatDocument, error)
	Close(ctx context.Context) error
}
