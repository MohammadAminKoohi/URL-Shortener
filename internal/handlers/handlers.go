package handlers

import (
	"UrlShortener/internal/util"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

var shortURLs = make(map[string]string)

func Shorten(c echo.Context) error {
	URL := c.QueryParam("url")
	if URL == "" {
		return c.String(http.StatusBadRequest, "URL parameter is required")
	}
	shortenedURL := util.ToBase62(URL)
	shortURLs[shortenedURL] = URL
	return c.String(http.StatusOK, "Shortened URL: http://localhost:1323/"+shortenedURL)
}

func Redirect(c echo.Context) error {
	shortenedURL := c.Param("shortenedURL")
	originalURL := shortURLs[shortenedURL]
	if originalURL == "" {
		fmt.Println("URL not found")
		return c.String(http.StatusNotFound, "URL not found")
	}
	return c.Redirect(http.StatusMovedPermanently, originalURL)
}

func List(c echo.Context) error {
	return c.JSON(http.StatusOK, shortURLs)
}
