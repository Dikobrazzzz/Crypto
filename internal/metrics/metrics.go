package metrics

import (
	"log/slog" // Подключаем пакет slog из стандартной библиотеки
	"net/http"
	"os"
	"time"

	"crypto/internal/cache" // замените на ваш корректный import-путь

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

	CacheMemoryUsage = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "cache_memory_usage_bytes",
			Help: "Approximate amount of memory occupied by the cache in bytes",
		},
		func() float64 {
			return float64(c.MemoryUsage())
		},
	)

	CacheSize = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "cache_items_count",
			Help: "Number of items stored in the cache",
		},
		func() float64 {
			return float64(c.Size())
		},
	)

	c *cache.CacheDecorator

	logger *slog.Logger
)

// InitMetrics регистрирует метрики и запускает сервер для /metrics
func InitMetrics(port string, cache *cache.CacheDecorator) {
	c = cache
	prometheus.MustRegister(RequestsTotalMetric)
	prometheus.MustRegister(HttpStatusMetric)
	prometheus.MustRegister(CacheMemoryUsage)
	prometheus.MustRegister(CacheSize)

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
			// slog не имеет метода Fatal, поэтому вручную выходим при критической ошибке
			logger.Error("Failed to start metrics server", "error", errors.Wrap(err, ""))
			os.Exit(1)
		}
	}()
}

// Увеличивает счётчик ответов с разбивкой по статусу и методу
func HttpStatusMetricInc(statusCode int, method string) {
	HttpStatusMetric.WithLabelValues(http.StatusText(statusCode), method).Inc()
}

// Счётчик общего числа запросов (можно вызывать из middleware)
func RequestsTotalInc() {
	RequestsTotalMetric.Inc()
}
