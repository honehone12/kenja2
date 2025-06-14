package main

import (
	"context"
	"kenja2"
	"kenja2/engine"
	kenja2lambda "kenja2/lambda"
	"kenja2/lambda/logs"

	"github.com/aws/aws-lambda-go/lambda"
)

var __ENGINE engine.Engine

func handler(ctx context.Context, req []byte) ([]byte, error) {
	if err := kenja2lambda.CheckBodyLimit(req); err != nil {
		logs.Error(err)
		return nil, err
	}

	res, err := __ENGINE.VectorSeach(ctx, req)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return res, nil
}

func main() {
	logs.FmtDefault()

	var err error
	__ENGINE, err = kenja2.ConnectAtlas(
		kenja2.NewJson(),
		kenja2.NewJson(),
	)
	if err != nil {
		logs.Fatal(err)
	}

	lambda.Start(handler)
}
