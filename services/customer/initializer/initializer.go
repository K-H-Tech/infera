package initializer

import (
	"context"

	"zarinpal-platform/core"
	"zarinpal-platform/core/logger"
	"zarinpal-platform/pkg/db/postgres"
	"zarinpal-platform/services/customer/api/grpc"
	pb "zarinpal-platform/services/customer/api/grpc/pb/src/golang"
	"zarinpal-platform/services/customer/config"
	"zarinpal-platform/services/customer/data/repository"
	"zarinpal-platform/services/customer/domain"

	gb "github.com/growthbook/growthbook-golang"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerService struct {
	pgxPool          *pgxpool.Pool
	growthbookClient *gb.Client
}

func (e *CustomerService) OnStart(service *core.Service) {
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
	onlineBusinessRepo := repository.NewOnlineBusinessRepository(pgxPool)

	// Initialize domain service
	customerService := domain.NewCustomerService(exampleRepo, onlineBusinessRepo)

	// Initialize gRPC handler
	handler := grpc.NewCustomerHandler(customerService)

	// Register gRPC server
	pb.RegisterCustomerServiceServer(service.Grpc.Server, handler)
	service.Grpc.Start()
}

func (e *CustomerService) OnStop() {
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
