package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	RabbitMQ  RabbitMQConfig
	Logger    LoggerConfig
	NewRelic  NewRelicConfig
	Minio     MinioConfig
	FTP       FTPConfig
	Scheduler SchedulerConfig
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

type NewRelicConfig struct {
	LicenseKey string
	Enabled    bool
}

type MinioConfig struct {
	AccessKey  string
	SecretKey  string
	BucketName string
	Endpoint   string
	UseSSL     bool
}

type FTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SchedulerConfig struct {
	Enabled        bool
	Interval       time.Duration
	MaxFilesPerDay int
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
		NewRelic: NewRelicConfig{
			LicenseKey: getEnv("NEW_RELIC_LICENSE_KEY", ""),
			Enabled:    getEnvAsBool("NEW_RELIC_ENABLED", false),
		},
		Minio: MinioConfig{
			Endpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:  getEnv("MINIO_ACCESS_KEY", "admin"),
			SecretKey:  getEnv("MINIO_SECRET_KEY", "password"),
			UseSSL:     false, // ou converter com strconv.ParseBool
			BucketName: getEnv("MINIO_BUCKET_NAME", "chargebacks"),
		},
		FTP: FTPConfig{
			Host:     getEnv("FTP_HOST", "localhost"),
			Port:     getEnvAsInt("FTP_PORT", 21),
			Username: getEnv("FTP_USER", "admin"),
			Password: getEnv("FTP_PASS", "admin"),
		},
		Scheduler: SchedulerConfig{
			Enabled: getEnvAsBool("SCHEDULER_ENABLED", true),
			Interval: getEnvAsDurationWithUnit(
				"SCHEDULER_INTERVAL_VALUE",
				"SCHEDULER_INTERVAL_UNIT",
				1*time.Minute, // default de 1 min
			),
			MaxFilesPerDay: getEnvAsInt("BATCH_MAX_FILES_PER_DAY", 4),
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

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func getEnvAsDurationWithUnit(valueKey string, unitKey string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(valueKey, "")
	unitStr := strings.ToLower(getEnv(unitKey, ""))

	if valueStr == "" || unitStr == "" {
		return defaultVal
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultVal
	}

	switch unitStr {
	case "seconds", "second", "s":
		return time.Duration(value) * time.Second
	case "minutes", "minute", "m":
		return time.Duration(value) * time.Minute
	case "hours", "hour", "h":
		return time.Duration(value) * time.Hour
	default:
		return defaultVal
	}
}
