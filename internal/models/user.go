package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey"`
	Username string    `json:"username" gorm:"unique"`
	Password string    `json:"password"`
}
