package provider

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/trace"
	"zarinpal-platform/pkg/common"
	"zarinpal-platform/services/notification/config"
)

type KavehNegarSmsProvider struct {
	url      string
	apiKey   string
	template string
}

func NewKavehNegarSmsProvider(config config.KavehNegarSection) *KavehNegarSmsProvider {
	return &KavehNegarSmsProvider{
		url:      config.Url,
		apiKey:   config.ApiKey,
		template: config.Template,
	}
}

func (p *KavehNegarSmsProvider) SendSMS(ctx context.Context, message string, mobile string) error {
	_, span := trace.GetTracer().Start(ctx, "KavehNegarSmsProvider.SendSMS")
	defer span.End()

	path := fmt.Sprintf(p.url+"?receptor=%v&token=123456&template=%v", p.apiKey, mobile, p.template)
	// todo fix it
	path = strings.ReplaceAll(path, "\n", "")
	res, err := common.Get(path)
	if err != nil {
		span.RecordError(err)
		span.SetString("mobile", mobile)
		logger.Log.Errorf(err.Error())
		return err
	}
	if res.StatusCode != 200 {
		// todo handle errors
	}
	logger.Log.Infof("successfuly send otp to %v", mobile)

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	if rand.Float64() < 0.5 {
		err := errors.New("error in connection to kaveh negar sms")
		span.RecordError(err)
		span.SetString("mobile", mobile)
		logger.Log.Errorf("error in connection to kaveh negar sms: %s", err.Error())
		return err
	}
	time.Sleep(1 * time.Second)
	logger.Log.Infof("sending sms message to %s on %s", mobile, p.url)
	return nil
}
