package services

import (
	"workflow-engine/internal/models"
	"workflow-engine/internal/repository"
)

type WorkflowTaskService struct {
	WorkflowTaskRepo *repository.WorkflowTasksRepository
}

func NewWorkflowTaskService(workflowTaskRepo *repository.WorkflowTasksRepository) *WorkflowTaskService {
	return &WorkflowTaskService{
		WorkflowTaskRepo: workflowTaskRepo,
	}
}

func (s *WorkflowTaskService) CreateTask(workflowId, taskName string, taskOrder int) (bool, error) {
	return s.WorkflowTaskRepo.CreateTask(workflowId, taskName, taskOrder)
}

func (s *WorkflowTaskService) GetTasksByWorkflow(workflowId string) (*[]models.WorkflowTask, error) {
	return s.WorkflowTaskRepo.GetTasksByWorkflow(workflowId)
}
