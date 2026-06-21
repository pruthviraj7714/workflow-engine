package services

import (
	"workflow-engine/internal/models"
	"workflow-engine/internal/repository"

	"github.com/google/uuid"
)

type WorkflowService struct {
	WorkflowRepo *repository.WorkflowRepository
}

func NewWorkflowService(workflowRepo *repository.WorkflowRepository) *WorkflowService {
	return &WorkflowService{
		WorkflowRepo: workflowRepo,
	}
}

func (s *WorkflowService) CreateWorkflow(workflowName string, tasks []string) (string, error) {
	return s.WorkflowRepo.CreateWorkflow(workflowName, tasks)
}

func (s *WorkflowService) GetWorkflow(workflowId uuid.UUID) (*models.WorkflowDefinition, error) {
	return s.WorkflowRepo.GetWorkflow(workflowId)
}

func (s *WorkflowService) ListWorkflows() ([]*models.WorkflowDefinition, error) {
	return s.WorkflowRepo.ListWorkflows()
}
