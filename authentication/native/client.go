package native

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-sdk/authentication"
	"github.com/TerraDharitri/drt-go-sdk/builders"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/workflows"
)

// ArgsNativeAuthClient is the DTO used in the native auth client constructor
type ArgsNativeAuthClient struct {
	Signer                 builders.Signer
	ExtraInfo              struct{}
	Proxy                  workflows.ProxyHandler
	CryptoComponentsHolder core.CryptoComponentsHolder
	TokenHandler           authentication.AuthTokenHandler
	TokenExpiryInSeconds   int64
	Host                   string
}

type authClient struct {
	signer                 builders.Signer
	extraInfo              []byte
	proxy                  workflows.ProxyHandler
	cryptoComponentsHolder core.CryptoComponentsHolder
	tokenExpiryInSeconds   int64
	host                   []byte
	token                  string
	tokenHandler           authentication.AuthTokenHandler
	tokenExpire            time.Time
	getTimeHandler         func() time.Time
}

// NewNativeAuthClient will create a new native client able to create authentication tokens
func NewNativeAuthClient(args ArgsNativeAuthClient) (*authClient, error) {
	if check.IfNil(args.Signer) {
		return nil, authentication.ErrNilSigner
	}

	extraInfoBytes, err := json.Marshal(args.ExtraInfo)
	if err != nil {
		return nil, fmt.Errorf("%w while marshaling args.extraInfo", err)
	}

	if check.IfNil(args.Proxy) {
		return nil, workflows.ErrNilProxy
	}

	if check.IfNil(args.TokenHandler) {
		return nil, authentication.ErrNilTokenHandler
	}

	if check.IfNil(args.CryptoComponentsHolder) {
		return nil, authentication.ErrNilCryptoComponentsHolder
	}

	return &authClient{
		signer:                 args.Signer,
		extraInfo:              extraInfoBytes,
		proxy:                  args.Proxy,
		cryptoComponentsHolder: args.CryptoComponentsHolder,
		host:                   []byte(args.Host),
		tokenHandler:           args.TokenHandler,
		tokenExpiryInSeconds:   args.TokenExpiryInSeconds,
		getTimeHandler:         time.Now,
	}, nil
}

// GetAccessToken returns an access token used for authentication into different Dharitri services
func (nac *authClient) GetAccessToken() (string, error) {
	now := nac.getTimeHandler()
	noToken := nac.tokenExpire.IsZero()
	tokenExpired := now.After(nac.tokenExpire)
	if noToken || tokenExpired {
		err := nac.createNewToken()
		if err != nil {
			return "", err
		}
	}
	return nac.token, nil
}

func (nac *authClient) createNewToken() error {
	nonce, err := nac.proxy.GetLatestHyperBlockNonce(context.Background())
	if err != nil {
		return err
	}

	lastHyperblock, err := nac.proxy.GetHyperBlockByNonce(context.Background(), nonce)
	if err != nil {
		return err
	}

	token := &AuthToken{
		ttl:       nac.tokenExpiryInSeconds,
		host:      nac.host,
		extraInfo: nac.extraInfo,
		blockHash: lastHyperblock.Hash,
		address:   []byte(nac.cryptoComponentsHolder.GetBech32()),
	}

	unsignedToken := nac.tokenHandler.GetUnsignedToken(token)
	signableMessage := nac.tokenHandler.GetSignableMessage(token.GetAddress(), unsignedToken)
	token.signature, err = nac.signer.SignMessage(signableMessage, nac.cryptoComponentsHolder.GetPrivateKey())
	if err != nil {
		return err
	}

	nac.token, err = nac.tokenHandler.Encode(token)
	if err != nil {
		return err
	}
	nac.tokenExpire = nac.getTimeHandler().Add(time.Duration(nac.tokenExpiryInSeconds))
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (nac *authClient) IsInterfaceNil() bool {
	return nac == nil
}
