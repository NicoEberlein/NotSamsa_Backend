package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strings"
)

func (h *Handler) Authenticator() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
			return
		}

		jwtToken := strings.Split(authHeader, " ")[1]

		userId, err := verifyToken(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		}

		c.Set("user", userId)
		c.Next()

	}
}

func (h *Handler) CheckCollectionOwnership(canBeParticipant bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectionId := c.Param("collectionId")
		if len(collectionId) == 0 {
			c.Next()
		}

		userId := c.GetString("user")
		user, _ := h.UserService.FindById(c, userId)

		collection, err := h.CollectionService.FindById(c, collectionId)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		userHasAccess := false

		if collection.OwnerId == user.Id {
			userHasAccess = true
		}
		if canBeParticipant && slices.Contains(collection.Participants, user) {
			userHasAccess = true
		}

		if userHasAccess {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}

	}
}
