package initializer

import (
	"zarinpal-platform/core"
	"zarinpal-platform/services/user-backoffice/api/grpc"
	pb "zarinpal-platform/services/user-backoffice/api/grpc/pb/src/golang"
	"zarinpal-platform/services/user-backoffice/client"
	"zarinpal-platform/services/user-backoffice/domain"
	"zarinpal-platform/services/user/config"
	userDomain "zarinpal-platform/services/user/domain"
)

type UserBackofficeService struct {
	configs     *config.Config
	userService userDomain.UserService
}

func NewUserBackofficeService(configs *config.Config, userService userDomain.UserService) *UserBackofficeService {
	return &UserBackofficeService{
		configs:     configs,
		userService: userService,
	}
}

func (i *UserBackofficeService) OnStart(service *core.Service) {
	userServiceClient := client.NewUserServiceClient(i.userService)

	// Initialize domain service
	internalService := domain.NewUserBackofficeService(userServiceClient)

	// Initialize gRPC handler
	handler := grpc.NewUserBackofficeHandler(internalService)

	// Register gRPC server
	pb.RegisterUserBackofficeServiceServer(service.Grpc.Server, handler)
}

func (i *UserBackofficeService) OnStop() {
}
