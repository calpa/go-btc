package exchanges

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchOKX(ch chan<- PriceResult) {
	type Ticker struct {
		Last string `json:"last"`
	}
	type Resp struct {
		Data []Ticker `json:"data"`
	}
	const url = "https://www.okx.com/api/v5/market/ticker?instId=BTC-USDT"

	resp, err := http.Get(url)
	if err != nil {
		ch <- PriceResult{Exchange: "OKX", Err: err}
		return
	}
	defer resp.Body.Close()

	var data Resp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- PriceResult{Exchange: "OKX", Err: err}
		return
	}
	if len(data.Data) == 0 {
		ch <- PriceResult{Exchange: "OKX", Err: fmt.Errorf("no ticker data")}
		return
	}

	var price float64
	if _, err := fmt.Sscanf(data.Data[0].Last, "%f", &price); err != nil {
		ch <- PriceResult{Exchange: "OKX", Err: err}
		return
	}

	ch <- PriceResult{Exchange: "OKX", Price: price}
}
