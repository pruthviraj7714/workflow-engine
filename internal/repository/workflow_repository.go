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

func CreateTask(db *gorm.DB, workflowId, taskName string, taskOrder int) (bool, error) {
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

func (r *WorkflowRepository) CreateWorkflow(ctx context.Context, workflowName string, tasks []string) (string, error) {
	if len(tasks) == 0 {
		return "", errors.New("No Tasks found for the workflow")
	}

	var workflow models.WorkflowDefinition

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		workflow := models.WorkflowDefinition{
			Name:      workflowName,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&workflow).Error; err != nil {
			return err
		}

		for idx, taskName := range tasks {
			if err := tx.Create(&models.WorkflowTask{
				WorkflowID: workflow.ID,
				TaskOrder:  idx + 1,
				TaskName:   taskName,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return workflow.ID, nil
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
