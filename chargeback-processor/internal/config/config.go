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

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "80"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DSN", "host=users.c0xiiyew0bah.us-east-1.rds.amazonaws.com port=5432 user=users password=userswellpass dbname=users sslmode=require timezone=UTC connect_timeout=5"),
		},
		Logger: LoggerConfig{
			URL: getEnv("LOGGER_URL", "http://logger-service/log"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqps://eirxeetr:S6jkk12c9SEoZoUuv7GoVgszBIkrhkbL@possum.lmq.cloudamqp.com/eirxeetr"),
		},
		SecretKeyConfig: SecretKeyConfig{
			SecretKey: []byte(getEnv("SECRET_KEY", "YSLjuEHpQIgYVaqOPo3Xxmq1iEhJ6msAdy0wO4yMWMbuGq8kGpDIeHDx99mW4smiFBPTSHIBE6NnMEBbAC2VJQ==")),
		},
		NewRelic: NewRelicConfig{
			LicenseKey: getEnv("NEW_RELIC_LICENSE_KEY", "a6c3120d4675abbb00eb86403160f723FFFFNRAL"),
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
