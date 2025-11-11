package engine

import (
	"context"
	"kenja2/documents"
)

type Stream[T any] interface {
	Next(ctx context.Context) bool
	Decode(val *T) error
	Err() error
	Close(ctx context.Context) error
}

type Engine interface {
	Ping(ctx context.Context) error
	Random(
		ctx context.Context, q documents.Query,
	) ([]documents.FlatDocument, error)
	TextSearch(
		ctx context.Context, q documents.Query,
	) ([]documents.FlatDocument, error)
	VectorSeach(
		ctx context.Context, q documents.Query,
	) ([]documents.FlatDocument, error)
	DirectVectorSearch(
		ctx context.Context, q documents.Query,
	) ([]documents.FlatDocument, error)
	MultiVectorSearch(
		ctx context.Context, q documents.Query,
	) ([]documents.FlatDocument, error)
	TextSearchStream(
		ctx context.Context, q documents.Query,
	) (Stream[documents.FlatDocument], error)
	VectorSearchStream(
		ctx context.Context, q documents.Query,
	) (Stream[documents.FlatDocument], error)
	DirectVectorSearchStream(
		ctx context.Context, q documents.Query,
	) (Stream[documents.FlatDocument], error)
	MultiVectorSearchStream(
		ctx context.Context, q documents.Query,
	) (Stream[documents.FlatDocument], error)
}
