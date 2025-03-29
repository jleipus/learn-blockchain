package blockchain

type BlockHash [32]byte

// Storage is an interface for a storage system that can store and retrieve blocks.
type Storage interface {
	// GetTip retrieves the hash of the last block in the blockchain.
	// It returns an empty hash if the blockchain is empty.
	GetTip() (BlockHash, error)

	// SetTip sets the hash of the last block in the blockchain.
	SetTip(hash BlockHash) error

	// GetBlock retrieves a block by its hash.
	GetBlock(hash BlockHash) (*Block, error)

	// AddBlock adds a new block to the blockchain.
	AddBlock(block *Block) error

	// Close closes the storage connection.
	Close() error
}

// ProofOfWorkFactory is an implementation of a proof-of-work algorithm.
// It is responsible for producing a hash for a block and validating it.
type ProofOfWorkFactory interface {
	// Create creates a new proof-of-work instance for the given block.
	Produce(block *Block) (hash BlockHash, powData []byte)

	// Validate validates the proof-of-work for the given block.
	Validate(block *Block) bool
}
