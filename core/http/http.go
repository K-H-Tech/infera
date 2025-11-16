package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"zarinpal-platform/core/docs"
	"zarinpal-platform/core/metric"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

type Http struct {
	Engine      *mux.Router
	GatewayMux  *runtime.ServeMux
	metrics     *metric.Metric
	serviceName string
	address     string
}

func NewHttp(serviceName string, address string, metrics *metric.Metric) *Http {
	router := mux.NewRouter()
	gwMux := runtime.NewServeMux(
		runtime.WithMetadata(forwardMetadata))

	h := &Http{Engine: router, GatewayMux: gwMux, serviceName: serviceName, address: address, metrics: metrics}

	router.Use(h.metricMiddleware)

	// todo change this with true prefix for grpc-gateway
	router.PathPrefix("/rest/").Handler(gwMux)

	h.initCommonRoutes()
	h.docsRoute()
	return h
}

func (h *Http) Start() {
	go func() {
		log.Printf("Starting http server on %s\n", h.address)

		if err := http.ListenAndServe(h.address, h.Engine); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	time.Sleep(200 * time.Millisecond)
}

func (h *Http) docsRoute() {
	h.Engine.HandleFunc(fmt.Sprintf("/docs/%s.swagger.json", h.serviceName), func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("./services/%s/docs/%s.swagger.json", h.serviceName, h.serviceName))
	})

	h.Engine.HandleFunc("/docs/swagger-ui/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(docs.GetSwaggerHtml(h.serviceName)))
	})

}

func (h *Http) initCommonRoutes() {
	h.Engine.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]int{"ok": 1}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	h.Engine.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]int{"ok": 1}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	h.Engine.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]int{"ok": 1}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})
}

func forwardMetadata(ctx context.Context, req *http.Request) metadata.MD {
	// grpc-gateway automatically converts headers to lowercase
	acceptLang := req.Header.Get("Accept-Language")
	if acceptLang != "" {
		return metadata.Pairs("accept-language", acceptLang)
	}
	return metadata.MD{}
}
