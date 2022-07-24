package database

import (
	"github.com/wjts/guess-number-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var dbError error

func Connect(connectionString string) {
	Instance, dbError = gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if dbError != nil {
		panic("Cannot connect to DB")
	}
}
func Migrate() {
	Instance.AutoMigrate(&models.User{}, &models.Guess{}, &models.Hint{})
}
