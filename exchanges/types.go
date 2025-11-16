package exchanges

// PriceResult is a unified result format for all exchanges.
type PriceResult struct {
	Exchange string
	Price    float64
	Err      error
}
