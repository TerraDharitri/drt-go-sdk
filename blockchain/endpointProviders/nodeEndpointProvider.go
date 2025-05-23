package endpointProviders

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-sdk/core"
)

const (
	nodeGetNodeStatusEndpoint      = "node/status"
	nodeRawBlockByHashEndpoint     = "internal/raw/block/by-hash/%s"
	nodeRawBlockByNonceEndpoint    = "internal/raw/block/by-nonce/%d"
	nodeRawMiniBlockByHashEndpoint = "internal/raw/miniblock/by-hash/%s/epoch/%d"
	nodeBlockByNonceEndpoint       = "block/by-nonce/%d"
	nodeGetBlockByHashEndpoint     = "block/by-hash/%s"
)

// nodeEndpointProvider is suitable to work with a Dharitri node (observer)
type nodeEndpointProvider struct {
	*baseEndpointProvider
}

// NewNodeEndpointProvider returns a new instance of a nodeEndpointProvider
func NewNodeEndpointProvider() *nodeEndpointProvider {
	return &nodeEndpointProvider{}
}

// GetNodeStatus returns the node status endpoint
func (node *nodeEndpointProvider) GetNodeStatus(_ uint32) string {
	return nodeGetNodeStatusEndpoint
}

// ShouldCheckShardIDForNodeStatus returns true as some extra check will need to be done when requesting from an observer
func (node *nodeEndpointProvider) ShouldCheckShardIDForNodeStatus() bool {
	return true
}

// GetRawBlockByHash returns the raw block by hash endpoint
func (node *nodeEndpointProvider) GetRawBlockByHash(_ uint32, hexHash string) string {
	return fmt.Sprintf(nodeRawBlockByHashEndpoint, hexHash)
}

// GetRawBlockByNonce returns the raw block by nonce endpoint
func (node *nodeEndpointProvider) GetRawBlockByNonce(_ uint32, nonce uint64) string {
	return fmt.Sprintf(nodeRawBlockByNonceEndpoint, nonce)
}

// GetRawMiniBlockByHash returns the raw miniblock by hash endpoint
func (node *nodeEndpointProvider) GetRawMiniBlockByHash(_ uint32, hexHash string, epoch uint32) string {
	return fmt.Sprintf(nodeRawMiniBlockByHashEndpoint, hexHash, epoch)
}

// GetRestAPIEntityType returns the observer node constant
func (node *nodeEndpointProvider) GetRestAPIEntityType() core.RestAPIEntityType {
	return core.ObserverNode
}

// GetBlockByNonce returns the block with the given nonce
func (node *nodeEndpointProvider) GetBlockByNonce(_ uint32, nonce uint64) string {
	return fmt.Sprintf(nodeBlockByNonceEndpoint, nonce)
}

// GetBlockByHash returns the block with the given hash
func (node *nodeEndpointProvider) GetBlockByHash(_ uint32, hash string) string {
	return fmt.Sprintf(nodeGetBlockByHashEndpoint, hash)
}

// IsInterfaceNil returns true if there is no value under the interface
func (node *nodeEndpointProvider) IsInterfaceNil() bool {
	return node == nil
}
