package db

import (
	"log"
	"workflow-engine/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.WorkflowDefinition{}, &models.WorkflowTask{}, &models.WorkflowExecution{})

	return db, nil
}
