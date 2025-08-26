package main

import (
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetPrefix("KENJA-PROXY")

	webUiUrl, err := url.Parse("http://localhost:5173")
	if err != nil {
		e.Logger.Fatal(err)
	}
	webUiBalancer := middleware.NewRoundRobinBalancer(
		[]*middleware.ProxyTarget{
			{
				Name: "web-ui",
				URL:  webUiUrl,
			},
		},
	)

	imgUrl, err := url.Parse("http://localhost:8081")
	if err != nil {
		e.Logger.Fatal(err)
	}
	imgBalancer := middleware.NewRoundRobinBalancer(
		[]*middleware.ProxyTarget{
			{
				Name: "img",
				URL:  imgUrl,
			},
		},
	)

	imgGroup := e.Group("/preview/*")
	imgGroup.Use(middleware.Proxy(imgBalancer))
	webUiGroup := e.Group("/*")
	webUiGroup.Use(middleware.Proxy(webUiBalancer))

	if err := e.Start("localhost:8082"); err != nil {
		e.Logger.Error(err)
	}
}
