package config

import (
	"log"

	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigStruct struct {
	ListenPort         string
	MongoUri           string
	DatabaseName       string
	CustomerCollection string
	AddressCollection  string
	CardCollection     string
	TokenSecret        string
	TokenExpiration    string
}

var (
	EnvVar      *ConfigStruct
	MongoClient *mongo.Client
)

func handleMissingEnv(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LoadEnvVars() {
	port, err := utils.GetEnvVar("PORT")
	handleMissingEnv(err)

	mongoUri, err := utils.GetEnvVar("MONGO_URI")
	handleMissingEnv(err)

	dbName, err := utils.GetEnvVar("DATABASE_NAME")
	handleMissingEnv(err)

	custColl, err := utils.GetEnvVar("CUSTOMER_COLLECTION")
	handleMissingEnv(err)

	addColl, err := utils.GetEnvVar("ADDRESS_COLLECTION")
	handleMissingEnv(err)

	cardColl, err := utils.GetEnvVar("CARD_COLLECTION")
	handleMissingEnv(err)

	tokEx, err := utils.GetEnvVar("TOKEN_EXPIRATION")
	handleMissingEnv(err)

	tokSec, err := utils.GetEnvVar("TOKEN_SECRET")
	handleMissingEnv(err)

	EnvVar = &ConfigStruct{
		ListenPort:         port,
		MongoUri:           mongoUri,
		DatabaseName:       dbName,
		CustomerCollection: custColl,
		AddressCollection:  addColl,
		CardCollection:     cardColl,
		TokenSecret:        tokSec,
		TokenExpiration:    tokEx,
	}
}
