package executor

import (
	"context"
	"fmt"
	"os"
	"time"
	"workflow-engine/internal/models"
	"workflow-engine/internal/repository"

	"github.com/google/uuid"
)

type WorkflowExecutor struct {
	Repo *repository.WorkflowRepository
}

func (e *WorkflowExecutor) ExecuteTask(task models.TaskExecution) error {
	fmt.Printf("Doing task: %s\n", task.TaskName)
	time.Sleep(1 * time.Second)
	fmt.Printf("Task Done: %s\n", task.TaskName)
	return nil
}

func (e *WorkflowExecutor) Execute(ctx context.Context, workflowExecutionId uuid.UUID) error {

	err := e.Repo.UpdateWorkflowExecutionStatus(ctx, workflowExecutionId, models.WorkflowRunning)

	if err != nil {
		e.Repo.UpdateWorkflowExecutionStatus(ctx, workflowExecutionId, models.WorkflowFailed)
		return err
	}

	workflowExecution, err := e.Repo.GetWorkflowExecutionById(workflowExecutionId)
	if err != nil {
		e.Repo.UpdateWorkflowExecutionStatus(ctx, workflowExecutionId, models.WorkflowFailed)
		return err
	}

	fmt.Fprintf(os.Stdout, "%+v\n", workflowExecution)

	for _, task := range workflowExecution.TaskExecutions {
		err := e.Repo.UpdateTaskExecutionStatus(ctx, task.ID, models.TaskRunning)
		if err != nil {
			fmt.Printf("Error updating task execution status: %v", err)
			return err
		}

		err = e.ExecuteTask(task)
		if err != nil {
			return fmt.Errorf("failed to execute task: %s", task.TaskName)
		}

		err = e.Repo.UpdateTaskExecutionStatus(ctx, task.ID, models.TaskCompleted)

		if err != nil {
			return fmt.Errorf("failed to update task execution status: %v", err)
		}
	}

	err = e.Repo.UpdateWorkflowExecutionStatus(ctx, workflowExecutionId, models.WorkflowCompleted)

	if err != nil {
		fmt.Printf("Error updating workflow execution status: %v", err)
		return err
	}

	return nil
}
