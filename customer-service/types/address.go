package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserId     primitive.ObjectID `bson:"userId"`
	Name       string             `bson:"name"`
	Address    string             `bson:"address"`
	PostalCode int64              `bson:"postalCode"`
	Phone      string             `bson:"phone"`
	Default    bool               `bson:"default"`
}
