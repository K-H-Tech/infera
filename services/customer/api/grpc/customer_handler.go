package grpc

import (
	"context"

	"zarinpal-platform/core/logger"
	"zarinpal-platform/services/customer/errors"

	"zarinpal-platform/core/trace"
	pb "zarinpal-platform/services/customer/api/grpc/pb/src/golang"
	"zarinpal-platform/services/customer/domain"
)

type CustomerHandler struct {
	pb.UnimplementedCustomerServiceServer
	service domain.CustomerService
}

func NewCustomerHandler(service domain.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		service: service,
	}
}

// validateField checks if a field is empty and returns an error message
func validateField(ctx context.Context, fieldValue, fieldName string) (string, error) {
	if fieldValue == "" {
		err := errors.NewAppError(ctx).InvalidArgumentError()
		logger.Log.Errorf("%s is required", fieldName)
		return fieldName + " is required", err
	}
	return "", nil
}

// CreateOnlineBusiness handles the creation of a new online business
func (h *CustomerHandler) CreateOnlineBusiness(ctx context.Context, req *pb.CreateOnlineBusinessRequest) (*pb.CreateOnlineBusinessResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "CustomerHandler.CreateOnlineBusiness")
	defer span.End()

	// Validate input
	if req.GetWebsiteName() == "" || req.GetUrl() == "" || req.GetUserId() == "" {
		err := errors.NewAppError(ctx).InvalidArgumentError()
		span.RecordError(err)
		logger.Log.Errorf("invalid request: %v", req)
		return nil, err
	}

	// Call domain method
	customerID, businessID, err := h.service.CreateOnlineBusiness(
		spannedCtx,
		req.GetWebsiteName(),
		req.GetUrl(),
		req.GetEnamadId(),
		req.GetUserId(),
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOnlineBusinessResponse{
		CustomerId: customerID,
		BusinessId: businessID,
	}, nil
}

// UpdateCorporateCustomerNationalID handles the update of corporate customer national ID
func (h *CustomerHandler) UpdateCorporateCustomerNationalID(ctx context.Context, req *pb.UpdateCorporateCustomerNationalIDRequest) (*pb.UpdateCorporateCustomerNationalIDResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "CustomerHandler.UpdateCorporateCustomerNationalID")
	defer span.End()

	// Validate input
	if errMsg, err := validateField(spannedCtx, req.GetCustomerId(), "customer_id"); err != nil {
		span.RecordError(err)
		return &pb.UpdateCorporateCustomerNationalIDResponse{Error: errMsg}, nil
	}

	if errMsg, err := validateField(spannedCtx, req.GetNationalId(), "national_id"); err != nil {
		span.RecordError(err)
		return &pb.UpdateCorporateCustomerNationalIDResponse{Error: errMsg}, nil
	}

	// TODO: Implement domain logic for updating corporate customer national ID
	logger.Log.Infof("Updating corporate customer (ID: %s) national ID to: %s", req.GetCustomerId(), req.GetNationalId())

	return &pb.UpdateCorporateCustomerNationalIDResponse{
		Error: "",
	}, nil
}

// SetCustomerAsIndividual handles setting a customer as individual type
func (h *CustomerHandler) SetCustomerAsIndividual(ctx context.Context, req *pb.SetCustomerAsIndividualRequest) (*pb.SetCustomerAsIndividualResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "CustomerHandler.SetCustomerAsIndividual")
	defer span.End()

	// Validate input
	if errMsg, err := validateField(spannedCtx, req.GetUserId(), "user_id"); err != nil {
		span.RecordError(err)
		return &pb.SetCustomerAsIndividualResponse{Error: errMsg}, nil
	}

	if errMsg, err := validateField(spannedCtx, req.GetCustomerId(), "customer_id"); err != nil {
		span.RecordError(err)
		return &pb.SetCustomerAsIndividualResponse{Error: errMsg}, nil
	}

	// TODO: Implement domain logic for setting customer as individual
	logger.Log.Infof("Setting customer (ID: %s) as individual for user (ID: %s)", req.GetCustomerId(), req.GetUserId())

	return &pb.SetCustomerAsIndividualResponse{
		Error: "",
	}, nil
}

// UpdateBusinessInfo handles updating business information
func (h *CustomerHandler) UpdateBusinessInfo(ctx context.Context, req *pb.UpdateBusinessInfoRequest) (*pb.UpdateBusinessInfoResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "CustomerHandler.UpdateBusinessInfo")
	defer span.End()

	// Validate input
	if req.GetGuild() == "" {
		err := errors.NewAppError(spannedCtx).InvalidArgumentError()
		span.RecordError(err)
		logger.Log.Errorf("guild is required")
		return &pb.UpdateBusinessInfoResponse{
			Error: "guild is required",
		}, nil
	}

	if req.GetLicenseId() == "" {
		err := errors.NewAppError(spannedCtx).InvalidArgumentError()
		span.RecordError(err)
		logger.Log.Errorf("license_id is required")
		return &pb.UpdateBusinessInfoResponse{
			Error: "license_id is required",
		}, nil
	}

	if req.GetPostalCode() == "" {
		err := errors.NewAppError(spannedCtx).InvalidArgumentError()
		span.RecordError(err)
		logger.Log.Errorf("postal_code is required")
		return &pb.UpdateBusinessInfoResponse{
			Error: "postal_code is required",
		}, nil
	}

	if req.GetPhoneTech() == "" {
		err := errors.NewAppError(spannedCtx).InvalidArgumentError()
		span.RecordError(err)
		logger.Log.Errorf("phone_tech is required")
		return &pb.UpdateBusinessInfoResponse{
			Error: "phone_tech is required",
		}, nil
	}

	if req.GetPhoneBusiness() == "" {
		err := errors.NewAppError(spannedCtx).InvalidArgumentError()
		span.RecordError(err)
		logger.Log.Errorf("phone_business is required")
		return &pb.UpdateBusinessInfoResponse{
			Error: "phone_business is required",
		}, nil
	}

	// TODO: Implement domain logic for updating business info
	logger.Log.Infof("Updating business info - Guild: %s, License: %s, Postal: %s, PhoneTech: %s, PhoneBusiness: %s",
		req.GetGuild(), req.GetLicenseId(), req.GetPostalCode(), req.GetPhoneTech(), req.GetPhoneBusiness())

	return &pb.UpdateBusinessInfoResponse{
		Error: "",
	}, nil
}

// UpdateFinancialData handles updating financial data
func (h *CustomerHandler) UpdateFinancialData(ctx context.Context, req *pb.UpdateFinancialDataRequest) (*pb.UpdateFinancialDataResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "CustomerHandler.UpdateFinancialData")
	defer span.End()

	// Validate input
	if errMsg, err := validateField(spannedCtx, req.GetTaxId(), "tax_id"); err != nil {
		span.RecordError(err)
		return &pb.UpdateFinancialDataResponse{Error: errMsg}, nil
	}

	if errMsg, err := validateField(spannedCtx, req.GetIban(), "iban"); err != nil {
		span.RecordError(err)
		return &pb.UpdateFinancialDataResponse{Error: errMsg}, nil
	}

	// TODO: Implement domain logic for updating financial data
	logger.Log.Infof("Updating financial data - Tax ID: %s, IBAN: %s", req.GetTaxId(), req.GetIban())

	return &pb.UpdateFinancialDataResponse{
		Error: "",
	}, nil
}

// ApproveBusiness handles approving a business registration
func (h *CustomerHandler) ApproveBusiness(ctx context.Context, req *pb.ApproveBusinessRequest) (*pb.ApproveBusinessResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "CustomerHandler.ApproveBusiness")
	defer span.End()

	// Validate input
	if req.GetBusinessId() == "" {
		err := errors.NewAppError(spannedCtx).InvalidArgumentError()
		span.RecordError(err)
		logger.Log.Errorf("business_id is required")
		return &pb.ApproveBusinessResponse{
			Error: "business_id is required",
		}, nil
	}

	// TODO: Implement domain logic for approving business
	logger.Log.Infof("Approving business (ID: %s)", req.GetBusinessId())

	return &pb.ApproveBusinessResponse{
		Error: "",
	}, nil
}
