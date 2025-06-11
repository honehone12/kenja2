package main

import (
	"context"
	"kenja2"
	"kenja2/endec"
	lib "kenja2/lambda"
	"kenja2/lambda/logs"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var __ENGINE kenja2.Engine[endec.Json, endec.Json]

func handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {

}

func main() {
	logs.FmtDefault()

	var err error
	__ENGINE, err = lib.ConnectAtlas[endec.Json, endec.Json]()
	if err != nil {
		logs.Fatal(err)
	}

	lambda.Start(handler)
}
