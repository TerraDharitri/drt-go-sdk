package disabled

// GenesisNodesSetupHandler is a disabled GenesisNodesSetupHandler
type GenesisNodesSetupHandler struct {
}

// MinShardHysteresisNodes returns 0
func (handler *GenesisNodesSetupHandler) MinShardHysteresisNodes() uint32 {
	return 0
}

// MinMetaHysteresisNodes returns 0
func (handler *GenesisNodesSetupHandler) MinMetaHysteresisNodes() uint32 {
	return 0
}

// IsInterfaceNil returns true if there is no value under the interface
func (handler *GenesisNodesSetupHandler) IsInterfaceNil() bool {
	return handler == nil
}