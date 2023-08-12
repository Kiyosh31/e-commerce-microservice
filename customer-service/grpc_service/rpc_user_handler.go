package grpcservice

import (
	"context"
	"fmt"

	grpcvalidators "github.com/Kiyosh31/e-commerce-microservice-common/grpc_validators"
	"github.com/Kiyosh31/e-commerce-microservice-common/middlewares"
	"github.com/Kiyosh31/e-commerce-microservice-common/token"
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/pb"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func validateUserRequest(in *pb.User) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := grpcvalidators.ValidateName(in.GetName()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("name", err))
	}

	if err := grpcvalidators.ValidateName(in.GetLastName()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("lastName", err))
	}

	if err := grpcvalidators.ValidateName(in.GetBirth()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("birth", err))
	}

	if err := grpcvalidators.ValidateEmail(in.GetEmail()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("email", err))
	}

	if err := grpcvalidators.ValidatePassword(in.GetPassword()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("password", err))
	}

	return violations
}

func validateSigninUser(req *pb.SigninUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := grpcvalidators.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("email", err))
	}

	if err := grpcvalidators.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("password", err))
	}

	return violations
}

func (svc *Service) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// violations := validateUserRequest(in.GetUser())
	// if violations != nil {
	// 	return nil, grpcvalidators.InvalidArgumentError(violations)
	// }

	hashedPassword, err := utils.HashPassword(in.GetUser().GetPassword())
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password: %v", err)
	}

	in.User.Password = hashedPassword

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

func (svc *Service) SigninUser(ctx context.Context, in *pb.SigninUserRequest) (*pb.SigninUserResponse, error) {
	violations := validateSigninUser(in)
	if violations != nil {
		return nil, grpcvalidators.InvalidArgumentError(violations)
	}

	user, err := svc.userStore.Signing(ctx, in.GetEmail())
	if err != nil {
		return nil, fmt.Errorf("could not signin: %v", err)
	}

	err = utils.CheckPassword(user.Password, in.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("password incorrect: %v", err)
	}

	tokenExpiration, err := utils.StringToTimeDuration(config.EnvVar.TokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("error parsing token expiration: %v", err)
	}

	token, err := token.GenerateToken(tokenExpiration, user.ID, config.EnvVar.TokenSecret)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %v", err)
	}

	res := &pb.SigninUserResponse{
		Token: token,
	}

	return res, nil
}

func (svc *Service) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	err := middlewares.ValidateTokenMatchesUser(ctx, in.GetUserId(), config.EnvVar.TokenSecret)
	if err != nil {
		return nil, err
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

func (svc *Service) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := middlewares.ValidateTokenMatchesUser(ctx, in.GetUserId(), config.EnvVar.TokenSecret)
	if err != nil {
		return nil, err
	}

	// violations := validateUserRequest(in.GetUser())
	// if violations != nil {
	// 	return nil, grpcvalidators.InvalidArgumentError(violations)
	// }

	hashedPassword, err := utils.HashPassword(in.GetUser().GetPassword())
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password: %v", err)
	}
	in.User.Password = hashedPassword

	userToUpdate, err := createUserTypeWithId(in.GetUserId(), in.GetUser())
	if err != nil {
		return nil, fmt.Errorf("Error parsing user to db: %v", err)
	}

	updatedUser, err := svc.userStore.Update(ctx, userToUpdate)
	if err != nil {
		return nil, fmt.Errorf("Error updating user in database: %v", err)
	}

	res := &pb.UpdateUserResponse{
		Result: &pb.UpdateResult{
			MatchedCount:  updatedUser.MatchedCount,
			ModifiedCount: updatedUser.ModifiedCount,
			UpsertedCount: updatedUser.UpsertedCount,
		},
	}

	return res, nil
}

func (svc *Service) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := middlewares.ValidateTokenMatchesUser(ctx, in.GetUserId(), config.EnvVar.TokenSecret)
	if err != nil {
		return nil, err
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
		Result: &pb.DeleteResult{
			DeletedCount: deletedUser.DeletedCount,
		},
	}

	return res, nil
}
