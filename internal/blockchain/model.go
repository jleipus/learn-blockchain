package blockchain

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain/block"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
)

type Storage interface {
	block.Storage
	wallet.Storage
}

// ProofOfWorkFactory is an implementation of a proof-of-work algorithm.
// It is responsible for producing a hash for a block and validating it.
type ProofOfWorkFactory interface {
	// Create creates a new proof-of-work instance for the given block.
	Produce(block *block.Block) (hash block.Hash, powData []byte)
	// Validate validates the proof-of-work for the given block.
	Validate(block *block.Block) bool
}
