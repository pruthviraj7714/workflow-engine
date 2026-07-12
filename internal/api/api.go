package api

import (
	"context"
	"log"
	"net/http"
	"workflow-engine/internal/config"
	"workflow-engine/internal/db"
	"workflow-engine/internal/executor"
	"workflow-engine/internal/handlers"
	"workflow-engine/internal/middlewares"
	"workflow-engine/internal/rabbitmq"
	"workflow-engine/internal/repository"
	"workflow-engine/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	cfg := config.LoadConfig()

	database, err := db.Connect(cfg.DBURL)

	if err != nil {
		log.Fatal(err)
	}

	mq, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQURL)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := mq.Close(); err != nil {
			log.Println("failed to close RabbitMQ:", err)
		}
	}()

	producer := rabbitmq.NewProducer(mq)
	consumer := rabbitmq.NewConsumer(mq)

	r.Use(cors.Default())

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

	workflowRepository := repository.NewWorkflowRepository(database)
	workflowService := services.NewWorkflowService(workflowRepository, producer)
	workflowHandler := handlers.NewWorkflowHandler(workflowService)

	executor := executor.WorkflowExecutor{
		Repo: workflowRepository,
	}

	go consumer.Start(context.Background(), executor.Execute)

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
