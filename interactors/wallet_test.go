package interactors

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWallet_GenerateMnemonicShouldWork(t *testing.T) {
	t.Parallel()

	w := NewWallet()
	mnemonic, err := w.GenerateMnemonic()
	assert.Nil(t, err)
	assert.Equal(t, 24, len(mnemonic.ToSplitMnemonicWords()))
}

func TestWallet_GetPrivateKeyFromMnemonic(t *testing.T) {
	t.Parallel()

	w := NewWallet()
	mnemonic := data.Mnemonic("acid twice post genre topic observe valid viable gesture fortune funny dawn around blood enemy page update reduce decline van bundle zebra rookie real")

	privKey := w.GetPrivateKeyFromMnemonic(mnemonic, 0, 0)
	expectedHexPrivKey := "0b7966138e80b8f3bb64046f56aea4250fd7bacad6ed214165cea6767fd0bc2c"
	assert.Equal(t, expectedHexPrivKey, hex.EncodeToString(privKey))

	privKey = w.GetPrivateKeyFromMnemonic(mnemonic, 0, 1)
	expectedHexPrivKey = "1648ad209d6b157a289884933e3bb30f161ec7113221ec16f87c3578b05830b0"
	assert.Equal(t, expectedHexPrivKey, hex.EncodeToString(privKey))
}

func TestWallet_CreateSeedFromMnemonicThenGetPrivateKeyFromSeed(t *testing.T) {
	t.Parallel()

	w := NewWallet()
	mnemonic := data.Mnemonic("acid twice post genre topic observe valid viable gesture fortune funny dawn around blood enemy page update reduce decline van bundle zebra rookie real")
	seed := w.CreateSeedFromMnemonic(mnemonic)

	privKey := w.GetPrivateKeyFromSeed(seed, 0, 0)
	expectedHexPrivKey := "0b7966138e80b8f3bb64046f56aea4250fd7bacad6ed214165cea6767fd0bc2c"
	assert.Equal(t, expectedHexPrivKey, hex.EncodeToString(privKey))

	privKey = w.GetPrivateKeyFromSeed(seed, 0, 1)
	expectedHexPrivKey = "1648ad209d6b157a289884933e3bb30f161ec7113221ec16f87c3578b05830b0"
	assert.Equal(t, expectedHexPrivKey, hex.EncodeToString(privKey))
}

func TestWallet_GetAddressFromMnemonicWalletIntegration(t *testing.T) {
	t.Parallel()

	w := NewWallet()
	mnemonic := data.Mnemonic("bid involve twenty cave offer life hello three walnut travel rare bike edit canyon ice brave theme furnace cotton swing wear bread fine latin")
	privKey := w.GetPrivateKeyFromMnemonic(mnemonic, 0, 0)
	fmt.Println(hex.EncodeToString(privKey))
	address, err := w.GetAddressFromPrivateKey(privKey)
	assert.Nil(t, err)
	expectedBech32Addr := "drt1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqspzc322"

	addressAsBech32String, err := address.AddressAsBech32String()
	assert.Nil(t, err)

	assert.Equal(t, expectedBech32Addr, addressAsBech32String)
}

func TestWallet_GetAddressFromPrivateKey(t *testing.T) {
	t.Parallel()

	hexPrivKey := "0b7966138e80b8f3bb64046f56aea4250fd7bacad6ed214165cea6767fd0bc2c"
	privKey, err := hex.DecodeString(hexPrivKey)
	require.Nil(t, err)

	w := NewWallet()
	address, err := w.GetAddressFromPrivateKey(privKey)
	assert.Nil(t, err)
	expectedBech32Addr := "drt1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsy839fr"

	addressAsBech32String, err := address.AddressAsBech32String()
	assert.Nil(t, err)

	assert.Equal(t, expectedBech32Addr, addressAsBech32String)
}

func TestWallet_LoadPrivateKeyFromJsonFile(t *testing.T) {
	t.Parallel()

	filename := "testdata/test.json"
	password := "pAssword1~"
	w := NewWallet()
	privkey, err := w.LoadPrivateKeyFromJsonFile(filename, password)
	require.Nil(t, err)
	expectedHexPrivKey := "15cfe2140ee9821f706423036ba58d1e6ec13dbc4ebf206732ad40b5236af403"
	assert.Equal(t, expectedHexPrivKey, hex.EncodeToString(privkey))
}

func TestWallet_LoadPrivateKeyFromJsonFileWithKind(t *testing.T) {
	t.Parallel()

	filename := "testdata/testWithKind.json"
	password := "password"
	w := NewWallet()
	privkey, err := w.LoadPrivateKeyFromJsonFile(filename, password)
	require.Nil(t, err)
	expectedHexPrivKey := "413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9"
	assert.Equal(t, expectedHexPrivKey, hex.EncodeToString(privkey))

	address, _ := w.GetAddressFromPrivateKey(privkey)

	addressAsBech32String, err := address.AddressAsBech32String()
	assert.Nil(t, err)

	require.Equal(t, "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf", addressAsBech32String)
}

func TestWallet_SavePrivateKeyToJsonFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fileName := path.Join(dir, "temp.json")

	file, err := os.Create(fileName)
	require.Nil(t, err)
	_ = file.Close() //close the file so the save can write in it

	defer func() {
		_ = os.Remove(file.Name())
	}()

	hexPrivKey := "15cfe2140ee9821f706423036ba58d1e6ec13dbc4ebf206732ad40b5236af403"
	privKey, err := hex.DecodeString(hexPrivKey)
	require.Nil(t, err)
	password := "pAssword1~"

	w := NewWallet()
	err = w.SavePrivateKeyToJsonFile(privKey, password, file.Name())
	require.Nil(t, err)

	recoveredSk, err := w.LoadPrivateKeyFromJsonFile(file.Name(), password)
	require.Nil(t, err)
	assert.Equal(t, privKey, recoveredSk)
}

func TestWallet_LoadPrivateKeyFromPemFile(t *testing.T) {
	t.Parallel()

	filename := "testdata/test.pem"
	w := NewWallet()
	privkey, err := w.LoadPrivateKeyFromPemFile(filename)
	require.Nil(t, err)

	address, err := w.GetAddressFromPrivateKey(privkey)
	require.Nil(t, err)

	expectedBech32Address := "drt1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsaa8s54"

	addressAsBech32String, err := address.AddressAsBech32String()
	require.Nil(t, err)

	assert.Equal(t, expectedBech32Address, addressAsBech32String)
}

func TestWallet_SavePrivateKeyToPemFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fileName := path.Join(dir, "temp.json")
	file, err := os.Create(fileName)
	require.Nil(t, err)
	_ = file.Close() //close the file so the save can write in it

	defer func() {
		_ = os.Remove(file.Name())
	}()

	hexPrivKey := "15cfe2140ee9821f706423036ba58d1e6ec13dbc4ebf206732ad40b5236af403"
	privKey, err := hex.DecodeString(hexPrivKey)
	require.Nil(t, err)

	w := NewWallet()
	err = w.SavePrivateKeyToPemFile(privKey, file.Name())
	require.Nil(t, err)

	recoveredSk, err := w.LoadPrivateKeyFromPemFile(file.Name())
	require.Nil(t, err)
	assert.Equal(t, privKey, recoveredSk)
}
