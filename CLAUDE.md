# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based monorepo for the Zarinpal Platform microservices. All services share common foundational libraries in `core/` and `pkg/` directories. Each microservice is self-contained within `services/` with a consistent architecture pattern.

## Essential Commands

### Development Setup
```bash
make setup              # Install tools (golangci-lint, goimports, etc.) and git hooks
make proto              # Generate gRPC code from .proto files (required after proto changes)
```

### Running Services
```bash
make run SERVICE=<service-name>    # Example: make run SERVICE=auth
```

### Code Quality
```bash
make format             # Format code with gofmt and organize imports with goimports
make lint               # Run golangci-lint (see .golangci.yml for configuration)
make lint-fix           # Auto-fix linting issues where possible
```

### Testing
```bash
make test               # Run all tests (2m timeout)
make test-race          # Run tests with race detector (3m timeout)
make test-verbose       # Run tests in verbose mode
make coverage           # Generate coverage report (minimum threshold: 70%)
make coverage-html      # Generate HTML coverage report (coverage.html)
```

### Service Generation
```bash
make new-service name=<service-name>    # Generate new service with boilerplate
```

## High-Level Architecture

### Monorepo Structure

**Core Libraries (`core/`)**: Shared foundational code used by all services
- `core/service.go`: Main service lifecycle management (initialization, startup, shutdown)
- `core/grpc/`: gRPC server utilities with interceptors for tracing, metrics, and error handling
- `core/http/`: HTTP server with gRPC-Gateway integration for REST endpoints
- `core/configuration/`: Config loading from YAML files and environment variables (Viper-based)
- `core/logger/`: Structured logging (Logrus-based)
- `core/trace/`: OpenTelemetry distributed tracing with Jaeger integration
- `core/metric/`: Prometheus metrics integration
- `core/prometheus/`: Prometheus server management
- `core/locale/`: i18n support for error messages and localization
- `core/boilerplate/`: Service code generator

**Shared Packages (`pkg/`)**: Additional utilities
- `pkg/cache/`: Caching utilities
- `pkg/db/`: Database utilities (PostgreSQL with pgx)
- `pkg/common/`: Common utilities

**Services (`services/`)**: Each microservice follows Clean Architecture pattern:
```
services/<service-name>/
├── api/grpc/              # gRPC handlers (thin layer)
│   └── pb/                # Proto definitions and generated code
├── domain/                # Business logic and domain services
├── data/
│   ├── model/            # Database models
│   └── repository/       # Data access layer
├── config/               # Service-specific configuration
├── errors/               # Custom error definitions
├── initializer/          # Service initialization
├── locales/              # Service-specific translations
└── config.yaml           # Service configuration file
```

### Service Lifecycle (core/service.go)

The `core.StartService()` function manages the complete lifecycle:

1. **Initialization Phase**:
   - Initialize locale/i18n system
   - Load configuration from `config.yaml` (supports env var overrides)
   - Initialize structured logger
   - Setup Prometheus metrics
   - Start HTTP server (includes health checks and gRPC-Gateway)
   - Initialize Jaeger tracer for distributed tracing
   - Create gRPC server with interceptors

2. **Start Phase**: Calls `IService.OnStart(service)` where each service:
   - Registers gRPC service implementations
   - Sets up domain services and repositories
   - Initializes database connections
   - Configures HTTP routes (if custom REST endpoints needed)

3. **Graceful Shutdown**: Waits for SIGINT/SIGTERM, then calls `IService.OnStop()`

### Protocol Buffers & Code Generation

- Proto files: `services/<service>/api/grpc/pb/*.proto`
- Generated code: `services/<service>/api/grpc/pb/src/golang/`
- OpenAPI docs: `services/<service>/docs/`
- Always run `make proto` after modifying `.proto` files
- Proto files use gRPC-Gateway annotations for REST API generation

### Configuration Management

Each service uses a `config.yaml` file with this structure:
```yaml
http:
  address: ":8080"
grpc:
  address: ":9090"
prometheus:
  address: ":2112"
jaeger:
  address: "http://localhost:14268/api/traces"
```

Environment variables override YAML (e.g., `HTTP_ADDRESS=:8081`).

### Observability

**Tracing**: All domain methods should create spans:
```go
ctx, span := trace.GetTracer().Start(ctx, "ServiceName.MethodName")
defer span.End()
if err != nil {
    span.RecordError(err)
}
```

**Logging**: Use structured logging with context:
```go
logger.GetLogger().WithContext(ctx).Info("message")
```

**Metrics**: Automatically collected via Prometheus interceptors for gRPC/HTTP.

### Error Handling Pattern

- Use custom error types in `services/<service>/errors/`
- Return appropriate gRPC status codes (InvalidArgument, NotFound, Internal, etc.)
- User-facing errors use i18n via locale system
- Always wrap errors with context in spans

### Database Access

- Use PostgreSQL with `pgx/v5` driver
- Repository pattern: database logic in `data/repository/`
- Always pass `context.Context` for cancellation and tracing
- Use parameterized queries to prevent SQL injection
- Connection pooling managed by pgxpool

## Git Workflow

### Commit Message Format
```
<type>/<subject>

Types: feat, fix, docs, style, refactor, test, chore, perf
Example: feat/add user authentication
Example: fix/update linter
```

### Git Hooks (Installed via `make install-hooks`)
- **pre-commit**: Checks formatting and runs linter on staged files
- **pre-push**: Runs tests, linter, and builds before push
- **commit-msg**: Validates commit message format

### Code Review Guidelines
See `docs/CODE_REVIEW.md` for comprehensive review checklist covering:
- Clean Architecture adherence
- Security (input validation, SQL injection prevention, auth/authz)
- Observability (tracing, logging, metrics)
- Testing (>70% coverage required)
- Error handling patterns
- Database query efficiency

## Important Development Notes

### Monorepo Rules & Workflow (.clinerules)
The `.clinerules` file defines the complete workflow protocol for working with this monorepo. **Key points:**

1. **Initial Setup**: At the start of every session:
   - Read `/CLAUDE.md` (this file) first
   - Read `/.clinerules` for monorepo rules
   - Read `/README.md` for project structure

2. **Working Modes**:
   - **Platform Mode**: Say "platform" to work on ALL files in the entire repository
   - **Service Mode** (Default): Work only within a selected service directory

3. **Core Principles from .clinerules**:
   - **Honesty & Accuracy**: Always be honest, admit uncertainty, correct mistakes immediately
   - **Ask for Clarification**: Never proceed with assumptions - use clarification protocol
   - **Verify First**: Always read actual files before answering, don't assume
   - **Core is Read-Only**: In service mode, shared code in `core/` and `pkg/` should not be modified

4. **Scope Rules**:
   - Platform mode: Full repository access
   - Service mode: Only `/services/[SERVICE_NAME]/` directory
   - Can READ from `/pkg/` and `/core/` in service mode (read-only)

**For complete details, see `/.clinerules`**

### Linting Configuration
See `.golangci.yml` for enabled linters. Key points:
- Skips generated `*.pb.go` files
- Enforces error checking, security scanning (gosec), complexity limits
- Disabled for test files: gocyclo, dupl, gosec
- Use `make lint-fix` for auto-fixable issues

### Testing Requirements
- Minimum 70% coverage (enforced by `make coverage-check`)
- Test naming: `Test<Function>_<Scenario>_<Expected>`
- Use mocks for external dependencies
- Integration tests for repositories
- Unit tests for business logic

### Import Organization (Automated by goimports)
```go
import (
    // Standard library
    "context"
    "fmt"

    // External dependencies
    "google.golang.org/grpc"

    // Internal packages (local prefix: zarinpal-platform)
    "zarinpal-platform/core/trace"
    "zarinpal-platform/services/user/domain"
)
```

## Key Dependencies

- **gRPC**: google.golang.org/grpc v1.75.0
- **gRPC-Gateway**: github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2
- **Database**: github.com/jackc/pgx/v5 v5.7.6
- **Tracing**: go.opentelemetry.io/otel v1.38.0
- **Metrics**: github.com/prometheus/client_golang v1.23.2
- **Logging**: github.com/sirupsen/logrus v1.9.3
- **Config**: github.com/spf13/viper v1.21.0
- **i18n**: github.com/leonelquinteros/gotext v1.7.2
- **Cache**: github.com/redis/go-redis/v9 v9.14.0
- **Feature Flags**: github.com/growthbook/growthbook-golang v0.2.4

## Common Patterns

### Creating a New Service
1. `make new-service name=myservice`
2. Define proto files in `services/myservice/api/grpc/pb/`
3. Run `make proto`
4. Implement domain logic in `services/myservice/domain/`
5. Implement repositories in `services/myservice/data/repository/`
6. Implement gRPC handlers in `services/myservice/api/grpc/`
7. Implement `IService` interface in service's main file
8. Add configuration to `services/myservice/config.yaml`

### Service Implementation Pattern
Each service must implement `core.IService`:
```go
type MyService struct {
    // dependencies
}

func (s *MyService) OnStart(service *core.Service) {
    // Register gRPC services
    pb.RegisterMyServiceServer(service.Grpc.Server, handler)

    // Start gRPC server
    service.Grpc.Start()
}

func (s *MyService) OnStop() {
    // Cleanup resources
}
```

## Architecture Principles

1. **Clean Architecture**: Domain layer is independent of infrastructure
2. **Dependency Injection**: Dependencies injected, not created
3. **Interface-Based Design**: Use interfaces for external dependencies
4. **Context Propagation**: Always pass `context.Context` as first parameter
5. **Separation of Concerns**: Handlers are thin, business logic in domain layer
6. **Repository Pattern**: All database access through repositories
7. **Observability First**: Trace every operation, log errors, measure metrics
