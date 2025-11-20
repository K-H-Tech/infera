package domain

import (
	"context"

	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/user-dashboard/client"
)

type UserDashboardService interface {
	ExampleMethod(ctx context.Context, param string) (string, error)
}

type userDashboardService struct {
	userServiceClient client.UserServiceClient
}

func NewUserDashboardService(userServiceClient client.UserServiceClient) UserDashboardService {
	return &userDashboardService{
		userServiceClient: userServiceClient,
	}
}

func (s *userDashboardService) ExampleMethod(ctx context.Context, param string) (string, error) {
	ctx, span := trace.GetTracer().Start(ctx, "user-dashboard.domain.ExampleMethod")
	defer span.End()

	err := s.userServiceClient.ExampleToDoMethod(ctx, param)
	if err != nil {
		return "", err
	}

	// Add your business logic here
	return "Example response", nil
}
