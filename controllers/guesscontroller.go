package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wjts/guess-number-api/database"
	"github.com/wjts/guess-number-api/models"
	"gorm.io/gorm"
)

type newGuess struct {
	Date  string `binding:"required"`
	Guess uint16 `binding:"required"`
}

func MakeGuess(context *gin.Context) {
	var guessRequest newGuess
	if err := context.ShouldBindJSON(&guessRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := time.Parse("2006-01-02", guessRequest.Date); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	var user models.User
	database.Instance.First(&user, "email = ?", context.GetString("UserEmail"))

	guess := models.Guess{Date: guessRequest.Date, Guess: guessRequest.Guess}
	if err := database.Instance.Create(&guess).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, guess)
}

type editGuess struct {
	Guess uint16 `binding:"required"`
}

func UpdateGuess(context *gin.Context) {
	var user models.User
	database.Instance.First(&user, "email = ?", context.GetString("UserEmail"))

	var guess models.Guess
	if err := database.Instance.First(&guess, "date = ?", context.Param("date")).Error; err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	var editGuess editGuess
	if err := context.ShouldBindJSON(&editGuess); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Model(&guess).Updates(models.Guess{Guess: editGuess.Guess}).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, guess)
}

func GetGuesses(context *gin.Context) {
	var guesses []models.Guess
	database.Instance.Find(&guesses)

	context.JSON(http.StatusOK, guesses)
}

func GetGuess(context *gin.Context) {
	var guess models.Guess
	err := database.Instance.Where("date = ?", context.Param("date")).First(&guess).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
			return
		}

		context.Status(http.StatusBadRequest)
		return
	}

	context.JSON(http.StatusOK, guess)
}
