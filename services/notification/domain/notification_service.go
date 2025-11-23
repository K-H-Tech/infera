package domain

import (
	"context"

	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/notification/domain/sms"
)

type NotificationService interface {
	SendOTP(ctx context.Context, otp string, mobile string) error
}

type notificationService struct {
	smsService sms.SmsService
}

func NewNotificationService(smsService sms.SmsService) NotificationService {
	return &notificationService{
		smsService: smsService,
	}
}

func (s *notificationService) SendOTP(ctx context.Context, otp string, mobile string) error {
	spannedCtx, span := trace.GetTracer().Start(ctx, "NotificationService.SendOTP")
	defer span.End()

	err := s.smsService.SendSMS(spannedCtx, otp, mobile)
	if err != nil {
		return err
	}

	return nil
}
