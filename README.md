# Zarinpal Platform

> 🗂️ **Monorepo for Zarinpal Platform Microservices**

[![Go Report Card](https://goreportcard.com/badge/gitlab.hamrah.in/big-bang/zarinpal-platform)](https://goreportcard.com/report/gitlab.hamrah.in/big-bang/zarinpal-platform)
[![Coverage](https://img.shields.io/badge/coverage-check-blue)](https://gitlab.hamrah.in/big-bang/zarinpal-platform/-/graphs/main/charts)
[![License](https://img.shields.io/badge/license-Proprietary-red.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.20+-blue.svg)](https://golang.org/dl/)

This repository is a monorepo containing multiple microservices for the Zarinpal platform. All services share common logic and libraries located in the `core/` directory. The `core` package provides essential utilities, configuration, logging, tracing, and other foundational components used across all services.

> 📁 **Note for Developers:**
>
> All microservices are located in the `services/` directory. To develop, update, or maintain a service, please work inside the relevant subdirectory under `services/`. The `core/` directory contains shared code and should not be modified for service-specific logic.

---

## 📑 Table of Contents
- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Development Workflow](#development-workflow)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

---

## Getting Started

Welcome! This project requires a few tools to be installed on your system before you can build and run services. Follow the steps below for a smooth setup.

### Prerequisites

Before you begin, make sure you have the following tools:

- 🛠️ **make**: For running build and automation commands.
  - On macOS, install via Xcode Command Line Tools:
    ```sh
    xcode-select --install
    ```
    Or using Homebrew:
    ```sh
    brew install make
    ```
- 📦 **protoc** (Protocol Buffers compiler): For generating gRPC and related code from `.proto` files.
  - Install using Homebrew:
    ```sh
    brew install protobuf
    ```
- 🧰 **Go plugins for protoc**:
  To generate Go and gRPC code from proto files, install these tools:
  ```sh
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
  ```
  💡 **Tip:** After installation, add the Go bin directory to your PATH. If you use zsh, add this line to the end of your `~/.zshrc` file:
  ```sh
  export PATH="$PATH:$(go env GOPATH)/bin"
  ```
  Then open a new terminal or run:
  ```sh
  source ~/.zshrc
  ```
  To verify installation, run:
  ```sh
  which protoc-gen-go
  which protoc-gen-go-grpc
  which protoc-gen-grpc-gateway
  which protoc-gen-openapiv2
  ```
  Each command should print the path to the corresponding binary.

---

## Setup

1. 📥 **Clone the repository:**
   ```sh
   git clone https://gitlab.hamrah.in/big-bang/zarinpal-platform.git
   cd zarinpal-platform
   ```

2. 🔧 **Run setup (recommended):**
   ```sh
   make setup
   ```
   This will install all required tools (golangci-lint, goimports, etc.) and git hooks.

3. ⚙️ **Generate proto files:**
   Run the following command to generate gRPC, gateway, and OpenAPI files for all services:
   ```sh
   make proto
   ```
   This will process all `.proto` files in each service and generate the necessary code and documentation.

   > ⚠️ **Important:** After making any changes to `.proto` files in any service, you must run `make proto` again to regenerate all the necessary gRPC and gateway code. This ensures your changes are properly reflected in the generated Go files.

4. 🔨 **Create a new service:**
   To generate a new service with all the necessary boilerplate code, use:
   ```sh
   make new-service name=<service-name>
   ```
   Replace `<service-name>` with your desired service name (e.g., `payment-service`). This will:
   - Create a new service directory under `services/`
   - Generate all required files and directory structure
   - Set up gRPC, configuration, and other necessary boilerplate code

5. 🚀 **Run a service:**
   To run a specific service, use:
   ```sh
   make run SERVICE=<service-name>
   ```
   Replace `<service-name>` with the name of the service directory (e.g., `auth`, `notification`).

---

## Development Workflow

### Quick Start for Developers

```bash
# Initial setup (one-time)
make setup              # Install tools and git hooks

# Daily workflow
make format             # Format your code
make lint               # Check for issues
make test               # Run tests
make coverage           # Check test coverage
```

### Available Commands

Run `make help` to see all available commands:

**Service Management:**
- `make new-service name=<service>` - Create a new service
- `make run SERVICE=<service>` - Run a service
- `make proto` - Generate protobuf files

**Code Quality:**
- `make lint` - Run linter
- `make lint-fix` - Auto-fix linting issues
- `make format` - Format code

**Testing:**
- `make test` - Run all tests
- `make test-race` - Run tests with race detector
- `make coverage` - Generate coverage report
- `make coverage-html` - Generate HTML coverage report

**Setup:**
- `make setup` - Complete development setup
- `make install-tools` - Install development tools
- `make install-hooks` - Install git hooks

### Git Hooks

This project uses git hooks to maintain code quality:

- **pre-commit**: Checks code formatting and runs linter on staged files
- **pre-push**: Runs tests, linter, and builds before pushing
- **commit-msg**: Validates commit message format

Install hooks with:
```bash
make install-hooks
```

### Commit Message Format

Use the following format:

```
<type>/<subject>
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style changes
- `refactor` - Code refactoring
- `test` - Adding or updating tests
- `chore` - Maintenance tasks
- `perf` - Performance improvements

**Examples:**
```
feat/add user service
fix/update linter
docs/update readme
refactor/simplify error handling
```

**Rules:**
- Use lowercase
- Spaces or hyphens allowed (both work)
- Be descriptive but concise
- Max 72 characters

### Code Review

Before submitting a merge request:

1. Run `make format` to format your code
2. Run `make lint` to check for issues
3. Run `make test` to ensure tests pass
4. Read [Code Review Guidelines](docs/CODE_REVIEW.md)

---

## Project Structure

```
zarinpal-platform/
├── core/                  # Shared libraries and utilities
│   ├── boilerplate/      # Service generator
│   ├── grpc/             # gRPC utilities
│   ├── http/             # HTTP utilities
│   ├── logger/           # Logging
│   ├── trace/            # Distributed tracing
│   └── locale/           # Internationalization
├── pkg/                   # Additional shared packages
│   ├── cache/            # Caching utilities
│   ├── db/               # Database utilities
│   └── common/           # Common utilities
├── services/              # Microservices
│   ├── user/             # User service
│   ├── auth/             # Authentication service
│   ├── notification/     # Notification service
│   └── ...
├── docs/                  # Documentation
│   └── CODE_REVIEW.md    # Code review guidelines
├── scripts/               # Utility scripts
│   └── hooks/            # Git hooks
├── .golangci.yml         # Linter configuration
└── Makefile              # Build commands
```

---

## Configuration

All configuration for each service should be managed by editing the `config.yaml` file located in the respective `services/<service-name>/config.yaml` directory.

---

## Troubleshooting

- If you encounter errors related to missing `protoc-gen-go` or `protoc-gen-go-grpc`, ensure you have installed the plugins and your PATH is set correctly (see Prerequisites).
- For proto generation issues, check that your `protoc` version is compatible (recommended: 3.20+).
- If you have permission errors, ensure binaries in your Go bin directory are executable (`chmod +x <binary>`).
- Run `make setup` to automatically install all required tools.
- For other issues, consult the documentation or reach out to the team.

---

## Contributing

We welcome contributions! Please:
- Work inside the appropriate `services/<service-name>/` directory for service-specific changes.
- For shared logic, discuss changes with the team before modifying `core/`.
- Follow [code review guidelines](docs/CODE_REVIEW.md)
- Follow commit message format (see above)
- Ensure all tests pass and coverage is adequate
- Open a merge request with a clear description of your changes.

---

## License

This project is proprietary and confidential. All rights reserved to Zarinpal.

---

## Resources

- [Code Review Guidelines](docs/CODE_REVIEW.md)
- [Architecture Documentation](docs/ARCHITECTURE.md) (coming soon)
- [API Documentation](docs/API.md) (coming soon)

---

## Support

For questions or issues:
1. Check existing documentation
2. Ask the team in your communication channel
3. Create an issue in GitLab

---

**Built with ❤️ by the Zarinpal Team**
