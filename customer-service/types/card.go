package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Card struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserId     primitive.ObjectID `bson:"userID"`
	Name       string             `bson:"name"`
	Number     int64              `bson:"number"`
	SecretCode string             `bson:"secretCode"`
	Expiration string             `bson:"expiration"`
	Type       string             `bson:"type"`
	Default    bool               `bson:"default"`
}
