package grpc

import (
	"context"

	"zarinpal-platform/core/trace"
	auth "zarinpal-platform/services/auth/api/grpc/pb/src/golang"
	"zarinpal-platform/services/auth/domain"
	"zarinpal-platform/services/auth/errors"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	loginWithOtpService domain.LoginWithOTPService
}

func NewAuthHandler(loginWithOtpService domain.LoginWithOTPService) *AuthHandler {
	return &AuthHandler{
		loginWithOtpService: loginWithOtpService,
	}
}

func (h *AuthHandler) AuthenticateWithOTP(ctx context.Context, req *auth.AuthenticateWithOTPRequest) (*auth.AuthenticateWithOTPResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "AuthHandler.AuthenticateWithOTP")
	defer span.End()

	if req.Mobile == "" {
		return nil, errors.NewAppError().MobileIsInvalidError()
	}

	return h.loginWithOtpService.LoginByMobile(spannedCtx, req.Mobile)
}

func (h *AuthHandler) VerifyOTP(ctx context.Context, req *auth.VerifyOTPRequest) (*auth.VerifyOTPResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "AuthHandler.VerifyOTP")
	defer span.End()

	span.SetString("mobile", req.Mobile)

	return h.loginWithOtpService.VerifyMobileAndOtp(spannedCtx, req.Mobile, req.Code)
}

func (h *AuthHandler) NewUser(ctx context.Context, req *auth.NewUserRequest) (*auth.NewUserResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "AuthHandler.NewUser")
	defer span.End()

	span.SetString("name", req.Name)

	return h.loginWithOtpService.NewUser(spannedCtx, req.Name)
}

func (h *AuthHandler) GetUser(ctx context.Context, req *auth.GetUserRequest) (*auth.GetUserResponse, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "AuthHandler.GetUser")
	defer span.End()

	span.SetInt("name", int(req.Id))

	return h.loginWithOtpService.GetUser(spannedCtx, req.Id)

}
