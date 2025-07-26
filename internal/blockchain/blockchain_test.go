package blockchain_test

import (
	"testing"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/mock"
	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	mockStorage := mock.NewStorage()
	mockPowFactory := mock.NewPoWFactory()

	wallets := wallet.NewCollection(mockStorage)

	var err error

	var address1 string
	t.Run("create wallet 1", func(t *testing.T) {
		address1, err = wallets.AddWallet()
		require.NoError(t, err, "failed to add wallet")
	})

	var address2 string
	t.Run("create wallet 2", func(t *testing.T) {
		address2, err = wallets.AddWallet()
		require.NoError(t, err, "failed to add wallet")
	})

	var address3 string
	t.Run("create wallet 3", func(t *testing.T) {
		address3, err = wallets.AddWallet()
		require.NoError(t, err, "failed to add wallet")
	})

	var bc *blockchain.Blockchain
	t.Run("create blockchain", func(t *testing.T) {
		err := blockchain.CreateBlockchain(mockStorage, mockPowFactory, address1)
		require.NoError(t, err, "failed to create blockchain")
		bc, err = blockchain.LoadBlockchain(mockStorage, mockPowFactory, wallets)
		require.NoError(t, err, "failed to load blockchain")
	})

	t.Run("check wallet 1 balance", func(t *testing.T) {
		balance := getBalance(t, bc, address1)
		assert.Equal(t, 10, balance)
	})

	t.Run("check wallet 2 balance", func(t *testing.T) {
		balance := getBalance(t, bc, address2)
		assert.Equal(t, 0, balance)
	})

	t.Run("check wallet 3 balance", func(t *testing.T) {
		balance := getBalance(t, bc, address3)
		assert.Equal(t, 0, balance)
	})

	t.Run("create transactions 1", func(t *testing.T) {
		tx, err := bc.NewUTXOTransaction(address1, address2, 7)
		require.NoError(t, err)

		err = bc.MineBlock([]*transaction.Tx{tx})
		require.NoError(t, err)
	})

	t.Run("create transactions 2", func(t *testing.T) {
		tx, err := bc.NewUTXOTransaction(address2, address3, 5)
		require.NoError(t, err)

		err = bc.MineBlock([]*transaction.Tx{tx})
		require.NoError(t, err)
	})

	t.Run("check new wallet 1 balance", func(t *testing.T) {
		balance := getBalance(t, bc, address1)
		assert.Equal(t, 3, balance)
	})

	t.Run("check new wallet 2 balance", func(t *testing.T) {
		balance := getBalance(t, bc, address2)
		assert.Equal(t, 2, balance)
	})

	t.Run("check new wallet 3 balance", func(t *testing.T) {
		balance := getBalance(t, bc, address3)
		assert.Equal(t, 5, balance)
	})
}

func getBalance(t *testing.T, bc *blockchain.Blockchain, address string) int {
	t.Helper()

	pubKeyHash, err := wallet.GetHashFromAddress([]byte(address))
	require.NoError(t, err, "failed to get public key hash from address")

	balance := 0
	for _, out := range bc.FindUnspentTxOutputs(pubKeyHash) {
		balance += int(out.Value)
	}

	return balance
}
