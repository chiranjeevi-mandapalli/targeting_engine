package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the application
type Metrics struct {
	HTTPRequests       *prometheus.CounterVec
	HTTPRequestLatency *prometheus.HistogramVec
	CacheHits          *prometheus.CounterVec
	CacheMisses        *prometheus.CounterVec
	CacheLatency       *prometheus.HistogramVec
	ActiveCampaigns    prometheus.Gauge
}

// NewMetrics initializes and returns a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total HTTP requests by method, path, and status",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency distribution",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),
		CacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total cache hits by type",
			},
			[]string{"type"},
		),
		CacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total cache misses by type",
			},
			[]string{"type"},
		),
		CacheLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "cache_request_duration_seconds",
				Help:    "Cache request latency distribution",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"type"},
		),
		ActiveCampaigns: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "campaigns_active_total",
				Help: "Current number of active campaigns",
			},
		),
	}
}

// Init initializes the metrics object
func Init() *Metrics {
	return NewMetrics()
}

// IncrementHTTPRequests increments the HTTP request counter for the given method, path, and status
func (m *Metrics) IncrementHTTPRequests(method, path, status string) {
	m.HTTPRequests.WithLabelValues(method, path, status).Inc()
}

// ObserveHTTPRequestLatency records the HTTP request duration
func (m *Metrics) ObserveHTTPRequestLatency(method, path string, duration float64) {
	m.HTTPRequestLatency.WithLabelValues(method, path).Observe(duration)
}

// IncrementCacheHits increments the cache hit counter for a given type
func (m *Metrics) IncrementCacheHits(cacheType string) {
	m.CacheHits.WithLabelValues(cacheType).Inc()
}

// IncrementCacheMisses increments the cache miss counter for a given type
func (m *Metrics) IncrementCacheMisses(cacheType string) {
	m.CacheMisses.WithLabelValues(cacheType).Inc()
}

// ObserveCacheLatency records the cache request latency for a given type
func (m *Metrics) ObserveCacheLatency(cacheType string, duration float64) {
	m.CacheLatency.WithLabelValues(cacheType).Observe(duration)
}

// SetActiveCampaigns sets the active campaigns gauge to a specific value
func (m *Metrics) SetActiveCampaigns(value float64) {
	m.ActiveCampaigns.Set(value)
}
