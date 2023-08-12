package httpservice

import (
	"github.com/Kiyosh31/e-commerce-microservice-common/middlewares"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
	"github.com/gin-gonic/gin"
)

type Service struct {
	userStore   store.UserStore
	cardStore   store.CardStore
	addresStore store.AddressStore
	listenAddr  string
	router      *gin.Engine
}

func NewService(userStore store.UserStore, cardStore store.CardStore, addressStore store.AddressStore, listenAddr string) (*Service, error) {
	server := &Service{
		userStore:   userStore,
		cardStore:   cardStore,
		addresStore: addressStore,
		listenAddr:  listenAddr,
	}

	server.registerRoutes()
	return server, nil
}

func (s *Service) registerRoutes() {
	router := gin.Default()
	router.Use(gin.Recovery())

	api := router.Group("/api")
	user := api.Group("/user")
	{
		user.POST("/signin", s.signinUser)
		user.POST("/", s.createUser)
		user.GET("/:userId", middlewares.AuthHttpMiddleware(config.EnvVar.TokenSecret), s.getUser)
		user.PUT("/:userId", middlewares.AuthHttpMiddleware(config.EnvVar.TokenSecret), s.updateUser)
		user.DELETE("/:userId", middlewares.AuthHttpMiddleware(config.EnvVar.TokenSecret), s.deleteUser)

		card := user.Group("/card").Use(middlewares.AuthHttpMiddleware(config.EnvVar.TokenSecret))
		{
			card.POST("/", s.createCard)
			card.GET("/:cardId", s.getCard)
			card.GET("/all/:userId", s.getAllCards)
			card.PUT("/:cardId", s.updateCard)
			card.DELETE("/:cardId", s.deleteCard)
		}

		address := user.Group("/address").Use(middlewares.AuthHttpMiddleware(config.EnvVar.TokenSecret))
		{
			address.POST("/", s.createAddress)
			address.GET("/:addressId", s.getAddress)
			address.GET("/all/:userId", s.getAllAddress)
			address.PUT("/:addressId", s.updateAddress)
			address.DELETE("/:addressId", s.deleteAddress)
		}
	}

	s.router = router
}

func (s *Service) Start() error {
	return s.router.Run(s.listenAddr)
}

func errorResponse(err error) *gin.H {
	return &gin.H{"Errors": err.Error()}
}
