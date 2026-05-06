package repository

import (
	"context"
	"fmt"

	"github.com/HackMateGolang/user-service/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
	redisClient *redis.Client
}

func NewUserRepository(db *gorm.DB, redisClient *redis.Client) *userRepository {
	return &userRepository{db: db, redisClient: redisClient}
}

func (r *userRepository) CreateUser (ctx context.Context, user *models.User) (string, error) {
	if err := r.db.Model(&models.User{}).Create(user).Error; err != nil {
		return "", fmt.Errorf("Repo: Create user failed: %w", err)
	}

	return user.Login, r.userCaching(ctx, user)
}

func (r *userRepository) ReadUser (req *models.ReadUserRequest) (string, error) {
	var login string
	return login, nil
}

func (r *userRepository) UpdateUser (req *models.UpdateUserRequest) (string, error) {
	var login string
	return login, nil
}

func (r *userRepository) PatchUser (req *models.PatchUserRequest) (string, error) {
	var login string
	return login, nil
}

func (r *userRepository) DeleteUser (req *models.DeleteUserRequest) (string, error) {
	var login string
	return login, nil
}

func (r *userRepository) userCaching(ctx context.Context, user *models.User) error {
	key := userCacheKey(user.Login)
	if err := r.redisClient.HSet(ctx, key, user).Err(); err != nil {
		return fmt.Errorf("Repo: user caching failed: %w", err)
	}
	return nil
}

func userCacheKey(login string) string {
	return fmt.Sprintf("user:%v", login)
}