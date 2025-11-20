package domain

import (
	"context"

	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/user/data/repository"
	"zarinpal-platform/services/user/domain/shahkar"
	"zarinpal-platform/services/user/errors"
)

// UserService defines the interface for domain operations
type UserService interface {
	ExampleMethod(ctx context.Context, param string) (string, error)
	ExampleMethodForDashboard(ctx context.Context, param string) (string, error)
	ExampleMethodForBackoffice(ctx context.Context, param string) (string, error)
	IsShahkarValid(ctx context.Context, mobile string, nationalCode string) (bool, error)
}

type userService struct {
	exampleRepo repository.ExampleRepository
	shahkar     shahkar.Shahkar
}

// NewUserService creates a new instance of the domain
func NewUserService(exampleRepo repository.ExampleRepository, shahkar shahkar.Shahkar) UserService {
	return &userService{
		exampleRepo: exampleRepo,
		shahkar:     shahkar,
	}
}

// Example implementation of domain method
func (s *userService) ExampleMethod(ctx context.Context, param string) (string, error) {
	ctx, span := trace.GetTracer().Start(ctx, "UserService.ExampleMethod")
	defer span.End()

	// Add your business logic here
	return "Example response: " + param, nil
}

func (s *userService) ExampleMethodForDashboard(ctx context.Context, param string) (string, error) {
	ctx, span := trace.GetTracer().Start(ctx, "UserService.ExampleMethodForDashboard")
	defer span.End()

	// Add your business logic here
	return "Example response: " + param, nil
}

func (s *userService) ExampleMethodForBackoffice(ctx context.Context, param string) (string, error) {
	ctx, span := trace.GetTracer().Start(ctx, "UserService.ExampleMethodForBackoffice")
	defer span.End()

	// Add your business logic here
	return "Example response: " + param, nil
}

func (s *userService) IsShahkarValid(ctx context.Context, mobile string, nationalCode string) (bool, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "UserService.IsShahkarValid")
	defer span.End()

	ok, err := s.shahkar.IsShahkarValid(spannedCtx, mobile, nationalCode)
	if err != nil {
		return false, errors.NewAppError(spannedCtx).FailedConnectToShahkarError()
	}
	return ok, nil
}
