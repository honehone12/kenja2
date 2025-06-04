package kenja2

import (
	"context"
	"kenja2/lib/documents"
)

type Engine interface {
	TextSearch(keywords string, rating documents.Rating)
	Close(ctx context.Context) error
}
