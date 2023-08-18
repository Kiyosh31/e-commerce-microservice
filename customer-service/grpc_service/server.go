package grpcservice

import (
	"fmt"

	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/customerPb"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
)

type Service struct {
	customerPb.UnimplementedCustomerServiceServer
	userStore   store.UserStore
	cardStore   store.CardStore
	addresStore store.AddressStore
	env         config.ConfigStruct
}

func NewService(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore, env config.ConfigStruct) (*Service, error) {
	service := &Service{
		userStore:   userStore,
		addresStore: addressStore,
		cardStore:   cardStore,
		env:         env,
	}

	return service, nil
}

func createUserTypeNoId(in *customerPb.User) types.User {
	user := types.User{
		Name:     in.GetName(),
		LastName: in.GetLastName(),
		Birth:    in.GetBirth(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
		Role:     in.GetRole(),
	}

	return user
}

func createUserTypeWithId(id string, in *customerPb.User) (types.User, error) {
	mongoId, err := utils.GetMongoId(id)
	if err != nil {
		return types.User{}, err
	}

	user := types.User{
		ID:       mongoId,
		Name:     in.GetName(),
		LastName: in.GetLastName(),
		Birth:    in.GetBirth(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
		Role:     in.GetRole(),
	}

	return user, nil
}

func createUserPbResponse(in types.User) customerPb.User {
	return customerPb.User{
		Id:       in.ID.Hex(),
		Name:     in.Name,
		LastName: in.LastName,
		Birth:    in.Birth,
		Email:    in.Email,
		Password: in.Password,
		Role:     in.Role,
	}
}

func createAddressTypeNoId(in *customerPb.Address) (types.Address, error) {
	mongoId, err := utils.GetMongoId(in.GetUserId())
	if err != nil {
		return types.Address{}, err
	}

	address := types.Address{
		UserId:     mongoId,
		Name:       in.GetName(),
		Address:    in.GetAddress(),
		PostalCode: in.GetPostalCode(),
		Phone:      in.GetPhone(),
		Default:    in.GetDefault(),
	}

	return address, nil
}

func createAddressTypeWithId(addressId string, in *customerPb.Address) (types.Address, error) {
	addressMongoId, err := utils.GetMongoId(addressId)
	if err != nil {
		return types.Address{}, fmt.Errorf("could not parse addressId: %v", err)
	}

	userMongoId, err := utils.GetMongoId(in.GetUserId())
	if err != nil {
		return types.Address{}, fmt.Errorf("Could not parse userId: %v", err)
	}

	address := types.Address{
		ID:         addressMongoId,
		UserId:     userMongoId,
		Name:       in.GetName(),
		Address:    in.GetAddress(),
		PostalCode: in.GetPostalCode(),
		Phone:      in.GetPhone(),
		Default:    in.GetDefault(),
	}

	return address, nil
}

func createAddressPbResponse(in types.Address) customerPb.Address {
	return customerPb.Address{
		Id:         in.ID.Hex(),
		UserId:     in.UserId.Hex(),
		Name:       in.Name,
		Address:    in.Address,
		PostalCode: in.PostalCode,
		Phone:      in.Phone,
		Default:    in.Default,
	}
}

func createAllAddressPbResponse(in []types.Address) []*customerPb.Address {
	var addressArray []*customerPb.Address

	for _, address := range in {
		addr := createAddressPbResponse(address)
		addressArray = append(addressArray, &addr)
	}

	return addressArray
}

func createCardTypeNoId(in *customerPb.Card) (types.Card, error) {
	userMongoId, err := utils.GetMongoId(in.GetUserId())
	if err != nil {
		return types.Card{}, fmt.Errorf("Could not parse userId: %v", err)
	}

	card := types.Card{
		UserId:     userMongoId,
		Name:       in.GetName(),
		Number:     in.GetNumber(),
		SecretCode: in.GetSecretCode(),
		Expiration: in.GetExpiration(),
		Type:       in.GetType(),
		Default:    in.GetDefault(),
	}

	return card, nil
}

func createCardTypeWithId(id string, in *customerPb.Card) (types.Card, error) {
	cardMongoId, err := utils.GetMongoId(id)
	if err != nil {
		return types.Card{}, fmt.Errorf("Could not parse userId: %v", err)
	}

	userMongoId, err := utils.GetMongoId(in.GetUserId())
	if err != nil {
		return types.Card{}, fmt.Errorf("Could not parse userId: %v", err)
	}

	card := types.Card{
		ID:         cardMongoId,
		UserId:     userMongoId,
		Name:       in.GetName(),
		Number:     in.GetNumber(),
		SecretCode: in.GetSecretCode(),
		Expiration: in.GetExpiration(),
		Type:       in.GetType(),
		Default:    in.GetDefault(),
	}

	return card, nil
}

func createCardPbResponse(in types.Card) customerPb.Card {
	return customerPb.Card{
		Id:         in.ID.Hex(),
		UserId:     in.UserId.Hex(),
		Name:       in.Name,
		Number:     in.Number,
		SecretCode: in.SecretCode,
		Expiration: in.Expiration,
		Type:       in.Type,
		Default:    in.Default,
	}
}

func createAllCardsPbresponse(in []types.Card) []*customerPb.Card {
	var cardsArray []*customerPb.Card

	for _, card := range in {
		crd := createCardPbResponse(card)
		cardsArray = append(cardsArray, &crd)
	}

	return cardsArray
}
