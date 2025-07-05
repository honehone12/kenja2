package lambda

import (
	"context"
	"errors"
	"kenja2/documents"
	"kenja2/engine"
)

const REQUEST_BODY_LIMIT = 1024

type Kenja struct {
	engine engine.Engine
}

func NewKenja(engine engine.Engine) *Kenja {
	return &Kenja{engine}
}

func (k *Kenja) Enter(
	ctx context.Context,
	q *documents.SearchQuery,
) ([]documents.FlatDocument, error) {
	size, ok := q.StringSize()
	if !ok || size > REQUEST_BODY_LIMIT {
		return nil, errors.New("invalid request size")
	}

	return nil, nil
}
