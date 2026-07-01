package executor

import (
	"context"
	"fmt"
	"workflow-engine/internal/models"
	"workflow-engine/internal/repository"

	"github.com/google/uuid"
)

type WorkflowExecutor struct {
	Repo *repository.WorkflowRepository
}

func (e *WorkflowExecutor) ExecuteTask(task models.WorkflowTask) bool {

	err := e.Repo.UpdateTaskExecutionStatus(context.Background(), task.ID, models.TaskRunning)
	if err != nil {
		fmt.Printf("Error updating task execution status: %v", err)
		return false
	}

	fmt.Printf("Doing task: %s", task.TaskName)
	fmt.Print(task)
	return true
}

func (e *WorkflowExecutor) Execute(workflowExecutionId uuid.UUID) error {

	workflowExecution, err := e.Repo.GetWorkflowExecutionById(workflowExecutionId)
	if err != nil {
		return err
	}

	fmt.Print(workflowExecution)

	workflowDefination, err := e.Repo.GetWorkflowDefinitionById(workflowExecution.WorkflowDefinationID)
	if err != nil {
		return err
	}

	for _, task := range workflowDefination.Tasks {
		success := e.ExecuteTask(task)
		if !success {
			return fmt.Errorf("failed to execute task: %s", task.TaskName)
		}
		err := e.Repo.UpdateTaskExecutionStatus(context.Background(), task.ID, models.TaskCompleted)
		if err != nil {
			return fmt.Errorf("failed to update task execution status: %v", err)
		}
	}

	return nil
}
