package grpc

import (
	"context"
	"zarinpal-platform/core/trace"
	pb "zarinpal-platform/services/user-backoffice/api/grpc/pb/src/golang"
	"zarinpal-platform/services/user-backoffice/domain"
)

type UserBackofficeHandler struct {
	pb.UnimplementedUserBackofficeServiceServer
	service domain.UserBackofficeService
}

func NewUserBackofficeHandler(service domain.UserBackofficeService) *UserBackofficeHandler {
	return &UserBackofficeHandler{
		service: service,
	}
}

// Example handler method
func (h *UserBackofficeHandler) ExampleMethod(ctx context.Context, req *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "UserBackofficeHandler.ExampleMethod")
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
