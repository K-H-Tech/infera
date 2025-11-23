package grpc

import (
	"context"
	"fmt"
	"os"
	"path"

	"zarinpal-platform/core/locale"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"log"
	"net"
	"time"

	"zarinpal-platform/core/logger"
	"zarinpal-platform/core/metric"

	"github.com/emicklei/proto"
)

type Grpc struct {
	serviceName string
	Address     string
	metrics     *metric.Metric
	Server      *grpc.Server
}

func NewGrpc(serviceName string, address string, metric *metric.Metric) *Grpc {
	grpcServer := &Grpc{
		serviceName: serviceName,
		Address:     address,
		metrics:     metric,
	}

	// todo R&D for use input validator as middleware globally (something like protoc-gen-validate)
	// todo another options is use in per service handler for validation

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcServer.metricInterceptor),
		),
		grpc.ChainStreamInterceptor(grpcmiddleware.ChainStreamServer()),
	}
	grpcServer.Server = grpc.NewServer(opts...)
	return grpcServer
}

func (g *Grpc) Start() {
	reflection.Register(g.Server)
	lis, err := net.Listen("tcp", g.Address)
	if err != nil {
		log.Fatalf("Grpc failed to listen: %v", err)
		return
	}
	log.Printf("Grpc listening on %s", g.Address)
	go func() {
		if err := g.Server.Serve(lis); err != nil {
			logger.Log.Fatalf("grpc failed to serve: %v", err)
		}
	}()

	g.addZeroValueMetrics()

}

func (g *Grpc) Stop() {
	g.Server.GracefulStop()
}

func (g *Grpc) metricInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		acceptLang := md.Get("Accept-Language")
		if len(acceptLang) > 0 {
			ctx = locale.WithLocale(ctx, locale.FromAcceptLang(acceptLang[0]))
		}
	}

	res, err := handler(ctx, req)

	st, _ := status.FromError(err)
	code := st.Code()

	duration := time.Since(startTime).Seconds()

	method := "/" + path.Base(info.FullMethod)

	g.metrics.MethodTotal.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
		"method": method}).Inc()

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	g.metrics.MethodDuration.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
		"method": method}).Observe(duration)

	if code == codes.OK {
		g.metrics.MethodSuccessTotal.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
			"method": method}).Inc()
	} else {
		g.metrics.MethodErrorDuration.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
			"method": method, "error": errMsg}).Observe(duration)
		if g.isUserError(err) {
			g.metrics.MethodUserErrorTotal.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
				"method": method, "error": errMsg}).Inc()
		} else {
			g.metrics.MethodServerErrorTotal.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
				"method": method, "error": errMsg}).Inc()
		}
	}

	return res, err
}

func (g *Grpc) isUserError(err error) bool {
	if err == nil {
		return false
	}

	code := status.Code(err)
	switch code {
	case codes.InvalidArgument,
		codes.FailedPrecondition,
		codes.OutOfRange,
		codes.Unauthenticated,
		codes.PermissionDenied,
		codes.AlreadyExists,
		codes.NotFound,
		codes.Aborted,
		codes.Canceled:
		return true
	default:
		return false
	}
}

func (g *Grpc) addZeroValueMetrics() {
	reader, err := os.Open("./" + g.serviceName + ".proto")
	if err != nil {
		logger.Log.Warn(fmt.Sprintf("failed to open proto file: %v", err))
		return
	}
	defer func(reader *os.File) {
		_ = reader.Close()
	}(reader)

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		logger.Log.Warn(fmt.Sprintf("failed to parse proto: %v", err))
		return
	}

	var methods []string

	// Walk through all top-level elements
	proto.Walk(
		definition,
		proto.WithService(func(s *proto.Service) {
			for _, element := range s.Elements {
				if rpc, ok := element.(*proto.RPC); ok {
					methods = append(methods, "/"+rpc.Name)
				}
			}
		}),
	)

	for _, method := range methods {
		g.metrics.MethodTotal.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
			"method": method}).Add(0)
		g.metrics.MethodSuccessTotal.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
			"method": method}).Add(0)
		g.metrics.MethodUserErrorTotal.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
			"method": method, "error": ""}).Add(0)
		g.metrics.MethodServerErrorTotal.With(prometheus.Labels{"service_name": g.serviceName, "type": "grpc",
			"method": method, "error": ""}).Add(0)
	}
}
