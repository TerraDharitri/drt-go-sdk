package native

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/ed25519"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-sdk/authentication"
	"github.com/TerraDharitri/drt-go-sdk/blockchain/cryptoProvider"
	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/TerraDharitri/drt-go-sdk/examples"
	"github.com/TerraDharitri/drt-go-sdk/interactors"
	"github.com/TerraDharitri/drt-go-sdk/testsCommon"
	"github.com/TerraDharitri/drt-go-sdk/workflows"
	"github.com/stretchr/testify/require"
)

var keyGen = signing.NewKeyGenerator(ed25519.NewEd25519())
var hrp = "drt"

func TestNativeserver_ClientServer(t *testing.T) {

	t.Run("valid token", func(t *testing.T) {
		t.Parallel()
		lastBlock := &data.HyperBlock{
			Timestamp: uint64(time.Now().Unix()),
			Hash:      "hash",
		}
		proxy := &testsCommon.ProxyStub{
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				return lastBlock, nil
			},
		}

		httpClientWrapper := &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				block := &data.Block{
					Timestamp: int(lastBlock.Timestamp),
				}
				buff, _ := json.Marshal(block)
				return buff, http.StatusOK, nil
			},
		}
		tokenHandler := NewAuthTokenHandler()
		server := createNativeServer(httpClientWrapper, tokenHandler)
		alice := createNativeClient(examples.AlicePemContents, proxy, tokenHandler, "host")

		authToken, _ := alice.GetAccessToken()

		tokenDecoded, err := tokenHandler.Decode(authToken)
		require.Nil(t, err)
		err = server.Validate(tokenDecoded)
		require.Nil(t, err)
	})
}

func createNativeClient(pem string, proxy workflows.ProxyHandler, tokenHandler authentication.AuthTokenHandler, host string) *authClient {
	w := interactors.NewWallet()
	privateKeyBytes, _ := w.LoadPrivateKeyFromPemData([]byte(pem))
	cryptoCompHolder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKeyBytes)

	clientArgs := ArgsNativeAuthClient{
		Signer:                 cryptoProvider.NewSigner(),
		ExtraInfo:              struct{}{},
		Proxy:                  proxy,
		CryptoComponentsHolder: cryptoCompHolder,
		TokenExpiryInSeconds:   60 * 60 * 24,
		TokenHandler:           tokenHandler,
		Host:                   host,
	}

	client, _ := NewNativeAuthClient(clientArgs)
	return client
}

func createNativeServer(httpClientWrapper authentication.HttpClientWrapper, tokenHandler authentication.AuthTokenHandler) *authServer {
	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, hrp)

	serverArgs := ArgsNativeAuthServer{
		HttpClientWrapper: httpClientWrapper,
		TokenHandler:      tokenHandler,
		Signer:            &testsCommon.SignerStub{},
		PubKeyConverter:   converter,
		KeyGenerator:      keyGen,
		TimestampsCacher:  testscommon.NewCacherMock(),
	}
	server, _ := NewNativeAuthServer(serverArgs)

	return server
}
