package http

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	r.body.Write(b) // capture body
	return r.ResponseWriter.Write(b)
}

func (h *Http) metricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &statusRecorder{ResponseWriter: w, status: 200, body: &bytes.Buffer{}} // default to 200

		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Seconds()

		route := mux.CurrentRoute(r)
		methodName := r.URL.Path
		if route != nil {
			pathTemplate, err := route.GetPathTemplate()
			if err == nil {
				methodName = pathTemplate
			}
		} else {
			methodName = "not_implemented_route"
		}

		if methodName == "/health" || methodName == "/readiness" {
			return
		}

		h.metrics.MethodTotal.With(prometheus.Labels{"service_name": h.serviceName, "type": "http",
			"method": methodName}).Inc()
		h.metrics.MethodDuration.With(prometheus.Labels{"service_name": h.serviceName, "type": "http",
			"method": methodName}).Observe(duration)

		if recorder.status >= 200 && recorder.status <= 399 {
			h.metrics.MethodSuccessTotal.With(prometheus.Labels{"service_name": h.serviceName, "type": "http",
				"method": methodName}).Inc()
		} else {

			errorMsg := recorder.body.String()
			if len(errorMsg) > 200 {
				errorMsg = errorMsg[:200] // truncate to avoid label explosion
			}

			h.metrics.MethodErrorDuration.With(prometheus.Labels{"service_name": h.serviceName, "type": "http",
				"method": methodName, "error": errorMsg}).Observe(duration)

			if recorder.status >= 400 && recorder.status <= 499 {
				h.metrics.MethodUserErrorTotal.With(prometheus.Labels{"service_name": h.serviceName, "type": "http",
					"method": methodName, "error": errorMsg}).Inc()
			} else {
				h.metrics.MethodServerErrorTotal.With(prometheus.Labels{"service_name": h.serviceName, "type": "http",
					"method": methodName, "error": errorMsg}).Inc()
			}
		}
	})
}
