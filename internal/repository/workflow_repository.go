package repository

import (
	"context"
	"errors"
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

func CreateTask(db *gorm.DB, workflowId uuid.UUID, taskName string, taskOrder int) (bool, error) {
	res := db.Create(&models.WorkflowTask{
		WorkflowID: workflowId,
		TaskOrder:  taskOrder,
		TaskName:   taskName,
	})

	if res.Error != nil {
		return false, res.Error
	}

	return true, nil
}

func GetTasksByWorkflow(db *gorm.DB, workflowId string) (*[]models.WorkflowTask, error) {
	var workflowTasks *[]models.WorkflowTask

	err := db.Find(&workflowTasks, "workflow_id = ?", workflowId).Error

	if err != nil {
		return nil, err
	}

	return workflowTasks, nil
}

func (r *WorkflowRepository) CreateWorkflow(ctx context.Context, workflowName string, tasks []string) (uuid.UUID, error) {
	if len(tasks) == 0 {
		return uuid.Nil, errors.New("No Tasks found for the workflow")
	}

	var createdWorkflow models.WorkflowDefinition

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&models.WorkflowDefinition{
			ID:        uuid.New(),
			Name:      workflowName,
			CreatedAt: time.Now(),
		}).Scan(&createdWorkflow).Error; err != nil {
			return err
		}

		for idx, taskName := range tasks {
			if err := tx.Create(&models.WorkflowTask{
				WorkflowID: createdWorkflow.ID,
				TaskOrder:  idx + 1,
				TaskName:   taskName,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return createdWorkflow.ID, nil
}

func (r *WorkflowRepository) GetWorkflow(ctx context.Context, workflowId uuid.UUID) (*models.WorkflowDefinition, error) {
	var workflow models.WorkflowDefinition
	res := r.DB.WithContext(ctx).Preload("Tasks").First(&workflow, workflowId)
	if res.Error != nil {
		return nil, res.Error
	}
	return &workflow, nil
}

func (r *WorkflowRepository) ListWorkflows(ctx context.Context) ([]*models.WorkflowDefinition, error) {
	var workflows []*models.WorkflowDefinition

	res := r.DB.WithContext(ctx).Preload("Tasks").Find(&workflows)

	if res.Error != nil {
		return nil, res.Error
	}

	return workflows, nil
}
