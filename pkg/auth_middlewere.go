package pkg

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorization"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) { return []byte(os.Getenv("SECRET")), nil })

		if !token.Valid || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorization, invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		c.Next()
	}
}
