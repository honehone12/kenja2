package kenja2

import (
	"context"
	"kenja2/ed"
)

type Engine[E ed.Encoder, D ed.Decoder] interface {
	TextSearch(ctx context.Context, input []byte) ([]byte, error)
	VectorSeach(ctx context.Context, input []byte) ([]byte, error)
	Close(ctx context.Context) error
}
