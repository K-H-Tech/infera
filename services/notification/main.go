package main

import (
	"zarinpal-platform/core"
	"zarinpal-platform/services/notification/initializer"
)

func main() {
	core.StartService("notification", initializer.NotificationInitializer{})
}
