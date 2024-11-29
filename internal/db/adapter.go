package db

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"usdt/config"
	"usdt/internal/models"
)

type DbAdapter struct {
	db *gorm.DB
}

func NewDB(cfg config.Config) (*DbAdapter, error) {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Database)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
		},
	)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к базе данных: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("Ошибка проверки соединения с базой данных: %w", err)
	}

	return &DbAdapter{db: db}, nil
}

func (adapter *DbAdapter) Close() error {
	sqlDB, err := adapter.db.DB()
	if err != nil {
		return fmt.Errorf("Ошибка получения SQL DB из GORM: %w", err)
	}
	return sqlDB.Close()
}

func (adapter *DbAdapter) CreateCurrencyRate(ctx context.Context, rate models.CurrencyRate) error {
	result := adapter.db.Create(&rate)
	return result.Error
}

func (adapter *DbAdapter) GetCurrencyRate(ctx context.Context, id int64) (*models.CurrencyRate, error) {
	var rate models.CurrencyRate
	result := adapter.db.Where("id = ?", id).First(&rate)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("Ошибка получения записи по ID: %w", result.Error)
	}
	return &rate, nil
}

func (adapter *DbAdapter) GetCurrencyRateByPair(ctx context.Context, pair string) (*models.CurrencyRate, error) {
	var rate models.CurrencyRate
	result := adapter.db.Where("pair = ?", pair).First(&rate)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("Ошибка получения записи по Pair: %w", result.Error)
	}
	return &rate, nil
}
func (adapter *DbAdapter) GetAllCurrencyRates(ctx context.Context) ([]models.CurrencyRate, error) {
	var rates []models.CurrencyRate
	result := adapter.db.Find(&rates)
	if result.Error != nil {
		return nil, fmt.Errorf("Ошибка получения всех записей: %w", result.Error)
	}
	return rates, nil
}

func (adapter *DbAdapter) UpdateCurrencyRate(ctx context.Context, rate models.CurrencyRate) error {
	result := adapter.db.Save(&rate)
	return result.Error
}

func (adapter *DbAdapter) DeleteCurrencyRate(ctx context.Context, id int64) error {
	result := adapter.db.Delete(&models.CurrencyRate{}, id)
	return result.Error
}
