package initializer

import (
	"context"
	"fmt"
	"time"

	"zarinpal-platform/core"
	"zarinpal-platform/core/logger"
	"zarinpal-platform/services/api-gateway/config"
	"zarinpal-platform/services/api-gateway/middleware"
	"zarinpal-platform/services/api-gateway/ratelimit"
	"zarinpal-platform/services/api-gateway/validator"
	auth "zarinpal-platform/services/auth/api/grpc/pb/src/golang"
	userbackoffice "zarinpal-platform/services/user-backoffice/api/grpc/pb/src/golang"
	userdashboard "zarinpal-platform/services/user-dashboard/api/grpc/pb/src/golang"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type APIGatewayService struct {
}

func (s APIGatewayService) OnStart(service *core.Service) {
	ctx := context.Background()
	mux := service.Http.GatewayMux
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Load configuration
	configs := config.GetConfig()

	// Register gRPC service handlers first
	var err error
	if err = auth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, configs.Clients.Auth.Address, opts); err != nil {
		logger.Log.Fatalf("failed to register auth domain grpc handler: %v", err)
	}

	if err = userbackoffice.RegisterUserBackofficeServiceHandlerFromEndpoint(ctx, mux, configs.Clients.User.Address, opts); err != nil {
		logger.Log.Fatalf("failed to register user-backoffice domain grpc handler: %v", err)
	}

	if err = userdashboard.RegisterUserDashboardServiceHandlerFromEndpoint(ctx, mux, configs.Clients.User.Address, opts); err != nil {
		logger.Log.Fatalf("failed to register user-dashboard domain grpc handler: %v", err)
	}

	// Setup rate limiting middleware if enabled (applies before auth)
	if configs.RateLimit.Enabled {
		logger.Log.Info("Rate limiting middleware is ENABLED")

		// Parse backoff configuration
		backoffConfig, err := parseBackoffConfig(configs.RateLimit.Backoff)
		if err != nil {
			logger.Log.Fatalf("failed to parse rate limit backoff config: %v", err)
		}

		// Create default limiter
		defaultLimiterConfig, err := parseRateLimitConfig(configs.RateLimit.Default)
		if err != nil {
			logger.Log.Fatalf("failed to parse default rate limit config: %v", err)
		}
		defaultLimiter := ratelimit.NewTokenBucketLimiter(defaultLimiterConfig, backoffConfig)

		// Create endpoint-specific limiters
		endpointLimiters := make(map[string]ratelimit.Limiter)
		for pattern, endpointConfig := range configs.RateLimit.Endpoints {
			limiterConfig, err := parseRateLimitConfig(endpointConfig)
			if err != nil {
				logger.Log.Warnf("failed to parse rate limit config for %s: %v, using default", pattern, err)
				continue
			}
			endpointLimiters[pattern] = ratelimit.NewTokenBucketLimiter(limiterConfig, backoffConfig)
			logger.Log.Infof("  - Rate limit for %s: %d req/%s (burst: %d)",
				pattern, endpointConfig.Requests, endpointConfig.Window, endpointConfig.Burst)
		}

		// Create rate limit middleware
		rateLimitMiddleware := middleware.NewRateLimitMiddleware(middleware.RateLimitConfig{
			Enabled:          true,
			DefaultLimiter:   defaultLimiter,
			EndpointLimiters: endpointLimiters,
		})

		// Apply rate limit middleware (executes before auth)
		service.Http.Engine.Use(rateLimitMiddleware.Middleware)

		logger.Log.Infof("Rate limiting configured: default=%d req/%s, %d endpoint-specific rules",
			configs.RateLimit.Default.Requests, configs.RateLimit.Default.Window, len(endpointLimiters))
	} else {
		logger.Log.Warn("Rate limiting middleware is DISABLED")
	}

	// Setup authentication middleware if enabled (after registering handlers)
	if configs.Auth.Enabled {
		logger.Log.Info("Authentication middleware is ENABLED")

		// Create JWT validator
		jwtValidator, err := validator.NewJWTValidator(configs.Auth.JWTSecret, configs.Auth.JWTAlgorithm)
		if err != nil {
			logger.Log.Fatalf("failed to create JWT validator: %v", err)
		}

		// Create auth middleware with configuration
		authMiddleware := middleware.NewAuthMiddleware(middleware.AuthMiddlewareConfig{
			Validator:    jwtValidator,
			PublicRoutes: configs.Auth.PublicRoutes,
			Enabled:      true,
		})

		// Apply auth middleware to Gorilla router using Use()
		service.Http.Engine.Use(authMiddleware.Middleware)

		logger.Log.Infof("Auth middleware configured with %d public routes", len(configs.Auth.PublicRoutes))
		for _, route := range configs.Auth.PublicRoutes {
			logger.Log.Infof("  - Public route: %s", route)
		}
	} else {
		logger.Log.Warn("Authentication middleware is DISABLED - all routes are public")
	}

	logger.Log.Info("API Gateway initialization complete")
}

func (s APIGatewayService) OnStop() {
}

// parseRateLimitConfig parses rate limit configuration from config structure
func parseRateLimitConfig(cfg config.RateLimitEndpointConfig) (*ratelimit.Config, error) {
	// Parse window duration
	window, err := time.ParseDuration(cfg.Window)
	if err != nil {
		return nil, fmt.Errorf("invalid window duration '%s': %w", cfg.Window, err)
	}

	return &ratelimit.Config{
		Requests:  cfg.Requests,
		Window:    window,
		BurstSize: cfg.Burst,
	}, nil
}

// parseBackoffConfig parses backoff configuration from config structure
func parseBackoffConfig(cfg config.RateLimitBackoffConfig) (*ratelimit.BackoffConfig, error) {
	if !cfg.Enabled {
		return &ratelimit.BackoffConfig{Enabled: false}, nil
	}

	baseDuration, err := time.ParseDuration(cfg.BaseDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid base duration '%s': %w", cfg.BaseDuration, err)
	}

	maxDuration, err := time.ParseDuration(cfg.MaxDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid max duration '%s': %w", cfg.MaxDuration, err)
	}

	return &ratelimit.BackoffConfig{
		Enabled:      true,
		BaseDuration: baseDuration,
		MaxDuration:  maxDuration,
		Multiplier:   cfg.Multiplier,
	}, nil
}
