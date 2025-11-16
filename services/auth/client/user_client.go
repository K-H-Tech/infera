package client

import (
	"context"
	"zarinpal-platform/core/grpc"
	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/trace"
	user "zarinpal-platform/services/user/api/grpc/pb/src/golang"
)

type UserServiceClient interface {
	IsShahkarValid(ctx context.Context, mobile string, nationalCode string) (bool, error)
}
type userServiceClient struct {
	client user.UserServiceClient
}

func NewUserServiceClient(address string) UserServiceClient {
	c := &userServiceClient{}
	cc := grpc.NewClient("user", address)
	c.client = user.NewUserServiceClient(cc)
	return c
}
func (c *userServiceClient) IsShahkarValid(ctx context.Context, mobile string, nationalCode string) (bool, error) {
	_, span := trace.GetTracer().Start(ctx, "UserServiceClient.IsShahkarValid")
	defer span.End()

	response, err := c.client.IsShahkarValid(ctx, &user.IsShahkarValidRequest{Mobile: mobile, NationalCode: nationalCode})
	if err != nil {
		span.RecordError(err)
		span.SetString("mobile", mobile)
		logger.Log.Errorf("check shahkar validation error: %v", err)
		return false, err
	}
	return response.IsValid, nil
}
