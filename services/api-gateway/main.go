package main

import (
	"zarinpal-platform/core"
	"zarinpal-platform/services/api-gateway/initializer"
)

func main() {
	core.StartService("api-gateway", initializer.APIGatewayService{})
}
