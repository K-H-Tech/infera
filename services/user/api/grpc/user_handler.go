package grpc

import (
	"context"
	"zarinpal-platform/core/trace"
	pb "zarinpal-platform/services/user/api/grpc/pb/src/golang"
	"zarinpal-platform/services/user/domain"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Example handler method
func (h *UserHandler) ExampleMethod(ctx context.Context, req *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "UserHandler.ExampleMethod")
	defer span.End()

	// Call domain method
	result, err := h.service.ExampleMethod(spannedCtx, req.GetParam())
	if err != nil {
		return nil, err
	}

	return &pb.ExampleResponse{
		Result: result,
	}, nil
}

func (h *UserHandler) IsShahkarValid(ctx context.Context, req *pb.IsShahkarValidRequest) (*pb.IsShahkarValidResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "UserHandler.IsShahkarValid")
	defer span.End()

	isValid, err := h.service.IsShahkarValid(spannedCtx, req.GetMobile(), req.GetNationalCode())
	if err != nil {
		return nil, err
	}

	return &pb.IsShahkarValidResponse{
		IsValid: isValid,
	}, nil
}
