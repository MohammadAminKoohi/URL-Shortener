package main

import (
	"UrlShortener/internal/DB"
	"UrlShortener/internal/handlers"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	h := handlers.Handler{
		Database: DB.DbInit(),
	}
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/shorten", h.Shorten)
	e.GET("/:shortenedURL", h.Redirect)
	e.GET("/urls", h.List)
	e.Logger.Fatal(e.Start(":1323"))
}
