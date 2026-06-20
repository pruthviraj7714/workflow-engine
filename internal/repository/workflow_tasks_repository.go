package repository

import (
	"workflow-engine/internal/models"

	"gorm.io/gorm"
)

type WorkflowTasksRepository struct {
	DB *gorm.DB
}

func NewWorkflowTasksRepository(db *gorm.DB) *WorkflowTasksRepository {
	return &WorkflowTasksRepository{
		DB: db,
	}
}

func (r *WorkflowTasksRepository) CreateTask(workflowId, taskName string, taskOrder int) (bool, error) {
	res := r.DB.Create(&models.WorkflowTask{
		WorkflowID: workflowId,
		TaskOrder:  taskOrder,
		TaskName:   taskName,
	})

	if res.Error != nil {
		return false, res.Error
	}

	return true, nil
}

func (r *WorkflowTasksRepository) GetTasksByWorkflow(workflowId string) (*[]models.WorkflowTask, error) {
	var workflowTasks *[]models.WorkflowTask

	err := r.DB.Find(&workflowTasks, "workflow_id = ?", workflowId).Error

	if err != nil {
		return nil, err
	}

	return workflowTasks, nil
}
