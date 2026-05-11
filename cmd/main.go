package main

import (
	"fmt"
	"log"
	"net"

	userpb "github.com/HackMateGolang/user-service/api/proto/v1"
	"github.com/HackMateGolang/user-service/internal/handlers"
	"github.com/HackMateGolang/user-service/internal/models"
	"github.com/HackMateGolang/user-service/internal/repository"
	"github.com/HackMateGolang/user-service/internal/service"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	cache := initCache()
	defer cache.Close()

	userRepo := repository.NewUserRepository(db, cache)

	userService := service.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(userService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Listen failed: %v", err)
	}

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, userHandler)

	fmt.Println("gRPC server started on :50051 !!!")

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("gRPC serve error: %v", err)
	}
}
//docker run -d --name userService -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=userService -p 5432:5432 postgres:16-alpine
func initDB() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=userService sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("Init db failed: %w", err)
	}

	if err := db.AutoMigrate(&models.Social{}, &models.Tech{} ,&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return db, nil
}

//docker run -d --name redis -p 6379:6379 redis:8.6.2
func initCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	return rdb
}
