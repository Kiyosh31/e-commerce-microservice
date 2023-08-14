package main

import (
	"context"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/Kiyosh31/e-commerce-microservice-common/database"
	"github.com/Kiyosh31/e-commerce-microservice-common/logger"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	grpcservice "github.com/Kiyosh31/e-commerce-microservice/customer/grpc_service"
	httpservice "github.com/Kiyosh31/e-commerce-microservice/customer/http_service"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/pb"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

	userStore := store.NewUserStore(config.MongoClient, env.DatabaseName, env.CustomerCollection)
	cardStore := store.NewCardStore(config.MongoClient, env.DatabaseName, env.CardCollection)
	addressStore := store.NewAddressStore(config.MongoClient, env.DatabaseName, env.AddressCollection)

	if env.ServiceMode == "grpc" {
		go runGatewayServer(*userStore, *addressStore, *cardStore, env)
		runGrpcServer(*userStore, *addressStore, *cardStore, env)
	} else if env.ServiceMode == "http" {
		runGinService(*userStore, *addressStore, *cardStore, env)
	} else {
		log.Fatal().Msg("To start [customer-service] You must provide an option in SERVICE_MODE env var")
	}
}

func runGatewayServer(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore, env config.ConfigStruct) {
	service, err := grpcservice.NewService(userStore, addressStore, cardStore, env)
	if err != nil {
		log.Fatal().Msgf("Error creating service: %v", err)
	}

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterCustomerServiceHandlerServer(ctx, grpcMux, service)
	if err != nil {
		log.Fatal().Msgf("Cannot register handler service: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	list, err := net.Listen(env.Protocol, env.HttpPort)
	if err != nil {
		log.Fatal().Msgf("Cannot create listener: %v", err)
	}

	log.Info().Msgf("Starting HTTP gateway service at: %v", list.Addr().String())
	// Logger
	loggerHandler := logger.HttpLogger(mux)
	err = http.Serve(list, loggerHandler)
	if err != nil {
		log.Fatal().Msgf("Cannot start HTTP gateway service: %v", err)
	}
}

func runGrpcServer(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore, env config.ConfigStruct) {
	service, err := grpcservice.NewService(userStore, addressStore, cardStore, env)
	if err != nil {
		log.Fatal().Msgf("Error creating service: %v", err)
	}

	// Logger
	logger := grpc.UnaryInterceptor(logger.GrpcLoggerInterceptor)

	// Server
	grpcServer := grpc.NewServer(logger)
	pb.RegisterCustomerServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	list, err := net.Listen(env.Protocol, env.GrpcPort)
	if err != nil {
		log.Fatal().Msgf("Cannot create listener: %v", err)
	}

	log.Info().Msgf("Starting gRPC server at: %v", list.Addr().String())
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatal().Msgf("Cannot start gRPC server: %v", err)
	}
}

func runGinService(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore, env config.ConfigStruct) {

	service, err := httpservice.NewService(userStore, cardStore, addressStore, env.HttpPort, env)
	if err != nil {
		log.Fatal().Msgf("Error creating service: %v", err)
	}

	err = service.Start()
	if err != nil {
		log.Fatal().Msgf("user-service could not listen: %v", err)
	}
}
