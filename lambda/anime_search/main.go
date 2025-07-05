package main

import (
	"context"
	"errors"
	"kenja2"
	"kenja2/documents"
	"kenja2/engine"
	"kenja2/lambda/logs"

	"github.com/aws/aws-lambda-go/lambda"
)

const REQUEST_BODY_LIMIT = 1024

var __engine engine.BufferingEngine

func entry(
	ctx context.Context,
	q *documents.SearchQuery,
) ([]documents.FlatDocument, error) {
	if __engine == nil {
		return nil, errors.New("engine is not initialized")
	}

	switch q.FunctionId {
	case documents.DOCUMENT:
		return __engine.Document(ctx, q)
	case documents.SEARCH_BUFFER_TEXT:
		return __engine.TextSearch(ctx, q)
	case documents.SEARCH_BUFFER_VECTOR:
		return __engine.VectorSeach(ctx, q)
	default:
		return nil, errors.New("invalid function id")
	}
}

func Handler(
	ctx context.Context,
	q documents.SearchQuery,
) (documents.SearchResult, error) {
	r := documents.SearchResult{}

	size, ok := q.StringSize()
	if !ok || size > REQUEST_BODY_LIMIT {
		err := errors.New("invalid request size")
		logs.Warn(err)
		return r, err
	}

	docs, err := entry(ctx, &q)
	if err != nil {
		logs.Error(err)
		return r, err
	}

	r.Documents = docs
	return r, nil
}

func main() {
	logs.FmtDefault()
	logs.SetLevel(logs.INFO)

	e, err := kenja2.ConnectAtlasAws(logs.Logger())
	if err != nil {
		logs.Panic(err)
	}

	__engine = e

	lambda.Start(Handler)
}
