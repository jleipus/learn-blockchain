package mock

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/utils"
)

type mockPoWFactory struct {
	counter int32
}

func NewPoWFactory() blockchain.ProofOfWorkFactory {
	return &mockPoWFactory{}
}

func (m *mockPoWFactory) Produce(_ *blockchain.Block) (blockchain.BlockHash, []byte) {
	// Mock implementation: return the counter as the hash and empty powData
	m.counter++
	counterHex, err := utils.IntToHex(m.counter)
	if err != nil {
		panic(err)
	}

	var hash blockchain.BlockHash
	copy(hash[:], counterHex)
	return hash, []byte{}
}

func (m *mockPoWFactory) Validate(_ *blockchain.Block) bool {
	// Mock implementation: always return true
	return true
}
