package api

import (
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
	"github.com/gin-gonic/gin"
)

type Service struct {
	userStore  store.UserStore
	cardStore  store.CardStore
	listenAddr string
	router     *gin.Engine
}

func NewService(userStore store.UserStore, cardStore store.CardStore, listenAddr string) *Service {
	server := &Service{
		userStore:  userStore,
		cardStore:  cardStore,
		listenAddr: listenAddr,
	}
	router := gin.Default()
	registerRoutes(router, server)
	server.router = router

	return server
}

func registerRoutes(router *gin.Engine, service *Service) {
	api := router.Group("/api")
	user := api.Group("/user")
	{
		user.POST("/", service.createUser)
		user.GET("/:id", service.getUser)
		user.PUT("/:id", service.updateUser)
		user.DELETE("/:id", service.deleteUser)

		card := user.Group("/card")
		{
			card.POST("/", service.createCard)
			card.GET("/:cardId", service.getCard)
			card.GET("/all/:userId", service.getAllCards)
			card.PUT("/:cardId", service.updateCard)
			card.DELETE("/:cardId", service.deleteCard)
		}
	}
}

func (s *Service) Start() error {
	return s.router.Run(s.listenAddr)
}

func errorResponse(err error) *gin.H {
	return &gin.H{"Errors": err.Error()}
}
