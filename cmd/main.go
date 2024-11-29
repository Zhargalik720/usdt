package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"usdt/config"
	"usdt/internal/db"
	migrate "usdt/internal/infrastructure/db"
	"usdt/internal/infrastructure/logger"
	"usdt/run"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	conf := config.NewConfig()
	logger, file := logger.NewLogger(conf)
	defer file.Close()
	defer logger.Sync()

	grpcServer := grpc.NewServer()
	err = migrate.RunMigrations(conf, logger)
	if err != nil {
		logger.Error(err.Error())
	}
	adapter, err := db.NewDB(conf)
	run.Run(adapter, logger, conf, grpcServer)
}
