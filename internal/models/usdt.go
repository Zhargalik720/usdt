package models

import "time"

type CurrencyRate struct {
	Pair      string    `json:"pair"`
	AskPrice  float64   `json:"ask_price"`
	BidPrice  float64   `json:"bid_price"`
	Timestamp time.Time `json:"timestamp"`
}
