package nonceHandlerV2

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	sdkCore "github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/interactors"
)

// NewAddressNonceHandlerWithPrivateAccess -
func NewAddressNonceHandlerWithPrivateAccess(proxy interactors.Proxy, address sdkCore.AddressHandler) (*addressNonceHandler, error) {
	if check.IfNil(proxy) {
		return nil, interactors.ErrNilProxy
	}
	if check.IfNil(address) {
		return nil, interactors.ErrNilAddress
	}
	return &addressNonceHandler{
		address:      address,
		proxy:        proxy,
		transactions: make(map[uint64]*transaction.FrontendTransaction),
	}, nil
}
