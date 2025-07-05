package lambda

import (
	"context"
	"errors"
	"kenja2"
	"kenja2/documents"
	"kenja2/engine"
	"kenja2/lambda/logs"
)

const REQUEST_BODY_LIMIT = 1024

type Kenja struct {
	engine engine.BufferingEngine
}

func NewKenja() (*Kenja, error) {
	engine, err := kenja2.ConnectAtlasAws(logs.Logger())
	if err != nil {
		return nil, err
	}

	return &Kenja{engine}, nil
}

func (k *Kenja) Enter(
	ctx context.Context,
	q *documents.SearchQuery,
) ([]documents.FlatDocument, error) {
	size, ok := q.StringSize()
	if !ok || size > REQUEST_BODY_LIMIT {
		return nil, errors.New("invalid request size")
	}

	switch q.FunctionId {
	case documents.DOCUMENT:
		return k.engine.Document(ctx, q)
	case documents.SEARCH_BUFFER_TEXT:
		return k.engine.TextSearch(ctx, q)
	case documents.SEARCH_BUFFER_VECTOR:
		return k.engine.VectorSeach(ctx, q)
	default:
		return nil, errors.New("invalid function id")
	}
}
