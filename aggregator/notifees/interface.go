package notifees

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
)

// TxBuilder defines the component able to build & sign a transaction
type TxBuilder interface {
	ApplyUserSignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
	IsInterfaceNil() bool
}

// Proxy holds the primitive functions that the dharitri proxy engine supports & implements
// dependency inversion: blockchain package is considered inner business logic, this package is considered "plugin"
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	SendTransactions(ctx context.Context, txs []*transaction.FrontendTransaction) ([]string, error)
	IsInterfaceNil() bool
}

// TransactionNonceHandler defines the component able to apply nonce for a given FrontendTransaction
type TransactionNonceHandler interface {
	ApplyNonceAndGasPrice(ctx context.Context, address core.AddressHandler, tx *transaction.FrontendTransaction) error
	SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	IsInterfaceNil() bool
}
