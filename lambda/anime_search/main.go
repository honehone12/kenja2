package main

import (
	"context"
	"errors"
	"kenja2/documents"
	kenjalambda "kenja2/lambda"
	"kenja2/lambda/logs"

	"github.com/aws/aws-lambda-go/lambda"
)

var __kenja *kenjalambda.Kenja

func Handler(
	ctx context.Context,
	q documents.SearchQuery,
) (documents.SearchResult, error) {
	r := documents.SearchResult{}

	if __kenja == nil {
		err := errors.New("kenja is not initialized")
		logs.Error(err)
		return r, err
	}

	docs, err := __kenja.Enter(ctx, &q)
	if err != nil {
		logs.Warn(err)
		return r, err
	}

	r.Documents = docs
	return r, nil
}

func main() {
	logs.FmtDefault()
	logs.SetLevel(logs.INFO)

	k, err := kenjalambda.NewKenja()
	if err != nil {
		logs.Panic(err)
	}

	__kenja = k

	lambda.Start(Handler)
}
