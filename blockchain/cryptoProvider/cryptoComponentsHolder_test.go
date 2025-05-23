package cryptoProvider

import (
	"encoding/hex"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/TerraDharitri/drt-go-chain/testscommon/cryptoMocks"
	"github.com/TerraDharitri/drt-go-sdk/testsCommon"
	"github.com/stretchr/testify/require"
)

func TestNewCryptoComponentsHolder(t *testing.T) {
	t.Parallel()

	t.Run("invalid privateKey bytes", func(t *testing.T) {
		t.Parallel()

		keyGenInstance := &cryptoMocks.KeyGenStub{
			PrivateKeyFromByteArrayStub: func(b []byte) (crypto.PrivateKey, error) {
				return nil, expectedError
			},
		}
		holder, err := NewCryptoComponentsHolder(keyGenInstance, []byte(""))
		require.Nil(t, holder)
		require.Equal(t, expectedError, err)
	})
	t.Run("invalid publicKey bytes", func(t *testing.T) {
		t.Parallel()

		publicKey := &testsCommon.PublicKeyStub{
			ToByteArrayCalled: func() ([]byte, error) {
				return nil, expectedError
			},
		}
		privateKey := &testsCommon.PrivateKeyStub{
			GeneratePublicCalled: func() crypto.PublicKey {
				return publicKey
			}}
		keyGenInstance := &cryptoMocks.KeyGenStub{
			PrivateKeyFromByteArrayStub: func(b []byte) (crypto.PrivateKey, error) {
				return privateKey, nil
			},
		}
		holder, err := NewCryptoComponentsHolder(keyGenInstance, []byte(""))
		require.Nil(t, holder)
		require.Equal(t, expectedError, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		privateKey := &testsCommon.PrivateKeyStub{
			GeneratePublicCalled: func() crypto.PublicKey {
				return &testsCommon.PublicKeyStub{
					ToByteArrayCalled: func() ([]byte, error) {
						return make([]byte, 32), nil
					},
				}
			},
		}
		keyGenInstance := &cryptoMocks.KeyGenStub{
			PrivateKeyFromByteArrayStub: func(b []byte) (crypto.PrivateKey, error) {
				return privateKey, nil
			},
		}
		holder, err := NewCryptoComponentsHolder(keyGenInstance, []byte(""))
		require.False(t, check.IfNil(holder))
		require.Nil(t, err)
		_ = holder.GetPublicKey()
		_ = holder.GetPrivateKey()
		_ = holder.GetBech32()
		_ = holder.GetAddressHandler()
	})
	t.Run("should work with real components", func(t *testing.T) {
		t.Parallel()

		sk, _ := hex.DecodeString("45f72e8b6e8d10086bacd2fc8fa1340f82a3f5d4ef31953b463ea03c606533a6")
		holder, err := NewCryptoComponentsHolder(keyGen, sk)
		require.False(t, check.IfNil(holder))
		require.Nil(t, err)
		publicKey := holder.GetPublicKey()
		privateKey := holder.GetPrivateKey()
		require.False(t, check.IfNil(publicKey))
		require.False(t, check.IfNil(privateKey))

		bech32Address := holder.GetBech32()
		addressHandler := holder.GetAddressHandler()

		addressAsBech32String, err := addressHandler.AddressAsBech32String()
		require.Nil(t, err)

		require.Equal(t, addressAsBech32String, bech32Address)
		require.Equal(t, "drt1j84k44nsqsme8r6e5aawutx0z2cd6cyx3wprkzdh73x2cf0kqvksqd8ss7", bech32Address)
	})
}
