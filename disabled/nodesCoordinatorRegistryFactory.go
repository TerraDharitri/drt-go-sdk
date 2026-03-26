package disabled

import "github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"

// NodesCoordinatorRegistryFactory is a disabled NodesCoordinatorRegistryFactory
type NodesCoordinatorRegistryFactory struct {
}

// CreateNodesCoordinatorRegistry returns a disabled NodesCoordinatorRegistryHandler
func (factory *NodesCoordinatorRegistryFactory) CreateNodesCoordinatorRegistry(_ []byte) (nodesCoordinator.NodesCoordinatorRegistryHandler, error) {
	return &NodesCoordinatorRegistryHandler{}, nil
}

// GetRegistryData returns an empty buff
func (factory *NodesCoordinatorRegistryFactory) GetRegistryData(_ nodesCoordinator.NodesCoordinatorRegistryHandler, _ uint32) ([]byte, error) {
	return make([]byte, 0), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (factory *NodesCoordinatorRegistryFactory) IsInterfaceNil() bool {
	return factory == nil
}