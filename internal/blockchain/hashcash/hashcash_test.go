package hashcash_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/jleipus/learn-blockchain/internal/blockchain/block"
	"github.com/jleipus/learn-blockchain/internal/blockchain/hashcash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	pow := hashcash.New()

	b := &block.Block{
		Timestamp:     1234567890,
		Transactions:  nil,
		PrevBlockHash: block.Hash{'0', '0', '0'},
	}

	hashcash.TimeNow = func() time.Time {
		return time.Unix(b.Timestamp, 0)
	}

	t.Run("produce", func(t *testing.T) {
		hash, nonce := pow.Produce(b)

		var hashInt big.Int
		hashInt.SetBytes(hash[:])

		expectedHash, err := hex.DecodeString("00003f60dc9c6a8c1840af11ebd071b0fe06df237854a49416dc1e8d9ec7d170")
		require.NoError(t, err)

		expectedNonce, err := hex.DecodeString("0000000000005e6c")
		require.NoError(t, err)

		assert.Equal(t, expectedHash, hash[:])
		assert.Equal(t, expectedNonce, nonce)

		b.Hash = hash
		b.PoW = nonce
	})

	t.Run("validate", func(t *testing.T) {
		valid := pow.Validate(b)
		assert.True(t, valid)
	})

	t.Run("validate invalid hash", func(t *testing.T) {
		invalidBlock := &block.Block{
			Timestamp:     1234567890,
			Transactions:  nil,
			PrevBlockHash: block.Hash{'0', '0', '0'},
			Hash:          block.Hash{'1', '2', '3', '4', '5', '6', '7', '8', '9', '1'}, // Invalid hash
			PoW:           []byte{'p', 'o', 'w', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'},
		}

		valid := pow.Validate(invalidBlock)
		assert.False(t, valid)
	})

	t.Run("validate nonce too short", func(t *testing.T) {
		invalidBlock := &block.Block{
			Timestamp:     1234567890,
			Transactions:  nil,
			PrevBlockHash: block.Hash{'0', '0', '0'},
			Hash:          block.Hash{'1', '2', '3', '4', '5', '6', '7', '8', '9', '1'}, // Invalid hash
			PoW:           []byte{'p', 'o', 'w'},
		}

		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "Invalid PoW length: expected 8, got 3", r.(string))
			}
		}()

		valid := pow.Validate(invalidBlock)
		assert.False(t, valid)

		t.Errorf("Test failed, panic was expected")
	})
}
