# API Gateway Service

Professional, optional, and pluggable API Gateway with JWT authentication and rate limiting for the Zarinpal Platform.

## Features

### Authentication
✅ **Optional Authentication** - Enable/disable via configuration
✅ **Pluggable Validators** - Interface-based design for custom token validators
✅ **Public Route Handling** - Wildcard pattern-based route whitelisting
✅ **JWT Token Validation** - Validates token signature, expiration, and claims
✅ **Per-Service Configuration** - Each proxied service can have custom auth requirements

### Rate Limiting
✅ **Optional Rate Limiting** - Enable/disable per endpoint
✅ **IP + Path Based** - Composite key for granular control
✅ **Per-Endpoint Limits** - Different limits for each route pattern
✅ **Token Bucket Algorithm** - Industry standard with burst support
✅ **Exponential Backoff** - Punishes repeat offenders automatically
✅ **429 Status Codes** - Standard HTTP Too Many Requests responses

### General
✅ **Proper Error Handling** - Returns 401/403/429 with descriptive JSON errors
✅ **Full Observability** - OpenTelemetry tracing, Prometheus metrics, structured logging
✅ **Comprehensive Testing** - Unit tests for all components

## Architecture

```
services/api-gateway/
├── config/
│   └── config.go              # Configuration management
├── middleware/
│   ├── auth.go                # JWT authentication middleware
│   ├── auth_test.go           # Middleware tests
│   ├── ratelimit.go           # Rate limiting middleware
│   └── ratelimit_test.go      # Rate limit tests
├── validator/
│   ├── validator.go           # Token validator interface
│   ├── jwt_validator.go       # JWT implementation
│   ├── jwt_validator_test.go  # Validator tests
│   └── mock_validator.go      # Mock for testing
├── ratelimit/
│   ├── limiter.go             # Rate limiter implementation
│   ├── limiter_test.go        # Limiter tests
│   ├── store.go               # In-memory token bucket store
│   ├── config.go              # Rate limit configuration
│   └── metrics.go             # Prometheus metrics
├── routes/
│   ├── routes.go              # Route matcher with wildcard support
│   └── routes_test.go         # Route matcher tests
├── initializer/
│   └── initiializer.go        # Service initialization
├── config.yaml                # Service configuration
├── main.go                    # Entry point
└── README.md                  # This file
```

### Request Flow

```
HTTP Request
    ↓
Rate Limit Middleware (check IP:Path limits)
    ↓ (if rate limit exceeded)
    → Return 429 with Retry-After header
    ↓ (if allowed)
Auth Middleware (check JWT for protected routes)
    ↓ (if public route)
    → Skip authentication
    ↓ (if protected route without valid token)
    → Return 401 Unauthorized
    ↓ (if valid)
Forward to Backend Services via gRPC
```

## Configuration

### config.yaml

```yaml
http:
  address: ":8088"
grpc:
  address: ":8087"
prometheus:
  address: ":8087"
jaeger:
  address: ":4318"

clients:
  auth:
    address: "localhost:8081"
  user:
    address: "localhost:8081"

# Authentication configuration
auth:
  # Enable/disable authentication middleware
  enabled: true

  # JWT secret key for token validation
  # Override with JWT_SECRET environment variable in production
  jwt_secret: "your-secret-key-change-in-production"

  # JWT signing algorithm (HS256, HS384, HS512, RS256, etc.)
  jwt_algorithm: "HS256"

  # Public routes that don't require authentication
  # Supports wildcard patterns with *
  public_routes:
    - "/rest/auth/otp/*"          # OTP authentication endpoints
    - "/auth/user"                 # New user registration
    - "/health"                    # Health check endpoint
    - "/metrics"                   # Prometheus metrics endpoint
```

### Environment Variables

You can override any configuration value using environment variables:

```bash
# Disable authentication
export AUTH_ENABLED=false

# Set JWT secret (IMPORTANT for production!)
export JWT_SECRET="your-production-secret-key"

# Set JWT algorithm
export JWT_ALGORITHM="HS256"

# Override HTTP port
export HTTP_ADDRESS=":8088"
```

## Usage

### Running the Gateway

```bash
# Using make
make run SERVICE=api-gateway

# Direct Go run
go run ./services/api-gateway/main.go
```

### Testing

```bash
# Run all tests
go test ./services/api-gateway/... -v

# Run specific package tests
go test ./services/api-gateway/middleware/... -v
go test ./services/api-gateway/validator/... -v
go test ./services/api-gateway/routes/... -v

# Run with coverage
go test ./services/api-gateway/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## How It Works

### 1. Request Flow

```
HTTP Request
    ↓
Auth Middleware (if enabled)
    ↓
Check if route is public
    ↓ (if public)
    → Skip authentication → Forward to backend
    ↓ (if protected)
Extract Bearer token
    ↓
Validate JWT token
    ↓
Extract claims (user_id, mobile)
    ↓
Forward claims in gRPC metadata
    ↓
Route to backend service
```

### 2. JWT Token Structure

The gateway expects JWT tokens with the following claims:

```json
{
  "user_id": "user123",
  "mobile": "09121234567",
  "exp": 1735123456
}
```

### 3. Public Route Patterns

The route matcher supports wildcard patterns:

- **Exact match**: `/rest/auth/login` - Only this specific path
- **Wildcard suffix**: `/rest/auth/*` - All paths under `/rest/auth/`
- **Deep wildcard**: `/rest/auth/otp/*` - Matches `/rest/auth/otp/verify`, `/rest/auth/otp/authenticate`, etc.

## Making Requests

### Public Endpoint (No Auth)

```bash
curl http://localhost:8088/rest/auth/otp/authenticate \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"mobile": "09121234567"}'
```

### Protected Endpoint (Requires Auth)

```bash
# First, get a token from auth service
TOKEN=$(curl http://localhost:8088/rest/auth/otp/verify \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"mobile": "09121234567", "code": "123456"}' \
  | jq -r '.access_token')

# Then use the token for protected endpoints
curl http://localhost:8088/rest/user/backoffice/profile \
  -H "Authorization: Bearer $TOKEN"
```

## Error Responses

### 401 Unauthorized - Missing Token

```json
{"error": "missing authorization header"}
```

### 401 Unauthorized - Invalid Token Format

```json
{"error": "invalid authorization header format"}
```

### 401 Unauthorized - Invalid Token

```json
{"error": "invalid token"}
```

### 401 Unauthorized - Expired Token

```json
{"error": "token has expired"}
```

### 401 Unauthorized - Invalid Claims

```json
{"error": "invalid token claims"}
```

### 429 Too Many Requests - Rate Limit Exceeded

```json
{
  "error": "rate limit exceeded",
  "retry_after_seconds": 60
}
```

**Headers returned:**
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Requests remaining in current window
- `X-RateLimit-Reset`: Unix timestamp when the limit resets
- `Retry-After`: Seconds to wait before retrying

## Rate Limiting

### How It Works

The API Gateway implements IP + Path-based rate limiting using the token bucket algorithm:

1. **Token Bucket**: Each IP:Path combination has a bucket of tokens
2. **Requests Consume Tokens**: Each request consumes 1 token
3. **Automatic Refill**: Tokens refill at a constant rate based on configuration
4. **Burst Support**: Allow temporary bursts beyond the base rate
5. **Exponential Backoff**: Repeat offenders get increasingly longer blocks

### Configuration

Rate limiting is configured per-endpoint in `config.yaml`:

```yaml
ratelimit:
  enabled: true

  # Default limit for all endpoints
  default:
    requests: 100       # 100 requests
    window: "1m"        # per minute
    burst: 20           # with 20 burst capacity

  # Per-endpoint limits (overrides default)
  endpoints:
    "/rest/auth/otp/authenticate":
      requests: 5
      window: "1m"
      burst: 2

    "/rest/user/*":     # Wildcard pattern
      requests: 200
      window: "1m"
      burst: 50

  # Exponential backoff for repeat violations
  backoff:
    enabled: true
    base_duration: "1m"   # First violation: 1 minute block
    max_duration: "1h"    # Maximum block duration
    multiplier: 2         # 1m → 2m → 4m → 8m → ...
```

### Rate Limit Headers

Every response includes rate limit information:

```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 73
X-RateLimit-Reset: 1735123456
```

When rate limited:

```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1735123456
Retry-After: 60
```

### Exponential Backoff

When a client exceeds their rate limit multiple times, backoff penalties increase:

| Violation | Backoff Duration |
|-----------|------------------|
| 1st       | 1 minute         |
| 2nd       | 2 minutes        |
| 3rd       | 4 minutes        |
| 4th       | 8 minutes        |
| 5th       | 16 minutes       |
| 6th+      | 1 hour (max)     |

After the backoff period expires, the violation count decreases.

### Testing Rate Limits

```bash
# Test rate limiting
for i in {1..10}; do
  curl -i http://localhost:8088/rest/auth/otp/authenticate \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"mobile": "09121234567"}'
  echo "---"
done

# Watch for 429 responses after exceeding the limit
```

### Disabling Rate Limiting

```yaml
# In config.yaml
ratelimit:
  enabled: false
```

Or via environment variable:

```bash
export RATELIMIT_ENABLED=false
```

### Prometheus Metrics

Rate limiting metrics are automatically tracked:

- `gateway_ratelimit_hits_total{endpoint, ip, status}` - Total rate limit checks
- `gateway_ratelimit_violations_total{endpoint, ip}` - Total violations
- `gateway_ratelimit_backoff_duration_seconds{endpoint, ip}` - Backoff durations
- `gateway_ratelimit_tokens_remaining{endpoint, ip}` - Current tokens available

View metrics at http://localhost:8087/metrics

## Extending the Gateway

### Adding a Custom Validator

Implement the `TokenValidator` interface:

```go
package validator

import "context"

type MyCustomValidator struct {
    // your fields
}

func (v *MyCustomValidator) Validate(ctx context.Context, token string) (*Claims, error) {
    // Your custom validation logic
    return &Claims{
        UserID:    "extracted-user-id",
        Mobile:    "extracted-mobile",
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }, nil
}
```

Then use it in the initializer:

```go
customValidator := &MyCustomValidator{}
authMiddleware := middleware.NewAuthMiddleware(middleware.AuthMiddlewareConfig{
    Validator:    customValidator,
    PublicRoutes: configs.Auth.PublicRoutes,
    Enabled:      true,
})
```

### Adding Public Routes Dynamically

```go
// In initializer or anywhere after middleware creation
authMiddleware.AddPublicRoute("/api/new-public/*")
```

### Disabling Auth at Runtime

```go
authMiddleware.SetEnabled(false)
```

## Security Considerations

### Production Checklist

- [ ] **Change JWT Secret**: Never use default secret in production
- [ ] **Use Environment Variables**: Store secrets in env vars, not config files
- [ ] **Use Strong Algorithms**: Prefer RS256 (asymmetric) over HS256 for production
- [ ] **Enable HTTPS**: Always use TLS in production (update `grpc.WithInsecure()`)
- [ ] **Set Token Expiration**: Keep token lifetime short (15-60 minutes)
- [ ] **Implement Token Refresh**: Add refresh token logic for better UX
- [ ] **Rate Limiting**: Add rate limiting to prevent brute force attacks
- [ ] **Logging**: Monitor authentication failures and suspicious patterns

### Security Best Practices

1. **Secrets Management**
   ```bash
   # ✅ Good
   export JWT_SECRET=$(cat /run/secrets/jwt_secret)

   # ❌ Bad
   jwt_secret: "hardcoded-secret"
   ```

2. **Token Storage** (Client-side)
   - Store tokens in `httpOnly` cookies (web)
   - Use secure storage (mobile)
   - Never store in localStorage

3. **Token Validation**
   - Always validate signature
   - Check expiration
   - Verify required claims
   - Validate issuer/audience (if using)

## Troubleshooting

### Auth is not working

1. Check if auth is enabled in config:
   ```yaml
   auth:
     enabled: true
   ```

2. Verify JWT secret matches between auth service and gateway

3. Check token format in request:
   ```
   Authorization: Bearer <token>
   ```

### Public routes still require auth

1. Verify route pattern in config.yaml matches request path
2. Check logs for "Public route, skipping auth" message
3. Test pattern matching:
   ```go
   matcher := routes.NewRouteMatcher([]string{"/rest/auth/*"})
   isPublic := matcher.IsPublic("/rest/auth/login") // should be true
   ```

### Token validation fails

1. Check token expiration
2. Verify signing algorithm matches
3. Ensure user_id claim is present
4. Check server logs for detailed error

## Monitoring

### Metrics

The gateway automatically exposes Prometheus metrics at `/metrics`:

- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request duration histogram
- `grpc_client_requests_total` - gRPC client requests
- `grpc_client_request_duration_seconds` - gRPC client duration

### Tracing

OpenTelemetry spans are automatically created for:
- `AuthMiddleware.Validate` - Token validation
- `JWTValidator.Validate` - JWT parsing and validation
- All gRPC calls to backend services

View traces in Jaeger UI at http://localhost:16686

### Logging

Structured logs include:
- Public route access (INFO)
- Token validation success (INFO)
- Missing/invalid tokens (WARN)
- Token validation failures (ERROR)

## Development

### Adding a New Backend Service

1. Add client configuration:
   ```yaml
   clients:
     myservice:
       address: "localhost:8082"
   ```

2. Register handler in initializer:
   ```go
   import myservice "zarinpal-platform/services/myservice/api/grpc/pb/src/golang"

   if err = myservice.RegisterMyServiceHandlerFromEndpoint(
       ctx, mux, configs.Clients.MyService.Address, opts); err != nil {
       logger.Log.Fatalf("failed to register myservice: %v", err)
   }
   ```

3. Run proto generation:
   ```bash
   make proto
   ```

## License

Proprietary - Zarinpal Platform

## Support

For issues or questions:
1. Check existing documentation
2. Review logs and traces
3. Contact the platform team
