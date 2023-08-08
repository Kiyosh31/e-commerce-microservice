package main

import (
	"log"

	"github.com/Kiyosh31/e-commerce-microservice-common/database"
	"github.com/Kiyosh31/e-commerce-microservice/customer/api"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
)

func main() {
	var err error

	config.MongoClient, err = database.ConnectToDB(config.MongoUri)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.DisconnectOfDB(config.MongoClient)

	user_store := store.NewUserStore(config.MongoClient, config.DatabaseName, config.CustomerCollection)
	card_store := store.NewCardStore(config.MongoClient, config.DatabaseName, config.CardCollection)
	address_store := store.NewAddressStore(config.MongoClient, config.DatabaseName, config.AddressCollection)

	service := api.NewService(*user_store, *card_store, *address_store, config.ListenPort)

	err = service.Start()
	if err != nil {
		log.Fatalf("user-service could not listen: %v", err)
	}
}
