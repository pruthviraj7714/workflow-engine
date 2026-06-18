package models

import (
	"time"
)

type WorkflowDefinition struct {
	ID        string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Name      string
	CreatedAt time.Time
}
