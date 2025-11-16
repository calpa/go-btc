package main

import (
	"fmt"

	"go-btc/exchanges"
)

func main() {
	ch := make(chan exchanges.PriceResult)

	go exchanges.FetchBinance(ch)
	go exchanges.FetchOKX(ch)
	go exchanges.FetchCoinbase(ch)
	go exchanges.FetchBybit(ch)
	go exchanges.FetchBitget(ch)

	results := make([]exchanges.PriceResult, 0, 5)
	for range 5 {
		results = append(results, <-ch)
	}

	fmt.Println("Exchange   Price (USD)")
	fmt.Println("---------------------------")

	bestBid := 1e12 // large starting value
	bestAsk := 0.0
	var bidEx, askEx string

	for _, r := range results {
		if r.Err != nil {
			fmt.Println(r.Exchange, "Error:", r.Err)
			continue
		}

		fmt.Printf("%-10s %.2f\n", r.Exchange, r.Price)

		if r.Price < bestBid {
			bestBid = r.Price
			bidEx = r.Exchange
		}
		if r.Price > bestAsk {
			bestAsk = r.Price
			askEx = r.Exchange
		}
	}

	fmt.Println("\nBest Bid:", bestBid, "(", bidEx, ")")
	fmt.Println("Best Ask:", bestAsk, "(", askEx, ")")
	fmt.Println("Spread:", bestAsk-bestBid)
}
