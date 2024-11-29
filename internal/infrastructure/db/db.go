package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	"usdt/config"
)

func RunMigrations(cfg config.Config, logger *zap.Logger) error {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Database)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("Ошибка при подключении к базе данных: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Ошибка при подключении к базе данных (migrate): %w", err)
	}

	migrationPath := "file:///app/migrations"
	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("Ошибка при создании мигратора: %w", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("Миграции уже выполнены.")
		} else {
			return fmt.Errorf("Ошибка при выполнении миграций: %w", err)
		}
	} else {
		logger.Info("Миграции успешно выполнены!")
	}

	return nil
}
