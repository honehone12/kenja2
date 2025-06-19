package main

import (
	"kenja2"
	"kenja2/engine"
	kenja2lambda "kenja2/lambda"
	"kenja2/lambda/logs"

	"github.com/aws/aws-lambda-go/lambda"
)

var __ENGINE engine.Engine

func main() {
	logs.FmtDefault()

	var err error
	__ENGINE, err = kenja2.ConnectAtlas(
		kenja2.NewMsgPack(),
		kenja2.NewMsgPack(),
	)
	if err != nil {
		logs.Fatal(err)
	}

	lambda.Start(kenja2lambda.Handler(__ENGINE.VectorSeach))
}
