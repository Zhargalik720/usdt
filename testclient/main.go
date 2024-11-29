package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"usdt/internal/proto/usdt_proto"
)

// RateController - структура для контроллера.
type RateController struct {
	grpcClient usdt_proto.AuthServiceClient
}

// NewRateController - функция для создания нового контроллера.
func NewRateController(grpcConn *grpc.ClientConn) *RateController {
	client := usdt_proto.NewAuthServiceClient(grpcConn)
	return &RateController{
		grpcClient: client,
	}
}

// GetRates - обработчик запроса для получения курсов валют.
func (c *RateController) GetRates(targetCurrency string) (*usdt_proto.GetRatesResponse, error) {
	req := &usdt_proto.GetRatesRequest{
		TargetCurrency: targetCurrency,
	}
	resp, err := c.grpcClient.GetRates(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("could not get rates: %v", err)
	}
	return resp, nil
}

// HealthCheck - функция для проверки HealthCheck
func (c *RateController) HealthCheck() (string, error) {
	req := &usdt_proto.HealthCheckRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := c.grpcClient.HealthCheck(ctx, req)
	if err != nil {
		return "", fmt.Errorf("HealthCheck failed: %v", err)
	}
	return resp.Status, nil
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	controller := NewRateController(conn)

	targetCurrency := "EUR"
	rates, err := controller.GetRates(targetCurrency)
	if err != nil {
		log.Fatalf("Error getting rates: %v", err)
	}

	log.Printf("Currency pair: %s", rates.GetRate().Pair)
	log.Printf("Ask Price: %f", rates.GetRate().AskPrice)
	log.Printf("Bid Price: %f", rates.GetRate().BidPrice)
	log.Printf("Timestamp: %s", rates.GetRate().Timestamp)

	status, err := controller.HealthCheck()
	if err != nil {
		log.Fatalf("HealthCheck failed: %v", err)
	}
	log.Printf("HealthCheck status: %s", status)
}
