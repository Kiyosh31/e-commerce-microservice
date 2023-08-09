package grpcserver

import (
	"context"

	"github.com/Kiyosh31/e-commerce-microservice-common/token"
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/pb"
	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
)

func (svc *Service) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := utils.HashPassword(in.User.GetPassword())
	if err != nil {
		return nil, errorResponse("Failed to hash password", err)
	}
	mongoId, err := utils.GetMongoId(in.GetUser().GetId())
	if err != nil {
		return nil, errorResponse("Failed to parse mongoId", err)
	}

	newUser := types.User{
		ID:       mongoId,
		Name:     in.GetUser().GetName(),
		LastName: in.GetUser().GetLasName(),
		Birth:    in.GetUser().GetBirth(),
		Email:    in.GetUser().GetEmail(),
		Password: hashedPassword,
	}

	createdUser, err := svc.userStore.CreateUser(ctx, newUser)
	if err != nil {
		return nil, errorResponse("Could not create user in database", err)
	}

	res := &pb.CreateUserResponse{
		InsertedID: createdUser.InsertedID.(string),
	}

	return res, nil
}

func (svc *Service) SigninUser(ctx context.Context, in *pb.SigninUserRequest) (*pb.SigninUserResponse, error) {
	user, err := svc.userStore.SigningUser(ctx, in.GetEmail())
	if err != nil {
		return nil, errorResponse("could not signin", err)
	}

	err = utils.CheckPassword(user.Password, in.GetPassword())
	if err != nil {
		return nil, errorResponse("password incorrect", err)
	}

	tokenExpiration, err := utils.StringToTimeDuration(config.EnvVar.TokenExpiration)
	if err != nil {
		return nil, errorResponse("error parsing token expiration", err)
	}

	token, err := token.GenerateToken(tokenExpiration, user.ID, config.EnvVar.TokenExpiration)
	if err != nil {
		return nil, errorResponse("error generating token", err)
	}

	res := &pb.SigninUserResponse{
		Token: token,
	}

	return res, nil
}
