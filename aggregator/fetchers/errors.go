package fetchers

import "errors"

var (
	errInvalidResponseData     = errors.New("invalid response data")
	errInvalidFetcherName      = errors.New("invalid fetcher name")
	errNilResponseGetter       = errors.New("nil response getter")
	errNilGraphqlGetter        = errors.New("nil graphql getter")
	errNilDharitriXTokensMap   = errors.New("nil DharitriX tokens map")
	errInvalidPair             = errors.New("invalid pair")
	errInvalidGraphqlResponse  = errors.New("invalid graphql response")
	errInvalidGasPriceSelector = errors.New("invalid gas price selector")
)
