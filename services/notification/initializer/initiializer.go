package initializer

import (
	"zarinpal-platform/core"
	"zarinpal-platform/core/logger"
	"zarinpal-platform/services/notification/api/grpc"
	notification "zarinpal-platform/services/notification/api/grpc/pb/src/golang"
	"zarinpal-platform/services/notification/config"
	"zarinpal-platform/services/notification/domain"
	"zarinpal-platform/services/notification/domain/sms"
)

type NotificationInitializer struct {
}

func (NotificationInitializer) OnStart(service *core.Service) {
	// configs
	configs := config.GetConfig()

	logger.Log.Info(configs.KavehNegar.ApiKey)

	// redis
	//cache.NewRedisCache(&cache.RedisConfig{
	//	Addr:     configs.Redis.Address,
	//	Password: configs.Redis.Password,
	//	DB:       configs.Redis.DB,
	//})

	// services
	smsService := sms.NewSmsService(configs.KavehNegar)
	notificationService := domain.NewNotificationService(smsService)

	// handlers
	notificationHandler := grpc.NewNotificationHandler(notificationService)

	// gRPC server
	notification.RegisterNotificationServiceServer(service.Grpc.Server, notificationHandler)
	service.Grpc.Start()
}

func (NotificationInitializer) OnStop() {
}
