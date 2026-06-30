package services

import (
	"context"
	"errors"
	"workflow-engine/internal/executor"
	"workflow-engine/internal/models"
	"workflow-engine/internal/repository"

	"github.com/google/uuid"
)

type WorkflowService struct {
	WorkflowRepo *repository.WorkflowRepository
	Executor     *executor.WorkflowExecutor
}

func NewWorkflowService(workflowRepo *repository.WorkflowRepository, executor *executor.WorkflowExecutor) *WorkflowService {
	return &WorkflowService{
		WorkflowRepo: workflowRepo,
		Executor:     executor,
	}
}

func (s *WorkflowService) CreateWorkflow(ctx context.Context, workflowName string, tasks []string) (uuid.UUID, error) {
	if len(tasks) == 0 {
		return uuid.Nil, errors.New("No Tasks found for the workflow")
	}

	return s.WorkflowRepo.CreateWorkflowDefinition(ctx, workflowName, tasks)
}

func (s *WorkflowService) GetWorkflow(ctx context.Context, workflowId uuid.UUID) (*models.WorkflowDefinition, error) {
	return s.WorkflowRepo.GetWorkflow(ctx, workflowId)
}

func (s *WorkflowService) ListWorkflows(ctx context.Context) ([]models.WorkflowDefinition, error) {
	return s.WorkflowRepo.ListWorkflows(ctx)
}

func (s *WorkflowService) CreateWorkflowExecution(ctx context.Context, workflowId uuid.UUID) (uuid.UUID, error) {

	//Validate workflow exists

	workflow, err := s.WorkflowRepo.GetWorkflow(ctx, workflowId)
	if err != nil {
		return uuid.Nil, err
	}

	// Create workflow execution
	executionId, err := s.WorkflowRepo.CreateWorkflowExecution(ctx, workflowId)
	if err != nil {
		return uuid.Nil, err
	}

	// Create task executions
	err = s.WorkflowRepo.CreateWorkflowExecution(ctx, executionId, workflow.Tasks)
	if err != nil {
		return uuid.Nil, err
	}

	// Start executor

	// Return execution ID

	return executionId, nil
}
