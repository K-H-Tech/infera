package routes

import (
	"path"
	"strings"
)

// RouteMatcher handles matching request paths against public route patterns
type RouteMatcher struct {
	publicRoutes []string
}

// NewRouteMatcher creates a new route matcher with the given public routes
func NewRouteMatcher(publicRoutes []string) *RouteMatcher {
	return &RouteMatcher{
		publicRoutes: publicRoutes,
	}
}

// IsPublic checks if the given path matches any public route pattern
// Supports wildcard patterns using * (e.g., "/api/auth/*" matches "/api/auth/login")
func (m *RouteMatcher) IsPublic(requestPath string) bool {
	for _, pattern := range m.publicRoutes {
		if m.matchPattern(pattern, requestPath) {
			return true
		}
	}
	return false
}

// matchPattern checks if a path matches a pattern with wildcard support
// Pattern examples:
//   - "/api/auth/login" - exact match
//   - "/api/auth/*" - matches any path starting with /api/auth/
//   - "/api/*/public" - matches /api/{anything}/public
func (m *RouteMatcher) matchPattern(pattern, requestPath string) bool {
	// Exact match
	if pattern == requestPath {
		return true
	}

	// No wildcard in pattern - no match
	if !strings.Contains(pattern, "*") {
		return false
	}

	// Clean paths to ensure consistent comparison
	pattern = path.Clean(pattern)
	requestPath = path.Clean(requestPath)

	// Split into segments
	patternSegments := strings.Split(pattern, "/")
	pathSegments := strings.Split(requestPath, "/")

	// If pattern has fewer segments (except for trailing *), can't match
	if len(patternSegments) > len(pathSegments) {
		return false
	}

	// Match each segment
	for i, patternSeg := range patternSegments {
		// Wildcard matches everything from this point
		if patternSeg == "*" {
			// If it's the last segment in pattern, it matches the rest
			if i == len(patternSegments)-1 {
				return true
			}
			// Otherwise, wildcard only matches this segment
			continue
		}

		// If we've exhausted path segments but pattern continues, no match
		if i >= len(pathSegments) {
			return false
		}

		// Exact segment match required
		if patternSeg != pathSegments[i] {
			return false
		}
	}

	// Pattern matched all segments
	// Allow trailing segments in path if pattern ends with *
	if len(patternSegments) > 0 && patternSegments[len(patternSegments)-1] == "*" {
		return true
	}

	// Exact length match required if no trailing wildcard
	return len(patternSegments) == len(pathSegments)
}

// AddPublicRoute adds a new public route pattern
func (m *RouteMatcher) AddPublicRoute(pattern string) {
	m.publicRoutes = append(m.publicRoutes, pattern)
}

// GetPublicRoutes returns all configured public routes
func (m *RouteMatcher) GetPublicRoutes() []string {
	return m.publicRoutes
}
