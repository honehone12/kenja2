package engine

import (
	"context"
	"kenja2/endec"
)

type Engine[E endec.Encoder, D endec.Decoder] interface {
	TextSearch(ctx context.Context, input []byte) ([]byte, error)
	VectorSeach(ctx context.Context, input []byte) ([]byte, error)
	Encoder() endec.Encoder
	Decoder() endec.Decoder
	Close(ctx context.Context) error
}
