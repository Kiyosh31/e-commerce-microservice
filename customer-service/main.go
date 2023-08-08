package main

import (
	"log"

	"github.com/Kiyosh31/e-commerce-microservice-common/database"
	"github.com/Kiyosh31/e-commerce-microservice/customer/api"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

func main() {
	var err error

	MongoClient, err = database.ConnectToDB(config.MongoUri)
	if err != nil {
		log.Fatalf("main: %v", err)
	}
	defer database.DisconnectOfDB(MongoClient)

	user_store := store.NewUserStore(MongoClient, config.DatabaseName, config.CustomerCollection)
	service := api.NewService(*user_store, config.ListenPort)

	err = service.Start()
	if err != nil {
		log.Fatalf("Service could not listen: %v", err)
	}
}
