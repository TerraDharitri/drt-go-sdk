package main

import (
	"context"
	"time"

	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-sdk/blockchain"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/TerraDharitri/drt-go-sdk/examples"
)

var log = logger.GetOrCreate("drt-go-sdk/examples/examplesVMQuery")

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

	vmRequest := &data.VmValueRequest{
		Address:    "drt1qqqqqqqqqqqqqpgqp699jngundfqw07d8jzkepucvpzush6k3wvqeyzkqc",
		FuncName:   "version",
		CallerAddr: "drt1rh5ws22jxm9pe7dtvhfy6j3uttuupkepferdwtmslms5fydtrh5smd3qya",
		CallValue:  "",
		Args:       nil,
	}
	response, err := ep.ExecuteVMQuery(context.Background(), vmRequest)
	if err != nil {
		log.Error("error executing vm query", "error", err)
		return
	}

	contractVersion := string(response.Data.ReturnData[0])
	log.Info("response", "contract version", contractVersion)
}
