package handlers

import (
	"net/http"
	"workflow-engine/internal/services"

	"github.com/gin-gonic/gin"
)

type WorkflowTaskHandler struct {
	WorkflowTaskService *services.WorkflowTaskService
}

func NewWorkflowTaskHandler(workflowTaskService *services.WorkflowTaskService) *WorkflowTaskHandler {
	return &WorkflowTaskHandler{
		WorkflowTaskService: workflowTaskService,
	}
}

func (h *WorkflowTaskHandler) CreateTask(c *gin.Context) {
	var req struct {
		WorkflowID string `json:"workflowId"`
		TaskOrder  int    `json:"taskOrder"`
		TaskName   string `json:"taskName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	done, err := h.WorkflowTaskService.CreateTask(req.WorkflowID, req.TaskName, req.TaskOrder)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
		return
	}

	if !done {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task created successfully",
	})
}

func (h *WorkflowTaskHandler) GetTasksByWorkflow(c *gin.Context) {
	workflowId := c.Param("workflowId")

	if workflowId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Workflow ID is missing",
		})
		return
	}

	tasks, err := h.WorkflowTaskService.GetTasksByWorkflow(workflowId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}
