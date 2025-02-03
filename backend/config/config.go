package config

import "os"

// Config содержит параметры конфигурации приложения.
type Config struct {
	DBHost      string
	DBUser      string
	DBPassword  string
	DBName      string
	KafkaBroker string
	RedisAddr   string
	Port        string
}

// LoadConfig читает параметры конфигурации из переменных окружения.
func LoadConfig() *Config {
	return &Config{
		DBHost:      os.Getenv("DB_HOST"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		Port:        os.Getenv("PORT"),
	}
}