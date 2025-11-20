package core

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"zarinpal-platform/core/configuration"
	"zarinpal-platform/core/grpc"
	"zarinpal-platform/core/http"
	"zarinpal-platform/core/locale"
	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/metric"
	"zarinpal-platform/core/prometheus"
	"zarinpal-platform/core/trace"
)

type Service struct {
	Name   string
	Http   *http.Http
	Grpc   *grpc.Grpc
	tracer *trace.Tracer
	metric *metric.Metric
}

type IService interface {
	OnStart(service *Service)
	OnStop()
}

func StartService(name string, iService IService) {
	service := initializeServices(name)
	service.start(iService)
	service.waitForOsSignal()
	service.stop(iService)
}

func initializeServices(name string) *Service {
	service := &Service{
		Name: name,
	}
	// Locale for translate
	locale.Init()

	// Config
	c := configuration.LoadConfig(name)

	// Logger
	logger.InitLogger()

	// metrics
	service.metric = metric.NewMetric()

	// Prometheus
	prometheus.RunPrometheus(c.Prometheus.Address)

	// Http Server
	httpServer := http.NewHttp(name, c.Http.Address, service.metric)
	service.Http = httpServer
	httpServer.Start()

	// Jaeger
	service.tracer = trace.NewTracer(name, c.Jaeger.Address)

	// Grpc
	service.Grpc = grpc.NewGrpc(name, c.Grpc.Address, service.metric)

	return service

}

func (s *Service) start(service IService) {
	log.Printf("Starting %v service...\n", s.Name)
	service.OnStart(s)

}

func (s *Service) waitForOsSignal() {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	<-osSignal
}

func (s *Service) stop(service IService) {
	log.Printf("Stopping %v service\n", s.Name)
	service.OnStop()
}
