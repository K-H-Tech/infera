package ratelimit

import (
	"sync"
	"time"
)

// Bucket represents a token bucket for rate limiting
type Bucket struct {
	Tokens        float64   // Current number of tokens
	LastRefill    time.Time // Last time tokens were refilled
	Violations    int       // Number of violations for exponential backoff
	LastViolation time.Time // Time of last violation
	BackoffUntil  time.Time // Time until which requests are blocked due to backoff
}

// Store manages in-memory storage of rate limit buckets
type Store struct {
	mu       sync.RWMutex
	buckets  map[string]*Bucket
	cleanup  time.Duration
	stopChan chan struct{}
}

// NewStore creates a new store with automatic cleanup
func NewStore(cleanupInterval time.Duration) *Store {
	store := &Store{
		buckets:  make(map[string]*Bucket),
		cleanup:  cleanupInterval,
		stopChan: make(chan struct{}),
	}

	// Start cleanup goroutine
	go store.cleanupLoop()

	return store
}

// Get retrieves a bucket for the given key, creating it if it doesn't exist
func (s *Store) Get(key string, config *Config) *Bucket {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucket, exists := s.buckets[key]
	if !exists {
		bucket = &Bucket{
			Tokens:     float64(config.Requests),
			LastRefill: time.Now(),
		}
		s.buckets[key] = bucket
	}

	return bucket
}

// Set updates a bucket for the given key
func (s *Store) Set(key string, bucket *Bucket) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buckets[key] = bucket
}

// Delete removes a bucket for the given key
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.buckets, key)
}

// cleanupLoop periodically removes expired buckets
func (s *Store) cleanupLoop() {
	ticker := time.NewTicker(s.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupExpired()
		case <-s.stopChan:
			return
		}
	}
}

// cleanupExpired removes buckets that haven't been accessed in a while
func (s *Store) cleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	expiry := 1 * time.Hour // Remove buckets not accessed in 1 hour

	for key, bucket := range s.buckets {
		if now.Sub(bucket.LastRefill) > expiry {
			delete(s.buckets, key)
		}
	}
}

// Stop stops the cleanup goroutine
func (s *Store) Stop() {
	close(s.stopChan)
}

// Size returns the number of buckets in the store
func (s *Store) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.buckets)
}
