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
	ID                   uuid.UUID           `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	WorkflowDefinition   *WorkflowDefinition `json:"-" gorm:"foreignKey:WorkflowID"`
	WorkflowDefinationID uuid.UUID           `json:"workflow_defination_id"`
	CurrentTaskOrder     int                 `json:"current_task_order"`
	Error                string
	Status               WorkflowStatus `json:"status" gorm:"default:WorkflowPending"`
	StartedAt            *time.Time     `json:"startedAt"`
	CompletedAt          *time.Time     `json:"completedAt"`

	TaskExecutions []TaskExecution
}

type TaskExecution struct {
	ID                  uuid.UUID          `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	WorkflowExecutionID uuid.UUID          `json:"-" gorm:"foreignKey:WorkflowExecutionID"`
	WorkflowExecution   *WorkflowExecution `json:"-"`
	ErrorMessage        string
	Status              TaskStatus `json:"status" gorm:"default:TaskPending"`
	CurrentTaskOrder    int
	TaskName            string     `json:"task_name"`
	StartedAt           *time.Time `json:"startedAt"`
	CompletedAt         *time.Time `json:"completedAt"`
}
