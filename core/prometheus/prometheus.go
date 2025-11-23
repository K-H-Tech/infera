package prometheus

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunPrometheus(address string) {
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus listening on " + address)

		srv := &http.Server{
			Addr:         address,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println("Error starting prometheus server:", err)
		}
	}()
}
