package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"usdt/internal/models"
)

// MockDbAdapter - mock для db.DbAdapter
type MockDbAdapter struct {
	mock.Mock
}

func (m *MockDbAdapter) CreateCurrencyRate(ctx context.Context, rate models.CurrencyRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockDbAdapter) UpdateCurrencyRate(ctx context.Context, rate models.CurrencyRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockDbAdapter) DeleteCurrencyRate(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDbAdapter) GetCurrencyRate(ctx context.Context, id int64) (*models.CurrencyRate, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return nil, nil
	}
	return args.Get(0).(*models.CurrencyRate), nil
}

func (m *MockDbAdapter) GetCurrencyRateByPair(ctx context.Context, pair string) (*models.CurrencyRate, error) {
	args := m.Called(ctx, pair)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return nil, nil
	}
	return args.Get(0).(*models.CurrencyRate), nil
}

func (m *MockDbAdapter) GetAllCurrencyRates(ctx context.Context) ([]models.CurrencyRate, error) {
	args := m.Called(ctx)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CurrencyRate), nil
}

func TestUsdtStorage_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAdapter := new(MockDbAdapter)
		storage := NewUsdtStorage(mockAdapter)
		rate := models.CurrencyRate{Pair: "USDT/USD", AskPrice: 10.1, BidPrice: 10.0, Timestamp: time.Now()}
		mockAdapter.On("CreateCurrencyRate", mock.Anything, rate).Return(nil)
		err := storage.Create(context.Background(), rate)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		mockAdapter := new(MockDbAdapter)
		storage := NewUsdtStorage(mockAdapter)
		rate := models.CurrencyRate{Pair: "USDT/USD", AskPrice: 10.1, BidPrice: 10.0, Timestamp: time.Now()}
		mockAdapter.On("CreateCurrencyRate", mock.Anything, rate).Return(errors.New("db error"))
		err := storage.Create(context.Background(), rate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestUsdtStorage_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAdapter := new(MockDbAdapter)
		storage := NewUsdtStorage(mockAdapter)
		rate := models.CurrencyRate{Pair: "USDT/USD", AskPrice: 10.1, BidPrice: 10.0, Timestamp: time.Now()}
		mockAdapter.On("UpdateCurrencyRate", mock.Anything, rate).Return(nil)
		err := storage.Update(context.Background(), rate)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		mockAdapter := new(MockDbAdapter)
		storage := NewUsdtStorage(mockAdapter)
		rate := models.CurrencyRate{Pair: "USDT/USD", AskPrice: 10.1, BidPrice: 10.0, Timestamp: time.Now()}
		mockAdapter.On("UpdateCurrencyRate", mock.Anything, rate).Return(errors.New("db error"))
		err := storage.Update(context.Background(), rate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestUsdtStorage_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAdapter := new(MockDbAdapter)
		storage := NewUsdtStorage(mockAdapter)
		testID := int64(1)
		mockAdapter.On("DeleteCurrencyRate", mock.Anything, testID).Return(nil)
		err := storage.Delete(context.Background(), testID)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		mockAdapter := new(MockDbAdapter)
		storage := NewUsdtStorage(mockAdapter)
		testID := int64(1)
		mockAdapter.On("DeleteCurrencyRate", mock.Anything, testID).Return(errors.New("db error"))
		err := storage.Delete(context.Background(), testID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestUsdtStorage_GetById(t *testing.T) {
	now := time.Now()
	t.Run("Success", func(t *testing.T) {
		mockAdapter := &MockDbAdapter{}
		storage := NewUsdtStorage(mockAdapter)
		testID := int64(1)
		expectedRate := models.CurrencyRate{Pair: "USDT/USD", AskPrice: 10.1, BidPrice: 10.0, Timestamp: now}
		mockAdapter.On("GetCurrencyRate", mock.Anything, testID).Return(&expectedRate, nil)
		rate, err := storage.GetById(context.Background(), testID)
		assert.NoError(t, err)
		assert.Equal(t, expectedRate, rate)
	})
	t.Run("Error", func(t *testing.T) {
		mockAdapter := &MockDbAdapter{}
		storage := NewUsdtStorage(mockAdapter)
		testID := int64(1)
		mockAdapter.On("GetCurrencyRate", mock.Anything, testID).Return(nil, errors.New("db error"))
		_, err := storage.GetById(context.Background(), testID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
	t.Run("NotFound", func(t *testing.T) {
		mockAdapter := &MockDbAdapter{}
		storage := NewUsdtStorage(mockAdapter)
		testID := int64(2)
		mockAdapter.On("GetCurrencyRate", mock.Anything, testID).Return(nil, nil)
		rate, err := storage.GetById(context.Background(), testID)
		assert.NoError(t, err)
		assert.Equal(t, models.CurrencyRate{}, rate)
	})
}

func TestUsdtStorage_GetByPair(t *testing.T) {
	now := time.Now()
	t.Run("Success", func(t *testing.T) {
		mockAdapter := &MockDbAdapter{}
		storage := NewUsdtStorage(mockAdapter)
		testPair := "USD"
		expectedRate := models.CurrencyRate{Pair: "USDT/USD", AskPrice: 10.1, BidPrice: 10.0, Timestamp: now}
		mockAdapter.On("GetCurrencyRateByPair", mock.Anything, testPair).Return(&expectedRate, nil)
		rate, err := storage.GetByPair(context.Background(), testPair)
		assert.NoError(t, err)
		assert.Equal(t, expectedRate, rate)
	})
	t.Run("Error", func(t *testing.T) {
		mockAdapter := &MockDbAdapter{}
		storage := NewUsdtStorage(mockAdapter)
		testPair := "USD"
		mockAdapter.On("GetCurrencyRateByPair", mock.Anything, testPair).Return(nil, errors.New("db error"))
		_, err := storage.GetByPair(context.Background(), testPair)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
	t.Run("NotFound", func(t *testing.T) {
		mockAdapter := &MockDbAdapter{}
		storage := NewUsdtStorage(mockAdapter)
		testPair := "NON_EXISTENT_PAIR"
		mockAdapter.On("GetCurrencyRateByPair", mock.Anything, testPair).Return(nil, nil)
		rate, err := storage.GetByPair(context.Background(), testPair)
		assert.NoError(t, err)
		assert.Equal(t, models.CurrencyRate{}, rate)
	})
}

func TestUsdtStorage_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAdapter := new(MockDbAdapter)
		storage := NewUsdtStorage(mockAdapter)
		expectedRates := []models.CurrencyRate{
			{Pair: "USDT/USD", AskPrice: 10.1, BidPrice: 10.0, Timestamp: time.Now()},
			{Pair: "USDT/EUR", AskPrice: 0.9, BidPrice: 0.8, Timestamp: time.Now()},
		}
		mockAdapter.On("GetAllCurrencyRates", mock.Anything).Return(expectedRates, nil)
		rates, err := storage.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedRates, rates)
	})
	t.Run("Error", func(t *testing.T) {
		mockAdapter := new(MockDbAdapter)
		storage := NewUsdtStorage(mockAdapter)
		mockAdapter.On("GetAllCurrencyRates", mock.Anything).Return(nil, errors.New("db error"))
		_, err := storage.GetAll(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}
