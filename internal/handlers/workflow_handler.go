package handlers

import (
	"workflow-engine/internal/services"
)

type WorkflowHandler struct {
	WorkflowService *services.WorkflowService
}

func NewWorkflowHandler(workflowService *services.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		WorkflowService: workflowService,
	}
}
