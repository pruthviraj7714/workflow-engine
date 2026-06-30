package repository

import (
	"context"
	"fmt"
	"time"
	"workflow-engine/internal/executor"
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

func (r *WorkflowRepository) CreateWorkflowDefinition(ctx context.Context, workflowName string, tasks []string) (uuid.UUID, error) {
	workflowDefination := &models.WorkflowDefinition{
		Name: workflowName,
	}

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&workflowDefination).Error; err != nil {
			return err
		}

		fmt.Print(workflowDefination)

		for idx, taskName := range tasks {
			if err := tx.Create(&models.WorkflowTask{
				WorkflowDefinationID: workflowDefination.ID,
				TaskOrder:            idx + 1,
				TaskName:             taskName,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return workflowDefination.ID, nil
}

func (r *WorkflowRepository) GetWorkflow(ctx context.Context, workflowId uuid.UUID) (*models.WorkflowDefinition, error) {
	var workflow models.WorkflowDefinition
	res := r.DB.WithContext(ctx).Preload("Tasks").First(&workflow, workflowId)
	if res.Error != nil {
		return nil, res.Error
	}
	return &workflow, nil
}

func (r *WorkflowRepository) ListWorkflows(ctx context.Context) ([]models.WorkflowDefinition, error) {
	var workflows []models.WorkflowDefinition

	res := r.DB.WithContext(ctx).Preload("Tasks").Find(&workflows)

	if res.Error != nil {
		return nil, res.Error
	}

	return workflows, nil
}

func (r *WorkflowRepository) CreateWorkflowExecution(ctx context.Context, workflowId uuid.UUID) (uuid.UUID, error) {

	var createdWorkflowExecution *models.WorkflowExecution

	now := time.Now()

	err := r.DB.WithContext(ctx).Create(&models.WorkflowExecution{
		WorkflowDefinationID: workflowId,
		Status:               models.WorkflowPending,
		ID:                   uuid.New(),
		StartedAt:            &now,
	}).Scan(&createdWorkflowExecution).Error

	if err != nil {
		return uuid.Nil, err
	}

	executor := executor.WorkflowExecutor{
		Repo: , 
	}

	return createdWorkflowExecution.ID, nil
}

func (r *WorkflowRepository) GetWorkflowExecutionById(workflowExecutionId uuid.UUID) (*models.WorkflowExecution, error) {
	var workflowExecution models.WorkflowExecution

	res := r.DB.Find(&workflowExecution, workflowExecutionId)

	if res.Error != nil {
		return nil, res.Error
	}

	return &workflowExecution, nil
}

func (r *WorkflowRepository) GetWorkflowDefinitionById(workflowDefinitionId uuid.UUID) (*models.WorkflowDefinition, error) {
	var workflow models.WorkflowDefinition
	res := r.DB.Preload("Tasks").First(&workflow, workflowDefinitionId)
	if res.Error != nil {
		return nil, res.Error
	}
	return &workflow, nil
}


func (r *WorkflowRepository) CreateTaskExecutions(ctx context.Context, workflowExecutionId uuid.UUID, tasks []models.WorkflowTask) error {


	for idx, task := range tasks {

	if err := r.DB.Create(&models.TaskExecution{
		WorkflowExecutionID: workflowExecutionId,
		Status: models.TaskPending,
		TaskName: task.TaskName,
		TaskOrder: idx + 1,
	}).Error; err != nil {
		return err
	}



	}



}
