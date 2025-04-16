package disabled

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// SimpleDCDTNFTStorageHandler is a disabled implementation of SimpleDCDTNFTStorageHandler interface
type SimpleDCDTNFTStorageHandler struct {
}

// SaveNFTMetaData returns nil
func (sns *SimpleDCDTNFTStorageHandler) SaveNFTMetaData(_ data.TransactionHandler) error {
	return nil
}

// GetDCDTNFTTokenOnDestination returns nil
func (sns *SimpleDCDTNFTStorageHandler) GetDCDTNFTTokenOnDestination(_ vmcommon.UserAccountHandler, _ []byte, _ uint64) (*dcdt.DCDigitalToken, bool, error) {
	return nil, false, nil
}

// SaveNFTMetaDataToSystemAccount returns nil
func (sns *SimpleDCDTNFTStorageHandler) SaveNFTMetaDataToSystemAccount(_ data.TransactionHandler) error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sns *SimpleDCDTNFTStorageHandler) IsInterfaceNil() bool {
	return sns == nil
}
