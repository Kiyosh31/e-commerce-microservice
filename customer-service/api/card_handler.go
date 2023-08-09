package api

import (
	"net/http"

	"github.com/Kiyosh31/e-commerce-microservice-common/utils"
	"github.com/Kiyosh31/e-commerce-microservice/customer/types"
	"github.com/gin-gonic/gin"
)

func (s *Service) createCard(c *gin.Context) {
	var req types.Card

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	card, err := s.cardStore.CreateCard(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, card)
}

func (s *Service) getCard(c *gin.Context) {
	mongoId, err := utils.GetMongoId(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	card, err := s.cardStore.GetCard(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, card)
}

func (s *Service) getAllCards(c *gin.Context) {
	mongoId, err := utils.GetMongoId(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cards, err := s.cardStore.GetAllCards(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (s *Service) updateCard(c *gin.Context) {
	cardMongoId, err := utils.GetMongoId(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req types.Card
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cardToUpdate := types.Card{
		ID:         cardMongoId,
		UserId:     req.UserId,
		Name:       req.Name,
		Number:     req.Number,
		SecretCode: req.SecretCode,
		Expiration: req.Expiration,
		Type:       req.Type,
		Default:    req.Default,
	}

	updatedCard, err := s.cardStore.UpdateCard(c, cardToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, updatedCard)
}

func (s *Service) deleteCard(c *gin.Context) {
	mongoId, err := utils.GetMongoId(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deletedCard, err := s.cardStore.DeleteCard(c, mongoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, deletedCard)
}
