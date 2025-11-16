package exchanges

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchMEXC(ch chan<- PriceResult) {
	type Resp struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	const url = "https://api.mexc.com/api/v3/ticker/price?symbol=BTCUSDT"

	resp, err := http.Get(url)
	if err != nil {
		ch <- PriceResult{Exchange: "MEXC", Err: err}
		return
	}
	defer resp.Body.Close()

	var data Resp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- PriceResult{Exchange: "MEXC", Err: err}
		return
	}

	var price float64
	if _, err := fmt.Sscanf(data.Price, "%f", &price); err != nil {
		ch <- PriceResult{Exchange: "MEXC", Err: err}
		return
	}

	ch <- PriceResult{Exchange: "MEXC", Price: price}
}
