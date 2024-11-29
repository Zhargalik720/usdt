package service

import (
	"context"
	"time"
	"usdt/internal/models"
)

type UsdtServicer interface {
	Create(ctx context.Context, rate models.CurrencyRate) error
	Update(ctx context.Context, rate models.CurrencyRate) error
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (models.CurrencyRate, error)
	GetByPair(ctx context.Context, pair string) (models.CurrencyRate, error)
	GetAll(ctx context.Context) ([]models.CurrencyRate, error)
}
type RequestAPI interface {
	GetRates(market string) (askPrice, bidPrice float64, timestamp time.Time, err error)
}
