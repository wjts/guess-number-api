package models

import (
	"gorm.io/gorm"
)

type Hint struct {
	gorm.Model `json:"-"`
	Date       string `gorm:"unique" json:"date"`
	Hint       uint16 `gorm:"not null" json:"hint"`
	CreatedBy  uint   `json:"-"`
	UpdatedBy  uint   `json:"-"`
}
