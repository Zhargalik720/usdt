package controller

import (
	"context"
	"usdt/internal/models"
)

type ControllerInterface interface {
	GetRates(ctx context.Context, pair string) (models.CurrencyRate, error)
}
