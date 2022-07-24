package models

import (
	"gorm.io/gorm"
)

type Guess struct {
	gorm.Model `json:"-"`
	Date       string `gorm:"unique" json:"date"`
	Guess      uint16 `gorm:"not null" json:"guess"`
	UserID     uint   `json:"-"`
}
