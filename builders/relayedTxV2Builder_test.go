package builders

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/stretchr/testify/require"
)

func TestRelayedTxV2Builder(t *testing.T) {
	t.Parallel()

	netConfig := &data.NetworkConfig{
		ChainID:               "T",
		MinTransactionVersion: 1,
		GasPerDataByte:        1500,
	}

	relayerAcc, relayerPrivKey := getAccount(t, testRelayerMnemonic)
	innerSenderAcc, innerSenderPrivKey := getAccount(t, testInnerSenderMnemonic)

	innerTx := &transaction.FrontendTransaction{
		Nonce:    innerSenderAcc.Nonce,
		Value:    "0",
		Receiver: "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls6prdez",
		Sender:   innerSenderAcc.Address,
		GasPrice: netConfig.MinGasPrice,
		GasLimit: 0,
		Data:     []byte("getContractConfig"),
		ChainID:  netConfig.ChainID,
		Version:  netConfig.MinTransactionVersion,
		Options:  0,
	}

	innerTxSig := signTx(t, innerSenderPrivKey, innerTx)
	innerTx.Signature = hex.EncodeToString(innerTxSig)

	txJson, _ := json.Marshal(innerTx)
	require.Equal(t,
		`{"nonce":37,"value":"0","receiver":"drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls6prdez","sender":"drt1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsy839fr","gasPrice":0,"gasLimit":0,"data":"Z2V0Q29udHJhY3RDb25maWc=","signature":"8bb517cd7dc360bd64b3901a57d174cf3d9fe054250b1ab94251b24d833354e88918b4b2bee21b195cfe73654c684f278b85b55738210b3f05da268a2b697806","chainID":"T","version":1}`,
		string(txJson),
	)

	relayedV2Builder := NewRelayedTxV2Builder()
	relayedV2Builder.SetInnerTransaction(innerTx)
	relayedV2Builder.SetRelayerAccount(relayerAcc)
	relayedV2Builder.SetNetworkConfig(netConfig)
	relayedV2Builder.SetGasLimitNeededForInnerTransaction(60_000_000)

	relayedTx, err := relayedV2Builder.Build()
	require.NoError(t, err)

	relayedTx.GasPrice = netConfig.MinGasPrice

	relayedTxSig := signTx(t, relayerPrivKey, relayedTx)
	relayedTx.Signature = hex.EncodeToString(relayedTxSig)

	txJson, _ = json.Marshal(relayedTx)
	require.Equal(t,
		`{"nonce":37,"value":"0","receiver":"drt1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsy839fr","sender":"drt1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqspzc322","gasPrice":0,"gasLimit":60364500,"data":"cmVsYXllZFR4VjJAMDAwMDAwMDAwMDAwMDAwMDAwMDEwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAyZmZmZkAyNUA2NzY1NzQ0MzZmNmU3NDcyNjE2Mzc0NDM2ZjZlNjY2OTY3QDhiYjUxN2NkN2RjMzYwYmQ2NGIzOTAxYTU3ZDE3NGNmM2Q5ZmUwNTQyNTBiMWFiOTQyNTFiMjRkODMzMzU0ZTg4OTE4YjRiMmJlZTIxYjE5NWNmZTczNjU0YzY4NGYyNzhiODViNTU3MzgyMTBiM2YwNWRhMjY4YTJiNjk3ODA2","signature":"e1866e4953f075642ba4ea82578e0f5569146e2a9f5ca0f2655ef1f7c2b0e36834a04d642967340225200eb69bc638aed8fd86b7e666270765c4387786f87205","chainID":"T","version":1}`,
		string(txJson),
	)
}
