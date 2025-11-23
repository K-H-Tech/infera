package validator

import (
	"context"
	"time"
)

// Claims represents the JWT token claims
type Claims struct {
	UserID    string
	Mobile    string
	ExpiresAt time.Time
}

// TokenValidator defines the interface for token validation
// This interface allows for pluggable implementations (JWT, OAuth, custom, etc.)
type TokenValidator interface {
	// Validate validates a token and returns the claims if valid
	// Returns error if token is invalid, expired, or malformed
	Validate(ctx context.Context, token string) (*Claims, error)
}
