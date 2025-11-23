package proxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"zarinpal-platform/services/api-gateway/config"

	"github.com/gorilla/mux"
)

func TestNewDocsProxy_Disabled(t *testing.T) {
	cfg := config.DocsSection{
		Enabled: false,
	}

	proxy, err := NewDocsProxy(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if proxy.IsEnabled() {
		t.Error("expected proxy to be disabled")
	}
}

func TestNewDocsProxy_InvalidURL(t *testing.T) {
	cfg := config.DocsSection{
		Enabled: true,
		Services: map[string]config.DocsService{
			"test": {URL: "://invalid-url"},
		},
	}

	_, err := NewDocsProxy(cfg)
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestNewDocsProxy_ValidConfig(t *testing.T) {
	cfg := config.DocsSection{
		Enabled: true,
		Services: map[string]config.DocsService{
			"auth":           {URL: "http://localhost:8080"},
			"user-dashboard": {URL: "http://localhost:8081"},
		},
	}

	proxy, err := NewDocsProxy(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !proxy.IsEnabled() {
		t.Error("expected proxy to be enabled")
	}

	names := proxy.GetServiceNames()
	if len(names) != 2 {
		t.Errorf("expected 2 services, got %d", len(names))
	}
}

func TestDocsProxy_DocsList(t *testing.T) {
	cfg := config.DocsSection{
		Enabled: true,
		Services: map[string]config.DocsService{
			"auth":           {URL: "http://localhost:8080"},
			"user-dashboard": {URL: "http://localhost:8081"},
		},
	}

	proxy, err := NewDocsProxy(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	router := mux.NewRouter()
	proxy.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "auth") {
		t.Error("expected response to contain 'auth' service")
	}
	if !strings.Contains(body, "user-dashboard") {
		t.Error("expected response to contain 'user-dashboard' service")
	}
	if !strings.Contains(body, "swagger-ui") {
		t.Error("expected response to contain swagger-ui links")
	}
}

func TestDocsProxy_ServiceNotFound(t *testing.T) {
	cfg := config.DocsSection{
		Enabled: true,
		Services: map[string]config.DocsService{
			"auth": {URL: "http://localhost:8080"},
		},
	}

	proxy, err := NewDocsProxy(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	router := mux.NewRouter()
	proxy.RegisterRoutes(router)

	// Test swagger UI for non-existent service
	req := httptest.NewRequest(http.MethodGet, "/docs/nonexistent/swagger-ui", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	// Test swagger JSON for non-existent service
	req = httptest.NewRequest(http.MethodGet, "/docs/nonexistent/swagger.json", nil)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestDocsProxy_SwaggerUIProxy(t *testing.T) {
	// Create a mock backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the path was rewritten correctly
		if r.URL.Path != "/docs/swagger-ui/" {
			t.Errorf("expected path /docs/swagger-ui/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("swagger ui content"))
	}))
	defer backend.Close()

	cfg := config.DocsSection{
		Enabled: true,
		Services: map[string]config.DocsService{
			"auth": {URL: backend.URL},
		},
	}

	proxy, err := NewDocsProxy(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	router := mux.NewRouter()
	proxy.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/docs/auth/swagger-ui", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestDocsProxy_SwaggerJSONProxy(t *testing.T) {
	// Create a mock backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the path was rewritten correctly
		if r.URL.Path != "/docs/auth.swagger.json" {
			t.Errorf("expected path /docs/auth.swagger.json, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"openapi": "3.0.0"}`))
	}))
	defer backend.Close()

	cfg := config.DocsSection{
		Enabled: true,
		Services: map[string]config.DocsService{
			"auth": {URL: backend.URL},
		},
	}

	proxy, err := NewDocsProxy(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	router := mux.NewRouter()
	proxy.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/docs/auth/swagger.json", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestDocsProxy_DisabledRoutes(t *testing.T) {
	cfg := config.DocsSection{
		Enabled: false,
	}

	proxy, err := NewDocsProxy(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	router := mux.NewRouter()
	proxy.RegisterRoutes(router)

	// Routes should not be registered when disabled
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should get 404 since routes were not registered
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d (routes not registered), got %d", http.StatusNotFound, rec.Code)
	}
}
