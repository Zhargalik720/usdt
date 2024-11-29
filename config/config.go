package config

import (
	"flag"
	"os"
)

const (
	AppName = "APP_NAME"
	LogLvl  = "LOG_LEVEL"
	Port    = "PORT"
)

type Config struct {
	AppName    string
	LogLvl     string
	Port       string
	Db         DB
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbDatabase string
}

type DB struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func NewConfig() Config {
	// Флаги командной строки
	flag.StringVar(&dbUser, "db-user", "", "Пользователь базы данных")
	flag.StringVar(&dbPassword, "db-password", "", "Пароль базы данных")
	flag.StringVar(&dbHost, "db-host", "", "Хост базы данных")
	flag.StringVar(&dbPort, "db-port", "", "Порт базы данных")
	flag.StringVar(&dbDatabase, "db-database", "", "Название базы данных")
	flag.Parse()

	return Config{
		AppName:    getEnvOrDefault(AppName, os.Getenv(AppName)),
		LogLvl:     getEnvOrDefault(LogLvl, os.Getenv(LogLvl)),
		Port:       getEnvOrDefault(Port, os.Getenv(Port)),
		dbUser:     dbUser,
		dbPassword: dbPassword,
		dbHost:     dbHost,
		dbPort:     dbPort,
		dbDatabase: dbDatabase,
		Db: DB{
			User:     getEnvOrDefault("DB_USER", dbUser),
			Password: getEnvOrDefault("DB_PASSWORD", dbPassword),
			Host:     getEnvOrDefault("DB_HOST", dbHost),
			Port:     getEnvOrDefault("DB_PORT", dbPort),
			Database: getEnvOrDefault("DB_DATABASE", dbDatabase),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
