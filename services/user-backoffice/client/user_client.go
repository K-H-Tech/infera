package client

import (
	"context"

	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/user/domain"
)

// UserServiceClient defines the interface for client operations
type UserServiceClient interface {
	ExampleToDoMethod(ctx context.Context, param string) error
}

type userServiceClient struct {
	userService domain.UserService
}

func NewUserServiceClient(userService domain.UserService) UserServiceClient {
	return &userServiceClient{
		userService: userService,
	}
}

func (c *userServiceClient) ExampleToDoMethod(ctx context.Context, param string) error {
	_, span := trace.GetTracer().Start(ctx, "ExampleServiceClient.ExampleToDoMethod")
	defer span.End()

	// Add implementation here
	value, err := c.userService.ExampleMethodForBackoffice(context.Background(), param)
	if err != nil {
		return err
	}
	logger.Log.Error("ExampleToDoMethod logger.", value)

	return nil
}
