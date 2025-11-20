package client

import (
	"context"

	"zarinpal-platform/core/grpc"
	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/trace"
	notification "zarinpal-platform/services/notification/api/grpc/pb/src/golang"
)

type NotificationServiceClient interface {
	SendOTP(ctx context.Context, mobile string, otp string) error
}
type notificationServiceClient struct {
	client notification.NotificationServiceClient
}

func NewNotificationServiceClient(address string) NotificationServiceClient {
	c := &notificationServiceClient{}
	cc := grpc.NewClient("notification", address)
	c.client = notification.NewNotificationServiceClient(cc)
	return c
}
func (c *notificationServiceClient) SendOTP(ctx context.Context, mobile string, otp string) error {
	_, span := trace.GetTracer().Start(ctx, "NotificationServiceClient.SendOTP")
	defer span.End()

	_, err := c.client.SendOTP(ctx, &notification.SendOTPRequest{Mobile: mobile, Otp: otp})
	if err != nil {
		span.RecordError(err)
		span.SetString("mobile", mobile)
		logger.Log.Errorf("Send OTP error: %v", err)
		return err
	}
	return nil
}
