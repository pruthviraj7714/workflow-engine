package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkflowDefinition struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Tasks      []WorkflowTask      `json:"tasks" gorm:"foreignKey:WorkflowDefinationID;constraint:OnDelete:CASCADE"`
	Executions []WorkflowExecution `json:"executions" gorm:"foreignKey:WorkflowDefinationID;constraint:OnDelete:CASCADE"`
}

type WorkflowTask struct {
	ID                   uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	WorkflowDefinationID uuid.UUID `json:"workflow_defination_id"`
	WorkflowDefination   *WorkflowDefinition
	TaskOrder            int    `json:"task_order"`
	TaskName             string `json:"task_name"`
}

type WorkflowStatus string

const (
	WorkflowPending   WorkflowStatus = "PENDING"
	WorkflowRunning   WorkflowStatus = "RUNNING"
	WorkflowCompleted WorkflowStatus = "COMPLETED"
	WorkflowFailed    WorkflowStatus = "FAILED"
)

type TaskStatus string

const (
	TaskPending   TaskStatus = "PENDING"
	TaskRunning   TaskStatus = "RUNNING"
	TaskCompleted TaskStatus = "COMPLETED"
	TaskFailed    TaskStatus = "FAILED"
)

type WorkflowExecution struct {
	ID                   uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	WorkflowDefinationID uuid.UUID
	WorkflowDefination   *WorkflowDefinition
	CurrentTaskOrder     int `json:"current_task_order"`
	Error                string
	Status               WorkflowStatus `json:"status" gorm:"default:PENDING"`
	StartedAt            *time.Time     `json:"startedAt"`
	CompletedAt          *time.Time     `json:"completedAt"`

	TaskExecutions []TaskExecution `json:"task_executions" gorm:"foreignKey:WorkflowExecutionID;constraint:OnDelete:CASCADE"`
}

type TaskExecution struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`

	WorkflowExecutionID uuid.UUID
	WorkflowExecution   *WorkflowExecution

	ErrorMessage     string
	Status           TaskStatus `json:"status" gorm:"default:PENDING"`
	CurrentTaskOrder int
	TaskName         string     `json:"task_name"`
	StartedAt        *time.Time `json:"startedAt"`
	CompletedAt      *time.Time `json:"completedAt"`
}
