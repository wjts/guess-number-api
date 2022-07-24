package models

import (
	"time"

	"gorm.io/gorm"
)

type Challenge struct {
	gorm.Model
	Date      time.Time `gorm:"unique"`
	Result    uint16
	CreatedBy uint
	UpdatedBy uint
}
