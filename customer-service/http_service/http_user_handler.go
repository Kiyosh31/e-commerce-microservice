package httpservice

import (
	"net/http"

	"github.com/Kiyosh31/e-commerce-microservice-common/token"
	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Service) signinUser(c *gin.Context) {
	var req types.SigninUserRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.userStore.Signing(c, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(user.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	tokenExpiration, err := utils.StringToTimeDuration(config.EnvVar.TokenExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, err := token.GenerateToken(tokenExpiration, user.ID, config.EnvVar.TokenSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := types.UserTokenResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, res)
}

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

	user, err := s.userStore.Create(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (s *Service) getUser(c *gin.Context) {
	mongoId, err := utils.GetMongoId(c.Param("userId"))
	if err != nil {
		log.Info().Msg("fue aqui")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.userStore.GetOne(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Service) updateUser(c *gin.Context) {
	mongoId, err := utils.GetMongoId(c.Param("userId"))
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

	updatedUser, err := s.userStore.Update(c, userToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (s *Service) deleteUser(c *gin.Context) {
	mongoId, err := utils.GetMongoId(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deletedUser, err := s.userStore.Delete(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, deletedUser)
}
