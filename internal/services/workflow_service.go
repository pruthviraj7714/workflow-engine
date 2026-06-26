package services

import (
	"context"
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

func (s *WorkflowService) CreateWorkflow(ctx context.Context, workflowName string, tasks []string) (uuid.UUID, error) {
	return s.WorkflowRepo.CreateWorkflow(ctx, workflowName, tasks)
}

func (s *WorkflowService) GetWorkflow(ctx context.Context, workflowId uuid.UUID) (*models.WorkflowDefinition, error) {
	return s.WorkflowRepo.GetWorkflow(ctx, workflowId)
}

func (s *WorkflowService) ListWorkflows(ctx context.Context) ([]*models.WorkflowDefinition, error) {
	return s.WorkflowRepo.ListWorkflows(ctx)
}
