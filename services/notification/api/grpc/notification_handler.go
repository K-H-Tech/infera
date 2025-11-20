package grpc

import (
	"context"

	"zarinpal-platform/core/trace"
	notification "zarinpal-platform/services/notification/api/grpc/pb/src/golang"

	"google.golang.org/protobuf/types/known/emptypb"

	"zarinpal-platform/services/notification/domain"
)

type NotificationHandler struct {
	notification.UnimplementedNotificationServiceServer
	notificationService domain.NotificationService
}

func NewNotificationHandler(notificationService domain.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) SendOTP(ctx context.Context, req *notification.SendOTPRequest) (*emptypb.Empty, error) {
	spannedCtx, span := trace.GetTracer().Start(ctx, "NotificationHandler.SendOTP")
	defer span.End()

	err := h.notificationService.SendOTP(spannedCtx, req.Otp, req.Mobile)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
