package domain

import (
	"context"

	"zarinpal-platform/services/customer/errors"

	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/customer/data/model"
	"zarinpal-platform/services/customer/data/repository"
)

// CustomerService defines the interface for domain operations
type CustomerService interface {
	CreateOnlineBusiness(ctx context.Context, websiteName, url, enamadID, userID string) (string, string, error)
}

type customerService struct {
	exampleRepo        repository.ExampleRepository
	onlineBusinessRepo repository.OnlineBusinessRepository
}

// NewCustomerService creates a new instance of the domain
func NewCustomerService(exampleRepo repository.ExampleRepository, onlineBusinessRepo repository.OnlineBusinessRepository) CustomerService {
	return &customerService{
		exampleRepo:        exampleRepo,
		onlineBusinessRepo: onlineBusinessRepo,
	}
}

// CreateOnlineBusiness creates a new online business profile
func (s *customerService) CreateOnlineBusiness(ctx context.Context, websiteName, url, enamadID, userID string) (string, string, error) {
	_, span := trace.GetTracer().Start(ctx, "CustomerService.CreateOnlineBusiness")
	defer span.End()

	// Create business model
	business := &model.OnlineBusiness{
		WebsiteName: websiteName,
		URL:         url,
		EnamadID:    enamadID,
		UserID:      userID,
	}

	// Save to database
	createdBusiness, err := s.onlineBusinessRepo.Create(ctx, business)
	if err != nil {
		span.RecordError(err)
		return "", "", errors.NewAppError(ctx).DatabaseError()
	}

	return createdBusiness.CustomerID, createdBusiness.ID, nil
}
