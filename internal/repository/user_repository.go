package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HackMateGolang/user-service/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserRepository(db *gorm.DB, redisClient *redis.Client) *UserRepository {
	return &UserRepository{db: db, redisClient: redisClient}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	if err := r.db.Model(&models.User{}).Create(user).Error; err != nil {
		return "", fmt.Errorf("Repo: Create user failed: %w", err)
	}

	return user.Login, r.userCaching(ctx, user)
}

func (r *UserRepository) ReadUser(ctx context.Context, req *models.ReadUserRequest) (*models.User, error) {
	if req.Login == "" {
		return nil, fmt.Errorf("Repo: Login is empty")
	}
	key := userCacheKey(req.Login)
	var user models.User
	data, err := r.redisClient.Get(ctx, key).Result()
	if err == nil && data != "" {
		if err := r.userUnmarshal(data, &user); err != nil {
			return nil, err
		}
		return &user, nil
	}
	

	if err := r.db.Preload("Stack").Preload("Contacts").Where("login = ?", req.Login).First(&user).Error; err != nil {
		return nil, fmt.Errorf("Repo: User not found: %w", err)
	}

	return &user, r.userCaching(ctx, &user)
}

func (r *UserRepository) ReplaceUser(ctx context.Context, req *models.UpdateUserRequest) (bool, error) {
	result := r.db.Model(&models.User{}).Where("login = ?", req.Login).Select("*").Updates(req)
	if result.Error != nil {
		return false, fmt.Errorf("Repo: replace user failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("Repo: user not found")
	}

	r.db.Where("user_login = ?", req.Login).Delete(&models.Tech{})
	r.db.Where("user_login = ?", req.Login).Delete(&models.Social{})
	
	if len(req.Stack) > 0 {
		r.db.Create(&req.Stack)
	}
	if len(req.Contacts) > 0 {
		r.db.Create(&req.Contacts)
	}

	var updatedUser models.User
	if err := r.db.Preload("Stack").Preload("Contacts").Where("login = ?", req.Login).First(&updatedUser).Error; err != nil {
		return false, fmt.Errorf("Repo: updated user not found: %w", err)
	}

	return true, r.userCaching(ctx, &updatedUser)
}

func (r *UserRepository) PatchUser(ctx context.Context, req *models.PatchUserRequest) (bool, error) {
	result := r.db.Model(&models.User{}).Where("login = ?", req.Login).Updates(req)
	if result.Error != nil {
		return false, fmt.Errorf("Repo: patch user failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("Repo: user not found")
	}

	r.db.Where("user_login = ?", req.Login).Delete(&models.Tech{})
	r.db.Where("user_login = ?", req.Login).Delete(&models.Social{})
	
	if len(req.Stack) > 0 {
		r.db.Create(&req.Stack)
	}
	if len(req.Contacts) > 0 {
		r.db.Create(&req.Contacts)
	}

	var patchedUser models.User
	if err := r.db.Preload("Stack").Preload("Contacts").Where("login = ?", req.Login).First(&patchedUser).Error; err != nil {
		return false, fmt.Errorf("Repo: patched user not found: %w", err)
	}

	return true, r.userCaching(ctx, &patchedUser)
}

func (r *UserRepository) DeleteUser(ctx context.Context, req *models.DeleteUserRequest) (bool, error) {
	if err := r.db.Where("login = ?", req.Login).Delete(&models.User{}).Error; err != nil {
		return false, fmt.Errorf("Repo: user not found: %w", err)
	}

	r.redisClient.Del(ctx, userCacheKey(req.Login))

	return true, nil
}

func (r *UserRepository) userCaching(ctx context.Context, user *models.User) error {
	key := userCacheKey(user.Login)
	
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Repo: JSON marshall failed: %w", err)
	}

	if err := r.redisClient.Set(ctx, key, string(jsonUser), 1*time.Hour).Err(); err != nil {
		return fmt.Errorf("Repo: user caching failed: %w", err)
	}
	return nil
}

func userCacheKey(login string) string {
	return fmt.Sprintf("user:%v", login)
}

func (r *UserRepository) userUnmarshal(jsonUser string, usModel *models.User) error{
	if err := json.Unmarshal([]byte(jsonUser), usModel); err != nil {
		return fmt.Errorf("Repo: JSON unmarshal failed: %w", err)
	}

	return nil
}