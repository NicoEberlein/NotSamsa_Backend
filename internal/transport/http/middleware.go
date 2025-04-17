package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Authenticatior() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		jwtToken := strings.Split(authHeader, " ")[1]

		userId, err := verifyToken(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, nil)
		}

		c.Set("user", userId)
		c.Next()

	}
}
