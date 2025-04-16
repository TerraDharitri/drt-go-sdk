package testsCommon

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	sdkCore "github.com/TerraDharitri/drt-go-sdk/core"
)

// TxBuilderStub -
type TxBuilderStub struct {
	ApplyUserSignatureCalled func(cryptoHolder sdkCore.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
}

// ApplyUserSignature -
func (stub *TxBuilderStub) ApplyUserSignature(cryptoHolder sdkCore.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
	if stub.ApplyUserSignatureCalled != nil {
		return stub.ApplyUserSignatureCalled(cryptoHolder, tx)
	}

	return nil
}

// IsInterfaceNil -
func (stub *TxBuilderStub) IsInterfaceNil() bool {
	return stub == nil
}
