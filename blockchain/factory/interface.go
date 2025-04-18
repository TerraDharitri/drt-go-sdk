package factory

import (
	"context"

	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
)

type proxy interface {
	GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error)
	GetRestAPIEntityType() core.RestAPIEntityType
	IsInterfaceNil() bool
}

// EndpointProvider is able to return endpoint routes strings
type EndpointProvider interface {
	GetNetworkConfig() string
	GetNetworkEconomics() string
	GetRatingsConfig() string
	GetEnableEpochsConfig() string
	GetAccount(addressAsBech32 string) string
	GetCostTransaction() string
	GetSendTransaction() string
	GetSendMultipleTransactions() string
	GetTransactionStatus(hexHash string) string
	GetTransactionInfo(hexHash string) string
	GetHyperBlockByNonce(nonce uint64) string
	GetHyperBlockByHash(hexHash string) string
	GetVmValues() string
	GetGenesisNodesConfig() string
	GetRawStartOfEpochMetaBlock(epoch uint32) string
	GetNodeStatus(shardID uint32) string
	ShouldCheckShardIDForNodeStatus() bool
	GetRawBlockByHash(shardID uint32, hexHash string) string
	GetRawBlockByNonce(shardID uint32, nonce uint64) string
	GetRawMiniBlockByHash(shardID uint32, hexHash string, epoch uint32) string
	GetGuardianData(address string) string
	GetRestAPIEntityType() core.RestAPIEntityType
	GetValidatorsInfo(epoch uint32) string
	GetProcessedTransactionStatus(hexHash string) string
	GetDCDTTokenData(addressAsBech32 string, tokenIdentifier string) string
	GetNFTTokenData(addressAsBech32 string, tokenIdentifier string, nonce uint64) string
	IsDataTrieMigrated(addressAsBech32 string) string
	GetBlockByNonce(shardID uint32, nonce uint64) string
	GetBlockByHash(shardID uint32, hash string) string
	IsInterfaceNil() bool
}

// FinalityProvider is able to check the shard finalization status
type FinalityProvider interface {
	CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error
	IsInterfaceNil() bool
}
