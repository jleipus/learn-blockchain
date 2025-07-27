package mock

import (
	"errors"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/block"
	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
)

type mockStorage struct {
	tip     block.Hash
	blocks  map[block.Hash]block.Block
	wallets map[string]wallet.Wallet
	utxos   map[transaction.TxID][]transaction.TxOutput
}

func NewStorage() blockchain.Storage {
	return &mockStorage{
		tip:     block.Hash{},
		blocks:  make(map[block.Hash]block.Block),
		wallets: make(map[string]wallet.Wallet),
		utxos:   make(map[transaction.TxID][]transaction.TxOutput),
	}
}

func (m *mockStorage) GetTip() (block.Hash, error) {
	return m.tip, nil
}

func (m *mockStorage) SetTip(tip block.Hash) error {
	m.tip = tip
	return nil
}

func (m *mockStorage) GetBlock(hash block.Hash) (*block.Block, error) {
	block, exists := m.blocks[hash]
	if !exists {
		return nil, errors.New("block not found")
	}
	return &block, nil
}

func (m *mockStorage) AddBlock(block block.Block) error {
	m.blocks[block.Hash] = block
	return nil
}

func (m *mockStorage) AddWallet(address string, wallet wallet.Wallet) error {
	m.wallets[address] = wallet
	return nil
}

func (m *mockStorage) GetAddresses() ([]string, error) {
	addresses := make([]string, 0, len(m.wallets))
	for address := range m.wallets {
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func (m *mockStorage) GetWallet(address string) (*wallet.Wallet, error) {
	if wallet, exists := m.wallets[address]; exists {
		return &wallet, nil
	}
	return nil, errors.New("wallet not found")
}

func (m *mockStorage) GetUTXOs() (map[transaction.TxID][]transaction.TxOutput, error) {
	return m.utxos, nil
}

func (m *mockStorage) SetUTXOs(txID transaction.TxID, outputs []transaction.TxOutput) error {
	m.utxos[txID] = outputs
	return nil
}

func (m *mockStorage) Close() error {
	return nil
}
