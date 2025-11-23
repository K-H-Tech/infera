package initializer

import (
	"context"

	"zarinpal-platform/core"
	"zarinpal-platform/core/logger"
	"zarinpal-platform/pkg/db/postgres"
	"zarinpal-platform/services/user/api/grpc"
	pb "zarinpal-platform/services/user/api/grpc/pb/src/golang"
	"zarinpal-platform/services/user/config"
	"zarinpal-platform/services/user/data/repository"
	"zarinpal-platform/services/user/domain"
	"zarinpal-platform/services/user/domain/shahkar"

	gb "github.com/growthbook/growthbook-golang"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	pgxPool          *pgxpool.Pool
	growthbookClient *gb.Client
}

func (e *UserService) OnStart(service *core.Service) {
	configs := config.GetConfig()
	var err error

	// growthbook
	e.growthbookClient, err = gb.NewClient(context.TODO(),
		gb.WithClientKey(configs.GrowthBook.ClientKey),
		gb.WithApiHost(configs.GrowthBook.APIHost),
		gb.WithPollDataSource(configs.GrowthBook.PollDatasourceTime))
	if err != nil {
		logger.Log.Fatalf("failed to create growthbook client: %v", err)
	}

	// postgres connection
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
	e.pgxPool = pgxPool

	exampleRepo := repository.NewExampleRepository(pgxPool)

	// modules
	shahkarModule := shahkar.NewShahkar()

	// Initialize domain domain
	userService := domain.NewUserService(exampleRepo, shahkarModule)

	// Initialize gRPC handler
	handler := grpc.NewUserHandler(userService)

	// Register gRPC server
	pb.RegisterUserServiceServer(service.Grpc.Server, handler)
	service.Grpc.Start()
}

func (e *UserService) OnStop() {
	if e.pgxPool != nil {
		e.pgxPool.Close()
	}

	if e.growthbookClient != nil {
		err := e.growthbookClient.Close()
		if err != nil {
			logger.Log.Errorf("failed to close growthbook client: %v", err)
		}
	}
}
