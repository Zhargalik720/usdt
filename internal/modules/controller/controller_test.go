package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"usdt/internal/models"
	"usdt/internal/proto/usdt_proto"
)

type MockControllerInterface struct {
	mock.Mock
}

func (m *MockControllerInterface) GetRates(ctx context.Context, pair string) (models.CurrencyRate, error) {
	args := m.Called(ctx, pair)
	return args.Get(0).(models.CurrencyRate), args.Error(1)
}

func TestNewController(t *testing.T) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	mockService := new(MockControllerInterface)
	controller := NewController(mockService, logger)
	assert.NotNil(t, controller)

}

func TestUsdtController_GetRates(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testQuery := "RUB"
		timeNow := time.Now()
		expectedResponse := models.CurrencyRate{
			Pair:      "USDT/RUB",
			AskPrice:  1.1,
			BidPrice:  2.2,
			Timestamp: timeNow,
		}

		mockController := new(MockControllerInterface)
		mockController.On("GetRates", context.Background(), testQuery).Return(expectedResponse, nil)

		controller := NewController(mockController, zap.NewNop())
		reqProto := &usdt_proto.GetRatesRequest{
			TargetCurrency: "RUB",
		}

		resp, err := controller.GetRates(context.Background(), reqProto)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Pair, resp.Rate.Pair)
		assert.Equal(t, expectedResponse.AskPrice, resp.Rate.AskPrice)
		assert.Equal(t, expectedResponse.BidPrice, resp.Rate.BidPrice)
		assert.Equal(t, expectedResponse.Timestamp.String(), resp.Rate.Timestamp)

	})

	t.Run("Error", func(t *testing.T) {
		testQuery := "EUR"
		expectedError := fmt.Errorf("что то %w", fmt.Errorf("some error"))

		mockController := new(MockControllerInterface)
		mockController.On("GetRates", context.Background(), testQuery).Return(models.CurrencyRate{}, expectedError)

		controller := NewController(mockController, zap.NewNop())
		reqProto := &usdt_proto.GetRatesRequest{
			TargetCurrency: "EUR",
		}

		_, err := controller.GetRates(context.Background(), reqProto)
		log.Println(err)
		assert.Error(t, err)
		assert.Equal(t, errors.Unwrap(expectedError), err)

	})
}
