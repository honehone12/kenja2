package engine

import (
	"context"
	"kenja2/documents"
	"kenja2/endec"
)

type Engine interface {
	TextSearch(ctx context.Context, q *documents.TextQuery) ([]byte, error)
	VectorSeach(ctx context.Context, q *documents.VectorQuery) ([]byte, error)
	Document(ctx context.Context, q *documents.DocumentQuery) ([]byte, error)
	Encoder() endec.Encoder
	Decoder() endec.Decoder
	Close(ctx context.Context) error
}
