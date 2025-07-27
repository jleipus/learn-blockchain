//go:build integration
// +build integration

package blockchain_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/badger"
	"github.com/jleipus/learn-blockchain/internal/blockchain/hashcash"
	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStorage(t *testing.T) (blockchain.Storage, func()) {
	dir := t.TempDir()
	db, err := badger.NewStorage(dir)
	require.NoError(t, err, "failed to create badger db")

	return db, func() {
		db.Close()
		os.RemoveAll(dir)
	}
}

func TestMain(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	t.Cleanup(cleanup)

	powFactory := hashcash.New()

	wallets := wallet.NewCollection(storage)

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
		err := blockchain.CreateBlockchain(storage, powFactory, address1)
		require.NoError(t, err, "failed to create blockchain")
		bc, err = blockchain.LoadBlockchain(storage, powFactory, wallets)
		require.NoError(t, err, "failed to load blockchain")
	})

	t.Run("reindex UTXO set", func(t *testing.T) {
		err = bc.ReindexUTXOSet()
		require.NoError(t, err, "failed to reindex UTXO set")
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

		cbTx, err := transaction.NewCoinbaseTX(address1, "")
		require.NoError(t, err, "failed to create coinbase transaction")

		b, err := bc.MineBlock([]*transaction.Tx{tx, cbTx})
		require.NoError(t, err)

		err = bc.Update(*b)
		require.NoError(t, err, "failed to update UTXO set after mining block")
	})

	t.Run("create transactions 2", func(t *testing.T) {
		tx, err := bc.NewUTXOTransaction(address2, address3, 5)
		require.NoError(t, err)

		cbTx, err := transaction.NewCoinbaseTX(address2, "")
		require.NoError(t, err, "failed to create coinbase transaction")

		b, err := bc.MineBlock([]*transaction.Tx{tx, cbTx})
		require.NoError(t, err)

		err = bc.Update(*b)
		require.NoError(t, err, "failed to update UTXO set after mining block")
	})

	t.Run("print blockchain", func(t *testing.T) {
		chain := ""
		for _, b := range bc.Blocks() {
			chain += fmt.Sprintf("============ Block %x ============\n", b.Hash)
			chain += fmt.Sprintf("Prev. block: %x\n", b.PrevBlockHash)
		}
		t.Logf("\n%s", chain)
	})

	t.Run("check new wallet 1 balance", func(t *testing.T) {
		balance := getBalance(t, bc, address1)
		assert.Equal(t, 13, balance)
	})

	t.Run("check new wallet 2 balance", func(t *testing.T) {
		balance := getBalance(t, bc, address2)
		assert.Equal(t, 12, balance)
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

	outputs, err := bc.FindUnspentTxOutputs(pubKeyHash)
	require.NoError(t, err, "failed to find unspent transaction outputs")

	balance := 0
	for _, out := range outputs {
		balance += int(out.Value)
	}

	return balance
}
