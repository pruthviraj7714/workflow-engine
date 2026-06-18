package services

import "workflow-engine/internal/repository"

type WorkflowService struct {
	WorkflowRepo *repository.WorkflowRepository
}

func NewWorkflowService(workflowRepo *repository.WorkflowRepository) *WorkflowService {
	return &WorkflowService{
		WorkflowRepo: workflowRepo,
	}
}

func (s *WorkflowService) CreateWorkflow(workflowName string) (string, error) {
	return s.WorkflowRepo.CreateWorkflow(workflowName)
}
