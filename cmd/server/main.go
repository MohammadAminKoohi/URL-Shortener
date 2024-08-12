package main

import (
	"UrlShortener/internal/handlers"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/shorten", handlers.Shorten)
	e.GET("/:shortenedURL", handlers.Redirect)
	e.GET("/urls", handlers.List)
	e.Logger.Fatal(e.Start(":1323"))
}
