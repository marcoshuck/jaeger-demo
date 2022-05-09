package main

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/marcoshuck/jaeger-demo/api"
	"github.com/marcoshuck/jaeger-demo/application/service"
	grpc_otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"log"
	"net"
	"os"
)

func main() {
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
	}

	// Connect to users service
	conn, err := grpc.Dial(os.Getenv("TWEETS_USERS_URL"), grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to initialize connection with Users V1 service", zap.Error(err))
	}
	client := api.NewUsersServiceClient(conn)

	// Initialize tweets service
	svc := service.NewTweetsV1(client, otel.Tracer("tweets_v1"))

	// Initialize gRPC server with respective interceptors.
	server := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			// grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			grpc_otel.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
		),
	),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_opentracing.UnaryServerInterceptor(),
				// grpc_prometheus.UnaryServerInterceptor,
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_otel.UnaryServerInterceptor(),
				grpc_recovery.UnaryServerInterceptor(),
				grpc_validator.UnaryServerInterceptor(),
			),
		),
	)

	// Register service in server
	api.RegisterTweetsServiceServer(server, svc)

	// Define listener
	listener, err := net.Listen("tcp", ":3031")
	if err != nil {
		logger.Fatal("Failed to initialize Tweets V1 listener on port 3031", zap.Error(err))
	}

	// Serve gRPC service
	log.Println("Listening incoming gRPC requests on port 3031")
	if err := server.Serve(listener); err != nil {
		logger.Fatal("Failed to listen on port 3031", zap.Error(err))
	}
}
