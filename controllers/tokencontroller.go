package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wjts/guess-number-api/auth"
	"github.com/wjts/guess-number-api/database"
	"github.com/wjts/guess-number-api/models"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Where("email = ?", request.Email).First(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	valid, credentialError := user.CheckPassword(request.Password)
	if !valid || credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	tokenString, err := auth.GenerateJWT(user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
