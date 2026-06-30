package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	userpb "github.com/HackMateGolang/proto-contracts/gen/go/user/v1"
	"github.com/HackMateGolang/user-service/config"
	"github.com/HackMateGolang/user-service/internal/handlers"
	"github.com/HackMateGolang/user-service/internal/models"
	pgpgx "github.com/HackMateGolang/user-service/internal/repository/postgres"
	"github.com/HackMateGolang/user-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	conf := config.New()

	db, err := initPgxPool(context.Background(), conf)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	cache := initCache(conf)
	defer cache.Close()

	userRepo := pgpgx.NewUserRepository(db)

	userService := service.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(userService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Listen failed: %v", err)
	}

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, userHandler)

	log.Printf("gRPC server started on :%v !!!\n", conf.Port)

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("gRPC serve error: %v", err)
	}
}

func initDB(conf *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.DBName)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("Init db failed: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Social{}, &models.Tech{}); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return db, nil
}

func initCache(conf *config.Config) *redis.Client {
	addr := fmt.Sprintf("%v:%v", conf.Cache.Host, conf.Cache.Port)
	db, _ := strconv.Atoi(conf.Cache.DB)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Cache.Password,
		DB:       db,
		Protocol: 2,
	})

	return rdb
}

func initPgxPool(ctx context.Context, conf *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.DBName))
	if err != nil {
		return nil, fmt.Errorf("Connpool creating error: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("DB connection error: %w", err)
	}
	return pool, nil
}
