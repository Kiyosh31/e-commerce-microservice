package main

import (
	"log"
	"net"

	"github.com/Kiyosh31/e-commerce-microservice-common/database"
	"github.com/Kiyosh31/e-commerce-microservice/customer/api"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	grpcserver "github.com/Kiyosh31/e-commerce-microservice/customer/grpc_server"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/pb"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
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

	runGrpcServer(*userStore, *addressStore, *cardStore)
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

	log.Println("Starting gRPC server")
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatalf("Cannot start gRPC server: %v", err)
	}
	log.Printf("gRPC server stared at: %v", list.Addr().String())
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
