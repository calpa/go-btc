package exchanges

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchKraken(ch chan<- PriceResult) {
	type Ticker struct {
		C []string `json:"c"`
	}
	type Resp struct {
		Error  []string          `json:"error"`
		Result map[string]Ticker `json:"result"`
	}

	const url = "https://api.kraken.com/0/public/Ticker?pair=XBTUSD"

	resp, err := http.Get(url)
	if err != nil {
		ch <- PriceResult{Exchange: "Kraken", Err: err}
		return
	}
	defer resp.Body.Close()

	var data Resp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- PriceResult{Exchange: "Kraken", Err: err}
		return
	}

	if len(data.Error) > 0 {
		ch <- PriceResult{Exchange: "Kraken", Err: fmt.Errorf("api error: %v", data.Error)}
		return
	}

	if len(data.Result) == 0 {
		ch <- PriceResult{Exchange: "Kraken", Err: fmt.Errorf("no ticker data")}
		return
	}

	var first Ticker
	for _, v := range data.Result {
		first = v
		break
	}

	if len(first.C) == 0 {
		ch <- PriceResult{Exchange: "Kraken", Err: fmt.Errorf("no last trade price")}
		return
	}

	var price float64
	if _, err := fmt.Sscanf(first.C[0], "%f", &price); err != nil {
		ch <- PriceResult{Exchange: "Kraken", Err: err}
		return
	}

	ch <- PriceResult{Exchange: "Kraken", Price: price}
}
