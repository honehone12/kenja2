package main

import (
	"context"
	"kenja2"
	"kenja2/documents"
	kenja "kenja2/lambda"
	"kenja2/lambda/logs"

	"github.com/aws/aws-lambda-go/lambda"
)

var __Kenja *kenja.Kenja

func Handler(ctx context.Context, q documents.SearchQuery) (documents.SearchResult, error) {
	docs, err := __Kenja.Enter(ctx, &q)
	if err != nil {
		logs.Warn(err)
		return documents.SearchResult{}, err
	}

	return documents.SearchResult{Documents: docs}, nil
}

func main() {
	logs.FmtDefault()
	logs.SetLevel(logs.INFO)

	engine, err := kenja2.ConnectAtlasAws(logs.Logger())
	if err != nil {
		logs.Panic(err)
	}

	__Kenja = kenja.NewKenja(engine)

	lambda.Start(Handler)
}
