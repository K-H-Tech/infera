package grpc

import (
	"zarinpal-platform/core/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(serviceName string, address string) *grpc.ClientConn {
	cc, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Fatalf("Faild to connect to %v: %v", serviceName, err.Error())
	}

	logger.Log.Infof("Successfuly GRPC connected to %v service: %v", serviceName, address)

	return cc
}
