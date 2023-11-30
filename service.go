package main

import (
	"context"
	"fmt"
	"time"
)

type PriceFetcher interface {
	FetchPrice(context.Context, string) (float64, error)
}

// implements PriceFetcher interface
type priceFetcher struct{}

func (s *priceFetcher) FetchPrice(ctx context.Context, ticker string) (float64, error) {
	return MockPriceFetcher(ctx, ticker)
}

var priceMocks = map[string]float64{
	"BTC": 47_000,
	"ETH": 3000,
	"TB":  3000_000,
}

func MockPriceFetcher(ctx context.Context, ticker string) (float64, error) {
	time.Sleep(100 * time.Millisecond)
	price, ok := priceMocks[ticker]

	if !ok {
		return price, fmt.Errorf("The given ticker (%s) is not supported", ticker)
	}

	return price, nil
}
