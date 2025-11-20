package middleware

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/api-gateway/ratelimit"
)

// RateLimitMiddleware provides rate limiting for the gateway
type RateLimitMiddleware struct {
	limiters       map[string]ratelimit.Limiter // path pattern -> limiter
	defaultLimiter ratelimit.Limiter
	enabled        bool
}

// RateLimitConfig holds configuration for rate limit middleware
type RateLimitConfig struct {
	Enabled          bool
	DefaultLimiter   ratelimit.Limiter
	EndpointLimiters map[string]ratelimit.Limiter // path pattern -> limiter
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(config RateLimitConfig) *RateLimitMiddleware {
	if config.EndpointLimiters == nil {
		config.EndpointLimiters = make(map[string]ratelimit.Limiter)
	}

	return &RateLimitMiddleware{
		limiters:       config.EndpointLimiters,
		defaultLimiter: config.DefaultLimiter,
		enabled:        config.Enabled,
	}
}

// Middleware returns an HTTP middleware function that enforces rate limits
func (m *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If rate limiting is disabled, skip
		if !m.enabled {
			next.ServeHTTP(w, r)
			return
		}

		ctx, span := trace.GetTracer().Start(r.Context(), "RateLimitMiddleware.Check")
		defer span.End()

		// Extract client IP
		ip := extractClientIP(r)

		// Get limiter for this path
		limiter := m.getLimiterForPath(ctx, r.URL.Path)
		if limiter == nil {
			// No rate limit configured for this path
			next.ServeHTTP(w, r)
			return
		}

		// Create composite key: IP:Path
		key := fmt.Sprintf("%s:%s", ip, r.URL.Path)

		// Check rate limit
		allowed, retryAfter := limiter.Allow(key)

		if !allowed {
			// Log violation
			logger.Log.Warnf("Rate limit exceeded for %s on %s - retry after %v", ip, r.URL.Path, retryAfter)
			span.RecordError(fmt.Errorf("rate limit exceeded"))

			// Record metrics
			ratelimit.RecordBlock(r.URL.Path, ip, retryAfter.Seconds())

			// Get limiter config for headers
			config := m.getLimiterConfig(limiter)

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Requests))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", limiter.GetMetrics(key).ResetAt.Unix()))
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))

			// Return 429 Too Many Requests
			http.Error(w, `{"error": "rate limit exceeded", "retry_after_seconds": `+fmt.Sprintf("%.0f", retryAfter.Seconds())+`}`, http.StatusTooManyRequests)
			return
		}

		// Get current metrics
		metrics := limiter.GetMetrics(key)

		// Record allowed request
		ratelimit.RecordAllow(r.URL.Path, ip, metrics.Remaining)

		// Get limiter config for headers
		config := m.getLimiterConfig(limiter)

		// Set rate limit headers for successful requests
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Requests))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", metrics.Remaining))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", metrics.ResetAt.Unix()))

		logger.Log.Infof("Rate limit check passed for %s on %s - %d remaining", ip, r.URL.Path, metrics.Remaining)

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// getLimiterForPath finds the appropriate limiter for a given path
func (m *RateLimitMiddleware) getLimiterForPath(ctx context.Context, requestPath string) ratelimit.Limiter {
	ctx, span := trace.GetTracer().Start(ctx, "RateLimitMiddleware.getLimiterForPath")
	defer span.End()
	// Clean the path
	requestPath = path.Clean(requestPath)

	// Check for exact match first
	if limiter, exists := m.limiters[requestPath]; exists {
		return limiter
	}

	// Check for wildcard pattern matches
	for pattern, limiter := range m.limiters {
		if m.matchPattern(ctx, pattern, requestPath) {
			return limiter
		}
	}

	// Return default limiter
	return m.defaultLimiter
}

// matchPattern checks if a path matches a pattern with wildcard support
func (m *RateLimitMiddleware) matchPattern(ctx context.Context, pattern, requestPath string) bool {
	_, span := trace.GetTracer().Start(ctx, "RateLimitMiddleware.matchPattern")
	defer span.End()

	// No wildcard - exact match already checked
	if !strings.Contains(pattern, "*") {
		return false
	}

	// Clean paths
	pattern = path.Clean(pattern)
	requestPath = path.Clean(requestPath)

	// Split into segments
	patternSegments := strings.Split(pattern, "/")
	pathSegments := strings.Split(requestPath, "/")

	// If pattern has more segments (except trailing *), can't match
	if len(patternSegments) > len(pathSegments) {
		return false
	}

	// Match each segment
	for i, patternSeg := range patternSegments {
		if patternSeg == "*" {
			// Wildcard matches rest of path
			return true
		}

		if i >= len(pathSegments) || patternSeg != pathSegments[i] {
			return false
		}
	}

	// Exact match on all segments
	return len(patternSegments) == len(pathSegments)
}

// getLimiterConfig extracts config from limiter
func (m *RateLimitMiddleware) getLimiterConfig(limiter ratelimit.Limiter) *ratelimit.Config {
	// Type assert to get config
	if tbl, ok := limiter.(*ratelimit.TokenBucketLimiter); ok {
		return tbl.GetConfig()
	}
	// Fallback default
	return &ratelimit.Config{Requests: 100, Window: 60}
}

// extractClientIP extracts the real client IP from the request
func extractClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (from proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}

// IsEnabled returns whether rate limiting is enabled
func (m *RateLimitMiddleware) IsEnabled() bool {
	return m.enabled
}

// SetEnabled enables or disables rate limiting
func (m *RateLimitMiddleware) SetEnabled(enabled bool) {
	m.enabled = enabled
}

// AddEndpointLimiter adds a rate limiter for a specific endpoint pattern
func (m *RateLimitMiddleware) AddEndpointLimiter(pattern string, limiter ratelimit.Limiter) {
	m.limiters[pattern] = limiter
}
