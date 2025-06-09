package kenja2

import (
	"context"
	"kenja2/endec"
)

type Engine[E endec.Encoder, D endec.Decoder] interface {
	TextSearch(ctx context.Context, input []byte) ([]byte, error)
	VectorSeach(ctx context.Context, input []byte) ([]byte, error)
	RequestContentType() string
	ResponseContentType() string
	Close(ctx context.Context) error
}
