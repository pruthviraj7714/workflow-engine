package executor

import (
	"fmt"
	"workflow-engine/internal/models"

	"github.com/google/uuid"
)

type WorkflowExecutor struct {
}

func ExecuteTask(task models.WorkflowTask) bool {
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
		ExecuteTask(task)
	}

	return nil
}
