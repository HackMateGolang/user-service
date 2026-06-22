package config

import "os"

type DBConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
}

type RedisConfig struct {
	Host string
	Port string
	Password string
	DB string
}

type Config struct {
	Port string
	DB DBConfig
	Cache RedisConfig
}

func New() *Config {
	return &Config{
		Port: getEnv("SERVICE_PORT", ""),
		DB: DBConfig{
			Host: getEnv("POSTGRES_HOST", ""),
			Port: getEnv("POSTGRES_PORT", ""),
			User: getEnv("POSTGRES_USER", ""),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			DBName: getEnv("POSTGRES_DB", ""),
		},
		Cache: RedisConfig{
			Host: getEnv("REDIS_HOST", ""),
			Port: getEnv("REDIS_PORT", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB: getEnv("REDIS_DB", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}


//docker run -d --name userService -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=userService -p 5432:5432 postgres:16-alpine
//docker run -d --name redis -p 6379:6379 redis:8.6.2
