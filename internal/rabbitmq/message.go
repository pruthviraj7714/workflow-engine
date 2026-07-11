package rabbitmq

import "github.com/google/uuid"

type WorkflowExecutionMessage struct {
	ExecutionID uuid.UUID `json:"execution_id"`
}
