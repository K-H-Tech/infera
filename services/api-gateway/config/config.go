package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Http       HttpSection       `yaml:"http"`
	Grpc       GrpcSection       `yaml:"grpc"`
	Prometheus PrometheusSection `yaml:"prometheus"`
	Jaeger     JaegerSection     `yaml:"jaeger"`
	Clients    ClientsSection    `yaml:"clients"`
	Auth       AuthSection       `yaml:"auth"`
	RateLimit  RateLimitSection  `yaml:"ratelimit"`
	Docs       DocsSection       `yaml:"docs"`
}

type HttpSection struct {
	Address string `yaml:"address"`
}

type GrpcSection struct {
	Address string `yaml:"address"`
}

type PrometheusSection struct {
	Address string `yaml:"address"`
}

type JaegerSection struct {
	Address string `yaml:"address"`
}

type ClientsSection struct {
	Auth struct {
		Address string `yaml:"address"`
	} `yaml:"auth"`
	User struct {
		Address string `yaml:"address"`
	} `yaml:"user"`
}

// AuthSection holds authentication configuration
type AuthSection struct {
	Enabled      bool     `yaml:"enabled"`
	JWTSecret    string   `yaml:"jwt_secret"`
	JWTAlgorithm string   `yaml:"jwt_algorithm"`
	PublicRoutes []string `yaml:"public_routes"`
}

// RateLimitSection holds rate limiting configuration
type RateLimitSection struct {
	Enabled   bool                               `yaml:"enabled"`
	Default   RateLimitEndpointConfig            `yaml:"default"`
	Endpoints map[string]RateLimitEndpointConfig `yaml:"endpoints"`
	Backoff   RateLimitBackoffConfig             `yaml:"backoff"`
}

// RateLimitEndpointConfig holds configuration for a single endpoint's rate limit
type RateLimitEndpointConfig struct {
	Requests int    `yaml:"requests"` // Number of requests allowed
	Window   string `yaml:"window"`   // Time window (e.g., "1m", "5m", "1h")
	Burst    int    `yaml:"burst"`    // Burst allowance
}

// RateLimitBackoffConfig holds exponential backoff configuration
type RateLimitBackoffConfig struct {
	Enabled      bool   `yaml:"enabled"`
	BaseDuration string `yaml:"base_duration"` // e.g., "1m"
	MaxDuration  string `yaml:"max_duration"`  // e.g., "1h"
	Multiplier   int    `yaml:"multiplier"`    // e.g., 2
}

// DocsSection holds documentation proxy configuration
type DocsSection struct {
	Enabled  bool                   `yaml:"enabled"`
	Services map[string]DocsService `yaml:"services"`
}

// DocsService holds the backend URL for a service's documentation
type DocsService struct {
	URL string `yaml:"url"` // Base HTTP URL of the service (e.g., "http://localhost:8080")
}

func GetConfig() *Config {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Unable to parse configuration: %s", err))
	}
	return &config
}
