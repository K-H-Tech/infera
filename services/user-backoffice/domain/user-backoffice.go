package domain

import (
	"context"

	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/user-backoffice/client"
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

	// Call user service client
	err := s.userServiceClient.ExampleToDoMethod(ctx, param)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	// Add your business logic here
	return "Example response", nil
}
