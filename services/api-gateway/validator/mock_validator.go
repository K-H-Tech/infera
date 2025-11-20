package validator

import (
	"context"
	"time"
)

// MockValidator is a simple mock validator for testing
// Always returns valid claims for any token
type MockValidator struct {
	ShouldFail bool
	Error      error
}

// NewMockValidator creates a new mock validator
func NewMockValidator() *MockValidator {
	return &MockValidator{
		ShouldFail: false,
	}
}

// Validate always returns valid claims for testing purposes
func (m *MockValidator) Validate(ctx context.Context, token string) (*Claims, error) {
	if m.ShouldFail {
		if m.Error != nil {
			return nil, m.Error
		}
		return nil, ErrInvalidToken
	}

	return &Claims{
		UserID:    "test-user-123",
		Mobile:    "09121234567",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}
