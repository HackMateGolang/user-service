package repository

import (
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

func (r *userRepository) CreateUser (req *models.CreateUserRequest) (string, error) {
	var login string
	return login, nil
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

