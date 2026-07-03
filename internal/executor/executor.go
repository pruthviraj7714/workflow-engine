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

func (e *WorkflowExecutor) ExecuteTask(task models.TaskExecution) error {
	err := e.Repo.UpdateTaskExecutionStatus(context.Background(), task.ID, models.TaskRunning)
	if err != nil {
		fmt.Printf("Error updating task execution status: %v", err)
		return err
	}

	fmt.Printf("Doing task: %s\n", task.TaskName)
	fmt.Println(task)
	return nil
}

func (e *WorkflowExecutor) Execute(workflowExecutionId uuid.UUID) error {

	e.Repo.UpdateWorkflowExecutionStatus(context.Background(), workflowExecutionId, models.WorkflowRunning)

	workflowExecution, err := e.Repo.GetWorkflowExecutionById(workflowExecutionId)
	if err != nil {
		e.Repo.UpdateWorkflowExecutionStatus(context.Background(), workflowExecutionId, models.WorkflowFailed)
		return err
	}

	fmt.Print(workflowExecution)

	for _, task := range workflowExecution.TaskExecutions {
		err := e.ExecuteTask(task)
		if err != nil {
			return fmt.Errorf("failed to execute task: %s", task.TaskName)
		}

		err = e.Repo.UpdateTaskExecutionStatus(context.Background(), task.ID, models.TaskCompleted)

		if err != nil {
			return fmt.Errorf("failed to update task execution status: %v", err)
		}
	}

	e.Repo.UpdateWorkflowExecutionStatus(context.Background(), workflowExecutionId, models.WorkflowCompleted)

	return nil
}
