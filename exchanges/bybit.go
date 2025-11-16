package exchanges

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchBybit(ch chan<- PriceResult) {
	type Ticker struct {
		LastPrice string `json:"lastPrice"`
	}
	type Result struct {
		List []Ticker `json:"list"`
	}
	type Resp struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  Result `json:"result"`
	}

	const url = "https://api.bybit.com/v5/market/tickers?category=spot&symbol=BTCUSDT"

	resp, err := http.Get(url)
	if err != nil {
		ch <- PriceResult{Exchange: "Bybit", Err: err}
		return
	}
	defer resp.Body.Close()

	var data Resp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- PriceResult{Exchange: "Bybit", Err: err}
		return
	}

	if data.RetCode != 0 {
		ch <- PriceResult{Exchange: "Bybit", Err: fmt.Errorf("api error: %s", data.RetMsg)}
		return
	}

	if len(data.Result.List) == 0 {
		ch <- PriceResult{Exchange: "Bybit", Err: fmt.Errorf("no ticker data")}
		return
	}

	var price float64
	if _, err := fmt.Sscanf(data.Result.List[0].LastPrice, "%f", &price); err != nil {
		ch <- PriceResult{Exchange: "Bybit", Err: err}
		return
	}

	ch <- PriceResult{Exchange: "Bybit", Price: price}
}
