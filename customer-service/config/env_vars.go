package config

import (
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ListenPort, _         = utils.GetEnvVar("PORT")
	MongoUri, _           = utils.GetEnvVar("MONGO_URI")
	DatabaseName, _       = utils.GetEnvVar("DATABASE_NAME")
	CustomerCollection, _ = utils.GetEnvVar("CUSTOMER_COLLECTION")
	AddressCollection, _  = utils.GetEnvVar("ADDRESS_COLLECTION")
	CardCollection, _     = utils.GetEnvVar("CARD_COLLECTION")
)

var (
	MongoClient *mongo.Client
	log         = utils.NewLogger()
)
