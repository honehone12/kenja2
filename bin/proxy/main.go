package main

import (
	"flag"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// transport headers are not set for local.
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// {
//   "Strict-Transport-Security": "max-age=86400; includeSubDomains; preload"
// }
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var __securityHeaders map[string]string = map[string]string{
	"X-XSS-Protection":             "0",
	"X-Frame-Options":              "DENY",
	"X-Content-Type-Options":       "nosniff",
	"Referrer-Policy":              "strict-origin-when-cross-origin",
	"Cross-Origin-Opener-Policy":   "same-origin",
	"Cross-Origin-Resource-Policy": "same-origin",
	"Permissions-Policy":           "geolocation=(), camera=(), microphone=()",
	"Content-Security-Policy":      "script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net/npm/ 'wasm-unsafe-eval'; font-src 'self'; img-src 'self' data: ; worker-src 'self' blob: ; connect-src 'self' https://huggingface.co/ https://cas-bridge.xethub.hf.co/ https://cdn.jsdelivr.net/npm/; base-uri 'self'; form-action 'self'; media-src 'self'; object-src 'self'; manifest-src 'self'; child-src 'none'; frame-src https://www.youtube.com/ https://www.youtube-nocookie.com/; frame-ancestors 'none';",
}

func securityHeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		h := ctx.Response().Header()
		for k, v := range __securityHeaders {
			h.Set(k, v)
		}
		return next(ctx)
	}
}

func hostSwapMiddleware(proxy *url.URL) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Request().Host = proxy.Hostname()
			return next(ctx)
		}
	}
}

func cacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		h := ctx.Response().Header()
		if len(h.Get("Cache-Control")) == 0 {
			h.Set("Cache-Control", "max-age=86400")
		}

		return next(ctx)
	}
}

func main() {
	dev := flag.Bool("dev", false, "development mode")
	flag.Parse()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetPrefix("KENJA-PROXY")

	err := godotenv.Load("../../.env")
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(securityHeaderMiddleware)

	imgRoot := os.Getenv("IMG_ROOT")
	if len(imgRoot) == 0 {
		e.Logger.Fatal("env for img root is not set")
	}
	imgGroup := e.Group("/preview")
	imgGroup.Use(cacheControlMiddleware)
	imgGroup.Use(middleware.Static(imgRoot))

	a := os.Getenv("AWS_API_URL")
	if len(a) == 0 {
		e.Logger.Fatal("env for api url is not set")
	}
	awsUrl, err := url.Parse(a)
	if err != nil {
		e.Logger.Fatal(err)
	}
	apiBalancer := middleware.NewRoundRobinBalancer(
		[]*middleware.ProxyTarget{
			{
				Name: "api",
				URL:  awsUrl,
			},
		},
	)
	apiGroup := e.Group("/api/*")
	apiGroup.Use(hostSwapMiddleware(awsUrl))
	apiGroup.Use(middleware.Proxy(apiBalancer))

	if *dev {
		e.Logger.Info("development mode")
		devUrl, err := url.Parse("http://localhost:5173")
		if err != nil {
			e.Logger.Fatal(err)
		}
		devBalancer := middleware.NewRoundRobinBalancer(
			[]*middleware.ProxyTarget{
				{
					Name: "dev",
					URL:  devUrl,
				},
			},
		)
		devGroup := e.Group("/*")
		devGroup.Use(middleware.Proxy(devBalancer))
	} else {
		e.Logger.Info("built app mode")
		appRoot := os.Getenv("APP_ROOT")
		if len(appRoot) == 0 {
			e.Logger.Fatal("env for app root is not set")
		}
		appGroup := e.Group("/*")
		appGroup.Use(cacheControlMiddleware)
		appGroup.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Skipper: middleware.DefaultSkipper,
			Root:    appRoot,
			Index:   "200.html",
			HTML5:   true,
		}))
	}

	if err := e.Start("0.0.0.0:8082"); err != nil {
		e.Logger.Error(err)
	}
}
