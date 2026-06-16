package repository

import (
	"errors"
	"workflow-engine/internal/models"
	"workflow-engine/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) RegisterUser(username, password string) (string, error) {

	res := r.DB.Model(&models.User{}).Where("username = ?", username)

	if res.RowsAffected > 0 {
		return "", errors.New("Username Already Exists")
	}

	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return "", res.Error
	}

	res = r.DB.Create(&models.User{
		ID:       uuid.New(),
		Username: username,
		Password: hashedPassword,
	})

	if res.Error != nil {
		return "", res.Error
	}

	return r.LoginUser(username, password)

}

func (r *UserRepository) LoginUser(username, password string) (string, error) {
	var existingUser models.User

	res := r.DB.Where("username = ?", username).First(&existingUser)

	if res.RowsAffected == 0 {
		return "", errors.New("User with given username not found")
	}

	isValid := utils.CheckPasswordHash(password, existingUser.Password)

	if !isValid {
		return "", errors.New("Incorrect Password")
	}

	token, err := utils.GenerateToken(existingUser.ID)

	if err != nil {
		return "", err
	}

	return token, err
}
