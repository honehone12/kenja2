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
const ENGINE_MAX_CONN = 1

var __engine engine.BufferingEngine

func Handler(
	ctx context.Context,
	q documents.SearchQuery,
) (documents.SearchResult, error) {
	r := documents.SearchResult{}

	if __engine == nil {
		return r, errors.New("engine is not initialized")
	}

	size, ok := q.StringSize()
	if !ok || size > REQUEST_BODY_LIMIT {
		err := errors.New("invalid request size")
		logs.Warn(err)
		return r, err
	}

	docs, err := kenja2.BufferedQuery(__engine, ctx, &q)
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

	e, err := kenja2.ConnectAtlasAws(ENGINE_MAX_CONN, logs.Logger())
	if err != nil {
		logs.Panic(err)
	}

	__engine = e

	lambda.Start(Handler)
}
