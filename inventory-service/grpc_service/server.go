package grpcservice

import (
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/config"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/proto/pb"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/store"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/types"
)

type Service struct {
	pb.UnimplementedInventoryServiceServer
	productStore store.ProductStore
	env          config.ConfigStruct
}

func NewService(productStore store.ProductStore, env config.ConfigStruct) (*Service, error) {
	service := &Service{
		productStore: productStore,
		env:          env,
	}

	return service, nil
}

func createProductTypeNoId(in *pb.Product) (types.Product, error) {
	sellerMongoId, err := utils.GetMongoId(in.GetSellerId())
	if err != nil {
		return types.Product{}, err
	}

	product := types.Product{
		SellerId:    sellerMongoId,
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Price:       in.GetPrice(),
		Brand:       in.GetBrand(),
		Stars:       in.GetStars(),
	}

	return product, nil
}

func createProductTypeWithId(id string, in *pb.Product) (types.Product, error) {
	productMongoId, err := utils.GetMongoId(id)
	if err != nil {
		return types.Product{}, err
	}

	sellerMongoId, err := utils.GetMongoId(in.GetSellerId())
	if err != nil {
		return types.Product{}, err
	}

	product := types.Product{
		ID:          productMongoId,
		SellerId:    sellerMongoId,
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Price:       in.GetPrice(),
		Brand:       in.GetBrand(),
		Stars:       in.GetStars(),
	}

	return product, nil
}

func createProductPbResponse(in types.Product) pb.Product {
	return pb.Product{
		Id:          in.ID.Hex(),
		SellerId:    in.SellerId.Hex(),
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Brand:       in.Brand,
		Stars:       in.Stars,
	}
}
