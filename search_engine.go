package kenja2

import "context"

type Engine interface {
	Close(ctx context.Context) error
}
