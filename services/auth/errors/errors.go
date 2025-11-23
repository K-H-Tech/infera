package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError interface {
	MobileIsInvalidError() error
	DatabaseError() error
	SendOTPError() error
}
type appError struct {
}

func NewAppError() AppError {
	return appError{}
}

func (e appError) MobileIsInvalidError() error {
	return status.Error(codes.InvalidArgument, "شماره ی وارد شده صحیح نیست")
}

func (e appError) DatabaseError() error {
	return status.Error(codes.Internal, "DatabaseError")
}

func (e appError) SendOTPError() error {
	return status.Error(codes.Internal, "failed to send otp error")
}
