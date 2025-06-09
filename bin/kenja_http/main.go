package main

import (
	"context"
	"flag"
	"fmt"
	"kenja2"
	"kenja2/endec"
	"kenja2/mongodb"

	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

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
	return e.JSON(http.StatusOK, `{"texts": ["hello"]}`)
}

func main() {
	e := echo.New()
	port := args()
	engineUri, err := env()
	if err != nil {
		e.Logger.Fatal(err)
	}

	ctx := context.Background()
	__ENGINE, err := mongodb.Connect(
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

	listenAt := fmt.Sprintf("localhost:%d", port)
	if err := e.Start(listenAt); err != nil {
		e.Logger.Error(err)
	}
}
