package data

// NetworkStatusResponse holds the network status response (for a specified shard)
type NetworkStatusResponse struct {
	Data struct {
		Status *NetworkStatus `json:"status"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// NetworkStatus holds the network status details of a specified shard
type NetworkStatus struct {
	CurrentRound               uint64 `json:"drt_current_round"`
	EpochNumber                uint64 `json:"drt_epoch_number"`
	Nonce                      uint64 `json:"drt_nonce"`
	NonceAtEpochStart          uint64 `json:"drt_nonce_at_epoch_start"`
	NoncesPassedInCurrentEpoch uint64 `json:"drt_nonces_passed_in_current_epoch"`
	RoundAtEpochStart          uint64 `json:"drt_round_at_epoch_start"`
	RoundsPassedInCurrentEpoch uint64 `json:"drt_rounds_passed_in_current_epoch"`
	RoundsPerEpoch             uint64 `json:"drt_rounds_per_epoch"`
	CrossCheckBlockHeight      string `json:"drt_cross_check_block_height"`
	HighestNonce               uint64 `json:"drt_highest_final_nonce"`
	ProbableHighestNonce       uint64 `json:"drt_probable_highest_nonce"`
	ShardID                    uint32 `json:"drt_shard_id"`
}
