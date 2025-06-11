package main

import (
	"context"
	"kenja2"
	"kenja2/endec"
	"kenja2/engine"
	"kenja2/lambda/logs"

	"github.com/aws/aws-lambda-go/lambda"
)

var __ENGINE engine.Engine[endec.Json, endec.Json]

func handler(
	ctx context.Context,
	req []byte,
) ([]byte, error) {

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
