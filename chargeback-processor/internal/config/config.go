package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	RabbitMQ   RabbitMQConfig
	Logger     LoggerConfig
	NewRelic   NewRelicConfig
	Chargeback ChargebackConfig
	Minio      MinioConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	CassandraHosts []string
	Keyspace       string
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

type ChargebackConfig struct {
	OutputDir   string
	MaxRecords  int
	MaxDuration time.Duration
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
			CassandraHosts: strings.Split(getEnv("CASSANDRA_HOSTS", "127.0.0.1"), ","),
			Keyspace:       getEnv("CASSANDRA_KEYSPACE", "chargebacks"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672"),
		},
		NewRelic: NewRelicConfig{
			LicenseKey: getEnv("NEW_RELIC_LICENSE_KEY", ""),
			Enabled:    getEnvAsBool("NEW_RELIC_ENABLED", false),
		},
		Chargeback: ChargebackConfig{
			OutputDir: getEnv("CHARGEBACK_OUTPUT_DIR", "/tmp/chargebacks"),
			MaxDuration: getEnvAsDurationWithUnit(
				"CHARGEBACK_MAX_DURATION_VALUE",
				"CHARGEBACK_MAX_DURATION_UNIT",
				30*time.Second,
			),
			MaxRecords: getEnvAsInt("CHARGEBACK_MAX_RECORDS", 5),
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
