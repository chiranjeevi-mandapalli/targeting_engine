package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics contains all monitoring metrics
type Metrics struct {
	HTTPRequests       *prometheus.CounterVec
	HTTPRequestLatency *prometheus.HistogramVec
	CacheHits          *prometheus.CounterVec
	CacheMisses        *prometheus.CounterVec
	ActiveCampaigns    prometheus.Gauge
}

// NewMetrics creates and registers all metrics
func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total HTTP requests by method, path and status",
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
		ActiveCampaigns: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "campaigns_active_total",
				Help: "Current number of active campaigns",
			},
		),
	}
}

// Register custom metrics
func Init() *Metrics {
	return NewMetrics()
}
