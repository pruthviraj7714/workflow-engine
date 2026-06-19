package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Username string    `json:"username" gorm:"unique"`
	Password string    `json:"password"`
}
