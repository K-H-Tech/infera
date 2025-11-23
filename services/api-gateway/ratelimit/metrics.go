package ratelimit

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RateLimitHitsTotal counts total rate limit checks
	RateLimitHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_ratelimit_hits_total",
			Help: "Total number of rate limit checks by endpoint and status",
		},
		[]string{"endpoint", "ip", "status"}, // status: allowed, blocked
	)

	// RateLimitViolationsTotal counts rate limit violations
	RateLimitViolationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_ratelimit_violations_total",
			Help: "Total number of rate limit violations by endpoint",
		},
		[]string{"endpoint", "ip"},
	)

	// RateLimitBackoffDuration tracks backoff duration histogram
	RateLimitBackoffDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_ratelimit_backoff_duration_seconds",
			Help:    "Duration of rate limit backoff periods",
			Buckets: prometheus.ExponentialBuckets(1, 2, 10), // 1s to 512s
		},
		[]string{"endpoint", "ip"},
	)

	// RateLimitTokensRemaining tracks remaining tokens gauge
	RateLimitTokensRemaining = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_ratelimit_tokens_remaining",
			Help: "Number of tokens remaining for rate limit keys",
		},
		[]string{"endpoint", "ip"},
	)
)

// RecordAllow records an allowed request
func RecordAllow(endpoint, ip string, remaining int) {
	RateLimitHitsTotal.WithLabelValues(endpoint, ip, "allowed").Inc()
	RateLimitTokensRemaining.WithLabelValues(endpoint, ip).Set(float64(remaining))
}

// RecordBlock records a blocked request
func RecordBlock(endpoint, ip string, backoffSeconds float64) {
	RateLimitHitsTotal.WithLabelValues(endpoint, ip, "blocked").Inc()
	RateLimitViolationsTotal.WithLabelValues(endpoint, ip).Inc()
	RateLimitBackoffDuration.WithLabelValues(endpoint, ip).Observe(backoffSeconds)
	RateLimitTokensRemaining.WithLabelValues(endpoint, ip).Set(0)
}
