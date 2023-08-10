package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/Kiyosh31/e-commerce-microservice-common/database"
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
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.DisconnectOfDB(config.MongoClient)

	userStore := store.NewUserStore(config.MongoClient, config.EnvVar.DatabaseName, config.EnvVar.CustomerCollection)
	cardStore := store.NewCardStore(config.MongoClient, config.EnvVar.DatabaseName, config.EnvVar.CardCollection)
	addressStore := store.NewAddressStore(config.MongoClient, config.EnvVar.DatabaseName, config.EnvVar.AddressCollection)

	go runGatewayServer(*userStore, *addressStore, *cardStore)
	runGrpcServer(*userStore, *addressStore, *cardStore)
}

func runGatewayServer(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore) {
	service, err := grpcserver.NewService(userStore, addressStore, cardStore)
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterCustomerServiceHandlerServer(ctx, grpcMux, service)
	if err != nil {
		log.Fatalf("Cannot register handler service: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	list, err := net.Listen(config.EnvVar.Protocol, config.EnvVar.HttpPort)
	if err != nil {
		log.Fatalf("Cannot create listener: %v", err)
	}

	log.Printf("Starting HTTP gateway service at: %v\n", list.Addr().String())
	err = http.Serve(list, mux)
	if err != nil {
		log.Fatalf("Cannot start HTTP gateway service: %v", err)
	}
}

func runGrpcServer(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore) {
	service, err := grpcserver.NewService(userStore, addressStore, cardStore)
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCustomerServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	list, err := net.Listen(config.EnvVar.Protocol, config.EnvVar.GrpcPort)
	if err != nil {
		log.Fatalf("Cannot create listener: %v", err)
	}

	log.Printf("Starting gRPC server at: %v", list.Addr().String())
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatalf("Cannot start gRPC server: %v", err)
	}
}

func runGinService(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore) {

	service, err := api.NewService(userStore, cardStore, addressStore, config.EnvVar.GrpcPort)
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	err = service.Start()
	if err != nil {
		log.Fatalf("user-service could not listen: %v", err)
	}
}
