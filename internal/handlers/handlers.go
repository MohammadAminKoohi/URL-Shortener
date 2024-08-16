package handlers

import (
	"UrlShortener/internal/util"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type Handler struct {
	Database *sql.DB
	Redis    *redis.Client
}

func (h *Handler) Shorten(c echo.Context) error {
	inputUrl := c.QueryParam("url")
	if inputUrl == "" {
		return c.String(http.StatusBadRequest, "URL parameter is required")
	}
	fmt.Println("Original URL: " + inputUrl)
	shortenedURL := util.ToBase62(inputUrl)
	fmt.Println("Shortened URL: http://localhost:1323/" + shortenedURL)
	_, err := h.Database.Exec(`INSERT INTO urls (original_url, shortened_url, count) VALUES ($1, $2, $3)`, inputUrl, shortenedURL, 0)
	if err != nil {
		fmt.Println("Error inserting into database " + err.Error())
		return c.String(http.StatusInternalServerError, "Error inserting into database")
	} else {
		fmt.Println("Inserted into database")
	}
	return c.String(http.StatusOK, "Shortened URL: http://localhost:1323/"+shortenedURL)
}

func (h *Handler) Redirect(c echo.Context) error {
	shortenedURL := c.Param("shortenedURL")
	var originalURL string
	originalURL, err := h.Redis.Get(c.Request().Context(), shortenedURL).Result()
	if err == redis.Nil {
		err := h.Database.QueryRow(`SELECT original_url FROM urls WHERE shortened_url = $1`, shortenedURL).Scan(&originalURL)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.String(http.StatusNotFound, "URL not found")
			}
			fmt.Println("Error querying database: " + err.Error())
			return c.String(http.StatusInternalServerError, "Error querying database")
		}
		_, err = h.Database.Exec(`UPDATE urls SET count = count + 1 WHERE shortened_url = $1`, shortenedURL)
		if err != nil {
			fmt.Println("Error updating count in database: " + err.Error())
			return c.String(http.StatusInternalServerError, "Error updating count in database")
		}
		err = h.Redis.Set(c.Request().Context(), shortenedURL, originalURL, 5*time.Minute).Err()
		if err != nil {
			fmt.Println("Error setting Redis cache: " + err.Error())
			return c.String(http.StatusInternalServerError, "Error setting Redis cache")
		}
	} else if err != nil {
		fmt.Println("Error retrieving from Redis: " + err.Error())
		return c.String(http.StatusInternalServerError, "Error retrieving from Redis")
	}
	return c.Redirect(http.StatusFound, originalURL)
}

func (h *Handler) List(c echo.Context) error {
	rows, err := h.Database.Query(`SELECT original_url, shortened_url FROM urls`)
	if err != nil {
		fmt.Println("Error querying database")
		return c.String(http.StatusInternalServerError, "Error querying database")
	}
	defer rows.Close()

	urls := make(map[string]string)
	for rows.Next() {
		var originalURL, shortenedURL string
		if err := rows.Scan(&originalURL, &shortenedURL); err != nil {
			fmt.Println("Error scanning row")
			return c.String(http.StatusInternalServerError, "Error scanning row")
		}
		urls[shortenedURL] = originalURL
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows")
		return c.String(http.StatusInternalServerError, "Error iterating rows")
	}

	return c.JSON(http.StatusOK, urls)
}
