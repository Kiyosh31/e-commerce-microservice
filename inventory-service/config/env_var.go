package config

import "github.com/Kiyosh31/e-commerce-microservice-common/utils"

type ConfigStruct struct {
	AppEnv            string
	ServiceMode       string
	HttpPort          string
	GrpcPort          string
	GrpcProtocol      string
	MongoUri          string
	DatabaseName      string
	ProductCollection string
	LoggerCollection  string
	TokenExpiration   string
}

func LoadEnvVars() (ConfigStruct, error) {
	appEnv, err := utils.GetEnvVar("APP_ENV")
	if err != nil {
		return ConfigStruct{}, err
	}

	serviceMode, err := utils.GetEnvVar("SERVICE_MODE")
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

	grpcProtoc, err := utils.GetEnvVar("GRPC_PROTOCOL")
	if err != nil {
		return ConfigStruct{}, err
	}

	tokenExp, err := utils.GetEnvVar("TOKEN_EXPIRATION")
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

	productCol, err := utils.GetEnvVar("PRODUCT_COLLECTION")
	if err != nil {
		return ConfigStruct{}, err
	}

	loggerCol, err := utils.GetEnvVar("LOGGER_COLLECTION")
	if err != nil {
		return ConfigStruct{}, err
	}

	env := ConfigStruct{
		AppEnv:            appEnv,
		ServiceMode:       serviceMode,
		HttpPort:          httpPort,
		GrpcPort:          grpcPort,
		GrpcProtocol:      grpcProtoc,
		MongoUri:          mongoUri,
		DatabaseName:      dbName,
		ProductCollection: productCol,
		LoggerCollection:  loggerCol,
		TokenExpiration:   tokenExp,
	}

	return env, nil
}
