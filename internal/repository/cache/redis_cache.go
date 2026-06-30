package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/HackMateGolang/user-service/internal/models"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	redisClient *redis.Client
	ttl         time.Duration
}

func NewUserCache(redisClient *redis.Client, ttl time.Duration) *UserCache {
	return &UserCache{redisClient: redisClient, ttl: ttl}
}

func (c *UserCache) AddUser(ctx context.Context, user *models.User) (bool, error) {
	key := userCacheKey(user.Login)

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return false, fmt.Errorf("CACHE: json marshal error: %w", err)
	}

	if err := c.redisClient.Set(ctx, key, string(jsonUser), c.ttl).Err(); err != nil {
		return false, fmt.Errorf("CACHE: user cache error: %w", err)
	}

	return true, nil
}

func (c *UserCache) DelUser(ctx context.Context, login string) (bool, error) {
	key := userCacheKey(login)

	if err := c.redisClient.Del(ctx, key).Err(); err != nil {
		return false, fmt.Errorf("CACHE: Delete user failed: %w", err)
	}

	return true, nil
}

// func (c *UserCache) UpdateUser(ctx context.Context, user *models.User) (bool, error) {
// 	ok, err := c.DelUser(ctx, user.Login)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	if !ok {
// 		log.Println("CACHE: User not cached")
// 	}

// 	ok, err = c.AddUser(ctx, user)
// 	if err != nil {
// 		return false, fmt.Errorf("CACHE: User cache failed: %w", err)
// 	}

// 	return true, nil
// }

func (c *UserCache) GetUser(ctx context.Context, login string) (*models.User, error) {
	key := userCacheKey(login)

	data, err := c.redisClient.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("CACHE: Get from cache error: %w", err)
	}

	if data == "" {
		c.DelUser(ctx, login)
		return nil, fmt.Errorf("CACHE: Cached user is empty")
	}

	var user models.User
	if err := userUnmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func userCacheKey(login string) string {
	return fmt.Sprintf("user:%v", login)
}

func userUnmarshal(jsonUser string, usModel *models.User) error {
	if err := json.Unmarshal([]byte(jsonUser), usModel); err != nil {
		return fmt.Errorf("CACHE: JSON unmarshal failed: %w", err)
	}

	return nil
}
