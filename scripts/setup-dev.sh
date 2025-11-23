#!/bin/bash

# Zarinpal Platform - Developer Environment Setup Script
# This script automates the setup of the development environment

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${CYAN}ℹ️  $1${NC}"
}

check_command() {
    if command -v "$1" &> /dev/null; then
        print_success "$1 is installed"
        return 0
    else
        print_warning "$1 is not installed"
        return 1
    fi
}

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "linux"
    else
        echo "unknown"
    fi
}

OS=$(detect_os)

# Welcome message
clear
print_header "Zarinpal Platform - Development Environment Setup"
echo "This script will set up your development environment for Zarinpal Platform."
echo "It will install necessary tools and configure git hooks."
echo ""
read -p "Press Enter to continue or Ctrl+C to cancel..."

# Check prerequisites
print_header "Checking Prerequisites"

# Check Go
if check_command go; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_info "Go version: $GO_VERSION"
else
    print_error "Go is not installed. Please install Go 1.20+ first."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

# Check Make
if ! check_command make; then
    print_error "Make is not installed."
    if [[ "$OS" == "macos" ]]; then
        echo "Install with: xcode-select --install"
        echo "Or: brew install make"
    else
        echo "Install with: sudo apt-get install make"
    fi
    exit 1
fi

# Check Git
if ! check_command git; then
    print_error "Git is not installed."
    exit 1
fi

# Check protoc
if ! check_command protoc; then
    print_warning "protoc is not installed."
    if [[ "$OS" == "macos" ]]; then
        read -p "Install protoc via Homebrew? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            brew install protobuf
            print_success "protoc installed"
        fi
    else
        echo "Please install protoc manually."
        echo "Visit: https://grpc.io/docs/protoc-installation/"
    fi
fi

# Install Go tools
print_header "Installing Go Tools"

print_info "Installing protoc plugins..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
print_success "Protoc plugins installed"

print_info "Installing goimports..."
go install golang.org/x/tools/cmd/goimports@latest
print_success "goimports installed"

print_info "Installing golangci-lint..."
if ! command -v golangci-lint &> /dev/null; then
    if [[ "$OS" == "macos" ]]; then
        brew install golangci-lint
    else
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    fi
    print_success "golangci-lint installed"
else
    print_success "golangci-lint already installed"
fi

print_info "Installing gosec (security scanner)..."
go install github.com/securego/gosec/v2/cmd/gosec@latest
print_success "gosec installed"

# Check PATH
print_header "Checking PATH Configuration"

GOPATH=$(go env GOPATH)
GOBIN="$GOPATH/bin"

if [[ ":$PATH:" != *":$GOBIN:"* ]]; then
    print_warning "Go bin directory is not in PATH"
    echo ""
    echo "Add this to your shell configuration file (~/.zshrc or ~/.bashrc):"
    echo ""
    echo -e "${YELLOW}export PATH=\"\$PATH:\$(go env GOPATH)/bin\"${NC}"
    echo ""

    if [[ "$SHELL" == *"zsh"* ]]; then
        SHELL_CONFIG="$HOME/.zshrc"
    else
        SHELL_CONFIG="$HOME/.bashrc"
    fi

    read -p "Add to $SHELL_CONFIG automatically? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "" >> "$SHELL_CONFIG"
        echo "# Added by Zarinpal Platform setup" >> "$SHELL_CONFIG"
        echo "export PATH=\"\$PATH:\$(go env GOPATH)/bin\"" >> "$SHELL_CONFIG"
        print_success "Added to $SHELL_CONFIG"
        print_warning "Please restart your terminal or run: source $SHELL_CONFIG"
    fi
else
    print_success "Go bin directory is in PATH"
fi

# Install git hooks
print_header "Installing Git Hooks"

if [ -d ".git" ]; then
    cp scripts/hooks/pre-commit .git/hooks/pre-commit
    cp scripts/hooks/pre-push .git/hooks/pre-push
    cp scripts/hooks/commit-msg .git/hooks/commit-msg
    chmod +x .git/hooks/pre-commit
    chmod +x .git/hooks/pre-push
    chmod +x .git/hooks/commit-msg
    print_success "Git hooks installed"

    echo ""
    print_info "Installed hooks:"
    echo "  - pre-commit:  Checks formatting and linter on staged files"
    echo "  - pre-push:    Runs tests and full linter before push"
    echo "  - commit-msg:  Validates commit message format"
else
    print_warning "Not a git repository - skipping git hooks"
fi

# Download dependencies
print_header "Downloading Go Dependencies"

print_info "Running go mod download..."
go mod download
print_success "Dependencies downloaded"

# Generate proto files
print_header "Generating Protocol Buffer Files"

read -p "Generate proto files now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    make proto
    print_success "Proto files generated"
else
    print_info "Skipped proto generation. Run 'make proto' when ready."
fi

# Setup complete
print_header "Setup Complete!"

echo ""
print_success "Your development environment is ready!"
echo ""
echo "Next steps:"
echo ""
echo "  1. If you modified shell config, restart your terminal or run:"
echo -e "     ${YELLOW}source ~/.zshrc${NC} (or ~/.bashrc)"
echo ""
echo "  2. Verify installation:"
echo -e "     ${YELLOW}which golangci-lint${NC}"
echo -e "     ${YELLOW}which protoc-gen-go${NC}"
echo ""
echo "  3. Try these commands:"
echo -e "     ${YELLOW}make help${NC}              # See all available commands"
echo -e "     ${YELLOW}make lint${NC}              # Run linter"
echo -e "     ${YELLOW}make test${NC}              # Run tests"
echo -e "     ${YELLOW}make new-service name=foo${NC}  # Create a new service"
echo ""
echo "Documentation:"
echo "  - Code Review Guidelines: docs/CODE_REVIEW.md"
echo "  - README: README.md"
echo ""
print_success "Happy coding! 🚀"
echo ""
