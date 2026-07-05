package services

import (
	"context"
	"errors"
	"workflow-engine/internal/repository"
	"workflow-engine/internal/utils"

	"github.com/google/uuid"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, username, password string) (uuid.UUID, error) {
	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return uuid.Nil, err
	}

	return s.UserRepo.CreateUser(ctx, username, hashedPassword)
}

func (s *UserService) LoginUser(ctx context.Context, username, password string) (string, error) {

	user, err := s.UserRepo.GetUserByUsername(ctx, username)

	if err != nil {
		return "", err
	}

	isValid := utils.CheckPasswordHash(password, user.Password)

	if !isValid {
		return "", errors.New("Incorrect Password")
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}
