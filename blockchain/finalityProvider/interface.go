package finalityProvider

import (
	"context"

	"github.com/TerraDharitri/drt-go-sdk/data"
)

type proxy interface {
	GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error)
	IsInterfaceNil() bool
}
