package controller

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"usdt/internal/proto/usdt_proto"
)

type UsdtInterface interface {
	GetRates(ctx context.Context, req *usdt_proto.GetRatesRequest) (*usdt_proto.GetRatesResponse, error)
	HealthCheck(ctx context.Context, req *usdt_proto.HealthCheckRequest) (*usdt_proto.HealthCheckResponse, error)
	usdt_proto.AuthServiceServer
}

type UsdtController struct {
	service ControllerInterface
	logger  *zap.Logger
	usdt_proto.AuthServiceServer
}

func NewController(service ControllerInterface, logger *zap.Logger) UsdtInterface {
	return &UsdtController{
		service: service,
		logger:  logger,
	}
}

func (s *UsdtController) GetRates(ctx context.Context, req *usdt_proto.GetRatesRequest) (*usdt_proto.GetRatesResponse, error) {
	rate, err := s.service.GetRates(ctx, req.TargetCurrency)
	if err != nil {
		s.logger.Error("Controller.GetRates error:", zap.Error(err))
		return nil, errors.Unwrap(err)
	}
	cr := usdt_proto.CurrencyRate{
		Pair:      rate.Pair,
		AskPrice:  rate.AskPrice,
		BidPrice:  rate.BidPrice,
		Timestamp: rate.Timestamp.String(),
	}
	resp := &usdt_proto.GetRatesResponse{
		Rate: &cr,
	}

	return resp, nil
}
func (s *UsdtController) HealthCheck(ctx context.Context, req *usdt_proto.HealthCheckRequest) (*usdt_proto.HealthCheckResponse, error) {
	return &usdt_proto.HealthCheckResponse{Status: "OK"}, nil
}
