package main

import (
	"context"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/Kiyosh31/e-commerce-microservice-common/database"
	"github.com/Kiyosh31/e-commerce-microservice-common/logger"
	"github.com/Kiyosh31/e-commerce-microservice/customer/api"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	grpcserver "github.com/Kiyosh31/e-commerce-microservice/customer/grpc_server"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/pb"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config.LoadEnvVars()
	var err error

	config.MongoClient, err = database.ConnectToDB(config.EnvVar.MongoUri)
	if err != nil {
		log.Fatal().Msgf("Could not connect to database: %v", err)
	}
	defer database.DisconnectOfDB(config.MongoClient)

	userStore := store.NewUserStore(config.MongoClient, config.EnvVar.DatabaseName, config.EnvVar.CustomerCollection)
	cardStore := store.NewCardStore(config.MongoClient, config.EnvVar.DatabaseName, config.EnvVar.CardCollection)
	addressStore := store.NewAddressStore(config.MongoClient, config.EnvVar.DatabaseName, config.EnvVar.AddressCollection)

	if config.EnvVar.AppMode == "grpc" {
		go runGatewayServer(*userStore, *addressStore, *cardStore)
		runGrpcServer(*userStore, *addressStore, *cardStore)
	} else {
		runGinService(*userStore, *addressStore, *cardStore)
	}
}

func runGatewayServer(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore) {
	service, err := grpcserver.NewService(userStore, addressStore, cardStore)
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

	list, err := net.Listen(config.EnvVar.Protocol, config.EnvVar.HttpPort)
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

func runGrpcServer(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore) {
	service, err := grpcserver.NewService(userStore, addressStore, cardStore)
	if err != nil {
		log.Fatal().Msgf("Error creating service: %v", err)
	}

	// Logger
	logger := grpc.UnaryInterceptor(logger.GrpcLogger)

	// Server
	grpcServer := grpc.NewServer(logger)
	pb.RegisterCustomerServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	list, err := net.Listen(config.EnvVar.Protocol, config.EnvVar.GrpcPort)
	if err != nil {
		log.Fatal().Msgf("Cannot create listener: %v", err)
	}

	log.Info().Msgf("Starting gRPC server at: %v", list.Addr().String())
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatal().Msgf("Cannot start gRPC server")
	}
}

func runGinService(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore) {

	service, err := api.NewService(userStore, cardStore, addressStore, config.EnvVar.HttpPort)
	if err != nil {
		log.Fatal().Msgf("Error creating service: %v", err)
	}

	err = service.Start()
	if err != nil {
		log.Fatal().Msgf("user-service could not listen: %v", err)
	}
}
