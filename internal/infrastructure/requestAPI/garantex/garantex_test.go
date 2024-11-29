package garantex

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRates(t *testing.T) {
	tests := []struct {
		name           string
		market         string
		mockResponse   string
		mockStatusCode int
		expectErr      bool
		expectAsk      float64
		expectBid      float64
	}{
		{
			name:           "Valid response",
			market:         "RUB",
			mockResponse:   `{"timestamp": 1698405000, "asks": [{"price": "100.5", "volume": "10", "amount": "1005", "type": "ask"}], "bids": [{"price": "99.5", "volume": "10", "amount": "995", "type": "bid"}]}`,
			mockStatusCode: http.StatusOK,
			expectErr:      false,
			expectAsk:      100.5,
			expectBid:      99.5,
		},
		{
			name:           "Market not exist",
			market:         "INVALID",
			mockResponse:   "",
			mockStatusCode: http.StatusOK,
			expectErr:      true,
		},
		{
			name:           "Invalid API response",
			market:         "RUB",
			mockResponse:   `invalid-json`,
			mockStatusCode: http.StatusOK,
			expectErr:      true,
		},
		{
			name:           "No asks or bids",
			market:         "RUB",
			mockResponse:   `{"timestamp": 1698405000, "asks": [], "bids": []}`,
			mockStatusCode: http.StatusOK,
			expectErr:      true,
		},
		{
			name:           "API returns non-200 status",
			market:         "RUB",
			mockResponse:   "",
			mockStatusCode: http.StatusInternalServerError,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			api := NewGrantexAPI(server.URL)
			ask, bid, ts, err := api.GetRates(tt.market)

			if (err != nil) != tt.expectErr {
				t.Fatalf("ожидали ошибку: %v, получили: %v", tt.expectErr, err)
			}

			if !tt.expectErr {
				if ask != tt.expectAsk {
					t.Errorf("ожидали ask: %f, получили: %f", tt.expectAsk, ask)
				}
				if bid != tt.expectBid {
					t.Errorf("ожидали bid: %f, получили: %f", tt.expectBid, bid)
				}
				if ts.IsZero() {
					t.Errorf("ожидали ненулевой timestamp, получили: %v", ts)
				}
			}
		})
	}
}
