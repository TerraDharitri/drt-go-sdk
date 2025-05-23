package endpointProviders

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/stretchr/testify/assert"
)

func TestNewNodeEndpointProvider(t *testing.T) {
	t.Parallel()

	provider := NewNodeEndpointProvider()
	assert.False(t, check.IfNil(provider))
}

func TestNodeEndpointProvider_GetNodeStatus(t *testing.T) {
	t.Parallel()

	provider := NewNodeEndpointProvider()
	assert.Equal(t, nodeGetNodeStatusEndpoint, provider.GetNodeStatus(2))
}

func TestNodeEndpointProvider_Getters(t *testing.T) {
	t.Parallel()

	provider := NewNodeEndpointProvider()
	assert.Equal(t, "internal/raw/block/by-hash/hex", provider.GetRawBlockByHash(2, "hex"))
	assert.Equal(t, "internal/raw/block/by-nonce/3", provider.GetRawBlockByNonce(2, 3))
	assert.Equal(t, "internal/raw/miniblock/by-hash/hex/epoch/4", provider.GetRawMiniBlockByHash(2, "hex", 4))
	assert.Equal(t, "block/by-nonce/5", provider.GetBlockByNonce(2, 5))
	assert.Equal(t, "block/by-hash/hex", provider.GetBlockByHash(2, "hex"))
	assert.Equal(t, core.ObserverNode, provider.GetRestAPIEntityType())
	assert.True(t, provider.ShouldCheckShardIDForNodeStatus())
}
