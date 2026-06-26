package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkflowDefinition struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	Tasks      []WorkflowTask      `json:"tasks" gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE"`
	Executions []WorkflowExecution `json:"executions" gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE"`
}

type WorkflowTask struct {
	ID         uuid.UUID           `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	WorkflowID uuid.UUID           `json:"workflow_id"`
	TaskOrder  int                 `json:"task_order"`
	TaskName   string              `json:"task_name"`
	Workflow   *WorkflowDefinition `json:"-" gorm:"foreignKey:WorkflowID"`
}

type WorkflowStatus string

const (
	Pending   WorkflowStatus = "PENDING"
	Running   WorkflowStatus = "RUNNING"
	Completed WorkflowStatus = "COMPLETED"
	Failed    WorkflowStatus = "FAILED"
)

type WorkflowExecution struct {
	ID         uuid.UUID           `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Workflow   *WorkflowDefinition `json:"-" gorm:"foreignKey:WorkflowID"`
	WorkflowID uuid.UUID           `json:"workflow_id"`

	Status      WorkflowStatus `json:"status" gorm:"default:PENDING"`
	StartedAt   *time.Time     `json:"startedAt"`
	CompletedAt *time.Time     `json:"completedAt"`
}
