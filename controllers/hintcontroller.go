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

type newHint struct {
	Date string `binding:"required"`
	Hint uint16 `binding:"required"`
}

func CreateHint(context *gin.Context) {
	var hintRequest newHint
	if err := context.ShouldBindJSON(&hintRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := time.Parse("2006-01-02", hintRequest.Date); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	var user models.User
	database.Instance.First(&user, "email = ?", context.GetString("UserEmail"))
	hint := models.Hint{Date: hintRequest.Date, Hint: hintRequest.Hint, CreatedBy: user.ID}
	if err := database.Instance.Create(&hint).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"ID": hint.ID})
}

type editHint struct {
	Hint uint16 `binding:"required"`
}

func UpdateHint(context *gin.Context) {
	var user models.User
	database.Instance.First(&user, "email = ?", context.GetString("UserEmail"))

	var hint models.Hint
	if err := database.Instance.First(&hint, "date = ?", context.Param("date")).Error; err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	var hintRequest editHint
	if err := context.ShouldBindJSON(&hintRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Model(&hint).Updates(models.Hint{Hint: hintRequest.Hint, UpdatedBy: user.ID}).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, hint)
}

func GetHints(context *gin.Context) {
	var hints []models.Hint
	database.Instance.Find(&hints)

	context.JSON(http.StatusOK, hints)
}

func GetHint(context *gin.Context) {
	var hint models.Hint
	err := database.Instance.Where("date = ?", context.Param("date")).First(&hint).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
			return
		}

		context.Status(http.StatusBadRequest)
		return
	}

	context.JSON(http.StatusOK, hint)
}
