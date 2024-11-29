package garantex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type GarantexDepth struct {
	Timestamp int64 `json:"timestamp"`
	Asks      []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
		Amount string `json:"amount"`
		Type   string `json:"type"`
	} `json:"asks"`
	Bids []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
		Amount string `json:"amount"`
		Type   string `json:"type"`
	} `json:"bids"`
}

type GrantexAPI struct {
	m       map[string]string
	baseURL string
}

func NewGrantexAPI(baseURL string) *GrantexAPI {
	if baseURL == "" {
		baseURL = "https://garantex.org/api/v2/depth"
	}
	m := map[string]string{
		"RUB": "usdtrub",
		"USD": "usdtusd",
		"EUR": "usdteur",
		"KGS": "usdtkgs",
	}

	return &GrantexAPI{
		m:       m,
		baseURL: baseURL,
	}
}

func (g *GrantexAPI) GetRates(market string) (askPrice, bidPrice float64, timestamp time.Time, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	market, ok := g.m[market]
	if !ok {
		return 0, 0, time.Time{}, fmt.Errorf("market not exist")
	}

	url := fmt.Sprintf("%s?market=%s", g.baseURL, market)
	resp, err := client.Get(url)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("не удалось выполнить запрос к API Garantex: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, time.Time{}, fmt.Errorf("неправильный статус ответа от API: %s", resp.Status)
	}

	var depth GarantexDepth
	if err := json.NewDecoder(resp.Body).Decode(&depth); err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if len(depth.Asks) == 0 || len(depth.Bids) == 0 {
		return 0, 0, time.Time{}, fmt.Errorf("не удалось найти данные о ценах в ответе")
	}

	askPrice = parsePrice(depth.Asks[0].Price)
	bidPrice = parsePrice(depth.Bids[0].Price)
	timestamp = time.Unix(depth.Timestamp, 0)

	return askPrice, bidPrice, timestamp, nil
}

func parsePrice(price string) float64 {
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0
	}
	return p
}
