package gormpg

import (
	"context"
	"fmt"

	"github.com/HackMateGolang/user-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	if err := r.db.WithContext(ctx).Model(&models.User{}).Create(user).Error; err != nil {
		return "", fmt.Errorf("Repo: Create user failed: %w", err)
	}

	return user.Login, nil
}

func (r *UserRepository) ReadUser(ctx context.Context, req *models.ReadUserRequest) (*models.User, error) {
	if req.Login == "" {
		return nil, fmt.Errorf("Repo: Login is empty")
	}

	var user models.User

	if err := r.db.WithContext(ctx).Preload("Stack").Preload("Contacts").Where("login = ?", req.Login).First(&user).Error; err != nil {
		return nil, fmt.Errorf("Repo: User not found: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) ReplaceUser(ctx context.Context, req *models.UpdateUserRequest) (bool, error) {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.User{}).Where("login = ?", req.Login).Select("*").Updates(req)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("Repo: user not found")
		}

		if err := tx.Where("user_login = ?", req.Login).Delete(&models.Tech{}).Error; err != nil {
			return err
		}

		if err := tx.Where("user_login = ?", req.Login).Delete(&models.Social{}).Error; err != nil {
			return err
		}

		if len(req.Stack) > 0 {
			if err := tx.Create(&req.Stack).Error; err != nil {
				return err
			}
		}

		if len(req.Contacts) > 0 {
			if err := tx.Create(&req.Contacts).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return false, fmt.Errorf("Repo: replace user failed %w", err)
	}

	return true, nil
}

func (r *UserRepository) PatchUser(ctx context.Context, req *models.PatchUserRequest) (bool, error) {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.User{}).Where("login = ?", req.Login).Updates(req)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("Repo: user not found")
		}

		if err := tx.Where("user_login = ?", req.Login).Delete(&models.Tech{}).Error; err != nil {
			return err
		}

		if err := tx.Where("user_login = ?", req.Login).Delete(&models.Social{}).Error; err != nil {
			return err
		}

		if len(req.Stack) > 0 {
			if err := tx.Create(&req.Stack).Error; err != nil {
				return err
			}
		}

		if len(req.Contacts) > 0 {
			if err := tx.Create(&req.Contacts).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return false, fmt.Errorf("Repo: Patch user failed %w", err)
	}

	return true, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, req *models.DeleteUserRequest) (bool, error) {
	tx := r.db.WithContext(ctx).Where("login = ?", req.Login).Delete(&models.User{})
	if tx.Error != nil {
		return false, fmt.Errorf("Repo: delete user failed %w", tx.Error)
	}

	if tx.RowsAffected == 0 {
		return false, fmt.Errorf("User doesnt exists")
	}

	return true, nil
}
