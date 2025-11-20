.PHONY: run proto new-service print
.PHONY: lint lint-fix format test test-race test-verbose coverage coverage-html coverage-check
.PHONY: install-tools install-hooks setup clean help

run:
	go run ./services/$(SERVICE)/main.go

new-service:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a service name. Usage: make new-service name=<service-name>"; \
		echo "Example: make new-service name=payment-service"; \
		exit 1; \
	fi
	@echo "Generating new service: $(name)"
	@if ! go run core/boilerplate/generate_service.go $(name); then \
		echo "Error: Failed to generate service."; \
		exit 1; \
	fi
	@echo "Service '$(name)' successfully created!"

SERVICES := $(shell find ./services -mindepth 1 -maxdepth 1 -type d)

proto:
	@for dir in $(SERVICES); do \
		if ls $$dir/api/grpc/pb/*.proto 1> /dev/null 2>&1; then \
			mkdir -p $$dir/api/grpc/pb/src/golang; \
			mkdir -p $$dir/docs; \
			protoc -I $$dir/api/grpc/pb -I core/grpc/proto/googleapis \
				-I core/grpc/proto \
				--go_out=paths=source_relative:$$dir/api/grpc/pb/src/golang \
				--go-grpc_out=paths=source_relative:$$dir/api/grpc/pb/src/golang \
				--grpc-gateway_out=paths=source_relative:$$dir/api/grpc/pb/src/golang \
				--openapiv2_out=$$dir/docs \
				$$dir/api/grpc/pb/*.proto; \
		fi; \
	done



print:
	@for dir in $(SERVICES); do \
		echo "Folder: $$dir"; \
	done

# ==============================================================================
# Code Quality
# ==============================================================================

## lint: Run golangci-lint on all code
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run --timeout=5m

## lint-fix: Run golangci-lint and auto-fix issues
lint-fix:
	@echo "🔧 Running linter with auto-fix..."
	@golangci-lint run --fix --timeout=5m

## format: Format all Go code
format:
	@echo "📐 Formatting code..."
	@gofmt -w .
	@if command -v goimports >/dev/null 2>&1; then \
		echo "📦 Organizing imports..."; \
		goimports -local zarinpal-platform -w .; \
	else \
		echo "⚠️  goimports not installed, skipping import organization"; \
	fi
	@echo "✅ Code formatted"

# ==============================================================================
# Testing
# ==============================================================================

## test: Run all tests
test:
	@echo "🧪 Running tests..."
	@go test ./... -timeout 2m

## test-race: Run tests with race detector
test-race:
	@echo "🏁 Running tests with race detector..."
	@go test ./... -race -timeout 3m

## test-verbose: Run tests in verbose mode
test-verbose:
	@echo "🧪 Running tests (verbose)..."
	@go test ./... -v -timeout 2m

## coverage: Generate test coverage report
coverage:
	@echo "📊 Generating coverage report..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out
	@echo ""
	@echo "📈 Coverage summary:"
	@go tool cover -func=coverage.out | grep total | awk '{print "  Total: " $$3}'

## coverage-html: Generate HTML coverage report
coverage-html:
	@echo "🌐 Generating HTML coverage report..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Report generated: coverage.html"
	@echo "📂 Open with: open coverage.html (macOS) or xdg-open coverage.html (Linux)"

## coverage-check: Check if coverage meets minimum threshold (70%)
coverage-check:
	@echo "✔️  Checking coverage threshold..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic > /dev/null 2>&1
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$coverage" ]; then \
		echo "⚠️  Unable to calculate coverage"; \
		exit 0; \
	fi; \
	echo "Coverage: $$coverage%"; \
	if [ $$(echo "$$coverage < 70" | bc 2>/dev/null || echo 0) -eq 1 ]; then \
		echo "❌ Coverage $$coverage% is below 70%"; \
		exit 1; \
	else \
		echo "✅ Coverage $$coverage% meets requirement"; \
	fi

# ==============================================================================
# Installation & Setup
# ==============================================================================

## install-tools: Install required development tools
install-tools:
	@echo "🔧 Installing development tools..."
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	else \
		echo "✅ golangci-lint already installed"; \
	fi
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing gosec (security scanner)..."
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "✅ All tools installed"

## install-hooks: Install git hooks
install-hooks:
	@echo "🎣 Installing git hooks..."
	@cp scripts/hooks/pre-commit .git/hooks/pre-commit
	@cp scripts/hooks/pre-push .git/hooks/pre-push
	@cp scripts/hooks/commit-msg .git/hooks/commit-msg
	@chmod +x .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-push
	@chmod +x .git/hooks/commit-msg
	@echo "✅ Git hooks installed"
	@echo ""
	@echo "Hooks installed:"
	@echo "  - pre-commit:  Checks formatting and runs linter on staged files"
	@echo "  - pre-push:    Runs tests, linter, and build before push"
	@echo "  - commit-msg:  Validates commit message format"

## setup: Complete development environment setup
setup: install-tools install-hooks
	@echo ""
	@echo "🎉 Development environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Generate proto files:  make proto"
	@echo "  2. Run tests:             make test"
	@echo "  3. Create a service:      make new-service name=my-service"
	@echo ""
	@echo "Useful commands:"
	@echo "  make help     - Show all available commands"
	@echo "  make lint     - Run code linter"
	@echo "  make format   - Format code"
	@echo "  make coverage - Check test coverage"

# ==============================================================================
# Utilities
# ==============================================================================

## clean: Clean build artifacts and temporary files
clean:
	@echo "🧹 Cleaning..."
	@rm -f coverage.out coverage.html
	@find . -type f -name '*.out' -delete
	@find . -type f -name '*.test' -delete
	@echo "✅ Cleaned"

## help: Show this help message
help:
	@echo "Zarinpal Platform - Available Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Service Management:"
	@echo "  run SERVICE=<name>  Run a specific service"
	@echo "  new-service         Create a new service (requires name=<service-name>)"
	@echo "  proto               Generate protobuf files for all services"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint                Run golangci-lint"
	@echo "  lint-fix            Run golangci-lint with auto-fix"
	@echo "  format              Format all Go code"
	@echo ""
	@echo "Testing:"
	@echo "  test                Run all tests"
	@echo "  test-race           Run tests with race detector"
	@echo "  test-verbose        Run tests in verbose mode"
	@echo "  coverage            Generate coverage report"
	@echo "  coverage-html       Generate HTML coverage report"
	@echo "  coverage-check      Check if coverage meets threshold"
	@echo ""
	@echo "Setup:"
	@echo "  install-tools       Install development tools (golangci-lint, etc.)"
	@echo "  install-hooks       Install git hooks"
	@echo "  setup               Complete development setup"
	@echo ""
	@echo "Utilities:"
	@echo "  clean               Clean build artifacts"
	@echo "  help                Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make new-service name=payment"
	@echo "  make run SERVICE=user"
	@echo "  make lint"
	@echo "  make test"