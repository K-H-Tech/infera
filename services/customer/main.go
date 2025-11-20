package main

import (
	"zarinpal-platform/core"
	"zarinpal-platform/services/customer/initializer"
)

func main() {
	core.StartService("customer", &initializer.CustomerService{})
}
