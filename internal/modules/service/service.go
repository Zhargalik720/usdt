package service

import (
	"context"
	"fmt"

	"usdt/internal/models"
)

type UsdtService struct {
	storage UsdtServicer
	api     RequestAPI
}

func NewUsdtService(storage UsdtServicer, api RequestAPI) *UsdtService {
	return &UsdtService{
		storage: storage,
		api:     api,
	}
}

func (u *UsdtService) GetRates(ctx context.Context, pair string) (models.CurrencyRate, error) {

	asc, bid, time, err := u.api.GetRates(pair)
	if err != nil {
		return models.CurrencyRate{}, fmt.Errorf("Service.GetRates: %w", err)
	}
	rates := models.CurrencyRate{
		"USDT/" + pair,
		asc,
		bid,
		time,
	}
	err = u.storage.Create(ctx, rates)
	if err != nil {
		return models.CurrencyRate{}, fmt.Errorf("Service.GetRates: %w", err)
	}
	return rates, nil
}
