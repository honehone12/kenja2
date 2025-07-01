package main

import (
	"kenja2/engine"
	kenja2lambda "kenja2/lambda"
	"kenja2/lambda/logs"

	"github.com/aws/aws-lambda-go/lambda"
)

var __ENGINE engine.Engine

func main() {
	en, err := kenja2lambda.Init()
	if err != nil {
		logs.Fatal(err)
	}

	__ENGINE = en
	lambda.Start(kenja2lambda.Handler(__ENGINE.Document))
}
