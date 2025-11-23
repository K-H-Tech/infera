package middleware

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"

	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/api-gateway/routes"
	"zarinpal-platform/services/api-gateway/validator"
)

// AuthMiddleware provides JWT authentication for the gateway
type AuthMiddleware struct {
	validator    validator.TokenValidator
	routeMatcher *routes.RouteMatcher
	enabled      bool
}

// AuthMiddlewareConfig holds configuration for the auth middleware
type AuthMiddlewareConfig struct {
	Validator    validator.TokenValidator
	PublicRoutes []string
	Enabled      bool
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(config AuthMiddlewareConfig) *AuthMiddleware {
	return &AuthMiddleware{
		validator:    config.Validator,
		routeMatcher: routes.NewRouteMatcher(config.PublicRoutes),
		enabled:      config.Enabled,
	}
}

// Middleware returns an HTTP middleware function that validates JWT tokens
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If auth is disabled, skip validation
		if !m.enabled {
			next.ServeHTTP(w, r)
			return
		}

		ctx, span := trace.GetTracer().Start(r.Context(), "AuthMiddleware.Validate")
		defer span.End()

		// Check if route is public
		if m.routeMatcher.IsPublic(r.URL.Path) {
			logger.Log.Infof("Public route, skipping auth: %s", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Log.Warnf("Missing authorization header for path: %s", r.URL.Path)
			span.RecordError(validator.ErrInvalidToken)
			http.Error(w, `{"error": "missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Parse Bearer token
		token := m.extractBearerToken(authHeader)
		if token == "" {
			logger.Log.Warnf("Invalid authorization header format for path: %s", r.URL.Path)
			span.RecordError(validator.ErrInvalidToken)
			http.Error(w, `{"error": "invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := m.validator.Validate(ctx, token)
		if err != nil {
			logger.Log.Errorf("Token validation failed for path %s: %v", r.URL.Path, err)
			span.RecordError(err)

			// Return appropriate error based on error type
			if err == validator.ErrExpiredToken {
				http.Error(w, `{"error": "token has expired"}`, http.StatusUnauthorized)
			} else if err == validator.ErrInvalidClaims {
				http.Error(w, `{"error": "invalid token claims"}`, http.StatusUnauthorized)
			} else {
				http.Error(w, `{"error": "invalid token"}`, http.StatusUnauthorized)
			}
			return
		}

		logger.Log.Infof("Token validated successfully for user_id %s, path: %s", claims.UserID, r.URL.Path)

		// Forward claims to gRPC backend via metadata
		ctx = m.forwardClaims(ctx, claims)
		r = r.WithContext(ctx)

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

// extractBearerToken extracts the token from "Bearer <token>" format
func (m *AuthMiddleware) extractBearerToken(authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	scheme := strings.ToLower(parts[0])
	if scheme != "bearer" {
		return ""
	}

	return parts[1]
}

// forwardClaims forwards user claims to gRPC backend as metadata
func (m *AuthMiddleware) forwardClaims(ctx context.Context, claims *validator.Claims) context.Context {
	md := metadata.Pairs(
		"user-id", claims.UserID,
		"user-mobile", claims.Mobile,
	)

	// Merge with existing metadata if any
	if existingMd, ok := metadata.FromOutgoingContext(ctx); ok {
		md = metadata.Join(existingMd, md)
	}

	return metadata.NewOutgoingContext(ctx, md)
}

// AddPublicRoute adds a new public route pattern at runtime
func (m *AuthMiddleware) AddPublicRoute(pattern string) {
	m.routeMatcher.AddPublicRoute(pattern)
}

// IsEnabled returns whether the middleware is enabled
func (m *AuthMiddleware) IsEnabled() bool {
	return m.enabled
}

// SetEnabled enables or disables the middleware
func (m *AuthMiddleware) SetEnabled(enabled bool) {
	m.enabled = enabled
}
