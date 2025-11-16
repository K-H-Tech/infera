package initializer

import (
	"context"
	"zarinpal-platform/core"
	"zarinpal-platform/core/logger"
	"zarinpal-platform/pkg/db/postgres"
	"zarinpal-platform/services/auth/api/grpc"
	auth "zarinpal-platform/services/auth/api/grpc/pb/src/golang"
	"zarinpal-platform/services/auth/client"
	"zarinpal-platform/services/auth/config"
	"zarinpal-platform/services/auth/data/repository"
	"zarinpal-platform/services/auth/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	pgxPool *pgxpool.Pool
}

func (AuthService) OnStart(service *core.Service) {
	// configs
	configs := config.GetConfig()

	// clients
	notificationClient := client.NewNotificationServiceClient(configs.Clients.Notification.Address)
	userClient := client.NewUserServiceClient(configs.Clients.User.Address)

	pgxPool, err := postgres.NewPgxPoolConnection(context.Background(), &postgres.Config{
		User:     configs.Postgres.User,
		Password: configs.Postgres.Password,
		Host:     configs.Postgres.Host,
		Port:     configs.Postgres.Port,
		Database: configs.Postgres.Database,
	})
	if err != nil {
		logger.Log.Fatalf("failed to connect to postgres: %v", err)
	}

	// repositories
	otpRepository := repository.NewOtpRepository(pgxPool)

	// services
	otpService := domain.NewLoginWithOTPService(notificationClient, userClient, otpRepository)

	// handlers
	authHandler := grpc.NewAuthHandler(otpService)

	// gRPC server
	auth.RegisterAuthServiceServer(service.Grpc.Server, authHandler)
	service.Grpc.Start()
}

func (i AuthService) OnStop() {
	if i.pgxPool != nil {
		i.pgxPool.Close()
	}
}
