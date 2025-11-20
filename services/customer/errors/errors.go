package errors

import (
	"context"

	"zarinpal-platform/core/locale"

	"github.com/leonelquinteros/gotext"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError interface {
	InvalidError() error
	DatabaseError() error
	InvalidArgumentError() error
}

type appError struct {
	locale *gotext.Locale
}

func NewAppError(ctx context.Context) AppError {
	return appError{
		locale: locale.FromContext(ctx),
	}
}

func (e appError) InvalidError() error {
	return status.Error(codes.PermissionDenied, "InvalidError")
}

func (e appError) DatabaseError() error {
	return status.Error(codes.Internal, e.locale.Get("Database Error"))
}

func (e appError) InvalidArgumentError() error {
	return status.Error(codes.InvalidArgument, e.locale.Get("Your argument are invalid"))
}
