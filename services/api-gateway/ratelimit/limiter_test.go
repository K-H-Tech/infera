package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucketLimiter_Allow(t *testing.T) {
	config := &Config{
		Requests:  5,
		Window:    1 * time.Minute,
		BurstSize: 2,
	}

	limiter := NewTokenBucketLimiter(config, nil)
	defer limiter.Stop()

	// Initial bucket has 5 tokens (requests), burst only adds to max capacity
	// Should allow first 5 requests
	for i := 0; i < 5; i++ {
		allowed, _ := limiter.Allow("test-key")
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be blocked
	allowed, retryAfter := limiter.Allow("test-key")
	if allowed {
		t.Error("Request 6 should be blocked")
	}
	if retryAfter == 0 {
		t.Error("RetryAfter should be > 0 when blocked")
	}
}

func TestTokenBucketLimiter_ExponentialBackoff(t *testing.T) {
	config := &Config{
		Requests:  1,
		Window:    1 * time.Second,
		BurstSize: 0,
	}

	backoffConfig := &BackoffConfig{
		Enabled:      true,
		BaseDuration: 2 * time.Second,
		MaxDuration:  10 * time.Second,
		Multiplier:   2,
	}

	limiter := NewTokenBucketLimiter(config, backoffConfig)
	defer limiter.Stop()

	// First request allowed
	allowed, _ := limiter.Allow("test-key")
	if !allowed {
		t.Error("First request should be allowed")
	}

	// Second request blocked (triggers first backoff)
	allowed, retryAfter := limiter.Allow("test-key")
	if allowed {
		t.Error("Second request should be blocked")
	}

	// First violation: 2s backoff
	expectedBackoff := 2 * time.Second
	if retryAfter < expectedBackoff || retryAfter > expectedBackoff+100*time.Millisecond {
		t.Errorf("Expected backoff ~%v, got %v", expectedBackoff, retryAfter)
	}

	// Immediate third request still blocked
	allowed, retryAfter2 := limiter.Allow("test-key")
	if allowed {
		t.Error("Third request should still be blocked")
	}
	if retryAfter2 >= retryAfter {
		t.Error("RetryAfter should decrease on subsequent checks")
	}
}

func TestTokenBucketLimiter_Reset(t *testing.T) {
	config := &Config{
		Requests:  1,
		Window:    1 * time.Minute,
		BurstSize: 0,
	}

	limiter := NewTokenBucketLimiter(config, nil)
	defer limiter.Stop()

	// Use up the token
	allowed, _ := limiter.Allow("test-key")
	if !allowed {
		t.Error("First request should be allowed")
	}

	// Second request blocked
	allowed, _ = limiter.Allow("test-key")
	if allowed {
		t.Error("Second request should be blocked")
	}

	// Reset the key
	limiter.Reset("test-key")

	// Should be allowed again
	allowed, _ = limiter.Allow("test-key")
	if !allowed {
		t.Error("Request after reset should be allowed")
	}
}

func TestTokenBucketLimiter_GetMetrics(t *testing.T) {
	config := &Config{
		Requests:  10,
		Window:    1 * time.Minute,
		BurstSize: 5,
	}

	limiter := NewTokenBucketLimiter(config, nil)
	defer limiter.Stop()

	// Use 3 tokens
	for i := 0; i < 3; i++ {
		limiter.Allow("test-key")
	}

	metrics := limiter.GetMetrics("test-key")
	if metrics.Total != 10 {
		t.Errorf("Expected total 10, got %d", metrics.Total)
	}
	if metrics.Remaining != 7 {
		t.Errorf("Expected remaining 7 (10 - 3 used), got %d", metrics.Remaining)
	}
}

func TestTokenBucketLimiter_Refill(t *testing.T) {
	config := &Config{
		Requests:  20, // 20 requests per second
		Window:    1 * time.Second,
		BurstSize: 0,
	}

	// No backoff for this test
	backoffConfig := &BackoffConfig{Enabled: false}

	limiter := NewTokenBucketLimiter(config, backoffConfig)
	defer limiter.Stop()

	// Use 15 tokens (leave 5)
	for i := 0; i < 15; i++ {
		allowed, _ := limiter.Allow("test-key")
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Wait for 500ms (should refill ~10 tokens at 20/sec rate)
	// Total available after refill: min(5 + 10, 20) = 15 tokens
	time.Sleep(550 * time.Millisecond)

	// Try to use more tokens - should get exactly 5 (remaining from before)
	// because max is capped at 20, and we had 5 left, no refill room
	successCount := 0
	for i := 0; i < 10; i++ {
		allowed, _ := limiter.Allow("test-key")
		if allowed {
			successCount++
		}
	}

	// Should get exactly 5 successful requests (the 5 that were remaining)
	if successCount != 5 {
		t.Errorf("Expected exactly 5 requests (remaining tokens), got %d", successCount)
	}
}
