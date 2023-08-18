package grpcservice

import (
	"regexp"

	grpcvalidators "github.com/Kiyosh31/e-commerce-microservice-common/grpc_validators"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/customerPb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

var (
	isValidName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidDate = regexp.MustCompile(`^\d{4}\-(0?[1-9]|1[012])\-(0?[1-9]|[12][0-9]|3[01])$`).MatchString
)

func validateCreateUserRequest(in *customerPb.User) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := grpcvalidators.ValidateName(in.GetName()); err != nil {

		violations = append(violations, grpcvalidators.FieldValidation("name", err))
	}

	if err := grpcvalidators.ValidateName(in.GetLastName()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("lastName", err))
	}

	if err := grpcvalidators.ValidateBirth(in.GetBirth()); err != nil {
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

func validateSigninUser(in *customerPb.SigninUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := grpcvalidators.ValidateEmail(in.GetEmail()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("email", err))
	}

	if err := grpcvalidators.ValidatePassword(in.GetPassword()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("password", err))
	}

	return violations
}

func validateGetUser(in *customerPb.GetUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := grpcvalidators.ValidateMongoId(in.GetUserId()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("password", err))
	}

	return violations
}

func validateUpdateUserRequest(in *customerPb.User) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := grpcvalidators.ValidateName(in.GetName()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("name", err))
	}

	if err := grpcvalidators.ValidateName(in.GetLastName()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("lastName", err))
	}

	if err := grpcvalidators.ValidateBirth(in.GetBirth()); err != nil {
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

func validateDeleteUserRequest(in *customerPb.DeleteUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := grpcvalidators.ValidateMongoId(in.GetUserId()); err != nil {
		violations = append(violations, grpcvalidators.FieldValidation("password", err))
	}

	return violations
}
