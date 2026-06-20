package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Username string    `json:"username" gorm:"unique"`
	Password string    `json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
