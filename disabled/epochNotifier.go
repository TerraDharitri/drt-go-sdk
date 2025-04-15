package disabled

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// EpochNotifier is a disabled implementation of EpochNotifier interface
type EpochNotifier struct {
}

// RegisterNotifyHandler does nothing
func (en *EpochNotifier) RegisterNotifyHandler(_ vmcommon.EpochSubscriberHandler) {
}

// CurrentEpoch returns 0
func (en *EpochNotifier) CurrentEpoch() uint32 {
	return 0
}

// CheckEpoch does nothing
func (en *EpochNotifier) CheckEpoch(_ data.HeaderHandler) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (en *EpochNotifier) IsInterfaceNil() bool {
	return en == nil
}
