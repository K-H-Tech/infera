package grpc

import (
	"context"
	"zarinpal-platform/core/trace"
	pb "zarinpal-platform/services/user-dashboard/api/grpc/pb/src/golang"
	"zarinpal-platform/services/user-dashboard/domain"
)

type UserDashboardHandler struct {
	pb.UnimplementedUserDashboardServiceServer
	service domain.UserDashboardService
}

func NewUserDashboardHandler(service domain.UserDashboardService) *UserDashboardHandler {
	return &UserDashboardHandler{
		service: service,
	}
}

// Example handler method
func (h *UserDashboardHandler) ExampleMethod(ctx context.Context, req *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "UserDashboardHandler.ExampleMethod")
	defer span.End()

	// Call domain service method
	result, err := h.service.ExampleMethod(spannedCtx, req.GetParam())
	if err != nil {
		return nil, err
	}

	return &pb.ExampleResponse{
		Result: result,
	}, nil
}
