package grpcservice

import (
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/pb"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
)

type Service struct {
	pb.UnimplementedCustomerServiceServer
	userStore   store.UserStore
	cardStore   store.CardStore
	addresStore store.AddressStore
}

func NewService(userStore store.UserStore, addressStore store.AddressStore, cardStore store.CardStore) (*Service, error) {
	service := &Service{
		userStore:   userStore,
		addresStore: addressStore,
		cardStore:   cardStore,
	}

	return service, nil
}

func createUserTypeNoId(in *pb.User) types.User {
	user := types.User{
		Name:     in.GetName(),
		LastName: in.GetLastName(),
		Birth:    in.GetBirth(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	return user
}

func createUserTypeWithId(id string, in *pb.User) (types.User, error) {
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
	}

	return user, nil
}

func createUserPbResponse(in types.User) pb.User {
	return pb.User{
		Id:       in.ID.Hex(),
		Name:     in.Name,
		LastName: in.LastName,
		Birth:    in.Birth,
		Email:    in.Email,
		Password: in.Password,
	}
}
