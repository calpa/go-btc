# go-btc

![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=for-the-badge&logo=go)
![Status](https://img.shields.io/badge/status-experimental-orange?style=for-the-badge)
![License](https://img.shields.io/badge/license-MIT-green?style=for-the-badge)

Simple Go program that fetches the current BTC price from multiple centralized exchanges (Binance, OKX, Coinbase, Bybit, Bitget, Hyperliquid, Kraken), prints them side by side, and shows the best bid/ask and spread.

## Features

- **Live price fetch** from:
  - Binance (`BTCUSDT`)
  - OKX (`BTC-USDT`)
  - Coinbase (`BTC-USD`)
  - Bybit (`BTCUSDT`, spot)
  - Bitget (`BTCUSDT`, spot)
  - Hyperliquid (BTC perp mid via `allMids`)
  - Kraken (`XBTUSD`)
- **Concurrent requests** to all exchanges using goroutines and channels.
- **Unified result type** so all exchanges report back in the same format.
- **Best bid / best ask / spread** calculation printed to the console.

## Requirements

- Go toolchain installed (Go 1.21+ is recommended; `go.mod` currently declares `go 1.24.1`).
- Internet connection (the program calls the public REST APIs of Binance, OKX, Coinbase, Bybit, Bitget, Hyperliquid, and Kraken).

## Getting started

Clone the repo and change into the directory:

```bash
git clone https://github.com/calpaliu/go-btc.git
cd go-btc
```

Build all packages:

```bash
go build ./...
```

Run the program:

```bash
go run .
# or
go run main.go
```

Example output:

```text
Exchange      Price (USD)
---------------------------
Bybit           95306.10
Coinbase        95294.12
Binance         95311.53
OKX             95312.30
Hyperliquid     95312.50
Bitget          95312.54
Kraken          95312.60

Best Bid: 95294.115 ( Coinbase )
Best Ask: 95312.60 ( Kraken )
Spread: 18.485
```

The order of the lines may change between runs because each exchange call is done concurrently.

## Exchanges support matrix

| #  | Exchange               | Symbol / Market   | Supported |
|----|------------------------|-------------------|-----------|
| 1  | Binance                | `BTCUSDT`         | ✅        |
| 2  | Coinbase               | `BTC-USD`         | ✅        |
| 3  | Upbit                  | `BTC/KRW`         | ❓        |
| 4  | Bybit                  | `BTCUSDT` (spot)  | ✅        |
| 5  | OKX                    | `BTC-USDT`        | ✅        |
| 6  | Bitget                 | `BTCUSDT` (spot)  | ✅        |
| 7  | Gate.io                | `BTCUSDT`         | ❓        |
| 8  | KuCoin                 | `BTCUSDT`         | ❓        |
| 9  | MEXC                   | `BTCUSDT`         | ❓        |
| 10 | HTX (Huobi)            | `BTCUSDT`         | ❓        |
| 11 | Crypto.com Exchange    | `BTCUSD`          | ❓        |
| 12 | Bitfinex               | `BTCUSD`          | ❓        |
| 13 | BingX                  | `BTCUSDT`         | ❓        |
| 14 | Kraken                 | `XBTUSD`          | ✅        |
| 15 | Binance TR             | `BTCUSDT`         | ❓        |
| 16 | BitMart                | `BTCUSDT`         | ❓        |
| 17 | LBank                  | `BTCUSDT`         | ❓        |
| 18 | Bitstamp               | `BTCUSD`          | ❓        |
| 19 | Bithumb                | `BTC/KRW`         | ❓        |
| 20 | XT.COM                 | `BTCUSDT`         | ❓        |
| 21 | Binance.US             | `BTCUSD`          | ❓        |
| 22 | Gemini                 | `BTCUSD`          | ❓        |
| 23 | Deepcoin               | `BTCUSDT`         | ❓        |
| 24 | Toobit                 | `BTCUSDT`         | ❓        |
| 25 | Biconomy.com           | `BTCUSDT`         | ❓        |
| 26 | KCEX                   | `BTCUSDT`         | ❓        |
| 27 | CoinW                  | `BTCUSDT`         | ❓        |
| 28 | WEEX                   | `BTCUSDT`         | ❓        |
| 29 | BTCC                   | `BTCUSDT`         | ❓        |
| 30 | DigiFinex              | `BTCUSDT`         | ❓        |
| 31 | Pionex                 | `BTCUSDT`         | ❓        |
| 32 | AscendEX               | `BTCUSDT`         | ❓        |
| 33 | P2B                    | `BTCUSDT`         | ❓        |
| 34 | Bitunix                | `BTCUSDT`         | ❓        |
| 35 | Hyperliquid            | BTC perp mid      | ✅        |

✅ = implemented in this repo (currently 7 CEX + 1 perp DEX-style venue). ❓ = not implemented yet.

## Project structure

```text
.
├── exchanges/
│   ├── binance.go      # FetchBinance: Binance REST API client
│   ├── coinbase.go     # FetchCoinbase: Coinbase REST API client
│   ├── okx.go          # FetchOKX: OKX REST API client
│   ├── bybit.go        # FetchBybit: Bybit REST API client
│   ├── bitget.go       # FetchBitget: Bitget REST API client
│   ├── hyperliquid.go  # FetchHyperliquid: Hyperliquid REST API client
│   ├── kraken.go       # FetchKraken: Kraken REST API client
│   └── types.go        # PriceResult type shared by all exchanges
├── main.go             # Orchestrates concurrent fetches and prints results
└── go.mod              # Go module definition (module "go-btc")
```

### `exchanges` package

All exchange-specific logic lives in the `exchanges` package:

- `PriceResult` (in `types.go`) is the unified struct returned via a channel:
  - `Exchange string`
  - `Price float64`
  - `Err error`
- `FetchBinance`, `FetchOKX`, `FetchCoinbase`, `FetchBybit`, `FetchBitget`, `FetchHyperliquid`, `FetchKraken` each:
  - Call the corresponding public REST endpoint.
  - Parse the JSON response.
  - Convert the price string to `float64`.
  - Send a `PriceResult` on the provided channel.

### `main` package

`main.go` is intentionally small and focuses on orchestration:

- Creates a `chan exchanges.PriceResult`.
- Starts one goroutine per exchange:
  - `go exchanges.FetchBinance(ch)`
  - `go exchanges.FetchOKX(ch)`
  - `go exchanges.FetchCoinbase(ch)`
  - `go exchanges.FetchBybit(ch)`
  - `go exchanges.FetchBitget(ch)`
  - `go exchanges.FetchHyperliquid(ch)`
  - `go exchanges.FetchKraken(ch)`
- Collects results from the channel (currently 7: Binance, OKX, Coinbase, Bybit, Bitget, Hyperliquid, Kraken).
- Prints the per-exchange prices.
- Computes and prints:
  - Best bid (lowest price)
  - Best ask (highest price)
  - Spread = bestAsk - bestBid

## Extending the project

Here are some ideas to extend this into a bigger engineering project:

- Add more exchanges (e.g., Bybit, Kraken, Bitfinex) by creating new `*.go` files in `exchanges/` with new `FetchX` functions.
- Introduce an `Exchange` interface so you can register exchanges dynamically.
- Add configurable quote/base pairs instead of hardcoding BTC/USDT and BTC/USD.
- Add unit tests for the JSON parsing and `PriceResult` handling.
- Implement timeouts and retry logic with `context.Context` and custom HTTP clients.
- Export the data via HTTP or gRPC instead of printing to stdout.

## License

This project is licensed under the **MIT License**.

Copyright (c) 2025 **Calpa Liu** <calpaliu@gmail.com>

See the [LICENSE](./LICENSE) file for full license text.
