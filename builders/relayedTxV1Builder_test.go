package builders

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/ed25519"
	"github.com/TerraDharitri/drt-go-sdk/blockchain/cryptoProvider"
	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/TerraDharitri/drt-go-sdk/interactors"
	"github.com/stretchr/testify/require"
)

const (
	testRelayerMnemonic     = "bid involve twenty cave offer life hello three walnut travel rare bike edit canyon ice brave theme furnace cotton swing wear bread fine latin"
	testInnerSenderMnemonic = "acid twice post genre topic observe valid viable gesture fortune funny dawn around blood enemy page update reduce decline van bundle zebra rookie real"
)

func TestRelayedTxV1Builder(t *testing.T) {
	t.Parallel()

	netConfig := &data.NetworkConfig{
		ChainID:               "T",
		MinTransactionVersion: 1,
		GasPerDataByte:        1500,
		MinGasLimit:           50000,
		MinGasPrice:           1000000000,
	}

	relayerAcc, relayerPrivKey := getAccount(t, testRelayerMnemonic)
	innerSenderAcc, innerSenderPrivKey := getAccount(t, testInnerSenderMnemonic)

	innerTx := &transaction.FrontendTransaction{
		Nonce:    innerSenderAcc.Nonce,
		Value:    "100000000",
		Receiver: "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf",
		Sender:   innerSenderAcc.Address,
		GasPrice: netConfig.MinGasPrice,
		GasLimit: netConfig.MinGasLimit,
		Data:     nil,
		ChainID:  netConfig.ChainID,
		Version:  netConfig.MinTransactionVersion,
		Options:  0,
	}

	innerTxSig := signTx(t, innerSenderPrivKey, innerTx)
	innerTx.Signature = hex.EncodeToString(innerTxSig)

	txJson, _ := json.Marshal(innerTx)
	require.Equal(t,
		`{"nonce":37,"value":"100000000","receiver":"drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf","sender":"drt1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsy839fr","gasPrice":1000000000,"gasLimit":50000,"signature":"b8d8988e0e3d15b628394fff2e375939e962b0951873492d649f8affa910def384f9601629e802d539acd89f50bd222a2f1abf8d1187f78ddf0ee2388904f507","chainID":"T","version":1}`,
		string(txJson),
	)

	relayedV1Builder := NewRelayedTxV1Builder()
	relayedV1Builder.SetInnerTransaction(innerTx)
	relayedV1Builder.SetRelayerAccount(relayerAcc)
	relayedV1Builder.SetNetworkConfig(netConfig)

	relayedTx, err := relayedV1Builder.Build()
	require.NoError(t, err)

	relayedTxSig := signTx(t, relayerPrivKey, relayedTx)
	relayedTx.Signature = hex.EncodeToString(relayedTxSig)

	txJson, _ = json.Marshal(relayedTx)
	require.Equal(t,
		`{"nonce":37,"value":"100000000","receiver":"drt1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsy839fr","sender":"drt1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqspzc322","gasPrice":1000000000,"gasLimit":1060000,"data":"cmVsYXllZFR4QDdiMjI2ZTZmNmU2MzY1MjIzYTMzMzcyYzIyNzY2MTZjNzU2NTIyM2EzMTMwMzAzMDMwMzAzMDMwMzAyYzIyNzI2NTYzNjU2OTc2NjU3MjIyM2EyMjQxNTQ2YzQ4NGM3NjM5NmY2ODZlNjM2MTZkNDMzODc3NjczOTcwNjQ1MTY4Mzg2Yjc3NzA0NzQyMzU2YTY5NDk0OTZmMzM0OTQ4NGI1OTRlNjE2NTQ1M2QyMjJjMjI3MzY1NmU2NDY1NzIyMjNhMjIzMzJiMmY2NzUyNTQ2ODQxMzU1YTQxMmY0YjM5NTU1YTMzNzA3MzRmMzE3NTZjNjk0ODZjNjY2OTY5MzY0MzMwNzc2MjQ2NTI0NjUxNmI2NDMxNzkzODNkMjIyYzIyNjc2MTczNTA3MjY5NjM2NTIyM2EzMTMwMzAzMDMwMzAzMDMwMzAzMDJjMjI2NzYxNzM0YzY5NmQ2OTc0MjIzYTM1MzAzMDMwMzAyYzIyNjM2ODYxNjk2ZTQ5NDQyMjNhMjI1NjQxM2QzZDIyMmMyMjc2NjU3MjczNjk2ZjZlMjIzYTMxMmMyMjczNjk2NzZlNjE3NDc1NzI2NTIyM2EyMjc1NGU2OTU5NmE2NzM0Mzk0NjYyNTk2ZjRmNTUyZjJmNGM2YTY0NWE0ZjY1NmM2OTczNGE1NTU5NjMzMDZiNzQ1YTRhMmI0YjJmMzY2YjUxMzM3NjRmNDUyYjU3NDE1NzRiNjU2NzQzMzE1NDZkNzMzMjRhMzk1MTc2NTM0OTcxNGM3ODcxMmY2YTUyNDc0ODM5MzQzMzY2NDQ3NTQ5MzQ2OTUxNTQzMTQyNzczZDNkMjI3ZA==","signature":"ab57f56ddb2a4669513d2c5e44a2c14b8c094a00c5662f694e4f4b1a65ebf30813eb448add483e7227c3bb15898ee03445651a82a595ffebb35980837e7b5f05","chainID":"T","version":1}`,
		string(txJson),
	)
}

func getAccount(t *testing.T, mnemonic string) (*data.Account, []byte) {
	wallet := interactors.NewWallet()

	privKey := wallet.GetPrivateKeyFromMnemonic(data.Mnemonic(mnemonic), 0, 0)
	address, err := wallet.GetAddressFromPrivateKey(privKey)
	require.NoError(t, err)

	addressAsBech32String, err := address.AddressAsBech32String()
	require.NoError(t, err)

	account := &data.Account{
		Nonce:   37,
		Address: addressAsBech32String,
	}

	return account, privKey
}

func signTx(t *testing.T, privKeyBytes []byte, tx *transaction.FrontendTransaction) []byte {
	keyGenInstance := signing.NewKeyGenerator(ed25519.NewEd25519())
	privKey, err := keyGenInstance.PrivateKeyFromByteArray(privKeyBytes)
	require.NoError(t, err)
	signer := cryptoProvider.NewSigner()

	signature, err := signer.SignTransaction(tx, privKey)
	require.NoError(t, err)

	return signature
}
