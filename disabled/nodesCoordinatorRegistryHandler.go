package disabled
import "github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"

// NodesCoordinatorRegistryHandler is a disabled NodesCoordinatorRegistryHandler
type NodesCoordinatorRegistryHandler struct {
}

// GetEpochsConfig returns an empty map
func (handler *NodesCoordinatorRegistryHandler) GetEpochsConfig() map[string]nodesCoordinator.EpochValidatorsHandler {
	return make(map[string]nodesCoordinator.EpochValidatorsHandler)
}

// GetCurrentEpoch returns 0
func (handler *NodesCoordinatorRegistryHandler) GetCurrentEpoch() uint32 {
	return 0
}

// SetCurrentEpoch does nothing
func (handler *NodesCoordinatorRegistryHandler) SetCurrentEpoch(_ uint32) {
}