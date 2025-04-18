package main

import (
	"context"
	"time"

	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-sdk/blockchain"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/examples"
	"github.com/TerraDharitri/drt-go-sdk/headerCheck"
)

var log = logger.GetOrCreate("drt-go-sdk/examples/examplesHeaderCheck")

func main() {
	args := blockchain.ArgsProxy{
		ProxyURL:            examples.TestnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	ep, err := blockchain.NewProxy(args)
	if err != nil {
		log.Error("error creating proxy", "error", err)
		return
	}

	enableEpochsConfig, err := ep.GetEnableEpochsConfig(context.Background())
	if err != nil {
		log.Error("error getting enable epochs config from proxy", "error", err)
		return
	}

	// set enable epochs config based on the environment
	if len(enableEpochsConfig.EnableEpochs.BLSMultiSignerEnableEpoch) == 0 {
		enableEpochsConfig.EnableEpochs.BLSMultiSignerEnableEpoch = []config.MultiSignerConfig{
			{
				EnableEpoch: 0,
				Type:        "no-KOSK",
			},
			{
				EnableEpoch: uint32(1000000),
				Type:        "KOSK",
			},
		}
	}

	headerVerifier, err := headerCheck.NewHeaderCheckHandler(ep, enableEpochsConfig)
	if err != nil {
		log.Error("error creating header check handler", "error", err)
		return
	}

	// set header headerHash and shard ID
	headerHash := "e0b29ef07f76b75ea9608eed37c588440113724077f57cda3bac84ea0de378ab"
	shardID := uint32(2)

	ok, err := headerVerifier.VerifyHeaderSignatureByHash(context.Background(), shardID, headerHash)
	if err != nil {
		log.Error("error verifying header signature", "error", err)
		return
	}
	if !ok {
		log.Info("header signature does not match")
		return
	}

	log.Info("header signature matches")
}
