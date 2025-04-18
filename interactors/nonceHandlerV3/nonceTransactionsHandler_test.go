package nonceHandlerV3

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/TerraDharitri/drt-go-sdk/testsCommon"
)

var testAddressAsBech32String = "drt1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsaa8s54"

func TestSendTransactionsOneByOne(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
	require.NoError(t, err, "failed to create transaction handler")

	var txs []*transaction.FrontendTransaction

	for i := 0; i < 100; i++ {
		tx := &transaction.FrontendTransaction{
			Sender:   testAddressAsBech32String,
			Receiver: testAddressAsBech32String,
			GasLimit: 50000,
			ChainID:  "T",
			Value:    "5000000000000000000",
			Nonce:    uint64(i),
			GasPrice: 1000000000,
			Version:  2,
		}
		txs = append(txs, tx)
	}

	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err, "failed to apply nonce")
	require.True(t, getAccountCalled, "get account was not called")

	var wg sync.WaitGroup
	for _, tt := range txs {
		wg.Add(1)
		go func(tt *transaction.FrontendTransaction) {
			defer wg.Done()
			h, err := transactionHandler.SendTransactions(context.Background(), tt)
			require.NoError(t, err, "failed to send transaction")
			require.Equal(t, []string{strconv.FormatUint(tt.Nonce, 10)}, h)
		}(tt)
	}
	wg.Wait()
}

func TestSendTransactionsBulk(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
	require.NoError(t, err, "failed to create transaction handler")

	var txs []*transaction.FrontendTransaction

	for i := 0; i < 100; i++ {
		tx := &transaction.FrontendTransaction{
			Sender:   testAddressAsBech32String,
			Receiver: testAddressAsBech32String,
			GasLimit: 50000,
			ChainID:  "T",
			Value:    "5000000000000000000",
			Nonce:    uint64(i),
			GasPrice: 1000000000,
			Version:  2,
		}
		txs = append(txs, tx)
	}

	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err, "failed to apply nonce")
	require.True(t, getAccountCalled, "get account was not called")

	txHashes, err := transactionHandler.SendTransactions(context.Background(), txs...)
	require.NoError(t, err, "failed to send transactions as bulk")
	require.Equal(t, mockedStrings(0, 100), txHashes)
}

func TestSendTransactionsCloseInstant(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
	require.NoError(t, err, "failed to create transaction handler")

	var txs []*transaction.FrontendTransaction

	// Create 1k transactions.
	for i := 0; i < 100; i++ {
		tx := &transaction.FrontendTransaction{
			Sender:   testAddressAsBech32String,
			Receiver: testAddressAsBech32String,
			GasLimit: 50000,
			ChainID:  "T",
			Value:    "5000000000000000000",
			Nonce:    uint64(i),
			GasPrice: 1000000000,
			Version:  2,
		}
		txs = append(txs, tx)
	}

	// Apply nonce to them in a bulk.
	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err, "failed to apply nonce")

	// We only do this once, we check if the bool has been modified.
	require.True(t, getAccountCalled, "get account was not called")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// make sure that the Close function is called before the send function
		time.Sleep(time.Second)

		hashes, errSend := transactionHandler.SendTransactions(context.Background(), txs...)

		var counter int
		// Since the close is almost instant there should be none or very few transactions that have been processed.
		for _, h := range hashes {
			if h != "" {
				counter++
			}
		}

		require.Equal(t, 0, counter)
		require.NotNil(t, errSend)
		wg.Done()
	}()

	// Close the processes related to the transaction handler.
	transactionHandler.Close()

	wg.Wait()
	require.NoError(t, err, "failed to send transactions as bulk")
}

func TestSendTransactionsCloseDelay(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	// Create another proxyStub that adds some delay when sending transactions.
	mockArgs := ArgsNonceTransactionsHandlerV3{
		Proxy: &testsCommon.ProxyStub{
			SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
				// Presume this operation is taking roughly 100 ms. Meaning 10 operations / second.
				time.Sleep(100 * time.Millisecond)
				return strconv.FormatUint(tx.Nonce, 10), nil
			},
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				getAccountCalled = true
				return &data.Account{}, nil
			},
		},
		IntervalToSend: time.Second * 5,
	}

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(mockArgs)
	require.NoError(t, err, "failed to create transaction handler")

	var txs []*transaction.FrontendTransaction

	// Create 1k transactions.
	for i := 0; i < 100; i++ {
		tx := &transaction.FrontendTransaction{
			Sender:   testAddressAsBech32String,
			Receiver: testAddressAsBech32String,
			GasLimit: 50000,
			ChainID:  "T",
			Value:    "5000000000000000000",
			Nonce:    uint64(i),
			GasPrice: 1000000000,
			Version:  2,
		}
		txs = append(txs, tx)
	}

	// Apply nonce to them in a bulk.
	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err, "failed to apply nonce")

	// We only do this once, we check if the bool has been modified.
	require.True(t, getAccountCalled, "get account was not called")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		hashes, errSend := transactionHandler.SendTransactions(context.Background(), txs...)

		// Since the close is not instant. There should be some hashes that have been processed.
		require.NotEmpty(t, hashes, "no transaction should be processed")
		require.Equal(t, "context canceled while sending transaction for address drt1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsaa8s54", errSend.Error())
		wg.Done()
	}()

	// Close the processes related to the transaction handler with a delay.
	time.AfterFunc(2*time.Second, func() {
		transactionHandler.Close()
	})

	wg.Wait()
	require.NoError(t, err, "failed to send transactions as bulk")
}

func TestApplyNonceAndGasPriceConcurrently(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
	require.NoError(t, err, "failed to create transaction handler")

	var txs []*transaction.FrontendTransaction
	for i := 0; i < 100; i++ {
		tx := &transaction.FrontendTransaction{
			Sender:   testAddressAsBech32String,
			Receiver: testAddressAsBech32String,
			GasLimit: 50000,
			ChainID:  "T",
			Value:    "5000000000000000000",
			Nonce:    uint64(i),
			GasPrice: 1000000000,
			Version:  2,
		}
		txs = append(txs, tx)
	}

	// we apply the nonce on the initial transaction list in batches of 20. in order to test that the nonce handler is
	// able to do it concurrently providing unique nonces for every transaction.
	var wg sync.WaitGroup
	indices := []int{0, 19, 39, 59, 79, 99}
	for i := 0; i < len(indices)-1; i++ {
		wg.Add(1)
		beginIdx := indices[i]
		endIdx := indices[i+1]
		go func() {
			defer wg.Done()
			err := transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs[beginIdx:endIdx]...)
			require.Nil(t, err, "error should be nil")
		}()
	}
	wg.Wait()

	// since we applied the nonces concurrently, the slice won't have all of them in order. therefore we sort them
	// before comparing them to the expected output.
	sort.SliceStable(txs, func(i, j int) bool {
		return txs[i].Nonce < txs[j].Nonce
	})
	mockedNonces := mockedStrings(0, 100)
	for idx := range txs {
		mockNonce, _ := strconv.ParseUint(mockedNonces[idx], 10, 64)
		require.Equal(t, mockNonce, txs[idx].Nonce)
	}
}

func TestSendDuplicateNonces(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
	require.NoError(t, err, "failed to create transaction handler")

	tx := &transaction.FrontendTransaction{
		Sender:   testAddressAsBech32String,
		Receiver: testAddressAsBech32String,
		GasLimit: 50000,
		ChainID:  "T",
		Value:    "5000000000000000000",
		Nonce:    0,
		GasPrice: 1000000000,
		Version:  2,
	}

	wg := sync.WaitGroup{}
	errCount := atomic.Uint32{}
	sentCount := atomic.Uint32{}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hashes, sendErr := transactionHandler.SendTransactions(context.Background(), tx)
			if sendErr != nil {
				errCount.Add(1)
			}

			if hashes[0] != "" {
				sentCount.Add(1)
			}
		}()
	}
	wg.Wait()

	require.Equal(t, uint32(1), errCount.Load())
	require.Equal(t, uint32(1), sentCount.Load())
}

func TestSendDuplicateNoncesBatch(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
	require.NoError(t, err, "failed to create transaction handler")

	nonce := uint64(0)
	txs := make([]*transaction.FrontendTransaction, 0)
	for i := 0; i < 100; i++ {
		tx := &transaction.FrontendTransaction{
			Sender:   testAddressAsBech32String,
			Receiver: testAddressAsBech32String,
			GasLimit: 50000,
			ChainID:  "T",
			Value:    "5000000000000000000",
			Nonce:    nonce,
			GasPrice: 1000000000,
			Version:  2,
		}
		txs = append(txs, tx)
	}

	hashes, err := transactionHandler.SendTransactions(context.Background(), txs...)
	require.Contains(t, hashes, strconv.FormatUint(nonce, 10), "no transaction has been sent")
	require.Error(t, errors.New("transaction with nonce: 0 has already been scheduled to send while sending transaction for address drt1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsaa8s54"), err)
}

func TestDoubleCycle(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
	require.NoError(t, err, "failed to create transaction handler")

	var txs []*transaction.FrontendTransaction
	for i := 0; i < 100; i++ {
		tx := &transaction.FrontendTransaction{
			Sender:   testAddressAsBech32String,
			Receiver: testAddressAsBech32String,
			GasLimit: 50000,
			ChainID:  "T",
			Value:    "5000000000000000000",
			Version:  2,
		}
		txs = append(txs, tx)
	}

	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err)

	hashes, err := transactionHandler.SendTransactions(context.Background(), txs...)
	require.NoError(t, err)

	mockedNonces := mockedStrings(0, 100)
	for idx := range txs {
		require.Equal(t, mockedNonces[idx], hashes[idx])
	}

	var txs2 []*transaction.FrontendTransaction
	for i := 100; i < 200; i++ {
		tx := &transaction.FrontendTransaction{
			Sender:   testAddressAsBech32String,
			Receiver: testAddressAsBech32String,
			GasLimit: 50000,
			ChainID:  "T",
			Value:    "5000000000000000000",
			Version:  2,
		}
		txs2 = append(txs2, tx)
	}

	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs2...)
	require.NoError(t, err)

	hashes, err = transactionHandler.SendTransactions(context.Background(), txs2...)
	require.NoError(t, err)

	mockedNonces = mockedStrings(100, 200)
	for idx := range txs {
		require.Equal(t, mockedNonces[idx], hashes[idx])
	}
}

func createMockArgsNonceTransactionsHandlerV3(getAccountCalled *bool) ArgsNonceTransactionsHandlerV3 {
	return ArgsNonceTransactionsHandlerV3{
		Proxy: &testsCommon.ProxyStub{
			SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
				return strconv.FormatUint(tx.Nonce, 10), nil
			},
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				*getAccountCalled = true
				return &data.Account{}, nil
			},
		},
		IntervalToSend: time.Millisecond * 1,
	}
}

func mockedStrings(beginIdx, endIdx int) []string {
	mock := make([]string, endIdx-beginIdx)
	for i := 0; i < endIdx-beginIdx; i++ {
		mock[i] = strconv.Itoa(beginIdx + i)
	}

	return mock
}
