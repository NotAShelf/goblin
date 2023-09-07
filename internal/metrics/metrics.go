package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	pasteLengthHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goblin",
			Subsystem: "paste",
			Name:      "length",
			Help:      "Length of pastes",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10), // Example buckets
		},
		[]string{"status"}, // Labels for success/failure
	)

	pasteCreatedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goblin",
			Subsystem: "paste",
			Name:      "created_total",
			Help:      "Total number of pastes created",
		},
		[]string{"status"}, // Labels for success/failure
	)
)

// ObservePasteLength records the length of pastes.
func ObservePasteLength(status string, length float64) {
	pasteLengthHistogram.WithLabelValues(status).Observe(length)
}

// IncrementPasteCreatedCounter increments the paste creation counter.
func IncrementPasteCreatedCounter(status string) {
	pasteCreatedCounter.WithLabelValues(status).Inc()
}

// InitPrometheus initializes Prometheus metrics and registers a handler for /metrics endpoint.
func InitPrometheus(router *mux.Router) {
	// Register Prometheus metrics collectors here
	// For example, you can use promauto to define and register metrics

	// Register Prometheus HTTP handler for /metrics endpoint with the Mux router
	router.Handle("/metrics", promhttp.Handler())
}
