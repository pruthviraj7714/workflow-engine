package executor

import (
	"workflow-engine/internal/repository"

	"github.com/google/uuid"
)

type WorkflowExecutor struct {
	repp *repository.WorkflowRepository
}

func (e *WorkflowExecutor) Execute(workflowId uuid.UUID) error {

	return nil
}
