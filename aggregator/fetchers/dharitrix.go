package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TerraDharitri/drt-go-sdk/aggregator"
)

const (
	// TODO EN-13146: extract this urls constants in a file
	dataApiUrl = "https://tools.dharitri.org/data-api/graphql"
	// TODO: update the query after data-api rebranding
	query = "query DurianPriceUrl($base: String!, $quote: String!) { trading { pair(first_token: $base, second_token: $quote) { price { last time } } } }"
)

type variables struct {
	BasePrice  string `json:"base"`
	QuotePrice string `json:"quote"`
}

type priceResponse struct {
	Last float64   `json:"last"`
	Time time.Time `json:"time"`
}

type graphqlResponse struct {
	Data struct {
		Trading struct {
			Pair struct {
				Price []priceResponse `json:"price"`
			} `json:"pair"`
		} `json:"trading"`
	} `json:"data"`
}

type DharitriX struct {
	aggregator.GraphqlGetter
	baseFetcher
	DharitriXTokensMap map[string]DharitriXTokensPair
}

// FetchPrice will fetch the price using the http client
func (x *DharitriX) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	if !x.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	DharitriXTokensPair, ok := x.fetchDharitriXTokensPair(base, quote)
	if !ok {
		return 0, errInvalidPair
	}

	vars, err := json.Marshal(variables{
		BasePrice:  DharitriXTokensPair.Base,
		QuotePrice: DharitriXTokensPair.Quote,
	})
	if err != nil {
		return 0, err
	}

	resp, err := x.GraphqlGetter.Query(ctx, dataApiUrl, query, string(vars))
	if err != nil {
		return 0, err
	}

	var graphqlResp graphqlResponse
	err = json.Unmarshal(resp, &graphqlResp)
	if err != nil {
		return 0, errInvalidGraphqlResponse
	}

	price := graphqlResp.Data.Trading.Pair.Price[0].Last

	if price <= 0 {
		return 0, errInvalidResponseData
	}
	return price, nil
}

func (x *DharitriX) fetchDharitriXTokensPair(base, quote string) (DharitriXTokensPair, bool) {
	pair := fmt.Sprintf("%s-%s", base, quote)
	mtp, ok := x.DharitriXTokensMap[pair]
	return mtp, ok
}

// Name returns the name
func (x *DharitriX) Name() string {
	return DharitriXName
}

// IsInterfaceNil returns true if there is no value under the interface
func (x *DharitriX) IsInterfaceNil() bool {
	return x == nil
}
