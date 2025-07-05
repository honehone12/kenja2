package engine

import (
	"context"
	"kenja2/documents"
)

type Engine interface {
	TextSearch(
		ctx context.Context,
		q *documents.SearchQuery,
	) ([]documents.FlatDocument, error)
	VectorSeach(
		ctx context.Context,
		q *documents.SearchQuery,
	) ([]documents.FlatDocument, error)
	Document(
		ctx context.Context,
		q *documents.SearchQuery,
	) ([]documents.FlatDocument, error)

	Close(ctx context.Context) error
}
