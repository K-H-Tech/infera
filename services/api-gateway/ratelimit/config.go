package ratelimit

import "time"

// Config holds rate limiting configuration for a specific limiter
type Config struct {
	Requests  int
	Window    time.Duration
	BurstSize int
}

// BackoffConfig holds exponential backoff configuration
type BackoffConfig struct {
	Enabled      bool
	BaseDuration time.Duration
	MaxDuration  time.Duration
	Multiplier   int
}

// LimitMetrics holds metrics for a specific rate limit key
type LimitMetrics struct {
	Remaining  int
	Total      int
	ResetAt    time.Time
	Violations int
}
