# ZarinpalPlatform

> 🗂️ **Monorepo for Zarinpal Platform Microservices**

This repository is a monorepo containing multiple microservices for the Zarinpal platform. All services share common logic and libraries located in the `core/` directory. The `core` package provides essential utilities, configuration, logging, tracing, and other foundational components used across all services.

> 📁 **Note for Developers:**
>
> All microservices are located in the `services/` directory. To develop, update, or maintain a service, please work inside the relevant subdirectory under `services/`. The `core/` directory contains shared code and should not be modified for service-specific logic.

---

## 📑 Table of Contents
- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
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

2. ⚙️ **Generate proto files:**
   Run the following command to generate gRPC, gateway, and OpenAPI files for all services:
   ```sh
   make proto
   ```
   This will process all `.proto` files in each service and generate the necessary code and documentation.

   > ⚠️ **Important:** After making any changes to `.proto` files in any service, you must run `make proto` again to regenerate all the necessary gRPC and gateway code. This ensures your changes are properly reflected in the generated Go files.

3. 🔨 **Create a new service:**
   To generate a new service with all the necessary boilerplate code, use:
   ```sh
   make new-domain name=<domain-name>
   ```
   Replace `<service-name>` with your desired service name (e.g., `payment-service`). This will:
   - Create a new service directory under `services/`
   - Generate all required files and directory structure
   - Set up gRPC, configuration, and other necessary boilerplate code

4. 🚀 **Run a service:**
   To run a specific service, use:
   ```sh
   make run SERVICE=<domain-name>
   ```
   Replace `<service-name>` with the name of the service directory (e.g., `auth`, `notification`).

---

## Configuration

All configuration for each service should be managed by editing the `config.yaml` file located in the respective `services/<service-name>/config.yaml` directory.

---

## Troubleshooting

- If you encounter errors related to missing `protoc-gen-go` or `protoc-gen-go-grpc`, ensure you have installed the plugins and your PATH is set correctly (see Prerequisites).
- For proto generation issues, check that your `protoc` version is compatible (recommended: 3.20+).
- If you have permission errors, ensure binaries in your Go bin directory are executable (`chmod +x <binary>`).
- For other issues, consult the documentation or reach out to the team.

---

## Contributing

We welcome contributions! Please:
- Work inside the appropriate `services/<service-name>/` directory for service-specific changes.
- For shared logic, discuss changes with the team before modifying `core/`.
- Follow code style and commit guidelines.
- Open a pull request with a clear description of your changes.

---

## License

This project is proprietary and confidential. All rights reserved to Zarinpal.

---

Continue to add more sections (Architecture, Usage examples, etc.) as needed.
