package grpcservice

import (
	"context"
	"fmt"

	grpcvalidators "github.com/Kiyosh31/e-commerce-microservice-common/grpc_validators"
	"github.com/Kiyosh31/e-commerce-microservice-common/middlewares"
	"github.com/Kiyosh31/e-commerce-microservice-common/token"
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/pb"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (svc *Service) CreateSeller(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// validate req
	violations := validateCreateUserRequest(in.GetUser())
	if violations != nil {
		return nil, grpcvalidators.InvalidArgumentError(violations)
	}

	// check is user with that email exists
	// if so you cannot create a new one
	existingUser, err := svc.userStore.GetOneByEmail(ctx, in.GetUser().GetEmail())
	if err == nil && &existingUser != nil {
		return nil, fmt.Errorf("User already exists")
	}

	hashedPassword, err := utils.HashPassword(in.GetUser().GetPassword())
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password: %v", err)
	}
	in.User.Password = hashedPassword
	in.User.Role = "seller"

	newUser := createUserTypeNoId(in.GetUser())

	createdUser, err := svc.userStore.Create(ctx, newUser)
	log.Printf("createdUser: %v", createdUser)
	if err != nil {
		return nil, fmt.Errorf("Could not create user in database: %v", err)
	}

	res := &pb.CreateUserResponse{
		Result: &pb.CreatedResult{
			InsertedId: createdUser.InsertedID.(primitive.ObjectID).Hex(),
		},
	}

	return res, nil
}

func (svc *Service) SigninSeller(ctx context.Context, in *pb.SigninUserRequest) (*pb.SigninUserResponse, error) {
	violations := validateSigninUser(in)
	if violations != nil {
		return nil, grpcvalidators.InvalidArgumentError(violations)
	}

	user, err := svc.userStore.Signing(ctx, in.GetEmail())
	if err != nil {
		return nil, fmt.Errorf("could not signin: %v", err)
	}

	if user.Role == "customer" {
		return nil, fmt.Errorf("This is a customer user please follow the right path /api/user/seller")
	}

	err = utils.CheckPassword(user.Password, in.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("password incorrect: %v", err)
	}

	tokenExpiration, err := utils.StringToTimeDuration(svc.env.TokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("error parsing token expiration: %v", err)
	}

	token, err := token.GenerateToken(tokenExpiration, user.ID, svc.env.TokenSecret)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %v", err)
	}

	res := &pb.SigninUserResponse{
		Token: token,
	}

	return res, nil
}

func (svc *Service) GetSeller(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	err := middlewares.ValidateTokenMatchesUser(ctx, in.GetUserId(), svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	violations := validateGetUser(in)
	if violations != nil {
		return nil, grpcvalidators.InvalidArgumentError(violations)
	}

	mongoId, err := utils.GetMongoId(in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("Could not parse string to mongoId: %v", err)
	}

	user, err := svc.userStore.GetOne(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("Error getting user: %v", err)
	}

	userPb := createUserPbResponse(user)
	res := &pb.GetUserResponse{
		User: &userPb,
	}

	return res, nil
}

func (svc *Service) UpdateSeller(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := middlewares.ValidateTokenMatchesUser(ctx, in.GetUserId(), svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	violations := validateUpdateUserRequest(in.GetUser())
	if violations != nil {
		return nil, grpcvalidators.InvalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(in.GetUser().GetPassword())
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password: %v", err)
	}
	in.User.Password = hashedPassword
	in.User.Role = "seller"

	userToUpdate, err := createUserTypeWithId(in.GetUserId(), in.GetUser())
	if err != nil {
		return nil, fmt.Errorf("Error parsing user to db: %v", err)
	}

	updatedUser, err := svc.userStore.Update(ctx, userToUpdate)
	if err != nil {
		return nil, fmt.Errorf("Error updating user in database: %v", err)
	}

	res := &pb.UpdateUserResponse{
		Result: &pb.UpdatedResult{
			MatchedCount:  updatedUser.MatchedCount,
			ModifiedCount: updatedUser.ModifiedCount,
			UpsertedCount: updatedUser.UpsertedCount,
		},
	}

	return res, nil
}

func (svc *Service) DeleteSeller(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := middlewares.ValidateTokenMatchesUser(ctx, in.GetUserId(), svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	violations := validateDeleteUserRequest(in)
	if violations != nil {
		return nil, grpcvalidators.InvalidArgumentError(violations)
	}

	mongoId, err := utils.GetMongoId(in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("Error parsing id to mongoId: %v", err)
	}

	deletedUser, err := svc.userStore.Delete(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("Error while deleting user from database: %v", err)
	}

	res := &pb.DeleteUserResponse{
		Result: &pb.DeletedResult{
			DeletedCount: deletedUser.DeletedCount,
		},
	}

	return res, nil
}
