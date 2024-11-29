package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"usdt/internal/models"
)

// MockUsdtStorage - mock для интерфейса UsdtServicer
type MockUsdtStorage struct {
	mock.Mock
}

func (m *MockUsdtStorage) Create(ctx context.Context, rate models.CurrencyRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockUsdtStorage) Update(ctx context.Context, rate models.CurrencyRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockUsdtStorage) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUsdtStorage) GetById(ctx context.Context, id int64) (models.CurrencyRate, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return models.CurrencyRate{}, args.Error(1)
	}
	return args.Get(0).(models.CurrencyRate), args.Error(1)
}

func (m *MockUsdtStorage) GetByPair(ctx context.Context, pair string) (models.CurrencyRate, error) {
	args := m.Called(ctx, pair)
	if args.Error(1) != nil {
		return models.CurrencyRate{}, args.Error(1)
	}
	return args.Get(0).(models.CurrencyRate), args.Error(1)
}

func (m *MockUsdtStorage) GetAll(ctx context.Context) ([]models.CurrencyRate, error) {
	args := m.Called(ctx)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CurrencyRate), args.Error(1)
}

type MockRequestAPI struct {
	mock.Mock
}

func (m *MockRequestAPI) GetRates(market string) (askPrice, bidPrice float64, timestamp time.Time, err error) {
	args := m.Called(market)
	return args.Get(0).(float64), args.Get(1).(float64), args.Get(2).(time.Time), args.Error(3)
}

func TestUsdtService_GetRates(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testMarket := "RUB"
		timeNow := time.Now()
		expectedAsk := 1.1
		expectedBid := 2.2

		mockStorage := new(MockUsdtStorage)
		mockAPI := new(MockRequestAPI)

		mockAPI.On("GetRates", testMarket).Return(expectedAsk, expectedBid, timeNow, nil)
		mockStorage.On("Create", mock.Anything, mock.MatchedBy(func(rate models.CurrencyRate) bool {
			return rate.Pair == "USDT/"+testMarket && rate.AskPrice == expectedAsk && rate.BidPrice == expectedBid
		})).Return(nil)

		service := NewUsdtService(mockStorage, mockAPI)
		rate, err := service.GetRates(context.Background(), testMarket)
		assert.NoError(t, err)
		assert.Equal(t, "USDT/"+testMarket, rate.Pair)
		assert.Equal(t, expectedAsk, rate.AskPrice)
		assert.Equal(t, expectedBid, rate.BidPrice)
		assert.Equal(t, timeNow.Format(time.RFC3339), rate.Timestamp.Format(time.RFC3339)) // Сравнение времени с учетом форматирования

	})

	t.Run("APIError", func(t *testing.T) {
		testMarket := "EUR"
		expectedError := errors.New("API error")

		mockStorage := new(MockUsdtStorage)
		mockAPI := new(MockRequestAPI)

		mockAPI.On("GetRates", testMarket).Return(0.0, 0.0, time.Time{}, expectedError)

		service := NewUsdtService(mockStorage, mockAPI)
		_, err := service.GetRates(context.Background(), testMarket)
		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("Service.GetRates: %w", expectedError), err)

	})

	t.Run("StorageError", func(t *testing.T) {
		testMarket := "GBP"
		expectedAsk := 1.1
		expectedBid := 2.2
		expectedError := errors.New("storage error")

		mockStorage := new(MockUsdtStorage)
		mockAPI := new(MockRequestAPI)

		mockAPI.On("GetRates", testMarket).Return(expectedAsk, expectedBid, time.Now(), nil)
		mockStorage.On("Create", mock.Anything, mock.Anything).Return(expectedError)

		service := NewUsdtService(mockStorage, mockAPI)
		_, err := service.GetRates(context.Background(), testMarket)
		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("Service.GetRates: %w", expectedError), err)
	})
}
