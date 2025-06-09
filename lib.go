package kenja2

import (
	"context"
	"kenja2/marshalers"
)

type Engine[E, D marshalers.Marshaler] interface {
	TextSearch(ctx context.Context, input []byte) ([]byte, error)
	VectorSeach(ctx context.Context)
	Close(ctx context.Context) error
}
