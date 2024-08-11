package main

import (
	"encoding/binary"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

var shortURLs = make(map[string]string)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func toBase62(url string) string {
	num := binary.BigEndian.Uint64([]byte(url))
	if num == 0 {
		return string(base62Chars[0])
	}
	var result []byte
	for num > 0 {
		remainder := num % 62
		result = append([]byte{base62Chars[remainder]}, result...)
		num = num / 62
	}
	return string(result)
}

func shorten(c echo.Context) error {
	URL := c.QueryParam("url")
	if URL == "" {
		return c.String(http.StatusBadRequest, "URL parameter is required")
	}
	shortenedURL := toBase62(URL)
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
