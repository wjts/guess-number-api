package models

import (
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model `json:"-"`
	Date       string `gorm:"unique" json:"date"`
	Answer     uint16 `gorm:"not null" json:"hint"`
	CreatedBy  uint   `json:"-"`
	UpdatedBy  uint   `json:"-"`
}
