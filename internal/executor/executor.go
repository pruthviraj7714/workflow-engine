package executor

import (
	"fmt"
	"workflow-engine/internal/repository"

	"github.com/google/uuid"
)

type WorkflowExecutor struct {
	repo *repository.WorkflowRepository
}

func (e *WorkflowExecutor) Execute(workflowExecutionId uuid.UUID) error {

	workflowExecution, err := e.repo.GetWorkflowExecutionById(workflowExecutionId)
	if err != nil {
		return err
	}

	fmt.Print(workflowExecution)

	workflowDefination, err := e.repo.GetWorkflowDefinitionById(workflowExecution.WorkflowDefinationID)
	if err != nil {
		return err
	}

	fmt.Print(workflowDefination.ID)

	return nil
}
