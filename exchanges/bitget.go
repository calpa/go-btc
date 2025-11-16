package exchanges

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchBitget(ch chan<- PriceResult) {
	type Ticker struct {
		LastPr string `json:"lastPr"`
	}
	type Resp struct {
		Code string   `json:"code"`
		Msg  string   `json:"msg"`
		Data []Ticker `json:"data"`
	}

	const url = "https://api.bitget.com/api/v2/spot/market/tickers?symbol=BTCUSDT"

	resp, err := http.Get(url)
	if err != nil {
		ch <- PriceResult{Exchange: "Bitget", Err: err}
		return
	}
	defer resp.Body.Close()

	var data Resp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- PriceResult{Exchange: "Bitget", Err: err}
		return
	}

	if data.Code != "00000" {
		ch <- PriceResult{Exchange: "Bitget", Err: fmt.Errorf("api error: %s", data.Msg)}
		return
	}

	if len(data.Data) == 0 {
		ch <- PriceResult{Exchange: "Bitget", Err: fmt.Errorf("no ticker data")}
		return
	}

	var price float64
	if _, err := fmt.Sscanf(data.Data[0].LastPr, "%f", &price); err != nil {
		ch <- PriceResult{Exchange: "Bitget", Err: err}
		return
	}

	ch <- PriceResult{Exchange: "Bitget", Price: price}
}
