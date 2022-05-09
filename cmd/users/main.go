package main

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/marcoshuck/jaeger-demo/api"
	"github.com/marcoshuck/jaeger-demo/application/service"
	"github.com/marcoshuck/jaeger-demo/collector"
	grpc_otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	environment := "production"
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("Failed to initialize Zap logger:", err)
	}

	// If development mode is enabled, use development logger.
	if os.Getenv("USERS_DEVELOPMENT_MODE") == "true" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalln("Failed to initialize Zap logger:", err)
		}

		environment = "staging"
	}

	traceProvider, err := collector.NewJaegerTracerProvider("users_v1", "http://localhost:14268/api/traces", environment)
	if err != nil {
		logger.Fatal("Failed to initialize Jaeger tracer provider", zap.Error(err))
	}

	// Initialize gRPC server with respective interceptors.
	server := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			// grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			grpc_otel.StreamServerInterceptor(grpc_otel.WithPropagators(jaeger.Jaeger{}), grpc_otel.WithTracerProvider(traceProvider)),
			grpc_recovery.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
		),
	),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				// grpc_prometheus.UnaryServerInterceptor,
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_otel.UnaryServerInterceptor(grpc_otel.WithPropagators(jaeger.Jaeger{}), grpc_otel.WithTracerProvider(traceProvider)),
				grpc_recovery.UnaryServerInterceptor(),
				grpc_validator.UnaryServerInterceptor(),
			),
		),
	)

	// Initialize tweets service
	svc := service.NewUsersV1()

	// Register service in server
	api.RegisterUsersServiceServer(server, svc)

	// Define listener
	listener, err := net.Listen("tcp", ":3030")
	if err != nil {
		logger.Fatal("Failed to initialize Tweets V1 listener on port 3030", zap.Error(err))
	}

	// Serve gRPC service
	logger.Info("Listening incoming gRPC requests on port 3030", zap.Int("port", 3030))
	if err := server.Serve(listener); err != nil {
		logger.Fatal("Failed to listen on port 3030", zap.Error(err), zap.Int("port", 3030))
	}
}
