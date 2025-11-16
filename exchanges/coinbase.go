package exchanges

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchCoinbase(ch chan<- PriceResult) {
	type AmountResp struct {
		Amount string `json:"amount"`
	}
	type Resp struct {
		Data AmountResp `json:"data"`
	}
	const url = "https://api.coinbase.com/v2/prices/BTC-USD/spot"

	resp, err := http.Get(url)
	if err != nil {
		ch <- PriceResult{Exchange: "Coinbase", Err: err}
		return
	}
	defer resp.Body.Close()

	var data Resp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- PriceResult{Exchange: "Coinbase", Err: err}
		return
	}

	var price float64
	if _, err := fmt.Sscanf(data.Data.Amount, "%f", &price); err != nil {
		ch <- PriceResult{Exchange: "Coinbase", Err: err}
		return
	}

	ch <- PriceResult{Exchange: "Coinbase", Price: price}
}
