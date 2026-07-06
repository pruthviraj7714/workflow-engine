package handlers

import (
	"net/http"
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
	userId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	parsedUUID, ok := userId.(uuid.UUID)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user_id type",
		})
		return
	}

	var req struct {
		WorkflowName string   `json:"workflowName"`
		Tasks        []string `json:"tasks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflowId, err := h.WorkflowService.CreateWorkflow(c.Request.Context(), parsedUUID, req.WorkflowName, req.Tasks)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workflowId": workflowId})
}

func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	userId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	parsedUUID, ok := userId.(uuid.UUID)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user_id type",
		})
		return
	}

	workflowId := c.Param("workflowId")

	uuid, err := uuid.Parse(workflowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow, err := h.WorkflowService.GetWorkflow(c.Request.Context(), parsedUUID, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	userId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	parsedUUID, ok := userId.(uuid.UUID)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user_id type",
		})
		return
	}

	if workflows, err := h.WorkflowService.ListWorkflows(c.Request.Context(), parsedUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, workflows)
	}
}

func (h *WorkflowHandler) CreateWorkflowExecution(c *gin.Context) {
	userId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	parsedUUID, ok := userId.(uuid.UUID)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user_id type",
		})
		return
	}
	var req struct {
		WorkflowID string `json:"workflowId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid Body",
			"error":   err.Error(),
		})
		return
	}

	workflowId, err := uuid.Parse(req.WorkflowID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid Workflow ID",
			"error":   err.Error(),
		})
		return
	}

	workflowExecutionId, err := h.WorkflowService.CreateWorkflowExecution(c.Request.Context(), parsedUUID, workflowId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Workflow execution started successfully",
		"executionId": workflowExecutionId,
	})
}
