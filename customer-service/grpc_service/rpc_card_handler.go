package grpcservice

import (
	"context"
	"fmt"

	"github.com/Kiyosh31/e-commerce-microservice-common/middlewares"
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (svc *Service) CreateCard(ctx context.Context, in *pb.CreateCardRequest) (*pb.CreateCardResponse, error) {
	userId, err := middlewares.AuthGrpcMiddleware(ctx, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateTokenMatchesUser(ctx, userId, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	cardType, err := createCardTypeNoId(in.GetCard())
	if err != nil {
		return nil, fmt.Errorf("Could not create card type: %v", err)
	}

	card, err := svc.cardStore.Create(ctx, cardType)
	if err != nil {
		return nil, fmt.Errorf("Could not create card in database: %v", err)
	}

	res := &pb.CreateCardResponse{
		Result: &pb.CreatedResult{
			InsertedId: card.InsertedID.(primitive.ObjectID).Hex(),
		},
	}

	return res, nil
}

func (svc *Service) GetCard(ctx context.Context, in *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	userId, err := middlewares.AuthGrpcMiddleware(ctx, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateTokenMatchesUser(ctx, userId, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	mongoId, err := utils.GetMongoId(in.GetCardId())
	if err != nil {
		return nil, fmt.Errorf("Could not parse string to mongoId: %v", err)
	}

	card, err := svc.cardStore.GetOne(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("Card does not exist in database: %v", err)
	}

	cardPb := createCardPbResponse(card)
	res := &pb.GetCardResponse{
		Card: &cardPb,
	}

	return res, nil
}

func (svc *Service) GetAllCard(ctx context.Context, in *pb.GetAllCardRequest) (*pb.GetAllCardResponse, error) {
	err := middlewares.ValidateTokenMatchesUser(ctx, in.GetUserId(), svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	mongoId, err := utils.GetMongoId(in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("Could not parse string to mongoId: %v", err)
	}

	cards, err := svc.cardStore.GetAll(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("No cards found for this user: %v", err)
	}

	cardsArray := createAllCardsPbresponse(cards)
	res := &pb.GetAllCardResponse{
		Card: cardsArray,
	}

	return res, nil
}

func (svc *Service) UpdateCard(ctx context.Context, in *pb.UpdateCardRequest) (*pb.UpdateCardResponse, error) {
	userId, err := middlewares.AuthGrpcMiddleware(ctx, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateTokenMatchesUser(ctx, userId, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	cardToUpdate, err := createCardTypeWithId(in.GetCardId(), in.GetCard())
	if err != nil {
		return nil, fmt.Errorf("Error creating card type: %v", err)
	}

	updatedCard, err := svc.cardStore.Update(ctx, cardToUpdate)
	if err != nil {
		return nil, fmt.Errorf("Error updating card: %v", err)
	}

	res := &pb.UpdateCardResponse{
		Result: &pb.UpdatedResult{
			MatchedCount:  updatedCard.MatchedCount,
			ModifiedCount: updatedCard.ModifiedCount,
			UpsertedCount: updatedCard.UpsertedCount,
		},
	}

	return res, nil
}

func (svc *Service) DeleteCard(ctx context.Context, in *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	userId, err := middlewares.AuthGrpcMiddleware(ctx, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateTokenMatchesUser(ctx, userId, svc.env.TokenSecret)
	if err != nil {
		return nil, err
	}

	mongoId, err := utils.GetMongoId(in.GetCardId())
	if err != nil {
		return nil, fmt.Errorf("Error parsing id to mongoId: %v", err)
	}

	deletedCard, err := svc.cardStore.Delete(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("Error deleting card from database: %v", err)
	}

	res := &pb.DeleteCardResponse{
		Result: &pb.DeletedResult{
			DeletedCount: deletedCard.DeletedCount,
		},
	}

	return res, nil
}
