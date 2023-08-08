package api

import (
	"net/http"

	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Service) createUser(c *gin.Context) {
	var req types.User

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse((err)))
		return
	}

	req.Password = hashedPassword

	user, err := s.userStore.CreateUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (s *Service) getUser(c *gin.Context) {
	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.userStore.GetOneUser(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Service) updateUser(c *gin.Context) {
	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req types.User

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userToUpdate := types.User{
		ID:       mongoId,
		Name:     req.Name,
		LastName: req.LastName,
		Birth:    req.Birth,
		Email:    req.Email,
		Password: req.Password,
	}

	updatedUser, err := s.userStore.UpdateUser(c, userToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (s *Service) deleteUser(c *gin.Context) {
	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deletedUser, err := s.userStore.DeleteUser(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, deletedUser)
}
