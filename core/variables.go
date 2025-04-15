package core

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
)

// AddressPublicKeyConverter represents the default address public key converter
var AddressPublicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(AddressBytesLen, core.DefaultAddressPrefix)
