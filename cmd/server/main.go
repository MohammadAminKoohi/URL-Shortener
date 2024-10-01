package main

import (
	"UrlShortener/internal/DB"
	"UrlShortener/internal/handlers"
	"UrlShortener/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	e := echo.New()
	h := handlers.NewHandler(
		DB.DbInit(),
		DB.RedisInit(),
	)
	e.Use(middleware.PrometheusMetrics)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to my Url Shortener")
	})
	e.POST("/shorten", h.Shorten)
	e.GET("/:shortenedURL", h.Redirect)
	e.GET("/urls", h.List)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.Logger.Fatal(e.Start(":1323"))
}
