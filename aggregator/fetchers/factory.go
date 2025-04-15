package fetchers

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-sdk/aggregator"
)

// DharitriXTokensPair defines a base-quote pair of ids used by DharitriX
type DharitriXTokensPair struct {
	Base  string
	Quote string
}

// ArgsPriceFetcher represents the arguments for the NewPriceFetcher function
type ArgsPriceFetcher struct {
	FetcherName        string
	ResponseGetter     aggregator.ResponseGetter
	GraphqlGetter      aggregator.GraphqlGetter
	DharitriXTokensMap map[string]DharitriXTokensPair
	EVMGasConfig       EVMGasPriceFetcherConfig
}

// NewPriceFetcher returns a new price fetcher of the type provided
func NewPriceFetcher(args ArgsPriceFetcher) (aggregator.PriceFetcher, error) {
	if args.ResponseGetter == nil {
		return nil, errNilResponseGetter
	}
	if args.GraphqlGetter == nil {
		return nil, errNilGraphqlGetter
	}
	if args.DharitriXTokensMap == nil && args.FetcherName == DharitriXName {
		return nil, errNilDharitriXTokensMap
	}

	return createFetcher(args)
}

func createFetcher(args ArgsPriceFetcher) (aggregator.PriceFetcher, error) {
	switch args.FetcherName {
	case BinanceName:
		return &binance{
			ResponseGetter: args.ResponseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case BitfinexName:
		return &bitfinex{
			ResponseGetter: args.ResponseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case CryptocomName:
		return &cryptocom{
			ResponseGetter: args.ResponseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case GeminiName:
		return &gemini{
			ResponseGetter: args.ResponseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case HitbtcName:
		return &hitbtc{
			ResponseGetter: args.ResponseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case HuobiName:
		return &huobi{
			ResponseGetter: args.ResponseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case KrakenName:
		return &kraken{
			ResponseGetter: args.ResponseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case OkxName:
		return &okx{
			ResponseGetter: args.ResponseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case DharitriXName:
		return &DharitriX{
			GraphqlGetter:      args.GraphqlGetter,
			baseFetcher:        newBaseFetcher(),
			DharitriXTokensMap: args.DharitriXTokensMap,
		}, nil
	case EVMGasPriceStation:
		return &evmGasPriceFetcher{
			ResponseGetter: args.ResponseGetter,
			config:         args.EVMGasConfig,
			baseFetcher:    newBaseFetcher(),
		}, nil
	}
	return nil, fmt.Errorf("%w, fetcherName %s", errInvalidFetcherName, args.FetcherName)
}
