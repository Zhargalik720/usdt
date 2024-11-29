package storage

import (
	"context"
	"fmt"
	"usdt/internal/models"
)

type UsdtStorage struct {
	adapter UsdtStorager
}

func NewUsdtStorage(adapter UsdtStorager) *UsdtStorage {
	return &UsdtStorage{adapter: adapter}
}

func (u *UsdtStorage) Create(ctx context.Context, rate models.CurrencyRate) error {
	err := u.adapter.CreateCurrencyRate(ctx, rate)
	if err != nil {
		return fmt.Errorf("Storage.Create.не удалось создать запись курса валют: %w", err)
	}
	return nil
}

func (u *UsdtStorage) Update(ctx context.Context, rate models.CurrencyRate) error {
	err := u.adapter.UpdateCurrencyRate(ctx, rate)
	if err != nil {
		return fmt.Errorf("Storage.Update.не удалось обновить запись курса валют: %w", err)
	}
	return nil
}

func (u *UsdtStorage) Delete(ctx context.Context, id int64) error {
	err := u.adapter.DeleteCurrencyRate(ctx, id)
	if err != nil {
		return fmt.Errorf("Storage.Delete.не удалось удалить запись курса валют: %w", err)
	}
	return nil
}

func (u *UsdtStorage) GetById(ctx context.Context, id int64) (models.CurrencyRate, error) {
	rate, err := u.adapter.GetCurrencyRate(ctx, id)
	if err != nil {
		return models.CurrencyRate{}, fmt.Errorf("Storage.GetById.не удалось получить запись по ID: %w", err)
	}
	if rate == nil {
		return models.CurrencyRate{}, nil
	}
	return *rate, nil
}

func (u *UsdtStorage) GetByPair(ctx context.Context, pair string) (models.CurrencyRate, error) {
	rate, err := u.adapter.GetCurrencyRateByPair(ctx, pair)
	if err != nil {
		return models.CurrencyRate{}, fmt.Errorf("Storage.GetByPair.не удалось получить запись по валютной паре: %w", err)
	}
	if rate == nil {
		return models.CurrencyRate{}, nil
	}
	return *rate, nil
}

func (u *UsdtStorage) GetAll(ctx context.Context) ([]models.CurrencyRate, error) {
	rates, err := u.adapter.GetAllCurrencyRates(ctx)
	if err != nil {
		return nil, fmt.Errorf("Storage.GetAll.не удалось получить все записи: %w", err)
	}
	return rates, nil
}
