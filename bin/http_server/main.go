package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"kenja2"
	"kenja2/endec"
	"kenja2/mongodb"
	"time"

	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const REQUEST_BODY_LIMIT = 1024

var __ENGINE kenja2.Engine[endec.Json, endec.Json]

func args() uint {
	port := flag.Uint("port", 8080, "listen port")
	flag.Parse()
	return *port
}

func env() (string, error) {
	godotenv.Load("../../.env")
	engine_uri := os.Getenv("SEARCHENGINE_URI")
	if len(engine_uri) == 0 {
		return "", nil
	}

	return engine_uri, nil
}

func text(e echo.Context) error {
	req := e.Request()
	if req.ContentLength > REQUEST_BODY_LIMIT {
		e.Logger().Error("content length over limit")
		return e.String(http.StatusBadRequest, "bad request")
	}
	contentType := req.Header.Get("Content-Type")
	if len(contentType) == 0 || contentType != __ENGINE.RequestContentType() {
		e.Logger().Error("unexpected content type header")
		return e.String(http.StatusBadRequest, "bad request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, "internal error")
	}

	b, err := __ENGINE.TextSearch(ctx, body)
	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, "internal error")
	}

	return e.Blob(
		http.StatusOK,
		__ENGINE.ResponseContentType(),
		b,
	)
}

func vector(e echo.Context) error {
	req := e.Request()
	if req.ContentLength > REQUEST_BODY_LIMIT {
		e.Logger().Error("content length over limit")
		return e.String(http.StatusBadRequest, "bad request")
	}
	contentType := req.Header.Get("Content-Type")
	if len(contentType) == 0 || contentType != __ENGINE.RequestContentType() {
		e.Logger().Error("unexpected content type header")
		return e.String(http.StatusBadRequest, "bad request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, "internal error")
	}

	b, err := __ENGINE.VectorSeach(ctx, body)
	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, "internal error")
	}

	return e.Blob(
		http.StatusOK,
		__ENGINE.ResponseContentType(),
		b,
	)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetPrefix("KENJA2")

	port := args()
	engineUri, err := env()
	if err != nil {
		e.Logger.Fatal(err)
	}

	ctx := context.Background()
	__ENGINE, err = mongodb.Connect(
		engineUri,
		endec.NewJson(),
		endec.NewJson(),
	)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer func() {
		if err := __ENGINE.Close(ctx); err != nil {
			e.Logger.Error(err)
		}
	}()

	e.GET("/text", text)
	e.GET("/vector", vector)

	listenAt := fmt.Sprintf("localhost:%d", port)
	if err := e.Start(listenAt); err != nil {
		e.Logger.Error(err)
	}
}
