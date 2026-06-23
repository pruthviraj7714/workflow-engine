package services

import (
	"context"
	"workflow-engine/internal/repository"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, username, password string) (string, error) {
	return s.UserRepo.RegisterUser(ctx, username, password)
}

func (s *UserService) LoginUser(ctx context.Context, username, password string) (string, error) {
	return s.UserRepo.LoginUser(ctx, username, password)
}
