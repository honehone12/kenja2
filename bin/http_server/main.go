package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"kenja2"
	"kenja2/documents"
	"kenja2/engine"
	"os"
	"time"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"lukechampine.com/blake3"
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
		e.Logger().Error(err)
		return responseBadRequest(e)
	}

	ctx, cancel := context.WithTimeout(req.Context(), REQUEST_TIME_LIMIT)
	defer cancel()

	docs, err := kenja2.BufferedQuery(__engine, ctx, &q)
	if err != nil {
		e.Logger().Error(err)
		return responseInternalError(e)
	}

	e.Logger().Warn("Don't forget to remove hashing")
	for i := range docs {
		url := docs[i].Url
		if len(url) == 0 {
			err := errors.New("url is empty")
			e.Logger().Error(err)
			return responseInternalError(e)
		}

		hash := blake3.Sum256([]byte(url))
		x := hex.EncodeToString(hash[:16])
		img := fmt.Sprintf("%s.jpg", x)
		docs[i].Img = img
	}

	return e.JSON(http.StatusOK, documents.SearchResult{Documents: docs})
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

	imgRoot := os.Getenv("IMG_ROOT")
	if len(imgRoot) == 0 {
		e.Logger.Fatal("env for img root is not set")
	}

	port := args()
	__engine, err = kenja2.ConnectAtlas(ENGINE_MAX_CONN, e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Static("/preview", imgRoot)
	e.POST("/anime-search", handler)

	listenAt := fmt.Sprintf("localhost:%d", port)
	if err := e.Start(listenAt); err != nil {
		e.Logger.Error(err)
	}
}
