package initializer

import (
	"context"
	"net/http"

	"zarinpal-platform/core"
	"zarinpal-platform/core/logger"
	auth "zarinpal-platform/services/auth/api/grpc/pb/src/golang"
	"zarinpal-platform/services/grpc-gateway/config"
	userbackoffice "zarinpal-platform/services/user-backoffice/api/grpc/pb/src/golang"
	userdashboard "zarinpal-platform/services/user-dashboard/api/grpc/pb/src/golang"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GrpcGatewayService struct {
}

func (s GrpcGatewayService) OnStart(service *core.Service) {
	ctx := context.Background()
	mux := service.Http.GatewayMux
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// configs
	configs := config.GetConfig()

	var err error
	if err = auth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, configs.Clients.Auth.Address, opts); err != nil {
		logger.Log.Fatalf("failed to register auth domain grpc handler: %v", err)
	}

	if err = userbackoffice.RegisterUserBackofficeServiceHandlerFromEndpoint(ctx, mux, configs.Clients.User.Address, opts); err != nil {
		logger.Log.Fatalf("failed to register user-backoffice domain grpc handler: %v", err)
	}
	if err = userdashboard.RegisterUserDashboardServiceHandlerFromEndpoint(ctx, mux, configs.Clients.User.Address, opts); err != nil {
		logger.Log.Fatalf("failed to register user-dashboard domain grpc handler: %v", err)
	}
}

func forwardMetadata(ctx context.Context, req *http.Request) metadata.MD {
	// grpc-gateway automatically converts headers to lowercase
	acceptLang := req.Header.Get("Accept-Language")
	if acceptLang != "" {
		return metadata.Pairs("accept-language", acceptLang)
	}
	return metadata.MD{}
}

func (s GrpcGatewayService) OnStop() {
}
