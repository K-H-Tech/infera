package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError interface {
	InvalidError() error
	DatabaseError() error
}
type appError struct {
}

func NewAppError() AppError {
	return appError{}
}

func (e appError) InvalidError() error {
	return status.Error(codes.PermissionDenied, "InvalidError")
}

func (e appError) DatabaseError() error {
	return status.Error(codes.Internal, "DatabaseError")
}