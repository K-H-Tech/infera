package routes

import (
	"testing"
)

func TestRouteMatcher_IsPublic(t *testing.T) {
	tests := []struct {
		name         string
		publicRoutes []string
		requestPath  string
		want         bool
	}{
		{
			name:         "Exact match",
			publicRoutes: []string{"/api/auth/login"},
			requestPath:  "/api/auth/login",
			want:         true,
		},
		{
			name:         "No match",
			publicRoutes: []string{"/api/auth/login"},
			requestPath:  "/api/auth/register",
			want:         false,
		},
		{
			name:         "Wildcard suffix match",
			publicRoutes: []string{"/api/auth/*"},
			requestPath:  "/api/auth/login",
			want:         true,
		},
		{
			name:         "Wildcard suffix match deep",
			publicRoutes: []string{"/api/auth/*"},
			requestPath:  "/api/auth/otp/verify",
			want:         true,
		},
		{
			name:         "Wildcard suffix no match",
			publicRoutes: []string{"/api/auth/*"},
			requestPath:  "/api/user/profile",
			want:         false,
		},
		{
			name:         "Multiple patterns first match",
			publicRoutes: []string{"/api/auth/*", "/api/public"},
			requestPath:  "/api/auth/register",
			want:         true,
		},
		{
			name:         "Multiple patterns second match",
			publicRoutes: []string{"/api/auth/*", "/api/public"},
			requestPath:  "/api/public",
			want:         true,
		},
		{
			name:         "Multiple patterns no match",
			publicRoutes: []string{"/api/auth/*", "/api/public"},
			requestPath:  "/api/private",
			want:         false,
		},
		{
			name:         "Empty public routes",
			publicRoutes: []string{},
			requestPath:  "/api/auth/login",
			want:         false,
		},
		{
			name:         "Root path wildcard",
			publicRoutes: []string{"/*"},
			requestPath:  "/anything/here",
			want:         true,
		},
		{
			name:         "Path with trailing slash",
			publicRoutes: []string{"/api/auth/*"},
			requestPath:  "/api/auth/login/",
			want:         true,
		},
		{
			name:         "Pattern with trailing slash",
			publicRoutes: []string{"/api/auth/*/"},
			requestPath:  "/api/auth/login",
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewRouteMatcher(tt.publicRoutes)
			if got := m.IsPublic(tt.requestPath); got != tt.want {
				t.Errorf("IsPublic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteMatcher_AddPublicRoute(t *testing.T) {
	m := NewRouteMatcher([]string{"/api/auth/login"})

	m.AddPublicRoute("/api/auth/register")

	if !m.IsPublic("/api/auth/register") {
		t.Error("Added route should be public")
	}

	if !m.IsPublic("/api/auth/login") {
		t.Error("Original route should still be public")
	}
}

func TestRouteMatcher_GetPublicRoutes(t *testing.T) {
	routes := []string{"/api/auth/login", "/api/public"}
	m := NewRouteMatcher(routes)

	got := m.GetPublicRoutes()
	if len(got) != len(routes) {
		t.Errorf("GetPublicRoutes() returned %d routes, want %d", len(got), len(routes))
	}

	for i, route := range routes {
		if got[i] != route {
			t.Errorf("GetPublicRoutes()[%d] = %s, want %s", i, got[i], route)
		}
	}
}
