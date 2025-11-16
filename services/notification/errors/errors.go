package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError interface {
	MobileIsInvalidError() error
	DatabaseError() error
}
type appError struct {
}

func NewAppError() AppError {
	return appError{}
}

func (e appError) MobileIsInvalidError() error {
	return status.Error(codes.PermissionDenied, "MobileIsInvalidError")
}

func (e appError) DatabaseError() error {
	return status.Error(codes.Internal, "DatabaseError")
}
