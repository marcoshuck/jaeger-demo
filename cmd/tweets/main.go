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
	"google.golang.org/grpc/credentials/insecure"

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
	if os.Getenv("TWEETS_DEVELOPMENT_MODE") == "true" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalln("Failed to initialize Zap logger:", err)
		}
		environment = "staging"
	}

	var propagator jaeger.Jaeger

	traceProvider, err := collector.NewJaegerTracerProvider("tweets_v1", "http://localhost:14268/api/traces", environment)
	if err != nil {
		logger.Fatal("Failed to initialize Jaeger tracer provider", zap.Error(err))
	}

	// Connect to users service
	url := os.Getenv("TWEETS_USERS_URL")
	logger.Info("Connecting to Users service", zap.String("service_url", url))
	conn, err := grpc.Dial(url,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_otel.UnaryClientInterceptor(grpc_otel.WithPropagators(propagator), grpc_otel.WithTracerProvider(traceProvider))),
		grpc.WithStreamInterceptor(grpc_otel.StreamClientInterceptor(grpc_otel.WithPropagators(propagator), grpc_otel.WithTracerProvider(traceProvider))),
	)
	if err != nil {
		logger.Fatal("Failed to initialize connection with Users V1 service", zap.Error(err))
	}
	client := api.NewUsersServiceClient(conn)

	server := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			// grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			grpc_otel.StreamServerInterceptor(grpc_otel.WithPropagators(propagator), grpc_otel.WithTracerProvider(traceProvider)),
			grpc_recovery.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
		),
	),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				// grpc_prometheus.UnaryServerInterceptor,
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_otel.UnaryServerInterceptor(grpc_otel.WithPropagators(propagator), grpc_otel.WithTracerProvider(traceProvider)),
				grpc_recovery.UnaryServerInterceptor(),
				grpc_validator.UnaryServerInterceptor(),
			),
		),
	)

	// Initialize tweets service
	svc := service.NewTweetsV1(client)

	// Register service in server
	api.RegisterTweetsServiceServer(server, svc)

	// Define listener
	listener, err := net.Listen("tcp", ":3031")
	if err != nil {
		logger.Fatal("Failed to initialize Tweets V1 listener", zap.Int("port", 3031), zap.Error(err))
	}

	// Serve gRPC service
	logger.Info("Listening gRPC requests", zap.Int("port", 3031))
	if err := server.Serve(listener); err != nil {
		logger.Fatal("Failed to listen incoming gRPC requests", zap.Int("port", 3031), zap.Error(err))
	}
}
