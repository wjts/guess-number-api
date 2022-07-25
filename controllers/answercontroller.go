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

type newAnswer struct {
	Date   string `binding:"required"`
	Answer uint16 `binding:"required"`
}

func CreateAnswer(context *gin.Context) {
	var answerRequest newAnswer
	if err := context.ShouldBindJSON(&answerRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := time.Parse("2006-01-02", answerRequest.Date); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	var user models.User
	database.Instance.First(&user, "email = ?", context.GetString("UserEmail"))
	answer := models.Answer{Date: answerRequest.Date, Answer: answerRequest.Answer, CreatedBy: user.ID}
	if err := database.Instance.Create(&answer).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, answer)
}

type editAnswer struct {
	Answer uint16 `binding:"required"`
}

func UpdateAnswer(context *gin.Context) {
	var user models.User
	database.Instance.First(&user, "email = ?", context.GetString("UserEmail"))

	var answer models.Answer
	if err := database.Instance.First(&answer, "date = ?", context.Param("date")).Error; err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	var answerRequest editAnswer
	if err := context.ShouldBindJSON(&answerRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Model(&answer).Updates(models.Answer{Answer: answerRequest.Answer, UpdatedBy: user.ID}).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, answer)
}

func GetAnswers(context *gin.Context) {
	var answers []models.Answer
	database.Instance.Find(&answers)

	context.JSON(http.StatusOK, answers)
}

func GetAnswer(context *gin.Context) {
	var answer models.Answer
	err := database.Instance.Where("date = ?", context.Param("date")).First(&answer).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
			return
		}

		context.Status(http.StatusBadRequest)
		return
	}

	context.JSON(http.StatusOK, answer)
}
