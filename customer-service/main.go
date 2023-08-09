package main

import (
	"log"

	"github.com/Kiyosh31/e-commerce-microservice-common/database"
	"github.com/Kiyosh31/e-commerce-microservice/customer/api"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
)

func main() {
	config.LoadEnvVars()
	var err error

	config.MongoClient, err = database.ConnectToDB(config.EnvVar.MongoUri)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.DisconnectOfDB(config.MongoClient)

	user_store := store.NewUserStore(config.MongoClient, config.EnvVar.DatabaseName, config.EnvVar.CustomerCollection)
	card_store := store.NewCardStore(config.MongoClient, config.EnvVar.DatabaseName, config.EnvVar.CardCollection)
	address_store := store.NewAddressStore(config.MongoClient, config.EnvVar.DatabaseName, config.EnvVar.AddressCollection)

	service, err := api.NewService(*user_store, *card_store, *address_store, config.EnvVar.ListenPort)
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	err = service.Start()
	if err != nil {
		log.Fatalf("user-service could not listen: %v", err)
	}
}
