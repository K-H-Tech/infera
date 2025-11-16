package prometheus

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunPrometheus(address string) {
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus listening on " + address)
		err := http.ListenAndServe(address, mux)
		if err != nil {
			fmt.Println("Error starting prometheus server:", err)
		}
	}()
}
