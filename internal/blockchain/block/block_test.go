package block_test

import (
	"bytes"
	"crypto/sha256"
	"testing"
	"time"

	"github.com/jleipus/learn-blockchain/internal/blockchain/block"
	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTransaction(id string) *transaction.Tx {
	txID := transaction.TxID{}
	copy(txID[:], []byte(id))

	return &transaction.Tx{
		ID: txID,
		Vin: []transaction.TxInput{
			{
				TxID:      transaction.TxID{},
				Vout:      0,
				Signature: []byte("test-signature"),
				PubKey:    []byte("test-pubkey"),
			},
		},
		Vout: []transaction.TxOutput{
			{
				Value:      100,
				PubKeyHash: []byte("test-pubkey-hash"),
			},
		},
	}
}

func getBlock() *block.Block {
	tx1 := getTransaction("tx1")
	tx2 := getTransaction("tx2")

	prevHash := block.Hash{}
	copy(prevHash[:], []byte("previous-block-hash"))

	blockHash := block.Hash{}
	copy(blockHash[:], []byte("current-block-hash"))

	return &block.Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  []*transaction.Tx{tx1, tx2},
		PrevBlockHash: prevHash,
		Hash:          blockHash,
		PoW:           []byte("proof-of-work-data"),
	}
}

func TestHashTransactions(t *testing.T) {
	t.Run("empty transactions", func(t *testing.T) {
		b := &block.Block{
			Transactions: []*transaction.Tx{},
		}

		result := b.HashTransactions()

		expected := sha256.Sum256([]byte{})

		assert.Equal(t, expected[:], result)
	})

	t.Run("single transaction", func(t *testing.T) {
		b := &block.Block{
			Transactions: []*transaction.Tx{
				getTransaction("single-tx"),
			},
		}

		result := b.HashTransactions()

		// Pad to 32 bytes to match TxID size
		txID := make([]byte, 32)
		copy(txID, []byte("single-tx"))
		expected := sha256.Sum256(txID)

		assert.Equal(t, expected[:], result)
	})

	t.Run("multiple transactions", func(t *testing.T) {
		b := getBlock()

		result := b.HashTransactions()

		tx1ID := make([]byte, 32)
		copy(tx1ID, []byte("tx1"))
		tx2ID := make([]byte, 32)
		copy(tx2ID, []byte("tx2"))

		combined := bytes.Join([][]byte{tx1ID, tx2ID}, []byte{})
		expected := sha256.Sum256(combined)
		assert.Equal(t, expected[:], result)
	})
}

func TestSerializeDeserialize(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		b := &block.Block{}
		serialized := b.Serialize()
		require.NotEmpty(t, serialized)

		var deserialized block.Block
		err := deserialized.Deserialize(serialized)
		require.NoError(t, err)

		assert.Equal(t, b, &deserialized)
	})

	t.Run("with basic data", func(t *testing.T) {
		b := &block.Block{
			Timestamp:     1234567890,
			Transactions:  nil,
			PrevBlockHash: block.Hash{},
			Hash:          block.Hash{},
			PoW:           []byte("test-pow"),
		}
		serialized := b.Serialize()
		require.NotEmpty(t, serialized)

		var deserialized block.Block
		err := deserialized.Deserialize(serialized)
		require.NoError(t, err)

		assert.Equal(t, b, &deserialized)
	})

	t.Run("with transactions", func(t *testing.T) {
		b := getBlock()
		serialized := b.Serialize()
		require.NotEmpty(t, serialized)

		var deserialized block.Block
		err := deserialized.Deserialize(serialized)
		require.NoError(t, err)

		assert.Equal(t, b, &deserialized)
	})
}
