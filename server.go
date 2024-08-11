package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

var shortURLs = make(map[string]string)

func generateHash(URL string) string {
	hash := sha256.New()
	hash.Write([]byte(URL))
	return hex.EncodeToString(hash.Sum(nil))[:8]
}

func shorten(c echo.Context) error {
	URL := c.QueryParam("url")
	if URL == "" {
		return c.String(http.StatusBadRequest, "URL parameter is required")
	}
	shortenedURL := generateHash(URL)
	shortURLs[shortenedURL] = URL
	return c.String(http.StatusOK, "Shortened URL: http://localhost:1323/"+shortenedURL)
}

func redirect(c echo.Context) error {
	shortenedURL := c.Param("shortenedURL")
	originalURL := shortURLs[shortenedURL]
	if originalURL == "" {
		fmt.Println("URL not found")
		return c.String(http.StatusNotFound, "URL not found")
	}
	return c.Redirect(http.StatusMovedPermanently, originalURL)
}

func list(c echo.Context) error {
	return c.JSON(http.StatusOK, shortURLs)
}
func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/shorten", shorten)
	e.GET("/:shortenedURL", redirect)
	e.GET("/urls", list)
	e.Logger.Fatal(e.Start(":1323"))
}
