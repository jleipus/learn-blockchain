package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Block represents a block in the blockchain.
type Block struct {
	// Timestamp is the time when the block was created.
	Timestamp int64
	// Transactions is a slice of transactions included in the block.
	Transactions []*Transaction
	// PrevBlockHash is the hash of the previous block in the chain.
	PrevBlockHash BlockHash
	// Hash is the hash of the current block.
	Hash BlockHash
	// PoW is the proof of work data for the block.
	// It contains any additional data needed to verify the block's validity.
	// e.g., nonce, difficulty target, etc.
	PoW []byte
}

func (b *Block) HashTransactions() [32]byte {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID[:])
	}

	return sha256.Sum256(bytes.Join(txHashes, []byte{}))
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}
