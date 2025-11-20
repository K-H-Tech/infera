package initializer

import (
	"zarinpal-platform/core"
	"zarinpal-platform/services/user-dashboard/api/grpc"
	pb "zarinpal-platform/services/user-dashboard/api/grpc/pb/src/golang"
	"zarinpal-platform/services/user-dashboard/client"
	"zarinpal-platform/services/user-dashboard/domain"
	"zarinpal-platform/services/user/config"
	userDomain "zarinpal-platform/services/user/domain"
)

type UserDashboardService struct {
	configs     *config.Config
	userService userDomain.UserService
}

func NewUserDashboardService(configs *config.Config, userService userDomain.UserService) *UserDashboardService {
	return &UserDashboardService{
		configs:     configs,
		userService: userService,
	}
}

func (i *UserDashboardService) OnStart(service *core.Service) {
	userServiceClient := client.NewUserServiceClient(i.userService)

	// Initialize domain service
	internalService := domain.NewUserDashboardService(userServiceClient)

	// Initialize gRPC handler
	handler := grpc.NewUserDashboardHandler(internalService)

	// Register gRPC server
	pb.RegisterUserDashboardServiceServer(service.Grpc.Server, handler)
}

func (i *UserDashboardService) OnStop() {
}
