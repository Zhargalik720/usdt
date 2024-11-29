package storage

import (
	"context"
	"usdt/internal/models"
)

type UsdtStorager interface {
	CreateCurrencyRate(ctx context.Context, rate models.CurrencyRate) error
	GetCurrencyRate(ctx context.Context, id int64) (*models.CurrencyRate, error)
	GetCurrencyRateByPair(ctx context.Context, pair string) (*models.CurrencyRate, error)
	GetAllCurrencyRates(ctx context.Context) ([]models.CurrencyRate, error)
	UpdateCurrencyRate(ctx context.Context, rate models.CurrencyRate) error
	DeleteCurrencyRate(ctx context.Context, id int64) error
}
