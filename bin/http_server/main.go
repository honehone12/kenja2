package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"kenja2"
	"kenja2/engine"
	"time"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const REQUEST_BODY_LIMIT = 1024
const REQUEST_TIME_LIMIT = time.Second * 3

var __ENGINE engine.Engine

func args() uint {
	port := flag.Uint("port", 8080, "listen port")
	flag.Parse()
	return *port
}

func responseBadRequest(e echo.Context) error {
	return e.String(http.StatusBadRequest, "bad request")
}

func responseInternalError(e echo.Context) error {
	return e.String(http.StatusInternalServerError, "internal error")
}

func checkBodyLimit(req *http.Request) error {
	if req.ContentLength > REQUEST_BODY_LIMIT {
		return errors.New("content length over limit")
	}
	return nil
}

func checkHeaders(req *http.Request) error {
	contentType := req.Header.Get("Content-Type")
	if len(contentType) == 0 || contentType != __ENGINE.Decoder().ContentType() {
		return errors.New("unexpected content type header")
	}
	return nil
}

func text(e echo.Context) error {
	req := e.Request()
	if err := checkBodyLimit(req); err != nil {
		e.Logger().Error(err)
		return responseBadRequest(e)
	}
	if err := checkHeaders(req); err != nil {
		e.Logger().Error(err)
		return responseBadRequest(e)
	}

	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_TIME_LIMIT)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		e.Logger().Error(err)
		return responseInternalError(e)
	}

	b, err := __ENGINE.TextSearch(ctx, body)
	if err != nil {
		e.Logger().Error(err)
		return responseInternalError(e)
	}

	return e.Blob(
		http.StatusOK,
		__ENGINE.Encoder().ContentType(),
		b,
	)
}

func vector(e echo.Context) error {
	req := e.Request()
	if err := checkBodyLimit(req); err != nil {
		e.Logger().Error(err)
		return responseBadRequest(e)
	}
	if err := checkHeaders(req); err != nil {
		e.Logger().Error(err)
		return responseBadRequest(e)
	}

	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_TIME_LIMIT)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		e.Logger().Error(err)
		return responseInternalError(e)
	}

	b, err := __ENGINE.VectorSeach(ctx, body)
	if err != nil {
		e.Logger().Error(err)
		return responseInternalError(e)
	}

	return e.Blob(
		http.StatusOK,
		__ENGINE.Encoder().ContentType(),
		b,
	)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetPrefix("KENJA2")

	err := godotenv.Load("../../.env")
	if err != nil {
		e.Logger.Fatal(err)
	}

	port := args()
	ctx := context.Background()
	__ENGINE, err = kenja2.ConnectAtlas(
		kenja2.NewJson(),
		kenja2.NewJson(),
	)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer func() {
		if err := __ENGINE.Close(ctx); err != nil {
			e.Logger.Error(err)
		}
	}()

	e.POST("/text", text)
	e.POST("/vector", vector)

	listenAt := fmt.Sprintf("localhost:%d", port)
	if err := e.Start(listenAt); err != nil {
		e.Logger.Error(err)
	}
}
