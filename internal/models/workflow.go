package models

import (
	"time"
)

type WorkflowDefinition struct {
	ID        string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	Tasks []WorkflowTask `json:"tasks" gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE"`
}

type WorkflowTask struct {
	ID         string              `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	WorkflowID string              `json:"workflow_id"`
	TaskOrder  int                 `json:"task_order"`
	TaskName   string              `json:"task_name"`
	Workflow   *WorkflowDefinition `json:"workflow"`
}
