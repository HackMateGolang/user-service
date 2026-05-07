package service

import (
	"context"
	"fmt"

	"github.com/HackMateGolang/user-service/internal/models"
	"github.com/HackMateGolang/user-service/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (string, error) {
	if _, err := s.repo.ReadUser(ctx, &models.ReadUserRequest{Login: req.Login}); err == nil {
		return "", fmt.Errorf("Service: User already exists: %w", err)
	}

	user := models.User{Login: req.Login, Username: req.Username}
	login, err := s.repo.CreateUser(ctx, &user)
	if err != nil {
		return "", fmt.Errorf("Service: Create user failed: %w", err)
	}

	return login, nil
}
