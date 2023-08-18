package grpcservice

import (
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/config"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/proto/inventoryPb"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/store"
	"github.com/Kiyosh31/e-commerce-microservice/inventory/types"
)

type Service struct {
	inventoryPb.UnimplementedInventoryServiceServer
	productStore        store.ProductStore
	productCommentStore store.ProductCommentStore
	env                 config.ConfigStruct
}

func NewService(productStore store.ProductStore, productCommentStore store.ProductCommentStore, env config.ConfigStruct) (*Service, error) {
	service := &Service{
		productStore:        productStore,
		productCommentStore: productCommentStore,
		env:                 env,
	}

	return service, nil
}

func createProductTypeNoId(in *inventoryPb.Product) (types.Product, error) {
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

func createProductTypeWithId(id string, in *inventoryPb.Product) (types.Product, error) {
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

func createProductPbResponse(in types.Product) inventoryPb.Product {
	return inventoryPb.Product{
		Id:          in.ID.Hex(),
		SellerId:    in.SellerId.Hex(),
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Brand:       in.Brand,
		Stars:       in.Stars,
	}
}

func createProductCommentTypeNoId(in *inventoryPb.ProductComment) (types.ProductComment, error) {
	productMongoId, err := utils.GetMongoId(in.GetProductId())
	if err != nil {
		return types.ProductComment{}, err
	}

	res := types.ProductComment{
		ProductId:  productMongoId,
		UserName:   in.GetUserName(),
		Comment:    in.GetComment(),
		RatingStar: in.GetRatingStar(),
	}

	return res, nil
}

func createProductCommentPbResponse(in types.ProductComment) inventoryPb.ProductComment {
	return inventoryPb.ProductComment{
		Id:         in.ID.Hex(),
		ProductId:  in.ProductId.Hex(),
		UserName:   in.UserName,
		Comment:    in.Comment,
		RatingStar: in.RatingStar,
	}
}

func createAllProductCommentPbResponse(in []types.ProductComment) []*inventoryPb.ProductComment {
	var productComments []*inventoryPb.ProductComment

	for _, prodComm := range in {
		comm := createProductCommentPbResponse(prodComm)
		productComments = append(productComments, &comm)
	}

	return productComments
}
