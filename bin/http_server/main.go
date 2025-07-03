package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
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

type handlerFunc = func(echo.Context) error
type bindFunc[Q documents.Query] = func(echo.Context) (Q, error)
type engineFunc[Q documents.Query] = func(context.Context, Q) ([]byte, error)

const REQUEST_BODY_LIMIT = 1024
const REQUEST_TIME_LIMIT = time.Second * 10

var __ENGINE engine.Engine

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

func handler[Q documents.Query](f engineFunc[Q], bind bindFunc[Q]) handlerFunc {
	return func(e echo.Context) error {
		req := e.Request()
		if err := checkBodyLimit(req); err != nil {
			e.Logger().Error(err)
			return responseBadRequest(e)
		}
		if err := checkHeaders(req); err != nil {
			e.Logger().Error(err)
			return responseBadRequest(e)
		}

		ctx, cancel := context.WithTimeout(req.Context(), REQUEST_TIME_LIMIT)
		defer cancel()

		q, err := bind(e)
		if err != nil {
			e.Logger().Error(err)
			return responseBadRequest(e)
		}

		b, err := f(ctx, q)
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
}

func args() uint {
	port := flag.Uint("port", 8080, "listen port")
	flag.Parse()
	return *port
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
		e.Logger,
	)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer func() {
		if err := __ENGINE.Close(ctx); err != nil {
			e.Logger.Error(err)
		}
	}()

	e.POST("/document", handler(
		__ENGINE.Document,
		func(e echo.Context) (*documents.DocumentQuery, error) {
			q := &documents.DocumentQuery{}
			if err := e.Bind(q); err != nil {
				return nil, err
			}
			return q, nil
		},
	))
	e.POST("/search/buffer/text", handler(
		__ENGINE.TextSearch,
		func(e echo.Context) (*documents.TextQuery, error) {
			q := &documents.TextQuery{}
			if err := e.Bind(q); err != nil {
				return nil, err
			}
			return q, nil
		},
	))
	e.POST("/search/buffer/vector", handler(
		__ENGINE.VectorSeach,
		func(e echo.Context) (*documents.VectorQuery, error) {
			q := &documents.VectorQuery{}
			if err := e.Bind(q); err != nil {
				return nil, err
			}
			return q, nil
		},
	))

	listenAt := fmt.Sprintf("localhost:%d", port)
	if err := e.Start(listenAt); err != nil {
		e.Logger.Error(err)
	}
}
