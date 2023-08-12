package api

import (
	"net/http"

	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
	"github.com/gin-gonic/gin"
)

func (s *Service) createAddress(c *gin.Context) {
	var req types.Address
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	card, err := s.addresStore.Create(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, card)
}

func (s *Service) getAddress(c *gin.Context) {
	mongoId, err := utils.GetMongoId(c.Param("addressId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	address, err := s.addresStore.GetOne(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, address)
}

func (s *Service) getAllAddress(c *gin.Context) {

	mongoId, err := utils.GetMongoId(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	address, err := s.addresStore.GetAll(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, address)
}

func (s *Service) updateAddress(c *gin.Context) {
	cardMongoId, err := utils.GetMongoId(c.Param("addressId"))
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

	updatedAddress, err := s.addresStore.Update(c, cardToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, updatedAddress)
}

func (s *Service) deleteAddress(c *gin.Context) {
	mongoId, err := utils.GetMongoId(c.Param("addressId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deletedAddress, err := s.addresStore.Delete(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, deletedAddress)
}
