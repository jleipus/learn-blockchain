package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
)

type Hash [32]byte

// Storage is an interface for a storage system that can store and retrieve blocks.
type Storage interface {
	// GetTip retrieves the hash of the last block in the blockchain.
	// It returns an empty hash if the blockchain is empty.
	GetTip() (Hash, error)
	// SetTip sets the hash of the last block in the blockchain.
	SetTip(hash Hash) error
	// GetBlock retrieves a block by its hash.
	GetBlock(hash Hash) (*Block, error)
	// AddBlock adds a new block to the blockchain.
	AddBlock(block Block) error
	// Close closes the storage connection.
	Close() error
}

// Block represents a block in the blockchain.
type Block struct {
	// Timestamp is the time when the block was created.
	Timestamp int64
	// Transactions is a slice of transactions included in the block.
	Transactions []*transaction.Tx
	// PrevBlockHash is the hash of the previous block in the chain.
	PrevBlockHash Hash
	// Hash is the hash of the current block.
	Hash Hash
	// PoW is the proof of work data for the block.
	// It contains any additional data needed to verify the block's validity.
	// e.g., nonce, difficulty target, etc.
	PoW []byte
}

// HashTransactions computes the hash of all transactions in the block.
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID[:])
	}

	txSHA256 := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txSHA256[:]
}

// Serialize serializes the block into a byte slice using gob encoding.
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		// Error will only occur if the input contains unsupported types.
		// In this case, we panic because the block should only contain serializable types.
		panic(err)
	}

	return result.Bytes()
}

// Deserialize deserializes a byte slice into a Block using gob encoding.
func (b *Block) Deserialize(d []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(b)
	if err != nil {
		return err
	}

	return nil
}
