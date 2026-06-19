package repository

import (
	"time"
	"workflow-engine/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkflowRepository struct {
	DB *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{
		DB: db,
	}
}

func (r *WorkflowRepository) CreateWorkflow(workflowName string) (string, error) {
	res := r.DB.Create(&models.WorkflowDefinition{
		Name:      workflowName,
		CreatedAt: time.Now(),
	})

	if res.Error != nil {
		return "", res.Error
	}

	return res.Statement.Dest.(*models.WorkflowDefinition).ID, nil
}

func (r *WorkflowRepository) GetWorkflow(workflowId uuid.UUID) (*models.WorkflowDefinition, error) {
	var workflow models.WorkflowDefinition
	res := r.DB.First(&workflow, workflowId)
	if res.Error != nil {
		return nil, res.Error
	}
	return &workflow, nil
}

func (r *WorkflowRepository) ListWorkflows() ([]*models.WorkflowDefinition, error) {
	var workflows []*models.WorkflowDefinition

	res := r.DB.Find(&workflows)

	if res.Error != nil {
		return nil, res.Error
	}

	return workflows, nil
}
