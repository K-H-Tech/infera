package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, opts ...SetOption) error
	Get(ctx context.Context, key string) (interface{}, bool, error)
}

// Option type for functional options
type SetOption func(*setOptions)

type setOptions struct {
	ttl time.Duration
}

// WithTTL sets the TTL option
func WithTTL(ttl time.Duration) SetOption {
	return func(o *setOptions) {
		o.ttl = ttl
	}
}
