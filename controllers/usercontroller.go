package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wjts/guess-number-api/database"
	"github.com/wjts/guess-number-api/models"
)

type newUser struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
	Admin    bool
}

func RegisterUser(context *gin.Context) {
	var userRequest newUser
	if err := context.ShouldBindJSON(&userRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Email: userRequest.Email, Admin: userRequest.Admin}
	if err := user.HashPassword(userRequest.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Create(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email})
}
