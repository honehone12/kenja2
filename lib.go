package kenja2

import (
	"context"
	"kenja2/marshalers"
)

type Engine[E, D marshalers.Marshaler] interface {
	TextSearch(input []byte) ([]byte, error)
	Close(ctx context.Context) error
}
