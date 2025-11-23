package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"zarinpal-platform/services/api-gateway/validator"
)

func TestAuthMiddleware_Middleware_Disabled(t *testing.T) {
	// Setup
	mockValidator := validator.NewMockValidator()
	middleware := NewAuthMiddleware(AuthMiddlewareConfig{
		Validator:    mockValidator,
		PublicRoutes: []string{},
		Enabled:      false, // Disabled
	})

	// Create handler
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	// Test request without auth header
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should pass through without auth check
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAuthMiddleware_Middleware_PublicRoute(t *testing.T) {
	// Setup
	mockValidator := validator.NewMockValidator()
	middleware := NewAuthMiddleware(AuthMiddlewareConfig{
		Validator:    mockValidator,
		PublicRoutes: []string{"/api/auth/*"},
		Enabled:      true,
	})

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	// Test public route without auth header
	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d for public route, got %d", http.StatusOK, rec.Code)
	}
}

func TestAuthMiddleware_Middleware_MissingAuthHeader(t *testing.T) {
	// Setup
	mockValidator := validator.NewMockValidator()
	middleware := NewAuthMiddleware(AuthMiddlewareConfig{
		Validator:    mockValidator,
		PublicRoutes: []string{},
		Enabled:      true,
	})

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test protected route without auth header
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestAuthMiddleware_Middleware_InvalidAuthHeader(t *testing.T) {
	// Setup
	mockValidator := validator.NewMockValidator()
	middleware := NewAuthMiddleware(AuthMiddlewareConfig{
		Validator:    mockValidator,
		PublicRoutes: []string{},
		Enabled:      true,
	})

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name       string
		authHeader string
	}{
		{"No Bearer prefix", "token123"},
		{"Wrong scheme", "Basic token123"},
		{"Empty token", "Bearer "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
			}
		})
	}
}

func TestAuthMiddleware_Middleware_ValidToken(t *testing.T) {
	// Setup
	mockValidator := validator.NewMockValidator()
	middleware := NewAuthMiddleware(AuthMiddlewareConfig{
		Validator:    mockValidator,
		PublicRoutes: []string{},
		Enabled:      true,
	})

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	// Test with valid token
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAuthMiddleware_Middleware_InvalidToken(t *testing.T) {
	// Setup
	mockValidator := validator.NewMockValidator()
	mockValidator.ShouldFail = true
	mockValidator.Error = validator.ErrInvalidToken

	middleware := NewAuthMiddleware(AuthMiddlewareConfig{
		Validator:    mockValidator,
		PublicRoutes: []string{},
		Enabled:      true,
	})

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test with invalid token
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestAuthMiddleware_Middleware_ExpiredToken(t *testing.T) {
	// Setup
	mockValidator := validator.NewMockValidator()
	mockValidator.ShouldFail = true
	mockValidator.Error = validator.ErrExpiredToken

	middleware := NewAuthMiddleware(AuthMiddlewareConfig{
		Validator:    mockValidator,
		PublicRoutes: []string{},
		Enabled:      true,
	})

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test with expired token
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestAuthMiddleware_AddPublicRoute(t *testing.T) {
	// Setup
	mockValidator := validator.NewMockValidator()
	middleware := NewAuthMiddleware(AuthMiddlewareConfig{
		Validator:    mockValidator,
		PublicRoutes: []string{"/api/auth/*"},
		Enabled:      true,
	})

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Add new public route at runtime
	middleware.AddPublicRoute("/api/public/*")

	// Test the newly added public route
	req := httptest.NewRequest(http.MethodGet, "/api/public/info", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d for dynamically added public route, got %d", http.StatusOK, rec.Code)
	}
}

func TestAuthMiddleware_SetEnabled(t *testing.T) {
	// Setup
	mockValidator := validator.NewMockValidator()
	middleware := NewAuthMiddleware(AuthMiddlewareConfig{
		Validator:    mockValidator,
		PublicRoutes: []string{},
		Enabled:      true,
	})

	// Initially enabled
	if !middleware.IsEnabled() {
		t.Error("Expected middleware to be enabled")
	}

	// Disable at runtime
	middleware.SetEnabled(false)

	if middleware.IsEnabled() {
		t.Error("Expected middleware to be disabled")
	}

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test without auth header - should pass since disabled
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d when disabled, got %d", http.StatusOK, rec.Code)
	}
}

func TestAuthMiddleware_ExtractBearerToken(t *testing.T) {
	middleware := &AuthMiddleware{}

	tests := []struct {
		name       string
		authHeader string
		want       string
	}{
		{"Valid Bearer token", "Bearer token123", "token123"},
		{"Valid bearer lowercase", "bearer token123", "token123"},
		{"Invalid format", "token123", ""},
		{"Wrong scheme", "Basic token123", ""},
		{"Empty token", "Bearer ", ""},
		{"No token", "Bearer", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := middleware.extractBearerToken(tt.authHeader)
			if got != tt.want {
				t.Errorf("extractBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthMiddleware_ForwardClaims(t *testing.T) {
	middleware := &AuthMiddleware{}

	claims := &validator.Claims{
		UserID:    "user123",
		Mobile:    "09121234567",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	ctx := context.Background()
	ctx = middleware.forwardClaims(ctx, claims)

	// Verify metadata was added
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("Expected metadata in context")
	}

	if userIDs := md.Get("user-id"); len(userIDs) == 0 || userIDs[0] != claims.UserID {
		t.Errorf("Expected user-id %s, got %v", claims.UserID, userIDs)
	}

	if mobiles := md.Get("user-mobile"); len(mobiles) == 0 || mobiles[0] != claims.Mobile {
		t.Errorf("Expected user-mobile %s, got %v", claims.Mobile, mobiles)
	}
}
