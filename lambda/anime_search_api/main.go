package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"kenja2"
	"kenja2/documents"
	"kenja2/engine"
	"kenja2/lambda/logs"
	"kenja2/lambda/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const REQUEST_BODY_LIMIT = 2048
const ENGINE_MAX_CONN = 1

var __engine engine.Engine

func handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	if __engine == nil {
		logs.Error("engine is not initialized")
		return response.InternalError(), nil
	}

	q := documents.SearchQuery{}
	{
		var b []byte
		if req.IsBase64Encoded {
			var err error
			b, err = base64.StdEncoding.DecodeString(req.Body)
			if err != nil {
				logs.Warn("could not decode base64 request body")
				return response.BadRequest(), nil
			}
		} else {
			b = []byte(req.Body)
		}

		err := json.Unmarshal(b, &q)
		if err != nil {
			logs.Warn("failed to unmarshal request body")
			return response.BadRequest(), nil
		}
	}

	size, ok := q.StringSize()
	if !ok || size > REQUEST_BODY_LIMIT {
		logs.Warn("invalid request body size")
		return response.BadRequest(), nil
	}

	docs, err := kenja2.BufferedQuery(__engine, ctx, &q)
	if err != nil {
		logs.Error(err)
		return response.InternalError(), nil
	}

	r := documents.SearchResult{Documents: docs}
	b, err := json.Marshal(r)
	if err != nil {
		logs.Error("failed to marshal request body")
		return response.BadRequest(), nil
	}
	return response.Json(string(b)), nil
}

func main() {
	logs.FmtDefault()
	logs.SetLevel(logs.INFO)

	var err error
	__engine, err = kenja2.ConnectAtlasAws(ENGINE_MAX_CONN, logs.Logger())
	if err != nil {
		logs.Panic(err)
	}

	lambda.Start(handler)
}
