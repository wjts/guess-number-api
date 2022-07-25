package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wjts/guess-number-api/controllers"
	"github.com/wjts/guess-number-api/database"
	"github.com/wjts/guess-number-api/middlewares"
)

func main() {
	database.Connect("db/sqlite.db")
	database.Migrate()

	router := initRouter()
	router.Run("0.0.0.0:8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))

	v1 := router.Group("/v1")
	{
		v1.GET("/ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		v1.GET("/secure-ping", middlewares.Auth(), func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"message": context.GetString("UserEmail"),
			})
		})

		v1.POST("/token", controllers.GenerateToken)
		v1.POST("/user/register", controllers.RegisterUser)

		v1.GET("/hints", middlewares.Auth(), controllers.GetHints)
		v1.GET("/hints/:date", middlewares.Auth(), controllers.GetHint)
		v1.POST("/hints", middlewares.Auth(), middlewares.RoleAdmin(), controllers.CreateHint)
		v1.PATCH("/hints/:date", middlewares.Auth(), middlewares.RoleAdmin(), controllers.UpdateHint)

		v1.GET("/guesses", middlewares.Auth(), controllers.GetGuesses)
		v1.GET("/guesses/:date", middlewares.Auth(), controllers.GetGuess)
		v1.POST("/guesses", middlewares.Auth(), controllers.MakeGuess)
		v1.PATCH("/guesses/:date", middlewares.Auth(), controllers.UpdateGuess)

		v1.GET("/answers", middlewares.Auth(), controllers.GetAnswers)
		v1.GET("/answers/:date", middlewares.Auth(), controllers.GetAnswer)
		v1.POST("/answers", middlewares.Auth(), middlewares.RoleAdmin(), controllers.CreateAnswer)
		v1.PATCH("/answers/:date", middlewares.Auth(), middlewares.RoleAdmin(), controllers.UpdateAnswer)

		v1.GET("/results", middlewares.Auth(), controllers.GetResults)
	}

	return router
}
