package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/api-gateway/config"

	"github.com/gorilla/mux"
)

// DocsProxy handles reverse proxy for service documentation
type DocsProxy struct {
	services map[string]*serviceProxy
	enabled  bool
}

type serviceProxy struct {
	name    string
	baseURL *url.URL
	proxy   *httputil.ReverseProxy
}

// NewDocsProxy creates a new documentation reverse proxy
func NewDocsProxy(cfg config.DocsSection) (*DocsProxy, error) {
	if !cfg.Enabled {
		return &DocsProxy{enabled: false}, nil
	}

	services := make(map[string]*serviceProxy)

	for name, svc := range cfg.Services {
		baseURL, err := url.Parse(svc.URL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL for service %s: %w", name, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(baseURL)

		// Custom director to modify the request
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = baseURL.Host
		}

		// Error handler for proxy failures
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Log.Errorf("docs proxy error for service %s: %v", name, err)
			http.Error(w, fmt.Sprintf("Service %s documentation unavailable", name), http.StatusBadGateway)
		}

		services[name] = &serviceProxy{
			name:    name,
			baseURL: baseURL,
			proxy:   proxy,
		}

		logger.Log.Infof("Registered docs proxy for service: %s -> %s", name, svc.URL)
	}

	return &DocsProxy{
		services: services,
		enabled:  true,
	}, nil
}

// RegisterRoutes registers documentation routes on the router
func (p *DocsProxy) RegisterRoutes(router *mux.Router) {
	if !p.enabled {
		logger.Log.Warn("Documentation proxy is disabled")
		return
	}

	// Route: /docs/{service}/swagger-ui -> proxies to backend's /docs/swagger-ui/
	router.PathPrefix("/docs/{service}/swagger-ui").HandlerFunc(p.handleSwaggerUI)

	// Route: /docs/{service}/swagger.json -> proxies to backend's /docs/{service}.swagger.json
	router.HandleFunc("/docs/{service}/swagger.json", p.handleSwaggerJSON)

	// Route: /docs -> list available services
	router.HandleFunc("/docs", p.handleDocsList).Methods(http.MethodGet)
	router.HandleFunc("/docs/", p.handleDocsList).Methods(http.MethodGet)

	logger.Log.Infof("Registered documentation routes for %d services", len(p.services))
}

// handleSwaggerUI proxies swagger UI requests to the backend service
func (p *DocsProxy) handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, span := trace.GetTracer().Start(ctx, "DocsProxy.handleSwaggerUI")
	defer span.End()

	vars := mux.Vars(r)
	serviceName := vars["service"]

	svc, exists := p.services[serviceName]
	if !exists {
		span.SetString("error", "service not found")
		http.Error(w, fmt.Sprintf("Documentation for service '%s' not available", serviceName), http.StatusNotFound)
		return
	}

	span.SetString("service", serviceName)
	span.SetString("target", svc.baseURL.String())

	// Rewrite the path: /docs/{service}/swagger-ui/* -> /docs/swagger-ui/*
	originalPath := r.URL.Path
	newPath := strings.Replace(originalPath, fmt.Sprintf("/docs/%s/swagger-ui", serviceName), "/docs/swagger-ui/", 1)
	r.URL.Path = newPath

	svc.proxy.ServeHTTP(w, r)
}

// handleSwaggerJSON proxies swagger JSON requests to the backend service
func (p *DocsProxy) handleSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, span := trace.GetTracer().Start(ctx, "DocsProxy.handleSwaggerJSON")
	defer span.End()

	vars := mux.Vars(r)
	serviceName := vars["service"]

	svc, exists := p.services[serviceName]
	if !exists {
		span.SetString("error", "service not found")
		http.Error(w, fmt.Sprintf("Documentation for service '%s' not available", serviceName), http.StatusNotFound)
		return
	}

	span.SetString("service", serviceName)
	span.SetString("target", svc.baseURL.String())

	// Rewrite the path: /docs/{service}/swagger.json -> /docs/{service}.swagger.json
	r.URL.Path = fmt.Sprintf("/docs/%s.swagger.json", serviceName)

	svc.proxy.ServeHTTP(w, r)
}

// handleDocsList returns a list of available service documentation
func (p *DocsProxy) handleDocsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, span := trace.GetTracer().Start(ctx, "DocsProxy.handleDocsList")
	defer span.End()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #007bff; padding-bottom: 10px; }
        .service-list { list-style: none; padding: 0; }
        .service-list li { margin: 15px 0; }
        .service-list a { display: inline-block; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; transition: background 0.2s; }
        .service-list a:hover { background: #0056b3; }
        .service-name { font-weight: bold; }
        .service-links { margin-left: 10px; font-size: 0.9em; }
        .service-links a { background: #6c757d; margin-left: 5px; padding: 8px 16px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>API Gateway Documentation</h1>
        <p>Select a service to view its API documentation:</p>
        <ul class="service-list">`

	for name := range p.services {
		html += fmt.Sprintf(`
            <li>
                <span class="service-name">%s</span>
                <span class="service-links">
                    <a href="/docs/%s/swagger-ui">Swagger UI</a>
                    <a href="/docs/%s/swagger.json">OpenAPI JSON</a>
                </span>
            </li>`, name, name, name)
	}

	html += `
        </ul>
    </div>
</body>
</html>`

	_, _ = w.Write([]byte(html))
}

// GetServiceNames returns a list of registered service names
func (p *DocsProxy) GetServiceNames() []string {
	names := make([]string, 0, len(p.services))
	for name := range p.services {
		names = append(names, name)
	}
	return names
}

// IsEnabled returns whether the docs proxy is enabled
func (p *DocsProxy) IsEnabled() bool {
	return p.enabled
}
