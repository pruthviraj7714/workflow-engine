package repository

import (
	"context"
	"errors"
	"workflow-engine/internal/models"
	"workflow-engine/internal/utils"

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

func (r *UserRepository) RegisterUser(ctx context.Context, username, password string) (string, error) {
	var user models.User

	res := r.DB.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).First(&user)

	if res.Error == nil {
		return "", errors.New("username already exists")
	}

	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return "", res.Error
	}

	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return "", err
	}

	res = r.DB.WithContext(ctx).Create(&models.User{
		Username: username,
		Password: hashedPassword,
	})

	if res.Error != nil {
		return "", res.Error
	}

	return r.LoginUser(ctx, username, password)

}

func (r *UserRepository) LoginUser(ctx context.Context, username, password string) (string, error) {
	var existingUser models.User

	res := r.DB.WithContext(ctx).Where("username = ?", username).First(&existingUser)

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
