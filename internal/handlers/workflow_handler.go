package handlers

import (
	"workflow-engine/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorkflowHandler struct {
	WorkflowService *services.WorkflowService
}

func NewWorkflowHandler(workflowService *services.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		WorkflowService: workflowService,
	}
}

func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {

	var req struct {
		WorkflowName string   `json:"workflowName"`
		Tasks        []string `json:"tasks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if workflowId, err := h.WorkflowService.CreateWorkflow(c.Request.Context(), req.WorkflowName, req.Tasks); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"workflowId": workflowId})
	}
}

func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	workflowId := c.Param("workflowId")

	uuid, err := uuid.Parse(workflowId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if workflow, err := h.WorkflowService.GetWorkflow(c.Request.Context(), uuid); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, workflow)
	}
}

func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	if workflows, err := h.WorkflowService.ListWorkflows(c.Request.Context()); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, workflows)
	}
}
