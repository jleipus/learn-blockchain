package mock

import (
	"errors"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
)

type mockStorage struct {
	tip    [32]byte
	blocks map[[32]byte]*blockchain.Block
}

func NewStorage() blockchain.Storage {
	return &mockStorage{
		tip:    [32]byte{},
		blocks: make(map[[32]byte]*blockchain.Block),
	}
}

func (m *mockStorage) GetTip() (blockchain.BlockHash, error) {
	return m.tip, nil
}

func (m *mockStorage) SetTip(tip blockchain.BlockHash) error {
	m.tip = tip
	return nil
}

func (m *mockStorage) GetBlock(hash blockchain.BlockHash) (*blockchain.Block, error) {
	block, exists := m.blocks[hash]
	if !exists {
		return nil, errors.New("block not found")
	}
	return block, nil
}

func (m *mockStorage) AddBlock(block *blockchain.Block) error {
	m.blocks[block.Hash] = block
	return nil
}

func (m *mockStorage) Close() error {
	return nil
}
