package main

import (
	"zarinpal-platform/core"
	"zarinpal-platform/services/user/initializer"
)

func main() {
	core.StartService("user", &initializer.UserService{})
}
