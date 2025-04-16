package data

import "github.com/TerraDharitri/drt-go-chain/config"

// EnableEpochsConfigResponse holds the enable epochs config endpoint response
type EnableEpochsConfigResponse struct {
	Data struct {
		Config *EnableEpochsConfig `json:"enableEpochs"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// MaxNodesChangeConfig holds the max nodes change config
type MaxNodesChangeConfig struct {
	EpochEnable            uint32 `json:"drt_epoch_enable"`
	MaxNumNodes            uint32 `json:"drt_max_num_nodes"`
	NodesToShufflePerShard uint32 `json:"drt_nodes_to_shuffle_per_shard"`
}

// EnableEpochsConfig holds the enable epochs configuration parameters
type EnableEpochsConfig struct {
	EnableEpochs config.EnableEpochs
}
