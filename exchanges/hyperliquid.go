package exchanges

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchHyperliquid(ch chan<- PriceResult) {
	// We assume BTC mid price is under key "BTC" in the allMids response.
	// If Hyperliquid uses a different symbol key, adjust btcKey accordingly.
	const (
		url    = "https://api.hyperliquid.xyz/info"
		btcKey = "BTC"
	)

	body := map[string]string{
		"type": "allMids",
	}

	payload, err := json.Marshal(body)
	if err != nil {
		ch <- PriceResult{Exchange: "Hyperliquid", Err: err}
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		ch <- PriceResult{Exchange: "Hyperliquid", Err: err}
		return
	}
	defer resp.Body.Close()

	// The response is a JSON object of string keys to string prices, e.g. {"BTC":"95448.5", ...}
	var data map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- PriceResult{Exchange: "Hyperliquid", Err: err}
		return
	}

	priceStr, ok := data[btcKey]
	if !ok {
		ch <- PriceResult{Exchange: "Hyperliquid", Err: fmt.Errorf("btc key %q not found in response", btcKey)}
		return
	}

	var price float64
	if _, err := fmt.Sscanf(priceStr, "%f", &price); err != nil {
		ch <- PriceResult{Exchange: "Hyperliquid", Err: err}
		return
	}

	ch <- PriceResult{Exchange: "Hyperliquid", Price: price}
}
