package main

import (
	"context"
	"encoding/json"
	"flag"
	"kenja2"
	"kenja2/documents"
	"kenja2/engine"
	"strings"
	"time"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type errorResponse = func(e echo.Context) error

const ENGINE_MAX_CONN = 3
const REQUEST_BODY_LIMIT = 2048
const REQUEST_TIME_LIMIT = time.Second * 10

var __engine engine.Engine

func responseBadRequest(e echo.Context) error {
	return e.String(http.StatusBadRequest, "bad request")
}

func responseInternalError(e echo.Context) error {
	return e.String(http.StatusInternalServerError, "internal error")
}

func check(e echo.Context) (*http.Request, errorResponse) {
	if __engine == nil {
		e.Logger().Error("engine is not initialized")
		return nil, responseInternalError
	}

	req := e.Request()

	if req.ContentLength > REQUEST_BODY_LIMIT {
		e.Logger().Warn("content length over limit")
		return nil, responseBadRequest
	}

	contentType := req.Header.Get("Content-Type")
	if len(contentType) == 0 || !strings.HasPrefix(contentType, "application/json") {
		e.Logger().Warn("unexpected content type header")
		return nil, responseBadRequest
	}

	return req, nil
}

func params(e echo.Context) (*documents.SearchQuery, errorResponse) {
	q := &documents.SearchQuery{}
	if err := e.Bind(q); err != nil {
		e.Logger().Warn(err)
		return nil, responseBadRequest
	}

	return q, nil
}

func animeSearch(e echo.Context) error {
	req, errRes := check(e)
	if errRes != nil {
		return errRes(e)
	}

	q, errRes := params(e)
	if errRes != nil {
		return errRes(e)
	}

	ctx, cancel := context.WithTimeout(req.Context(), REQUEST_TIME_LIMIT)
	defer cancel()

	docs, err := kenja2.BufferedQuery(__engine, ctx, q)
	if err != nil {
		e.Logger().Error(err)
		return responseInternalError(e)
	}

	return e.JSON(http.StatusOK, documents.SearchResult{Documents: docs})
}

func animeSearchStream(e echo.Context) error {
	req, errRes := check(e)
	if errRes != nil {
		return errRes(e)
	}

	q, errRes := params(e)
	if errRes != nil {
		return errRes(e)
	}

	ctx, cancel := context.WithTimeout(req.Context(), REQUEST_TIME_LIMIT)
	defer cancel()

	stream, err := kenja2.StreamQuery(__engine, ctx, q)
	if err != nil {
		e.Logger().Error(err)
		return responseInternalError(e)
	}
	defer stream.Close(ctx)

	res := e.Response()
	res.Header().Set("Content-Type", "application/json; charset=utf8")
	jsonEnc := json.NewEncoder(res)

	for stream.Next(ctx) {
		doc := documents.FlatDocument{}
		if err := stream.Decode(&doc); err != nil {
			return responseInternalError(e)
		}
		if err := jsonEnc.Encode(doc); err != nil {
			return responseInternalError(e)
		}

		res.Flush()
	}
	if err := stream.Err(); err != nil {
		return responseInternalError(e)
	}

	return nil
}

func main() {
	cors := flag.Bool("cors", false, "-cors")
	flag.Parse()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetPrefix("KENJA-ANIME-SEARCH")

	err := godotenv.Load("../../.env")
	if err != nil {
		e.Logger.Fatal(err)
	}

	__engine, err = kenja2.ConnectAtlas(ENGINE_MAX_CONN, e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}

	if *cors {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:5173"},
			AllowMethods: []string{http.MethodPost},
		}))
	}

	e.POST("/anime-search", animeSearch)
	e.POST("/anime-search/stream", animeSearchStream)

	if err := e.Start("localhost:8080"); err != nil {
		e.Logger.Error(err)
	}
}
