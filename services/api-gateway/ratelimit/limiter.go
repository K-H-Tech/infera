package ratelimit

import (
	"math"
	"sync"
	"time"
)

// Limiter defines the interface for rate limiting
type Limiter interface {
	Allow(key string) (allowed bool, retryAfter time.Duration)
	Reset(key string)
	GetMetrics(key string) *LimitMetrics
}

// TokenBucketLimiter implements rate limiting using token bucket algorithm
type TokenBucketLimiter struct {
	store         *Store
	config        *Config
	backoffConfig *BackoffConfig
	mu            sync.RWMutex
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
func NewTokenBucketLimiter(config *Config, backoffConfig *BackoffConfig) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		store:         NewStore(5 * time.Minute), // Cleanup every 5 minutes
		config:        config,
		backoffConfig: backoffConfig,
	}
}

// Allow checks if a request should be allowed for the given key
func (l *TokenBucketLimiter) Allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	bucket := l.store.Get(key, l.config)
	now := time.Now()

	// Check if in exponential backoff period
	if l.backoffConfig != nil && l.backoffConfig.Enabled && now.Before(bucket.BackoffUntil) {
		retryAfter := bucket.BackoffUntil.Sub(now)
		return false, retryAfter
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(bucket.LastRefill)
	refillRate := float64(l.config.Requests) / float64(l.config.Window)
	tokensToAdd := refillRate * elapsed.Seconds()

	// Calculate max tokens (requests + burst)
	maxTokens := float64(l.config.Requests + l.config.BurstSize)

	bucket.Tokens = math.Min(bucket.Tokens+tokensToAdd, maxTokens)
	bucket.LastRefill = now

	// Check if we have tokens available
	if bucket.Tokens >= 1.0 {
		bucket.Tokens -= 1.0
		l.store.Set(key, bucket)
		return true, 0
	}

	// Request blocked - record violation
	bucket.Violations++
	bucket.LastViolation = now

	// Calculate exponential backoff if enabled
	if l.backoffConfig != nil && l.backoffConfig.Enabled {
		backoff := l.calculateBackoff(bucket.Violations)
		bucket.BackoffUntil = now.Add(backoff)
		l.store.Set(key, bucket)
		return false, backoff
	}

	// Calculate retry after based on when next token will be available
	tokensNeeded := 1.0 - bucket.Tokens
	retryAfter := time.Duration(tokensNeeded/refillRate) * time.Second

	l.store.Set(key, bucket)
	return false, retryAfter
}

// calculateBackoff calculates exponential backoff duration
func (l *TokenBucketLimiter) calculateBackoff(violations int) time.Duration {
	if l.backoffConfig == nil || !l.backoffConfig.Enabled || violations == 0 {
		return 0
	}

	// Exponential: base * (multiplier ^ violations)
	multiplier := float64(l.backoffConfig.Multiplier)
	exponent := float64(violations - 1) // violations start at 1, so we use violations-1 for exponent
	backoff := time.Duration(float64(l.backoffConfig.BaseDuration) * math.Pow(multiplier, exponent))

	// Cap at max duration
	if backoff > l.backoffConfig.MaxDuration {
		backoff = l.backoffConfig.MaxDuration
	}

	return backoff
}

// Reset resets the rate limit for a given key
func (l *TokenBucketLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.store.Delete(key)
}

// GetMetrics returns current metrics for a given key
func (l *TokenBucketLimiter) GetMetrics(key string) *LimitMetrics {
	l.mu.RLock()
	defer l.mu.RUnlock()

	bucket := l.store.Get(key, l.config)
	remaining := int(math.Floor(bucket.Tokens))

	resetAt := bucket.LastRefill.Add(l.config.Window)

	return &LimitMetrics{
		Remaining:  remaining,
		Total:      l.config.Requests,
		ResetAt:    resetAt,
		Violations: bucket.Violations,
	}
}

// GetConfig returns the limiter configuration
func (l *TokenBucketLimiter) GetConfig() *Config {
	return l.config
}

// Stop stops the limiter and cleanup processes
func (l *TokenBucketLimiter) Stop() {
	l.store.Stop()
}
