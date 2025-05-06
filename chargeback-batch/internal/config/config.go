package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server          ServerConfig
	Database        DatabaseConfig
	RabbitMQ        RabbitMQConfig
	Logger          LoggerConfig
	SecretKeyConfig SecretKeyConfig
	NewRelic        NewRelicConfig
	Minio           MinioConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	DSN string
}

type LoggerConfig struct {
	URL string
}

type RabbitMQConfig struct {
	URL string
}

type SecretKeyConfig struct {
	SecretKey []byte
}

type NewRelicConfig struct {
	LicenseKey string
}

type MinioConfig struct {
	AccessKey  string
	SecretKey  string
	BucketName string
	Endpoint   string
	UseSSL     bool
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "81"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DSN", "host=localhost port=5432 user=admin password=admin dbname=batch sslmode=disable timezone=UTC connect_timeout=5"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672"),
		},
		SecretKeyConfig: SecretKeyConfig{
			SecretKey: []byte(getEnv("SECRET_KEY", "YSLjuEHpQIgYVaqOPo3Xxmq1iEhJ6msAdy0wO4yMWMbuGq8kGpDIeHDx99mW4smiFBPTSHIBE6NnMEBbAC2VJQ==")),
		},
		NewRelic: NewRelicConfig{
			LicenseKey: getEnv("NEW_RELIC_LICENSE_KEY", ""),
		},
		Minio: MinioConfig{
			Endpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:  getEnv("MINIO_ACCESS_KEY", "admin"),
			SecretKey:  getEnv("MINIO_SECRET_KEY", "password"),
			UseSSL:     false, // ou converter com strconv.ParseBool
			BucketName: getEnv("MINIO_BUCKET_NAME", "chargebacks"),
		},
	}
}

// Funções auxiliares para carregar variáveis de ambiente
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	if valueStr, exists := os.LookupEnv(name); exists {
		value, err := time.ParseDuration(valueStr)
		if err == nil {
			return value
		}
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		value, err := strconv.Atoi(valueStr)
		if err == nil {
			return value
		}
	}
	return defaultVal
}
