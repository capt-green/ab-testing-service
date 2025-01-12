package proxy

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	RequestsTotal       *prometheus.CounterVec
	LatencyHistogram    *prometheus.HistogramVec
	BytesSentTotal      *prometheus.CounterVec
	BytesReceivedTotal  *prometheus.CounterVec
	ResponseStatusTotal *prometheus.CounterVec
	ActiveConnections   *prometheus.GaugeVec
	RequestErrors       *prometheus.CounterVec
}

func newProxyMetrics(proxyID string) *Metrics {
	return &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "ab_test_requests_total",
				Help:        "Total number of requests per target",
				ConstLabels: prometheus.Labels{"proxy_id": proxyID},
			},
			[]string{"target"},
		),
		LatencyHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "ab_test_request_duration_seconds",
				Help:        "Request duration in seconds",
				ConstLabels: prometheus.Labels{"proxy_id": proxyID},
				Buckets:     []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"target"},
		),
		BytesSentTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "ab_test_bytes_sent_total",
				Help:        "Total number of bytes sent to clients",
				ConstLabels: prometheus.Labels{"proxy_id": proxyID},
			},
			[]string{"target"},
		),
		BytesReceivedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "ab_test_bytes_received_total",
				Help:        "Total number of bytes received from targets",
				ConstLabels: prometheus.Labels{"proxy_id": proxyID},
			},
			[]string{"target"},
		),
		ResponseStatusTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "ab_test_response_status_total",
				Help:        "Total number of responses by status code",
				ConstLabels: prometheus.Labels{"proxy_id": proxyID},
			},
			[]string{"target", "status"},
		),
		ActiveConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        "ab_test_active_connections",
				Help:        "Number of active connections",
				ConstLabels: prometheus.Labels{"proxy_id": proxyID},
			},
			[]string{"target"},
		),
		RequestErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "ab_test_request_errors_total",
				Help:        "Total number of request errors",
				ConstLabels: prometheus.Labels{"proxy_id": proxyID},
			},
			[]string{"target", "error_type"},
		),
	}
}
