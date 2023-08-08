package config

import "github.com/Kiyosh31/e-commerce-microservice-common/utils"

var (
	ListenPort, _         = utils.GetEnvVar("PORT")
	MongoUri, _           = utils.GetEnvVar("MONGO_URI")
	DatabaseName, _       = utils.GetEnvVar("DATABASE_NAME")
	CustomerCollection, _ = utils.GetEnvVar("CUSTOMER_COLLECTION")
	AddressCollection, _  = utils.GetEnvVar("ADDRESS_COLLECTION")
	CardCollection, _     = utils.GetEnvVar("CARD_COLLECTION")
)
