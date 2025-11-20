package sms

import (
	"context"

	"zarinpal-platform/services/notification/config"
	"zarinpal-platform/services/notification/domain/sms/provider"
)

type SmsService interface {
	SendSMS(ctx context.Context, text string, mobile string) error
}

func NewSmsService(config config.KavehNegarSection) SmsService {
	return provider.NewKavehNegarSmsProvider(config)
}
