package validator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/trace"
)

var (
	// ErrInvalidToken is returned when token is malformed or invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when token has expired
	ErrExpiredToken = errors.New("token has expired")
	// ErrInvalidClaims is returned when token claims are invalid or missing
	ErrInvalidClaims = errors.New("invalid token claims")
)

// JWTValidator implements TokenValidator interface using JWT tokens
type JWTValidator struct {
	secret        []byte
	signingMethod jwt.SigningMethod
}

// JWTClaims represents the custom JWT claims structure
type JWTClaims struct {
	UserID string `json:"user_id"`
	Mobile string `json:"mobile"`
	jwt.RegisteredClaims
}

// NewJWTValidator creates a new JWT validator with the given secret
// secret: The secret key used to sign/verify tokens
// algorithm: The signing algorithm (HS256, HS384, HS512, RS256, etc.)
func NewJWTValidator(secret string, algorithm string) (*JWTValidator, error) {
	if secret == "" {
		return nil, errors.New("jwt secret cannot be empty")
	}

	var method jwt.SigningMethod
	switch algorithm {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS512":
		method = jwt.SigningMethodHS512
	case "RS256":
		method = jwt.SigningMethodRS256
	case "RS384":
		method = jwt.SigningMethodRS384
	case "RS512":
		method = jwt.SigningMethodRS512
	default:
		method = jwt.SigningMethodHS256 // Default to HS256
	}

	return &JWTValidator{
		secret:        []byte(secret),
		signingMethod: method,
	}, nil
}

// Validate validates the JWT token and returns claims if valid
func (v *JWTValidator) Validate(ctx context.Context, tokenString string) (*Claims, error) {
	_, span := trace.GetTracer().Start(ctx, "JWTValidator.Validate")
	defer span.End()

	// Parse token with custom claims
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if token.Method != v.signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.secret, nil
	})

	if err != nil {
		span.RecordError(err)
		logger.Log.Errorf("Failed to parse token: %v", err)

		// Check for specific errors
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		span.RecordError(ErrInvalidClaims)
		logger.Log.Error("Invalid token claims")
		return nil, ErrInvalidClaims
	}

	// Validate required claims
	if claims.UserID == "" {
		span.RecordError(ErrInvalidClaims)
		logger.Log.Error("Missing user_id in token claims")
		return nil, ErrInvalidClaims
	}

	// Check expiration manually (additional safety check)
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		span.RecordError(ErrExpiredToken)
		logger.Log.Errorf("Token has expired: %v", claims.ExpiresAt.Time)
		return nil, ErrExpiredToken
	}

	logger.Log.Infof("Token validated successfully for user_id: %s", claims.UserID)

	return &Claims{
		UserID:    claims.UserID,
		Mobile:    claims.Mobile,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
