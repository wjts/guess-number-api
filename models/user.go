package models

import (
	"github.com/wjts/guess-number-api/argon2id"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email             string      `gorm:"unique;not null" json:"email"`
	Password          string      `gorm:"not null" json:"-"`
	Admin             bool        `gorm:"not null;default:false" json:"-"`
	Guess             []Guess     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ChallengesCreated []Challenge `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ChallengesUpdated []Challenge `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := argon2id.DefaultParams().Hash(password)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) (bool, error) {
	ok, err := argon2id.DefaultParams().Verify(providedPassword, user.Password)
	if err != nil {
		return false, err
	}
	return ok, nil
}
