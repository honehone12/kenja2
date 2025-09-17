package engine

import (
	"context"
	"kenja2/documents"
)

type BufferingEngine interface {
	Ping(ctx context.Context) error
	Random(
		ctx context.Context,
		q documents.Query,
	) ([]documents.FlatDocument, error)
	TextSearch(
		ctx context.Context,
		q documents.Query,
	) ([]documents.FlatDocument, error)
	VectorSeach(
		ctx context.Context,
		q documents.Query,
	) ([]documents.FlatDocument, error)
}
