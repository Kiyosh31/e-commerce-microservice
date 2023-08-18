package grpcservice

import (
	"context"
	"fmt"

	"github.com/Kiyosh31/e-commerce-microservice-common/middlewares"
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/customerPb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (svc *Service) CreateAddress(ctx context.Context, in *customerPb.CreateAddressRequest) (*customerPb.CreateAddressResponse, error) {
	err := middlewares.ValidateTokenMatchesUser(ctx, in.GetAddress().GetUserId(), svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	newAddress, err := createAddressTypeNoId(in.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("Error creating address type: %v", err)
	}

	createdAddress, err := svc.addresStore.Create(ctx, newAddress)
	if err != nil {
		return nil, fmt.Errorf("Could not create address: %v", err)
	}

	res := &customerPb.CreateAddressResponse{
		Result: &customerPb.CreatedResult{
			InsertedId: createdAddress.InsertedID.(primitive.ObjectID).Hex(),
		},
	}

	return res, nil
}

func (svc *Service) GetAddress(ctx context.Context, in *customerPb.GetAddressRequest) (*customerPb.GetAddressResponse, error) {
	userId, err := middlewares.AuthGrpcMiddleware(ctx, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateTokenMatchesUser(ctx, userId, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	mongoId, err := utils.GetMongoId(in.GetAddressId())
	if err != nil {
		return nil, fmt.Errorf("Could not parse string to mongoId: %v", err)
	}

	address, err := svc.addresStore.GetOne(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("Could not found the address: %v", err)
	}

	addressPb := createAddressPbResponse(address)
	res := &customerPb.GetAddressResponse{
		Address: &addressPb,
	}

	return res, nil
}

func (svc *Service) GetAllAddress(ctx context.Context, in *customerPb.GetAllAddressRequest) (*customerPb.GetAllAddressResponse, error) {
	err := middlewares.ValidateTokenMatchesUser(ctx, in.GetUserId(), svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	mongoId, err := utils.GetMongoId(in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("Could not parse string to mongoId: %v", err)
	}

	addresses, err := svc.addresStore.GetAll(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("No address for this user: %v", err)
	}

	addressArray := createAllAddressPbResponse(addresses)
	res := &customerPb.GetAllAddressResponse{
		Address: addressArray,
	}

	return res, nil
}

func (svc *Service) UpdateAddress(ctx context.Context, in *customerPb.UpdateAddressRequest) (*customerPb.UpdateAddressResponse, error) {
	userId, err := middlewares.AuthGrpcMiddleware(ctx, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateTokenMatchesUser(ctx, userId, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	addressToUpdate, err := createAddressTypeWithId(in.GetAddressId(), in.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("Could not create address type: %v", err)
	}

	updatedAddress, err := svc.addresStore.Update(ctx, addressToUpdate)
	if err != nil {
		return nil, fmt.Errorf("Error updating address: %v", err)
	}

	res := &customerPb.UpdateAddressResponse{
		Result: &customerPb.UpdatedResult{
			MatchedCount:  updatedAddress.MatchedCount,
			ModifiedCount: updatedAddress.ModifiedCount,
			UpsertedCount: updatedAddress.UpsertedCount,
		},
	}

	return res, nil
}

func (svc *Service) DeleteAddress(ctx context.Context, in *customerPb.DeleteAddressRequest) (*customerPb.DeleteAddressResponse, error) {
	mongoId, err := utils.GetMongoId(in.GetAddressId())
	if err != nil {
		return nil, fmt.Errorf("Error parsing id to mongoId: %v", err)
	}

	deletedAddress, err := svc.addresStore.Delete(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("Error while deleting address from database: %v", err)
	}

	res := &customerPb.DeleteAddressResponse{
		Result: &customerPb.DeletedResult{
			DeletedCount: deletedAddress.DeletedCount,
		},
	}

	return res, nil
}
