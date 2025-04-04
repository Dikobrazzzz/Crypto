package metrics

import (
	"log/slog"
	"net/http"
	"time"

	"crypto/internal/cache"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotalMetric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of incoming HTTP requests",
		},
	)

	HttpStatusMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_responses_total",
			Help: "Count of HTTP responses, labeled by status code and method",
		},
		[]string{"status", "method"},
	)

	logger *slog.Logger
)

func InitMetrics(port string, cache *cache.CacheDecorator) {
	prometheus.MustRegister(RequestsTotalMetric)
	prometheus.MustRegister(HttpStatusMetric)

	cacheMemoryUsage := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "cache_memory_usage_bytes",
			Help: "Approximate amount of memory occupied by the cache in bytes",
		},
		func() float64 {
			return float64(cache.MemoryUsage())
		},
	)
	prometheus.MustRegister(cacheMemoryUsage)

	cacheSize := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "cache_items_count",
			Help: "Number of items stored in the cache",
		},
		func() float64 {
			return float64(cache.Size())
		},
	)
	prometheus.MustRegister(cacheSize)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		logger.Info("starting metrics server", "port", port)
		server := &http.Server{
			Addr:         ":" + port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			logger.Error("Failed to start metrics server", "error", errors.Wrap(err, ""))
		}
	}()
}

func HttpStatusMetricInc(statusCode int, method string) {
	HttpStatusMetric.WithLabelValues(http.StatusText(statusCode), method).Inc()
}

func RequestsTotalInc() {
	RequestsTotalMetric.Inc()
}
