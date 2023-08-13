package config

import (
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
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
	AuthHeaderKey      string
	AuthTypeBearer     string
	AuthPayloadKey     string
}

func LoadEnvVars() (ConfigStruct, error) {
	appEnv, err := utils.GetEnvVar("APP_ENV")
	if err != nil {
		return ConfigStruct{}, err
	}

	appMode, err := utils.GetEnvVar("APP_MODE")
	if err != nil {
		return ConfigStruct{}, err
	}

	httpPort, err := utils.GetEnvVar("HTTP_PORT")
	if err != nil {
		return ConfigStruct{}, err
	}

	grpcPort, err := utils.GetEnvVar("GRPC_PORT")
	if err != nil {
		return ConfigStruct{}, err
	}

	grpcProtoc, err := utils.GetEnvVar("PROTOCOL")
	if err != nil {
		return ConfigStruct{}, err
	}

	mongoUri, err := utils.GetEnvVar("MONGO_URI")
	if err != nil {
		return ConfigStruct{}, err
	}

	dbName, err := utils.GetEnvVar("DATABASE_NAME")
	if err != nil {
		return ConfigStruct{}, err
	}

	custColl, err := utils.GetEnvVar("CUSTOMER_COLLECTION")
	if err != nil {
		return ConfigStruct{}, err
	}

	addColl, err := utils.GetEnvVar("ADDRESS_COLLECTION")
	if err != nil {
		return ConfigStruct{}, err
	}

	cardColl, err := utils.GetEnvVar("CARD_COLLECTION")
	if err != nil {
		return ConfigStruct{}, err
	}

	tokEx, err := utils.GetEnvVar("TOKEN_EXPIRATION")
	if err != nil {
		return ConfigStruct{}, err
	}

	tokSec, err := utils.GetEnvVar("TOKEN_SECRET")
	if err != nil {
		return ConfigStruct{}, err
	}

	loggCol, err := utils.GetEnvVar("lOGGER_COLLECTION")
	if err != nil {
		return ConfigStruct{}, err
	}

	env := ConfigStruct{
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
		AuthHeaderKey:      "Authorization",
		AuthTypeBearer:     "bearer",
		AuthPayloadKey:     "userId",
	}

	return env, nil
}
