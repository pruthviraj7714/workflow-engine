package repository

import (
	"context"
	"fmt"
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
	now := time.Now()

	res := r.DB.WithContext(ctx).Create(&models.WorkflowExecution{
		WorkflowDefinationID: workflowId,
		Status:               models.WorkflowPending,
		ID:                   uuid.New(),
		StartedAt:            &now,
	})

	if res.Error != nil {
		return uuid.Nil, res.Error
	}

	return res.Statement.ReflectValue.Interface().(models.WorkflowExecution).ID, nil
}

func (r *WorkflowRepository) GetWorkflowExecutionById(workflowExecutionId uuid.UUID) (*models.WorkflowExecution, error) {
	var workflowExecution models.WorkflowExecution

	res := r.DB.Preload("TaskExecutions").Find(&workflowExecution, workflowExecutionId)

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

	for _, task := range tasks {
		if err := r.DB.Create(&models.TaskExecution{
			WorkflowExecutionID: workflowExecutionId,
			Status:              models.TaskPending,
			TaskName:            task.TaskName,
			CurrentTaskOrder:    task.TaskOrder,
		}).Error; err != nil {
			return err
		}

	}

	return nil
}

func (r *WorkflowRepository) UpdateWorkflowExecutionStatus(ctx context.Context, executionId uuid.UUID, workflowStatus models.WorkflowStatus) error {

	err := r.DB.Model(&models.WorkflowExecution{}).Where("id = ?", executionId).UpdateColumn("status", workflowStatus)
	if err.Error != nil {
		return err.Error
	}

	return nil
}

func (r *WorkflowRepository) UpdateTaskExecutionStatus(ctx context.Context, taskId uuid.UUID, taskStatus models.TaskStatus) error {

	err := r.DB.Model(&models.TaskExecution{}).Where("id = ?", taskId).UpdateColumn("status", taskStatus)
	if err.Error != nil {
		return err.Error
	}

	return nil
}
