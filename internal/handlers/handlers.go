package handlers

import (
	"UrlShortener/internal/util"
	"context"
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

const (
	maxKeys  = 10000
	keyList  = "key_list"
	cacheTTL = 5 * time.Minute
)

func NewHandler(db *sql.DB, redis *redis.Client) *Handler {
	return &Handler{
		Database: db,
		Redis:    redis,
	}
}

func (h *Handler) Shorten(c echo.Context) error {
	inputURL := c.QueryParam("url")
	if inputURL == "" {
		return c.String(http.StatusBadRequest, "URL parameter is required")
	}

	shortenedURL := util.ToBase62(inputURL)
	_, err := h.Database.Exec(`INSERT INTO urls (original_url, shortened_url, count) VALUES ($1, $2, $3)`, inputURL, shortenedURL, 0)
	if err != nil {
		fmt.Printf("Database insertion error: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error inserting into database")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Shortened URL: http://localhost:1323/%s", shortenedURL))
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
			return c.String(http.StatusInternalServerError, "Error querying database")
		}

		_, err = h.Database.Exec(`UPDATE urls SET count = count + 1 WHERE shortened_url = $1`, shortenedURL)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error updating count in database")
		}

		if err := h.manageCache(c.Request().Context(), shortenedURL, originalURL); err != nil {
			return c.String(http.StatusInternalServerError, "Error managing Redis cache")
		}
	} else if err != nil {
		return c.String(http.StatusInternalServerError, "Error retrieving from Redis")
	}

	return c.Redirect(http.StatusFound, originalURL)
}

func (h *Handler) manageCache(ctx context.Context, key, value string) error {
	_, err := h.Redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, key, value, cacheTTL).Err()
		pipe.RPush(ctx, keyList, key).Err()

		listLength, _ := pipe.LLen(ctx, keyList).Result()
		if listLength > maxKeys {
			oldestKey, _ := pipe.LPop(ctx, keyList).Result()
			pipe.Del(ctx, oldestKey).Err()
		}
		return nil
	})
	return err
}

func (h *Handler) List(c echo.Context) error {
	rows, err := h.Database.Query(`SELECT original_url, shortened_url FROM urls`)
	if err != nil {
		fmt.Printf("Database query error: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error querying database")
	}
	defer rows.Close()

	urls := make(map[string]string)
	for rows.Next() {
		var originalURL, shortenedURL string
		if err := rows.Scan(&originalURL, &shortenedURL); err != nil {
			fmt.Printf("Row scanning error: %v\n", err)
			return c.String(http.StatusInternalServerError, "Error scanning row")
		}
		urls[shortenedURL] = originalURL
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Row iteration error: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error iterating rows")
	}

	return c.JSON(http.StatusOK, urls)
}
