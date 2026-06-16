package services

import (
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

func (s *UserService) RegisterUser(username, password string) (string, error) {
	return s.UserRepo.RegisterUser(username, password)
}

func (s *UserService) LoginUser(username, password string) (string, error) {
	return s.UserRepo.LoginUser(username, password)
}
