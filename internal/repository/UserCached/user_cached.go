package UserCached

import (
	"context"
	"log"

	"github.com/HackMateGolang/user-service/internal/models"
	"github.com/HackMateGolang/user-service/internal/repository/cache"
	"github.com/HackMateGolang/user-service/internal/service"
)

type CachedUser struct {
	db    service.UserRepository
	cache *cache.UserCache
}

func NewCachedUser(db service.UserRepository, cache *cache.UserCache) *CachedUser {
	return &CachedUser{db: db, cache: cache}
}

func (c *CachedUser) CreateUser(ctx context.Context, user *models.User) (string, error) {
	login, err := c.db.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	ok, err := c.cache.AddUser(ctx, user)
	if err != nil {
		log.Println(err.Error())
	}

	if ok {
		log.Println("REPO: User cached")
	}

	return login, nil
}

func (c *CachedUser) ReadUser(ctx context.Context, req *models.ReadUserRequest) (*models.User, error) {
	user, err := c.cache.GetUser(ctx, req.Login)
	if err != nil {
		log.Println(err.Error())
	} else if user != nil {
		return user, nil
	}

	log.Println("CACHE: User not found")

	user, err = c.db.ReadUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *CachedUser) ReplaceUser(ctx context.Context, req *models.UpdateUserRequest) (bool, error) {
	ok, err := c.db.ReplaceUser(ctx, req)
	if err != nil {
		return false, err
	}

	if ok {
		if ok, err := c.cache.DelUser(ctx, req.Login); err != nil {
			log.Println(err.Error())
		} else if ok {
			log.Println("CACHE: Old user removed")
		}
	}

	return ok, nil
}

func (c *CachedUser) PatchUser(ctx context.Context, req *models.PatchUserRequest) (bool, error) {
	ok, err := c.db.PatchUser(ctx, req)
	if err != nil {
		return false, err
	}

	if ok {
		if ok, err := c.cache.DelUser(ctx, req.Login); err != nil {
			log.Println(err.Error())
		} else if ok {
			log.Println("CACHE: Old user removed")
		}
	}

	return ok, nil
}

func (c *CachedUser) DeleteUser(ctx context.Context, req *models.DeleteUserRequest) (bool, error) {
	ok, err := c.db.DeleteUser(ctx, req)
	if err != nil {
		return false, err
	}

	if ok {
		ok, err := c.cache.DelUser(ctx, req.Login)
		if err != nil {
			log.Println(err.Error())
		}
		if !ok {
			log.Println("CACHE: delete user from cache failed")
		}
	}

	return ok, nil
}

// readReq := &models.ReadUserRequest{Login: req.Login}
// user, err := c.db.ReadUser(ctx, readReq)
// if err != nil {
// 	log.Println(err.Error())
// } else {
// 	ok, err := c.cache.UpdateUser(ctx, user)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	log.Printf("CACHE: User cached: %v", ok)
// }
