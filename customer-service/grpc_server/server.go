package grpcserver

import (
	"github.com/Kiyosh31/e-commerce-microservice/customer/proto/pb"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func errorResponse(message string, err error) error {
	return status.Errorf(codes.Internal, message+": ", err)
}
