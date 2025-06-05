package kenja2

import (
	"context"
	"kenja2/marshalers"
)

type Engine[M marshalers.Marshaler] interface {
	TextSearch(input []byte)
	Close(ctx context.Context) error
}
