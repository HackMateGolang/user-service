package repository

import (
	"context"
	"fmt"

	"github.com/HackMateGolang/user-service/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserRepository(db *gorm.DB, redisClient *redis.Client) *userRepository {
	return &userRepository{db: db, redisClient: redisClient}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	if err := r.db.Model(&models.User{}).Create(user).Error; err != nil {
		return "", fmt.Errorf("Repo: Create user failed: %w", err)
	}

	return user.Login, r.userCaching(ctx, user)
}

func (r *userRepository) ReadUser(ctx context.Context, req *models.ReadUserRequest) (*models.User, error) {
	if req.Login == "" {
		return nil, fmt.Errorf("Repo: Login is empty")
	}
	key := userCacheKey(req.Login)
	var user models.User
	err := r.redisClient.HGetAll(ctx, key).Scan(&user)
	if err == nil && user.Login != "" {
		return &user, nil
	}

	if err := r.db.Where("login = ?", req.Login).First(&user).Error; err != nil {
		return nil, fmt.Errorf("Repo: User not found: %w", err)
	}

	return &user, r.userCaching(ctx, &user)
}

func (r *userRepository) ReplaceUser(ctx context.Context, req *models.UpdateUserRequest) (bool, error) {
	result := r.db.Model(&models.User{}).Where("login = ?", req.Login).Select("*").Updates(req)
	if result.Error != nil {
		return false, fmt.Errorf("Repo: replace user failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("Repo: user not found")
	}

	var updatedUser models.User
	if err := r.db.Where("login = ?", req.Login).First(&updatedUser).Error; err != nil {
		return false, fmt.Errorf("Repo: updated user not found: %w", err)
	}

	return true, r.userCaching(ctx, &updatedUser)
}

func (r *userRepository) PatchUser(ctx context.Context, req *models.PatchUserRequest) (bool, error) {
	result := r.db.Model(&models.User{}).Where("login = ?", req.Login).Updates(req)
	if result.Error != nil {
		return false, fmt.Errorf("Repo: patch user failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("Repo: user not found")
	}

	var patchedUser models.User
	if err := r.db.Where("login = ?", req.Login).First(&patchedUser).Error; err != nil {
		return false, fmt.Errorf("Repo: patched user not found: %w", err)
	}

	return true, r.userCaching(ctx, &patchedUser)
}

func (r *userRepository) DeleteUser(req *models.DeleteUserRequest) (string, error) {
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
