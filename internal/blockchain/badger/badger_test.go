package badger_test

import (
	"os"
	"testing"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/badger"
	"github.com/jleipus/learn-blockchain/internal/blockchain/block"
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

func TestSetAndGetTip(t *testing.T) {
	db, cleanup := setupTestStorage(t)
	t.Cleanup(cleanup)

	t.Run("not found", func(t *testing.T) {
		tip, err := db.GetTip()
		require.NoError(t, err)
		var zero block.Hash
		assert.Equal(t, zero, tip)
	})

	t.Run("ok", func(t *testing.T) {
		var hash block.Hash
		copy(hash[:], []byte("testhash"))

		err := db.SetTip(hash)
		require.NoError(t, err)

		tip, err := db.GetTip()
		require.NoError(t, err)
		assert.Equal(t, hash, tip)
	})
}

func TestAddAndGetBlock(t *testing.T) {
	db, cleanup := setupTestStorage(t)
	t.Cleanup(cleanup)

	b := &block.Block{
		Timestamp:     1234567890,
		Transactions:  nil,
		PrevBlockHash: block.Hash{'0', '0', '0'},
		Hash:          block.Hash{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'},
		PoW:           []byte{'p', 'o', 'w'},
	}

	err := db.AddBlock(*b)
	require.NoError(t, err)

	t.Run("ok", func(t *testing.T) {
		retrievedBlock, err := db.GetBlock(b.Hash)
		require.NoError(t, err)
		assert.Equal(t, b, retrievedBlock)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := db.GetBlock(block.Hash{'n', 'o', 't', 'f', 'o', 'u', 'n', 'd'})
		assert.Error(t, err)
	})
}

func TestAddAndGetWallet(t *testing.T) {
	db, cleanup := setupTestStorage(t)
	t.Cleanup(cleanup)

	collection := wallet.NewCollection(db)

	address1, err := collection.AddWallet()
	require.NoError(t, err)
	wlt1, err := collection.GetWallet(address1)
	require.NoError(t, err)
	address2, err := collection.AddWallet()
	require.NoError(t, err)
	wlt2, err := collection.GetWallet(address2)
	require.NoError(t, err)

	err = db.AddWallet(address1, *wlt1)
	require.NoError(t, err)
	err = db.AddWallet(address2, *wlt2)
	require.NoError(t, err)

	t.Run("ok", func(t *testing.T) {
		retrievedWallet, err := db.GetWallet(address1)
		require.NoError(t, err)
		assert.Equal(t, wlt1, retrievedWallet)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := db.GetWallet("nonexistent")
		assert.Error(t, err)
	})

	t.Run("get addresses", func(t *testing.T) {
		addresses := db.GetAddresses()
		assert.Contains(t, addresses, address1)
		assert.Contains(t, addresses, address2)
		assert.Len(t, addresses, 2)
	})
}
