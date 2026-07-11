package services

import (
	"context"
	"errors"
	"workflow-engine/internal/models"
	"workflow-engine/internal/rabbitmq"
	"workflow-engine/internal/repository"

	"github.com/google/uuid"
)

type WorkflowService struct {
	WorkflowRepo *repository.WorkflowRepository
	Producer     *rabbitmq.Producer
}

func NewWorkflowService(workflowRepo *repository.WorkflowRepository, producer *rabbitmq.Producer) *WorkflowService {
	return &WorkflowService{
		WorkflowRepo: workflowRepo,
		Producer:     producer,
	}
}

func (s *WorkflowService) CreateWorkflow(ctx context.Context, userId uuid.UUID, workflowName string, tasks []string) (uuid.UUID, error) {
	if len(tasks) == 0 {
		return uuid.Nil, errors.New("No Tasks found for the workflow")
	}

	return s.WorkflowRepo.CreateWorkflowDefinition(ctx, workflowName, tasks, userId)
}

func (s *WorkflowService) GetWorkflow(ctx context.Context, userId uuid.UUID, workflowId uuid.UUID) (*models.WorkflowDefinition, error) {
	return s.WorkflowRepo.GetWorkflow(ctx, userId, workflowId)
}

func (s *WorkflowService) ListWorkflows(ctx context.Context, userId uuid.UUID) ([]models.WorkflowDefinition, error) {
	return s.WorkflowRepo.ListWorkflows(ctx, userId)
}

func (s *WorkflowService) CreateWorkflowExecution(ctx context.Context, userId uuid.UUID, workflowId uuid.UUID) (uuid.UUID, error) {

	//Validate workflow exists
	workflow, err := s.WorkflowRepo.GetWorkflow(ctx, userId, workflowId)
	if err != nil {
		return uuid.Nil, err
	}

	// Create workflow execution
	executionId, err := s.WorkflowRepo.CreateWorkflowExecution(ctx, workflowId)
	if err != nil {
		return uuid.Nil, err
	}

	// Create task executions
	err = s.WorkflowRepo.CreateTaskExecutions(ctx, executionId, workflow.Tasks)
	if err != nil {
		return uuid.Nil, err
	}

	// Start executor
	err = s.Producer.PublishWorkflowExecution(executionId)
	if err != nil {
		return uuid.Nil, err
	}
	// Return execution ID
	return executionId, nil
}
