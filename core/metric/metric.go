package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	MethodDuration         *prometheus.HistogramVec
	MethodErrorDuration    *prometheus.HistogramVec
	MethodTotal            *prometheus.CounterVec
	MethodSuccessTotal     *prometheus.CounterVec
	MethodUserErrorTotal   *prometheus.CounterVec
	MethodServerErrorTotal *prometheus.CounterVec
}

func NewMetric() *Metric {
	methodLabels := []string{"service_name", "type", "method"}
	errorLabels := []string{"service_name", "type", "method", "error"}
	buckets := []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1, 1.5,
		2, 2.5, 3, 4, 5, 10, 15, 20, 30, 45, 60, 80, 100, 150, 200}
	metrics := &Metric{
		MethodDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "method_duration_seconds",
				Help:    "A histogram of latencies for each method",
				Buckets: buckets,
			}, methodLabels,
		),
		MethodErrorDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "method_error_duration_seconds",
				Help:    "A histogram of latencies for each method with errors",
				Buckets: buckets,
			}, errorLabels,
		),
		MethodTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "method_total",
			Help: "A counter for each method",
		}, methodLabels),
		MethodSuccessTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "method_success_total",
			Help: "A counter for each method",
		}, methodLabels),
		MethodUserErrorTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "method_user_error_total",
			Help: "A counter for each method with errors",
		}, errorLabels),
		MethodServerErrorTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "method_server_error_total",
			Help: "A counter for each method with errors",
		}, errorLabels),
	}

	prometheus.MustRegister(metrics.MethodDuration, metrics.MethodErrorDuration, metrics.MethodTotal,
		metrics.MethodSuccessTotal, metrics.MethodUserErrorTotal, metrics.MethodServerErrorTotal)

	return metrics
}
