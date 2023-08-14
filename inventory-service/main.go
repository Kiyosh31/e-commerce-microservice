package main

import (
	"context"
	"net"
	"net/http"

	"github.com/Kiyosh31/e-commerce-microservice-common/database"
	"github.com/Kiyosh31/e-commerce-microservice-common/logger"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/config"
	grpcservice "github.com/Kiyosh31/e-commerce-microservice/inventory/grpc_service"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/proto/pb"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/store"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	env, err := config.LoadEnvVars()
	if err != nil {
		log.Fatal().Msgf("Could not load env vars: %v", err)
	}

	config.MongoClient, err = database.ConnectToDB(env.MongoUri)
	if err != nil {
		log.Fatal().Msgf("Could not connect to database: %v", err)
	}
	defer database.DisconnectOfDB(config.MongoClient)

	productStore := store.NewProductStore(config.MongoClient, env.DatabaseName, env.ProductCollection)

	if env.ServiceMode == "grpc" {
		go runGrpcGateway(*productStore, env)
		runGrpcServer(*productStore, env)
	} else {
		log.Fatal().Msg("To start [inventory-service] You must provide an option in SERVICE_MODE env var")
	}
}

func runGrpcGateway(productStore store.ProductStore, env config.ConfigStruct) {
	service, err := grpcservice.NewService(productStore, env)
	if err != nil {
		log.Fatal().Msgf("Error creating service: %v", err)
	}

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterInventoryServiceHandlerServer(ctx, grpcMux, service)
	if err != nil {
		log.Fatal().Msgf("Cannot register handler service: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	list, err := net.Listen(env.GrpcProtocol, env.HttpPort)
	if err != nil {
		log.Fatal().Msgf("Cannot create listener: %v", err)
	}

	log.Info().Msgf("Starting HTTP gateway service at: %v", list.Addr().String())
	loggerHandler := logger.HttpLogger(mux)
	err = http.Serve(list, loggerHandler)
	if err != nil {
		log.Fatal().Msgf("Cannot start HTTP gateway service: %v", err)
	}
}

func runGrpcServer(productStore store.ProductStore, env config.ConfigStruct) {
	service, err := grpcservice.NewService(productStore, env)
	if err != nil {
		log.Fatal().Msgf("Error creating service: %v", err)
	}

	// Logger
	logger := grpc.UnaryInterceptor(logger.GrpcLoggerInterceptor)

	grpcServer := grpc.NewServer(logger)
	pb.RegisterInventoryServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	list, err := net.Listen(env.GrpcProtocol, env.GrpcPort)
	if err != nil {
		log.Fatal().Msgf("Cannot create listener: %v", err)
	}

	log.Info().Msgf("Starting gRPC server at: %v", list.Addr().String())
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatal().Msgf("Cannot start gRPC server: %v", err)
	}
}
