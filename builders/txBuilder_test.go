package builders

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/ed25519"
	"github.com/TerraDharitri/drt-go-sdk/blockchain/cryptoProvider"
	"github.com/TerraDharitri/drt-go-sdk/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
)

func TestNewTxBuilder(t *testing.T) {
	t.Parallel()

	t.Run("nil signer should error", func(t *testing.T) {
		t.Parallel()

		tb, err := NewTxBuilder(nil)
		assert.True(t, check.IfNil(tb))
		assert.Equal(t, ErrNilSigner, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		tb, err := NewTxBuilder(&testsCommon.SignerStub{})
		assert.False(t, check.IfNil(tb))
		assert.Nil(t, err)
	})
}

func TestTxBuilder_ApplySignature(t *testing.T) {
	t.Parallel()

	sk, err := hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c40d248169f4dd7c90537f05be1c49772ddbf8f7948b507ed17fb23284cf218b7d")
	require.Nil(t, err)
	cryptoHolder, err := cryptoProvider.NewCryptoComponentsHolder(keyGen, sk)
	require.Nil(t, err)
	value := big.NewInt(999)
	tx := transaction.FrontendTransaction{
		Value:    value.String(),
		Receiver: "drt1l20m7kzfht5rhdnd4zvqr82egk7m4nvv3zk06yw82zqmrt9kf0zs5ewnr7",
		GasPrice: 10,
		GasLimit: 100000,
		Data:     []byte(""),
		ChainID:  "integration test chain id",
		Version:  uint32(1),
	}

	t.Run("tx signer errors when signing should error", func(t *testing.T) {
		t.Parallel()

		txCopy := tx
		expectedErr := errors.New("expected error")
		tb, _ := NewTxBuilder(&testsCommon.SignerStub{
			SignByteSliceCalled: func(_ []byte, _ crypto.PrivateKey) ([]byte, error) {
				return nil, expectedErr
			},
		})

		errGenerate := tb.ApplyUserSignature(cryptoHolder, &txCopy)
		assert.Empty(t, txCopy.Signature)
		assert.Equal(t, expectedErr, errGenerate)
	})

	signer := cryptoProvider.NewSigner()
	tb, err := NewTxBuilder(signer)
	require.Nil(t, err)

	t.Run("sign on all tx bytes should work", func(t *testing.T) {
		t.Parallel()

		txCopy := tx
		errGenerate := tb.ApplyUserSignature(cryptoHolder, &txCopy)
		require.Nil(t, errGenerate)

		assert.Equal(t, "drt1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7sk3pvxc", txCopy.Sender)
		assert.Equal(t, "e5df0318f3b9769d556f769e343e66d42254513996626fa01045c394d88e186645acf496b37fb5726ddb52dbe20068b130b9853bf86373e5645409f8c131f605",
			txCopy.Signature)
	})
	t.Run("sign on tx hash should work", func(t *testing.T) {
		t.Parallel()

		txCopy := tx
		txCopy.Version = 2
		txCopy.Options = 1

		errGenerate := tb.ApplyUserSignature(cryptoHolder, &txCopy)
		require.Nil(t, errGenerate)

		assert.Equal(t, "drt1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7sk3pvxc", txCopy.Sender)
		assert.Equal(t, "0dfba6d199ee742602da106fdbc126384419f8b530f5f6fcd47b6a1fb7886bb1571b85c52efb0b345a90fb4f2a818fac55dd6dd8a02fcbe8d0375136c78c390f",
			txCopy.Signature)
	})
}

func TestTxBuilder_ApplySignatureAndGenerateTxHash(t *testing.T) {
	t.Parallel()

	sk, err := hex.DecodeString("28654d9264f55f18d810bb88617e22c117df94fa684dfe341a511a72dfbf2b68")
	require.Nil(t, err)
	cryptoHolder, err := cryptoProvider.NewCryptoComponentsHolder(keyGen, sk)
	require.Nil(t, err)

	t.Run("fails if the signature is missing", func(t *testing.T) {
		t.Parallel()

		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		txHash, errGenerate := tb.ComputeTxHash(&transaction.FrontendTransaction{})
		assert.Nil(t, txHash)
		assert.Equal(t, ErrMissingSignature, errGenerate)
	})

	t.Run("should generate tx hash", func(t *testing.T) {
		t.Parallel()

		tx := &transaction.FrontendTransaction{
			Nonce:    1,
			Value:    "11500313000000000000",
			Receiver: "drt1p72ru5zcdsvgkkcm9swtvw2zy5epylwgv8vwquptkw7ga7pfvk7qlz8sps",
			GasPrice: 1000000000,
			GasLimit: 60000,
			Data:     []byte(""),
			ChainID:  "T",
			Version:  uint32(1),
		}
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())

		_ = tb.ApplyUserSignature(cryptoHolder, tx)
		assert.Equal(t, "19d8333ed6aee1fa3d2e2a05e410a16092e9b81b27b162b51c782c3f570b9bed900783d92c4b273c12f6698939814536cfe43d60d3e319d845dd8a5c8220bb0b", tx.Signature)

		txHash, errGenerate := tb.ComputeTxHash(tx)
		assert.Nil(t, errGenerate)
		assert.Equal(t, "712f95bac04a4898d668281a7ace94105afdb4283e7b64986eb87e2906b0b19e", hex.EncodeToString(txHash))
	})
}

func TestTxBuilder_ApplyUserSignatureAndGenerateWithTxGuardian(t *testing.T) {
	t.Parallel()

	guardianAddress := "drt1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7sk3pvxc"
	skGuardian, err := hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c40d248169f4dd7c90537f05be1c49772ddbf8f7948b507ed17fb23284cf218b7d")
	require.Nil(t, err)
	cryptoHolderGuardian, err := cryptoProvider.NewCryptoComponentsHolder(keyGen, skGuardian)
	require.Nil(t, err)

	senderAddress := "drt1lta2vgd0tkeqqadkvgef73y0efs6n3xe5ss589ufhvmt6tcur8kqvfh4da"
	sk, err := hex.DecodeString("28654d9264f55f18d810bb88617e22c117df94fa684dfe341a511a72dfbf2b68")
	require.Nil(t, err)
	cryptoHolder, err := cryptoProvider.NewCryptoComponentsHolder(keyGen, sk)
	require.Nil(t, err)

	tx := transaction.FrontendTransaction{
		Nonce:        1,
		Value:        "11500313000000000000",
		Receiver:     "drt1p72ru5zcdsvgkkcm9swtvw2zy5epylwgv8vwquptkw7ga7pfvk7qlz8sps",
		GasPrice:     1000000000,
		GasLimit:     60000,
		Data:         []byte(""),
		ChainID:      "T",
		Version:      uint32(2),
		Options:      transaction.MaskGuardedTransaction,
		GuardianAddr: guardianAddress,
	}

	t.Run("no guardian option should fail", func(t *testing.T) {
		txLocal := tx
		txLocal.Options = 0
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		_ = tb.ApplyUserSignature(cryptoHolder, &txLocal)
		err = tb.ApplyGuardianSignature(cryptoHolderGuardian, &txLocal)

		require.Equal(t, ErrMissingGuardianOption, err)
	})

	t.Run("no guardian address should fail", func(t *testing.T) {
		txLocal := tx
		txLocal.GuardianAddr = ""
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		_ = tb.ApplyUserSignature(cryptoHolder, &txLocal)

		err = tb.ApplyGuardianSignature(cryptoHolderGuardian, &txLocal)
		require.NotNil(t, err)
	})

	t.Run("different guardian address should fail", func(t *testing.T) {
		txLocal := tx
		txLocal.GuardianAddr = senderAddress
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		_ = tb.ApplyUserSignature(cryptoHolder, &txLocal)

		err = tb.ApplyGuardianSignature(cryptoHolderGuardian, &txLocal)
		require.Equal(t, ErrGuardianDoesNotMatch, err)
	})
	t.Run("correct guardian ok", func(t *testing.T) {
		txLocal := tx
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		_ = tb.ApplyUserSignature(cryptoHolder, &txLocal)

		err = tb.ApplyGuardianSignature(cryptoHolderGuardian, &txLocal)
		require.Nil(t, err)
	})
	t.Run("correct guardian and sign with hash ok", func(t *testing.T) {
		txLocal := tx
		txLocal.Options |= transaction.MaskSignedWithHash
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		_ = tb.ApplyUserSignature(cryptoHolder, &txLocal)

		err = tb.ApplyGuardianSignature(cryptoHolderGuardian, &txLocal)
		require.Nil(t, err)
	})
}

func TestTxBuilder_ApplyRelayerSignature(t *testing.T) {
	t.Parallel()

	relayerAddress := "drt1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7sk3pvxc"
	skRelayer, err := hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c40d248169f4dd7c90537f05be1c49772ddbf8f7948b507ed17fb23284cf218b7d")
	require.Nil(t, err)
	cryptoHolderRelayer, err := cryptoProvider.NewCryptoComponentsHolder(keyGen, skRelayer)
	require.Nil(t, err)

	senderAddress := "drt1lta2vgd0tkeqqadkvgef73y0efs6n3xe5ss589ufhvmt6tcur8kqvfh4da"
	sk, err := hex.DecodeString("28654d9264f55f18d810bb88617e22c117df94fa684dfe341a511a72dfbf2b68")
	require.Nil(t, err)
	cryptoHolder, err := cryptoProvider.NewCryptoComponentsHolder(keyGen, sk)
	require.Nil(t, err)

	tx := transaction.FrontendTransaction{
		Nonce:       1,
		Value:       "1000000000000000000",
		Receiver:    "drt1p72ru5zcdsvgkkcm9swtvw2zy5epylwgv8vwquptkw7ga7pfvk7qlz8sps",
		GasPrice:    1000000000,
		GasLimit:    100000,
		Data:        []byte("gift"),
		ChainID:     "T",
		Version:     uint32(2),
		RelayerAddr: relayerAddress,
	}

	t.Run("no relayer address should fail", func(t *testing.T) {
		txLocal := tx
		txLocal.RelayerAddr = ""
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		_ = tb.ApplyUserSignature(cryptoHolder, &txLocal)

		err = tb.ApplyRelayerSignature(cryptoHolderRelayer, &txLocal)
		require.NotNil(t, err)
	})
	t.Run("different relayer address should fail", func(t *testing.T) {
		txLocal := tx
		txLocal.RelayerAddr = senderAddress
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		_ = tb.ApplyUserSignature(cryptoHolder, &txLocal)

		err = tb.ApplyRelayerSignature(cryptoHolderRelayer, &txLocal)
		require.Equal(t, ErrRelayerDoesNotMatch, err)
	})
	t.Run("correct relayer should work", func(t *testing.T) {
		txLocal := tx
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		_ = tb.ApplyUserSignature(cryptoHolder, &txLocal)

		err = tb.ApplyRelayerSignature(cryptoHolderRelayer, &txLocal)
		require.Nil(t, err)
	})
	t.Run("correct relayer and sign with hash ok", func(t *testing.T) {
		txLocal := tx
		txLocal.Options |= transaction.MaskSignedWithHash
		tb, _ := NewTxBuilder(cryptoProvider.NewSigner())
		_ = tb.ApplyUserSignature(cryptoHolder, &txLocal)

		err = tb.ApplyRelayerSignature(cryptoHolderRelayer, &txLocal)
		require.Nil(t, err)
	})
}
