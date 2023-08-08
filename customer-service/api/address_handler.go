package api

import (
	"net/http"

	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Service) createAddress(c *gin.Context) {
	var req types.Address
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	card, err := s.addresStore.CreateAddress(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, card)
}

func (s *Service) getAddress(c *gin.Context) {
	mongoId, err := primitive.ObjectIDFromHex(c.Param("addressId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	address, err := s.addresStore.GetAddress(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, address)
}

func (s *Service) getAllAddress(c *gin.Context) {
	mongoId, err := primitive.ObjectIDFromHex(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	address, err := s.addresStore.GetAllAddress(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, address)
}

func (s *Service) updateAddress(c *gin.Context) {
	cardMongoId, err := primitive.ObjectIDFromHex(c.Param("addressId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req types.Address
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cardToUpdate := types.Address{
		ID:         cardMongoId,
		UserId:     req.UserId,
		Name:       req.Name,
		Address:    req.Address,
		PostalCode: req.PostalCode,
		Phone:      req.Phone,
		Default:    req.Default,
	}

	updatedAddress, err := s.addresStore.UpdateAddress(c, cardToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, updatedAddress)
}

func (s *Service) deleteAddress(c *gin.Context) {
	mongoId, err := primitive.ObjectIDFromHex(c.Param("addressId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deletedAddress, err := s.addresStore.DeleteAddress(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, deletedAddress)
}
