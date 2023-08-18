package grpcservice

import (
	"context"
	"fmt"

	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/proto/inventoryPb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (svc *Service) CreateProductComment(ctx context.Context, in *inventoryPb.CreateProductCommentRequest) (*inventoryPb.CreateProductCommentResponse, error) {
	productCommentToCreate, err := createProductCommentTypeNoId(in.GetProductComment())
	if err != nil {
		return nil, fmt.Errorf("Error creating product comment type: %v", err)
	}

	createdProductComment, err := svc.productCommentStore.Create(ctx, productCommentToCreate)
	if err != nil {
		return nil, fmt.Errorf("Erorr creating product comment in database: %v", err)
	}

	res := &inventoryPb.CreateProductCommentResponse{
		Result: &inventoryPb.CreatedResult{
			InsertedId: createdProductComment.InsertedID.(primitive.ObjectID).Hex(),
		},
	}

	return res, nil
}

func (svc *Service) GetProductComment(ctx context.Context, in *inventoryPb.GetProductCommentRequest) (*inventoryPb.GetProductCommentRespone, error) {
	mongoId, err := utils.GetMongoId(in.GetCommentId())
	if err != nil {
		return nil, fmt.Errorf("Could not parse string to mongoId: %v", err)
	}

	productComment, err := svc.productCommentStore.GetOne(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("Product comment not found: %v", err)
	}

	productCommentPb := createProductCommentPbResponse(productComment)
	res := &inventoryPb.GetProductCommentRespone{
		ProductComment: &productCommentPb,
	}

	return res, nil
}

func (svc *Service) GetAllProductComment(ctx context.Context, in *inventoryPb.GetAllProductCommentRequest) (*inventoryPb.GetAllProductCommentRespone, error) {
	mongoId, err := utils.GetMongoId(in.GetProductId())
	if err != nil {
		return nil, fmt.Errorf("Could not parse string to mongoId: %v", err)
	}

	productComments, err := svc.productCommentStore.GetAll(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("No comments found for this product in database: %v", err)
	}

	commentsArray := createAllProductCommentPbResponse(productComments)
	res := &inventoryPb.GetAllProductCommentRespone{
		ProductComment: commentsArray,
	}

	return res, nil
}
