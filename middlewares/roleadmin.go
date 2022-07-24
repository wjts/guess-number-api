package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wjts/guess-number-api/database"
	"github.com/wjts/guess-number-api/models"
)

func RoleAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user models.User
		if err := database.Instance.Where("email = ?", context.GetString("UserEmail")).First(&user).Error; err != nil {
			context.Status(http.StatusUnauthorized)
			context.Abort()
			return
		}

		if !user.Admin {
			context.Status(http.StatusUnauthorized)
			context.Abort()
			return
		}

		context.Next()
	}
}
