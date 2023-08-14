package grpcservice

import (
	"context"
	"fmt"

	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/proto/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (svc *Service) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	productToCreate, err := createProductTypeNoId(in.GetProduct())
	if err != nil {
		return nil, fmt.Errorf("Could not create product type: %v", err)
	}

	createdUser, err := svc.productStore.Create(ctx, productToCreate)
	if err != nil {
		return nil, fmt.Errorf("Could not create product into database: %v", err)
	}

	res := &pb.CreateProductResponse{
		Result: &pb.CreatedResult{
			InsertedId: createdUser.InsertedID.(primitive.ObjectID).Hex(),
		},
	}

	return res, nil
}

func (svc *Service) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	mongoId, err := utils.GetMongoId(in.GetProductId())
	if err != nil {
		return nil, fmt.Errorf("Could not parse string to mongoId: %v", err)
	}

	product, err := svc.productStore.GetOne(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("Product not found: %v", err)
	}

	productPb := createProductPbResponse(product)
	res := &pb.GetProductResponse{
		Product: &productPb,
	}

	return res, nil
}

func (svc *Service) UpdateProduct(ctx context.Context, in *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	productType, err := createProductTypeWithId(in.GetProductId(), in.GetProduct())
	if err != nil {
		return nil, fmt.Errorf("Could not create product type: %v", err)
	}

	productToUpdate, err := svc.productStore.Update(ctx, productType)
	if err != nil {
		return nil, fmt.Errorf("Could not update product in database: %v", err)
	}

	res := &pb.UpdateProductResponse{
		Result: &pb.UpdatedResult{
			MatchedCount:  productToUpdate.MatchedCount,
			ModifiedCount: productToUpdate.ModifiedCount,
			UpsertedCount: productToUpdate.UpsertedCount,
		},
	}

	return res, nil
}

func (svc *Service) DeleteProduct(ctx context.Context, in *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	mongoId, err := utils.GetMongoId(in.GetProductId())
	if err != nil {
		return nil, fmt.Errorf("Error parsing id to mongoId: %v", err)
	}

	deletedProduct, err := svc.productStore.Delete(ctx, mongoId)
	if err != nil {
		return nil, fmt.Errorf("Could  ot delete product from database: %v", err)
	}

	res := &pb.DeleteProductResponse{
		Result: &pb.DeletedResult{
			DeletedCount: deletedProduct.DeletedCount,
		},
	}

	return res, nil
}
