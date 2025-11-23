# ZarinpalPlatform - Monorepo for Zarinpal Platform Microservices

## Project Overview

This is a Go-based monorepo containing multiple microservices for the Zarinpal platform. The repository follows a microservices architecture where all services share common logic and libraries located in the `core/` directory. The `core` package provides essential utilities including configuration, logging, tracing, gRPC setup, HTTP server, and other foundational components used across all services.

### Architecture
- **Monorepo Structure**: Contains multiple microservices in the `services/` directory
- **Shared Core Library**: Common components in `core/` directory for all services
- **Service Structure**: Each service follows a standardized domain-driven design (DDD) structure
- **gRPC Communication**: Services communicate using gRPC with Protocol Buffers
- **OpenTelemetry**: Built-in tracing with Jaeger and metrics with Prometheus

### Service Structure
Each service follows this directory structure:
```
services/<service-name>/
├── api/                   # API definitions (gRPC, HTTP)
├── client/                # Service clients for inter-service communication
├── config/                # Configuration files and loaders
├── data/                  # Data access layer (repositories, database connection)
├── domain/                # Business logic and domain services
├── errors/                # Custom error definitions
├── initializer/           # Service initialization logic
├── config.yaml            # Service configuration
├── Dockerfile             # Containerization
├── main.go                # Entry point
└── <service-name>.ci.yml  # CI configuration
```

### Core Components
The `core/` directory contains:
- `configuration/` - Configuration management using Viper
- `grpc/` - gRPC server setup and utilities
- `http/` - HTTP server and routing
- `logger/` - Logging utilities
- `metric/` - Metrics collection
- `prometheus/` - Prometheus metrics exporter
- `trace/` - OpenTelemetry tracing with Jaeger integration
- `locale/` - Localization and translation support
- `service.go` - Main service initialization logic

## Building and Running

### Prerequisites
- **Go 1.23.5** or higher
- **make** - For build automation
- **protoc** - Protocol Buffers compiler
- **Go plugins for protoc**: 
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`
  - `protoc-gen-grpc-gateway`
  - `protoc-gen-openapiv2`

### Key Commands
```bash
# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Generate proto files (run after changing .proto files)
make proto

# Create a new service with boilerplate
make new-service name=<service-name>

# Run a specific service
make run SERVICE=<service-name>
```

### Development Workflow
1. **Service Generation**: Use `make new-service name=<service-name>` to create a new service with all necessary boilerplate
2. **Proto Generation**: Use `make proto` to generate gRPC code from `.proto` files after any changes
3. **Running Services**: Use `make run SERVICE=<service-name>` to run individual services

## Development Conventions

### Code Structure
- Follow Domain-Driven Design (DDD) principles
- Separate concerns into layers: API, Client, Config, Data, Domain
- Use interfaces for dependency injection and testing
- Implement the `IService` interface for service lifecycle management

### Configuration Management
- Store service-specific configuration in `services/<service-name>/config.yaml`
- Use viper for configuration management
- Support for multiple configuration sources (file, environment, etc.)

### Logging and Tracing
- Use the centralized logger from the core package
- Automatic OpenTelemetry tracing with Jaeger
- Prometheus metrics collection built-in
- Structured logging with consistent format

### Service Lifecycle
- Each service must implement the `IService` interface (`OnStart` and `OnStop` methods)
- Proper resource cleanup in `OnStop` method
- Graceful shutdown handling

### Inter-Service Communication
- gRPC is the primary communication protocol
- Generated clients in the `client/` directory for each service
- Automatic gateway generation for REST endpoints

## Project Dependencies

The project uses several key Go libraries:
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/redis/go-redis/v9` - Redis client
- `google.golang.org/grpc` - gRPC framework
- `github.com/grpc-ecosystem/grpc-gateway/v2` - gRPC to JSON proxy
- `go.opentelemetry.io/otel` - OpenTelemetry for observability
- `github.com/sirupsen/logrus` - Logging
- `github.com/spf13/viper` - Configuration management

## Testing Strategy
While not explicitly shown in the basic files, the architecture suggests:
- Domain logic testing in the `domain/` packages
- Integration tests for gRPC and HTTP endpoints
- Repository pattern allows for easier testing with mocks

## Troubleshooting

- **Proto generation issues**: Ensure `protoc` and all Go plugins are properly installed and in PATH
- **Import errors**: Ensure Go modules are properly initialized and dependencies are downloaded
- **Runtime issues**: Check config.yaml files for each service are properly configured
- **Service communication**: Verify service addresses and network connectivity between services