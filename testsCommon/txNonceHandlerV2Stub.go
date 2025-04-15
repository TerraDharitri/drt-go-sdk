package testsCommon

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/core"
)

// TxNonceHandlerV2Stub -
type TxNonceHandlerV2Stub struct {
	ApplyNonceAndGasPriceCalled func(ctx context.Context, address core.AddressHandler, tx *transaction.FrontendTransaction) error
	SendTransactionCalled       func(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	ForceNonceReFetchCalled     func(address core.AddressHandler) error
	CloseCalled                 func() error
}

// ApplyNonceAndGasPrice -
func (stub *TxNonceHandlerV2Stub) ApplyNonceAndGasPrice(ctx context.Context, address core.AddressHandler, tx *transaction.FrontendTransaction) error {
	if stub.ApplyNonceAndGasPriceCalled != nil {
		return stub.ApplyNonceAndGasPriceCalled(ctx, address, tx)
	}

	return nil
}

// SendTransaction -
func (stub *TxNonceHandlerV2Stub) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	if stub.SendTransactionCalled != nil {
		return stub.SendTransactionCalled(ctx, tx)
	}

	return "", nil
}

// Close -
func (stub *TxNonceHandlerV2Stub) Close() error {
	if stub.CloseCalled != nil {
		return stub.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (stub *TxNonceHandlerV2Stub) IsInterfaceNil() bool {
	return stub == nil
}
