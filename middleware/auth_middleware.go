package middleware

import (
	"jwt-authentication/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		tokenData, err := c.Cookie("session")

		if err != nil {
			c.JSON(404, gin.H{"message": "No cookie available"})
			c.Abort()
			return
		}

		claims := &utils.Claims{}

		token, err1 := jwt.ParseWithClaims(tokenData, claims, func(t *jwt.Token) (any, error) {
			return utils.SecretKey, nil
		})

		if err1 != nil || !token.Valid {
			c.JSON(404, gin.H{"message": "Token error"})
			c.Abort()
			return
		}

		c.Next()
	}

}
