package api

import (
	"log"
	"net/http"
	"workflow-engine/internal/config"
	"workflow-engine/internal/db"
	"workflow-engine/internal/executor"
	"workflow-engine/internal/handlers"
	"workflow-engine/internal/middlewares"
	"workflow-engine/internal/repository"
	"workflow-engine/internal/services"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	cfg := config.LoadConfig()

	database, err := db.Connect(cfg.DBURL)

	if err != nil {
		log.Fatal(err)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	userRepository := repository.NewUserRepository(database)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	authRouter := r.Group("/auth")
	{
		authRouter.POST("/register", userHandler.Register)
		authRouter.POST("/login", userHandler.Login)
	}

	workflowExecutor := executor.WorkflowExecutor{
		Repo: repository.NewWorkflowRepository(database),
	}

	workflowRepository := repository.NewWorkflowRepository(database)
	workflowService := services.NewWorkflowService(workflowRepository, &workflowExecutor)
	workflowHandler := handlers.NewWorkflowHandler(workflowService)

	workflowRouter := r.Group("/workflows")
	{
		workflowRouter.Use(middlewares.AuthMiddleware())
		workflowRouter.POST("/", workflowHandler.CreateWorkflow)
		workflowRouter.GET("/", workflowHandler.ListWorkflows)
		workflowRouter.GET("/:workflowId", workflowHandler.GetWorkflow)
		workflowRouter.POST("/workflow-executions", workflowHandler.CreateWorkflowExecution)
	}
	r.Run(":" + cfg.Port)
}
