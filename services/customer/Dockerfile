# --- Builder stage ---
FROM registry.hamrah.in/zpldocker/library/golang:1.24-alpine AS builder
RUN apk add --no-cache git make protobuf-dev protoc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


# Install Go plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.3
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.27.3

ENV PATH="$PATH:$(go env GOPATH)/bin"

# Generate protobufs
RUN make proto

WORKDIR /app/services/customer
RUN go build -o customer .


# --- Final stage ---
FROM registry.hamrah.in/zpldocker/library/alpine:3.18
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY ./services/customer/api/grpc/pb/customer.proto .

COPY --from=builder /app/services/customer/customer .

ENTRYPOINT ["./customer"]
