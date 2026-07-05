package repository

import (
	"context"
	"errors"
	"workflow-engine/internal/models"

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

func (r *UserRepository) CreateUser(ctx context.Context, username, password string) (uuid.UUID, error) {
	user := models.User{
		Username: username,
		Password: password,
	}

	err := r.DB.WithContext(ctx).Create(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return uuid.Nil, errors.New("user name already exists")
		} else {
			return uuid.Nil, err
		}
	}

	return user.ID, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	res := r.DB.WithContext(ctx).Where("username = ?", username).First(&user)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user with given username not found")
		}
		return nil, res.Error
	}

	return &user, nil

}
