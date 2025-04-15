package data

// NetworkConfigResponse holds the network config endpoint response
type NetworkConfigResponse struct {
	Data struct {
		Config *NetworkConfig `json:"config"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// NetworkConfig holds the network configuration parameters
type NetworkConfig struct {
	ChainID                  string  `json:"drt_chain_id"`
	Denomination             int     `json:"drt_denomination"`
	GasPerDataByte           uint64  `json:"drt_gas_per_data_byte"`
	LatestTagSoftwareVersion string  `json:"drt_latest_tag_software_version"`
	MetaConsensusGroup       uint32  `json:"drt_meta_consensus_group_size"`
	MinGasLimit              uint64  `json:"drt_min_gas_limit"`
	MinGasPrice              uint64  `json:"drt_min_gas_price"`
	MinTransactionVersion    uint32  `json:"drt_min_transaction_version"`
	NumMetachainNodes        uint32  `json:"drt_num_metachain_nodes"`
	NumNodesInShard          uint32  `json:"drt_num_nodes_in_shard"`
	NumShardsWithoutMeta     uint32  `json:"drt_num_shards_without_meta"`
	RoundDuration            int64   `json:"drt_round_duration"`
	ShardConsensusGroupSize  uint64  `json:"drt_shard_consensus_group_size"`
	StartTime                int64   `json:"drt_start_time"`
	Adaptivity               bool    `json:"drt_adaptivity,string"`
	Hysteresys               float32 `json:"drt_hysteresis,string"`
	RoundsPerEpoch           uint32  `json:"drt_rounds_per_epoch"`
	ExtraGasLimitGuardedTx   uint64  `json:"drt_extra_gas_limit_guarded_tx"`
}
