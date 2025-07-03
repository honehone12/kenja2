package engine

import (
	"context"
	"kenja2/documents"
)

type Engine interface {
	TextSearch(
		ctx context.Context,
		q *documents.TextQuery,
	) (*documents.SearchResult, error)
	VectorSeach(
		ctx context.Context,
		q *documents.VectorQuery,
	) (*documents.SearchResult, error)
	Document(
		ctx context.Context,
		q *documents.DocumentQuery,
	) (*documents.SingleDocument, error)

	Close(ctx context.Context) error
}
