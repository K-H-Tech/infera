package domain

import (
	"context"

	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/user-backoffice/client"
	"zarinpal-platform/services/user/errors"
)

// UserBackofficeService defines the interface for service operations
type UserBackofficeService interface {
	// Define your service methods here
	ExampleMethod(ctx context.Context, param string) (string, error)
}

type UserBackofficeServiceImpl struct {
	userServiceClient client.UserServiceClient
}

// NewUserBackofficeService creates a new instance of the service
func NewUserBackofficeService(userServiceClient client.UserServiceClient) UserBackofficeService {
	return &UserBackofficeServiceImpl{
		userServiceClient: userServiceClient,
	}
}

// Example implementation of service method
func (s *UserBackofficeServiceImpl) ExampleMethod(ctx context.Context, param string) (string, error) {
	ctx, span := trace.GetTracer().Start(ctx, "user-backoffice.domain.ExampleMethod")
	defer span.End()

	return "", errors.NewAppError(ctx).InvalidArgumentError()

	err := s.userServiceClient.ExampleToDoMethod(ctx, param)
	if err != nil {
		return "", err
	}

	// Add your business logic here
	return "Example response", nil
}
