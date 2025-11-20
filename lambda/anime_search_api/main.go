package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
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

func logHeaders(h map[string]string) error {
	var err error

	userAgent, ok := h["user-agent"]
	if !ok {
		errors.Join(err, errors.New("empty user agent header"))
	}
	referer, ok := h["referer"]
	if !ok {
		errors.Join(err, errors.New("empty referer header"))
	}
	host, ok := h["host"]
	if !ok {
		errors.Join(err, errors.New("empty host header"))
	}
	origin := h["origin"]
	forwarded, ok := h["x-forwarded-for"]
	if !ok {
		errors.Join(err, errors.New("empty forwarded header"))
	}
	address, ok := h["cloudfront-viewer-address"]
	if !ok {
		errors.Join(err, errors.New("empty address header"))
	}
	asn, ok := h["cloudfront-viewer-asn"]
	if !ok {
		errors.Join(err, errors.New("empty asn header"))
	}
	country, ok := h["cloudfront-viewer-country"]
	if !ok {
		errors.Join(err, errors.New("empty country header"))
	}
	city := h["cloudfront-viewer-city"]

	logs.Infof(
		"user-agent=%s referer=%s host=%s origin=%s forwarded=%s address=%s asn=%s country=%s city=%s",
		userAgent,
		referer,
		host,
		origin,
		forwarded,
		address,
		asn,
		country,
		city,
	)

	return err
}

func handler(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	if __engine == nil {
		logs.Error("engine is not initialized")
		return response.InternalError(), nil
	}

	if err := logHeaders(req.Headers); err != nil {
		logs.Warn(err)
		return response.BadRequest(), nil
	}

	q := documents.SearchQuery{}
	{
		p, ok := req.QueryStringParameters["payload"]
		if !ok {
			logs.Warn("request without payload")
			return response.BadRequest(), nil
		}

		b, err := base64.RawURLEncoding.DecodeString(p)
		if err != nil {
			logs.Warn("failed to decode base64 payload")
			return response.BadRequest(), nil
		}

		logs.Infof("payload=%s", string(b))

		err = json.Unmarshal(b, &q)
		if err != nil {
			logs.Warn("failed to unmarshal json payload")
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
		logs.Error("failed to marshal json response body")
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
