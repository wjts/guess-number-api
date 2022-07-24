package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wjts/guess-number-api/auth"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		context.Set("UserEmail", claims.Email)
		context.Next()
	}
}
