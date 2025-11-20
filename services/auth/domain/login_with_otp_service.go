package domain

import (
	"context"
	"time"

	"zarinpal-platform/core/trace"
	auth "zarinpal-platform/services/auth/api/grpc/pb/src/golang"
	"zarinpal-platform/services/auth/client"
	"zarinpal-platform/services/auth/data/repository"
	"zarinpal-platform/services/auth/errors"
)

type LoginWithOTPService interface {
	LoginByMobile(ctx context.Context, mobile string) (*auth.AuthenticateWithOTPResponse, error)
	VerifyMobileAndOtp(ctx context.Context, mobile string, code string) (*auth.VerifyOTPResponse, error)
	NewUser(ctx context.Context, name string) (*auth.NewUserResponse, error)
	GetUser(ctx context.Context, id int32) (*auth.GetUserResponse, error)
}

type loginWithOTPService struct {
	notificationServiceClient client.NotificationServiceClient
	userServiceClient         client.UserServiceClient
	otpRepository             repository.OtpRepository
}

func NewLoginWithOTPService(notificationServiceClient client.NotificationServiceClient,
	userServiceClient client.UserServiceClient, otpRepository repository.OtpRepository) LoginWithOTPService {
	return &loginWithOTPService{
		notificationServiceClient: notificationServiceClient,
		userServiceClient:         userServiceClient,
		otpRepository:             otpRepository,
	}
}

func (s *loginWithOTPService) LoginByMobile(ctx context.Context, mobile string) (*auth.AuthenticateWithOTPResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "LoginWithOTPService.LoginByMobile")
	defer span.End()

	otp := "1234"
	err := s.otpRepository.GenerateOTP(spannedCtx, mobile)
	if err != nil {
		return nil, err
	}

	err = s.notificationServiceClient.SendOTP(spannedCtx, mobile, otp)
	if err != nil {
		return nil, errors.NewAppError().SendOTPError()
	}

	return &auth.AuthenticateWithOTPResponse{
		ExpireAt: "1234",
	}, nil
}

func (s *loginWithOTPService) VerifyMobileAndOtp(ctx context.Context, mobile string, code string) (*auth.VerifyOTPResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "LoginWithOTPService.VerifyMobileAndOtp")
	defer span.End()

	err := s.otpRepository.VerifyOTP(spannedCtx, mobile, code)
	if err != nil {
		return nil, err
	}

	return &auth.VerifyOTPResponse{
		AccessToken: "1234",
	}, nil
}

func (s *loginWithOTPService) NewUser(ctx context.Context, name string) (*auth.NewUserResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "LoginWithOTPService.NewUser")
	defer span.End()

	user, err := s.otpRepository.NewUser(spannedCtx, name)
	if err != nil {
		return nil, err
	}

	return &auth.NewUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *loginWithOTPService) GetUser(ctx context.Context, id int32) (*auth.GetUserResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "LoginWithOTPService.GetUser")
	defer span.End()

	user, err := s.otpRepository.GetUser(spannedCtx, id)
	if err != nil {
		return nil, err
	}

	return &auth.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}
