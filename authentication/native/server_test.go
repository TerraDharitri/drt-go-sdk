package native

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	genesisMock "github.com/TerraDharitri/drt-go-chain/genesis/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-sdk/authentication"
	"github.com/TerraDharitri/drt-go-sdk/authentication/native/mock"
	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/TerraDharitri/drt-go-sdk/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedErr = errors.New("expected error")
var httpExpectedErr = authentication.CreateHTTPStatusError(http.StatusInternalServerError, expectedErr)

func TestNativeserver_NewNativeAuthServer(t *testing.T) {
	t.Parallel()

	t.Run("nil http server wrapper should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.HttpClientWrapper = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilHttpClientWrapper, err)
	})
	t.Run("nil signer should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.Signer = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilSigner, err)
	})
	t.Run("nil KeyGenerator should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.KeyGenerator = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, crypto.ErrNilKeyGenerator, err)
	})
	t.Run("nil pubKeyConverter should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.PubKeyConverter = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, core.ErrNilPubkeyConverter, err)
	})
	t.Run("nil token handler should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilTokenHandler, err)
	})
	t.Run("nil cacher should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TimestampsCacher = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilCacher, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		server, err := NewNativeAuthServer(args)
		require.NotNil(t, server)
		require.False(t, server.IsInterfaceNil())
		require.Nil(t, err)
	})
}
func TestNativeserver_Validate(t *testing.T) {
	t.Parallel()

	tokenTtl := int64(20)
	blockTimestamp := int64(10)
	providedBlockHash := "provided block hash"

	t.Run("invalid cached value should return error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TimestampsCacher = &testscommon.CacherStub{
			GetCalled: func(key []byte) (value interface{}, ok bool) {
				assert.Equal(t, []byte(providedBlockHash), key)
				return "invalid value", true
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp, 1)
		}

		err := server.Validate(&AuthToken{
			ttl:       tokenTtl,
			blockHash: providedBlockHash,
		})
		assert.True(t, errors.Is(err, authentication.ErrInvalidValue))
	})
	t.Run("httpClientWrapper returns error should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.HttpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				return nil, http.StatusInternalServerError, expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, httpExpectedErr, err)
	})
	t.Run("token expired should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.HttpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				buff, _ := json.Marshal(block)
				return buff, http.StatusOK, nil
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp+tokenTtl+1, 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, authentication.ErrTokenExpired, err)
	})
	t.Run("pubKeyConverter errors should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.HttpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				buff, _ := json.Marshal(block)
				return buff, http.StatusOK, nil
			},
		}
		args.PubKeyConverter = &testscommon.PubkeyConverterStub{
			DecodeCalled: func(humanReadable string) ([]byte, error) {
				return nil, expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp, 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, expectedErr, err)
	})
	t.Run("keyGenerator errors should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.HttpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				buff, _ := json.Marshal(block)
				return buff, http.StatusOK, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, expectedErr
			},
		}
		args.PubKeyConverter = &testscommon.PubkeyConverterStub{
			DecodeCalled: func(humanReadable string) ([]byte, error) {
				return nil, nil
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp, 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, expectedErr, err)
	})
	t.Run("invalid http result should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.HttpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				return nil, http.StatusOK, nil
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp, 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		assert.True(t, strings.Contains(err.Error(), "unexpected end of JSON input"))
	})
	t.Run("verification errors should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			GetUnsignedTokenCalled: func(authToken authentication.AuthToken) []byte {
				return []byte("token")
			},
		}
		args.HttpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				buff, _ := json.Marshal(block)
				return buff, http.StatusOK, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, nil
			},
		}
		args.Signer = &testsCommon.SignerStub{
			VerifyMessageCalled: func(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
				return expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp, 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, expectedErr, err)
	})
	t.Run("should work - token not cached", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			GetUnsignedTokenCalled: func(authToken authentication.AuthToken) []byte {
				return []byte("token")
			},
		}
		args.HttpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				buff, _ := json.Marshal(block)
				return buff, http.StatusOK, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, nil
			},
		}
		args.Signer = &testsCommon.SignerStub{
			VerifyMessageCalled: func(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
				return nil
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp, 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		assert.Nil(t, err)
	})
	t.Run("should work - token cached", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			GetUnsignedTokenCalled: func(authToken authentication.AuthToken) []byte {
				return []byte("token")
			},
		}
		args.HttpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				assert.Fail(t, "should have not been called")
				return []byte{}, http.StatusOK, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, nil
			},
		}
		args.Signer = &testsCommon.SignerStub{
			VerifyMessageCalled: func(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
				return nil
			},
		}
		args.TimestampsCacher = &testscommon.CacherStub{
			GetCalled: func(key []byte) (value interface{}, ok bool) {
				assert.Equal(t, []byte(providedBlockHash), key)
				return blockTimestamp, true
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp, 1)
		}

		err := server.Validate(&AuthToken{
			ttl:       tokenTtl,
			blockHash: providedBlockHash,
		})
		assert.Nil(t, err)
	})
}

func createMockArgsNativeAuthServer() ArgsNativeAuthServer {
	return ArgsNativeAuthServer{
		HttpClientWrapper: &testsCommon.HTTPClientWrapperStub{},
		TokenHandler:      &mock.AuthTokenHandlerStub{},
		Signer:            &testsCommon.SignerStub{},
		PubKeyConverter:   &testscommon.PubkeyConverterStub{},
		KeyGenerator:      &genesisMock.KeyGeneratorStub{},
		TimestampsCacher:  &testscommon.CacherStub{},
	}
}
