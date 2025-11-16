package main

import (
	"zarinpal-platform/core"
	"zarinpal-platform/services/grpc-gateway/initializer"
)

func main() {
	core.StartService("grpc-gateway", initializer.GrpcGatewayService{})
}
