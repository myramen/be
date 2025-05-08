package config

import (
	"os"
)

type Config struct {
	AppPort       string
	AdminPassword string
	MySQLHost     string
	MySQLPort     string
	MySQLDatabase string
	MySQLUsername string
	MySQLPassword string
}

var AppConfig Config

func Load() {
	AppConfig = Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin1234"),
		MySQLHost:     getEnv("DB_HOST", "localhost"),
		MySQLPort:     getEnv("DB_PORT", "3306"),
		MySQLDatabase: getEnv("DB_DATABASE", "myramen"),
		MySQLUsername: getEnv("DB_USER", "root"),
		MySQLPassword: getEnv("DB_PASSWORD", "password"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
