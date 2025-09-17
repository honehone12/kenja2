package main

import (
	"context"
	"kenja2"
	"kenja2/documents"
	"kenja2/engine"
	"time"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const ENGINE_MAX_CONN = 3
const REQUEST_BODY_LIMIT = 1024
const REQUEST_TIME_LIMIT = time.Second * 10

var __engine engine.BufferingEngine

func responseBadRequest(e echo.Context) error {
	return e.String(http.StatusBadRequest, "bad request")
}

func responseInternalError(e echo.Context) error {
	return e.String(http.StatusInternalServerError, "internal error")
}

func handler(e echo.Context) error {
	if __engine == nil {
		e.Logger().Error("engine is not initialized")
		return responseInternalError(e)
	}

	req := e.Request()

	if req.ContentLength > REQUEST_BODY_LIMIT {
		e.Logger().Warn("content length over limit")
		return responseBadRequest(e)
	}

	contentType := req.Header.Get("Content-Type")
	if len(contentType) == 0 || contentType != "application/json; charset=utf-8" {
		e.Logger().Warn("unexpected content type header")
		return responseBadRequest(e)
	}

	q := documents.SearchQuery{}
	if err := e.Bind(&q); err != nil {
		e.Logger().Warn(err)
		return responseBadRequest(e)
	}

	ctx, cancel := context.WithTimeout(req.Context(), REQUEST_TIME_LIMIT)
	defer cancel()

	docs, err := kenja2.BufferedQuery(__engine, ctx, &q)
	if err != nil {
		e.Logger().Error(err)
		return responseInternalError(e)
	}

	return e.JSON(http.StatusOK, documents.SearchResult{Documents: docs})
}

func main() {
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

	e.POST("/anime-search", handler)

	if err := e.Start("localhost:8080"); err != nil {
		e.Logger.Error(err)
	}
}
