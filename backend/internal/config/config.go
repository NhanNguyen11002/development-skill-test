package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type JWTConfig struct {
	SecretKey string
	Duration  int // in hours
}

const (
	// Server defaults
	DefaultServerPort = "8080"
	DefaultServerHost = "0.0.0.0"

	// Database defaults
	DefaultDBHost     = "localhost"
	DefaultDBPort     = "5432"
	DefaultDBUser     = "postgres"
	DefaultDBPassword = "password"
	DefaultDBName     = "smart_city_surveillance"
	DefaultDBSSLMode  = "disable"

	// Redis defaults
	DefaultRedisHost     = "localhost"
	DefaultRedisPort     = "6379"
	DefaultRedisPassword = ""
	DefaultRedisDB       = 0

	// Kafka defaults
	DefaultKafkaBrokers = "localhost:9092"
	DefaultKafkaTopic   = "surveillance-alerts"

	// JWT defaults
	DefaultJWTSecretKey     = "your-secret-key"
	DefaultJWTDurationHours = 24
)

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Continue without .env file
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", DefaultServerPort),
			Host: getEnv("SERVER_HOST", DefaultServerHost),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", DefaultDBHost),
			Port:     getEnv("DB_PORT", DefaultDBPort),
			User:     getEnv("POSTGRES_USER", DefaultDBUser),
			Password: getEnv("POSTGRES_PASSWORD", DefaultDBPassword),
			DBName:   getEnv("POSTGRES_DB", DefaultDBName),
			SSLMode:  getEnv("DB_SSLMODE", DefaultDBSSLMode),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", DefaultRedisHost),
			Port:     getEnv("REDIS_PORT", DefaultRedisPort),
			Password: getEnv("REDIS_PASSWORD", DefaultRedisPassword),
			DB:       getEnvAsInt("REDIS_DB", DefaultRedisDB),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", DefaultKafkaBrokers)},
			Topic:   getEnv("KAFKA_TOPIC", DefaultKafkaTopic),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", DefaultJWTSecretKey),
			Duration:  getEnvAsInt("JWT_DURATION_HOURS", DefaultJWTDurationHours),
		},
	}

	return config, nil
}


func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
} 