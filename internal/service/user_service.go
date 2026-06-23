package service

import (
	"context"
	"fmt"

	"github.com/HackMateGolang/user-service/internal/models"
)

type UserService struct {
	repo UserRepository
}

type UserRepository interface {
	CreateUser(context.Context, *models.User) (string, error)
	ReadUser(context.Context, *models.ReadUserRequest) (*models.User, error)
	ReplaceUser(context.Context, *models.UpdateUserRequest) (bool, error)
	PatchUser(context.Context, *models.PatchUserRequest) (bool, error)
	DeleteUser(context.Context, *models.DeleteUserRequest) (bool, error)
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (string, error) {
	if _, err := s.repo.ReadUser(ctx, &models.ReadUserRequest{Login: req.Login}); err == nil {
		return "", fmt.Errorf("Service: User already exists")
	}

	user := models.User{Login: req.Login, Username: req.Username}
	login, err := s.repo.CreateUser(ctx, &user)
	if err != nil {
		return "", fmt.Errorf("Service: Create user failed: %w", err)
	}

	return login, nil
}

func (s *UserService) ReadUser(ctx context.Context, req *models.ReadUserRequest) (*models.User, error) {
	return s.repo.ReadUser(ctx, req)
}

func (s *UserService) ReplaceUser(ctx context.Context, req *models.UpdateUserRequest) (bool, error) {
	return s.repo.ReplaceUser(ctx, req)
}

func (s *UserService) PatchUser(ctx context.Context, req *models.PatchUserRequest) (bool, error) {
	return s.repo.PatchUser(ctx, req)
}

func (s *UserService) DeleteUser(ctx context.Context, req *models.DeleteUserRequest) (bool, error) {
	return s.repo.DeleteUser(ctx, req)
}
