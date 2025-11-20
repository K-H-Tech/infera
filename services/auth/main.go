package main

import (
	"zarinpal-platform/core"
	"zarinpal-platform/services/auth/initializer"
)

func main() {
	core.StartService("auth", initializer.AuthService{})
}
