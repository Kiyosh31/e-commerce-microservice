package config

import (
	"log"

	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigStruct struct {
	AppEnv             string
	AppMode            string
	HttpPort           string
	GrpcPort           string
	MongoUri           string
	DatabaseName       string
	CustomerCollection string
	AddressCollection  string
	CardCollection     string
	TokenSecret        string
	TokenExpiration    string
	Protocol           string
	LoggerCollection   string
}

var (
	EnvVar      *ConfigStruct
	MongoClient *mongo.Client
)

const (
	AuthHeaderKey  = "Authorization"
	AuthTypeBearer = "bearer"
	AuthPayloadKey = "userId"
)

func handleMissingEnv(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LoadEnvVars() {
	appEnv, err := utils.GetEnvVar("APP_ENV")
	handleMissingEnv(err)

	appMode, err := utils.GetEnvVar("APP_MODE")
	handleMissingEnv(err)

	httpPort, err := utils.GetEnvVar("HTTP_PORT")
	handleMissingEnv(err)

	grpcPort, err := utils.GetEnvVar("GRPC_PORT")
	handleMissingEnv(err)

	grpcProtoc, err := utils.GetEnvVar("PROTOCOL")
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

	loggCol, err := utils.GetEnvVar("lOGGER_COLLECTION")
	handleMissingEnv(err)

	EnvVar = &ConfigStruct{
		AppEnv:             appEnv,
		AppMode:            appMode,
		HttpPort:           httpPort,
		GrpcPort:           grpcPort,
		MongoUri:           mongoUri,
		DatabaseName:       dbName,
		CustomerCollection: custColl,
		AddressCollection:  addColl,
		CardCollection:     cardColl,
		TokenSecret:        tokSec,
		TokenExpiration:    tokEx,
		Protocol:           grpcProtoc,
		LoggerCollection:   loggCol,
	}
}
