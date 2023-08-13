package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name" binding:"required"`
	LastName string             `bson:"lastName" binding:"required"`
	Birth    string             `bson:"birth" binding:"required"`
	Email    string             `bson:"email" binding:"required"`
	Password string             `bson:"password" binding:"required"`
	Role     string             `bson:"role" binding:"required"`
}

type SigninUserRequest struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type UserTokenResponse struct {
	Token string `bson:"token"`
}
