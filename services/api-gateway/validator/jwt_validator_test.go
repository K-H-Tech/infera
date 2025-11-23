package validator

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewJWTValidator(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		algorithm string
		wantError bool
	}{
		{"Valid HS256", "secret123", "HS256", false},
		{"Valid HS384", "secret123", "HS384", false},
		{"Valid HS512", "secret123", "HS512", false},
		{"Empty secret", "", "HS256", true},
		{"Unknown algorithm defaults to HS256", "secret123", "UNKNOWN", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tt.secret, tt.algorithm)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantError && validator == nil {
				t.Error("Expected validator but got nil")
			}
		})
	}
}

func TestJWTValidator_Validate_ValidToken(t *testing.T) {
	secret := "test-secret-key"
	validator, err := NewJWTValidator(secret, "HS256")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Create a valid token
	claims := &JWTClaims{
		UserID: "user123",
		Mobile: "09121234567",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Validate token
	ctx := context.Background()
	parsedClaims, err := validator.Validate(ctx, tokenString)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if parsedClaims.UserID != claims.UserID {
		t.Errorf("Expected UserID %s, got %s", claims.UserID, parsedClaims.UserID)
	}

	if parsedClaims.Mobile != claims.Mobile {
		t.Errorf("Expected Mobile %s, got %s", claims.Mobile, parsedClaims.Mobile)
	}
}

func TestJWTValidator_Validate_ExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	validator, err := NewJWTValidator(secret, "HS256")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Create an expired token
	claims := &JWTClaims{
		UserID: "user123",
		Mobile: "09121234567",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Validate token
	ctx := context.Background()
	_, err = validator.Validate(ctx, tokenString)
	if err != ErrExpiredToken {
		t.Errorf("Expected ErrExpiredToken, got %v", err)
	}
}

func TestJWTValidator_Validate_InvalidSignature(t *testing.T) {
	secret := "test-secret-key"
	validator, err := NewJWTValidator(secret, "HS256")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Create token with different secret
	claims := &JWTClaims{
		UserID: "user123",
		Mobile: "09121234567",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Validate token
	ctx := context.Background()
	_, err = validator.Validate(ctx, tokenString)
	if err == nil {
		t.Error("Expected error for invalid signature")
	}
}

func TestJWTValidator_Validate_MissingUserID(t *testing.T) {
	secret := "test-secret-key"
	validator, err := NewJWTValidator(secret, "HS256")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Create token without UserID
	claims := &JWTClaims{
		Mobile: "09121234567",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Validate token
	ctx := context.Background()
	_, err = validator.Validate(ctx, tokenString)
	if err != ErrInvalidClaims {
		t.Errorf("Expected ErrInvalidClaims, got %v", err)
	}
}

func TestJWTValidator_Validate_MalformedToken(t *testing.T) {
	secret := "test-secret-key"
	validator, err := NewJWTValidator(secret, "HS256")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Test with malformed tokens
	tests := []struct {
		name  string
		token string
	}{
		{"Empty string", ""},
		{"Random string", "not-a-jwt-token"},
		{"Incomplete JWT", "header.payload"},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validator.Validate(ctx, tt.token)
			if err == nil {
				t.Error("Expected error for malformed token")
			}
		})
	}
}

func TestJWTValidator_Validate_WrongSigningMethod(t *testing.T) {
	secret := "test-secret-key"
	validator, err := NewJWTValidator(secret, "HS256")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Create token with HS512 but validator expects HS256
	claims := &JWTClaims{
		UserID: "user123",
		Mobile: "09121234567",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Validate token
	ctx := context.Background()
	_, err = validator.Validate(ctx, tokenString)
	if err == nil {
		t.Error("Expected error for wrong signing method")
	}
}
