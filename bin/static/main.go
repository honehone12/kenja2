package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetPrefix("KENJA-STATIC")

	err := godotenv.Load("../../.env")
	if err != nil {
		e.Logger.Fatal(err)
	}

	imgRoot := os.Getenv("IMG_ROOT")
	if len(imgRoot) == 0 {
		e.Logger.Fatal("env for img root is not set")
	}

	e.Static("/", imgRoot)

	if err := e.Start("localhost:8081"); err != nil {
		e.Logger.Error(err)
	}
}
