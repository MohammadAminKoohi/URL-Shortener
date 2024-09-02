package main

import (
	"UrlShortener/internal/DB"
	"UrlShortener/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"path", "method"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func main() {
	e := echo.New()
	h := handlers.Handler{
		Database: DB.DbInit(),
		Redis:    DB.RedisInit(),
	}
	e.Use(prometheusMiddleware)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/shorten", h.Shorten)
	e.GET("/:shortenedURL", h.Redirect)
	e.GET("/urls", h.List)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.Logger.Fatal(e.Start(":1323"))
}
func prometheusMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			httpRequestDuration.WithLabelValues(c.Path(), c.Request().Method).Observe(v)
		}))
		httpRequestsTotal.WithLabelValues(c.Path(), c.Request().Method).Inc()
		err := next(c)
		timer.ObserveDuration()

		return err
	}
}
