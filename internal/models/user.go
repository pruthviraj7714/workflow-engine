package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Username string    `json:"username" gorm:"unique"`
	Password string    `json:"password"`

	Workflows []WorkflowDefinition `json:"workflows" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
