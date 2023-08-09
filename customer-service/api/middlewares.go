package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Kiyosh31/e-commerce-microservice-common/token"
	"github.com/Kiyosh31/e-commerce-microservice/customer/config"
	"github.com/Kiyosh31/e-commerce-microservice/customer/store"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	authHeaderKey  = "Authorization"
	authTypeBearer = "bearer"
	authPayloadKey = "userId"
)

func authMiddleware(userStore store.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get(authHeaderKey)
		fields := strings.Fields(authHeader)

		// validate token is provided
		if len(fields) == 0 {
			err := errors.New("Authorization is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// validate auth has Bearer token
		if len(fields) < 2 {
			err := errors.New("Invalid authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// validate auth is Bearer token type
		authType := strings.ToLower(fields[0])
		if authType != authTypeBearer {
			err := fmt.Errorf("Unsuported auth type: %s", authType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// extracting the token code
		// checking token is valid
		tokenCode := fields[1]
		tokenUserId, err := token.ValidateToken(tokenCode, config.EnvVar.TokenSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// checking the userId matched token id this means the user is modifyng it's own info
		// otherwise this request has no permission to continue
		mongoId, err := primitive.ObjectIDFromHex(fmt.Sprint(tokenUserId))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		user, err := userStore.GetOneUser(context.Background(), mongoId)
		if tokenUserId != user.ID.Hex() {
			err := errors.New("Cannot modify the information, does not belong to you")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(authPayloadKey, tokenUserId)
		c.Next()
	}
}
