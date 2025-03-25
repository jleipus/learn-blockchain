package pow

type ProofOfWorkFactory interface {
	Produce(block Block) (nonce int64, hash [32]byte)
	Validate(block Block) bool
}

type Block interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
