package blockchain_test

import (
	"testing"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindUnspentTransactions(t *testing.T) {
	mockStorage := mock.NewStorage()
	mockPowFactory := mock.NewPoWFactory()
	address := "test-address"

	// Create a blockchain with a genesis block
	err := blockchain.CreateBlockchain(mockStorage, mockPowFactory, address)
	require.NoError(t, err)

	bc, err := blockchain.LoadBlockchain(mockStorage, mockPowFactory)
	require.NoError(t, err)

	// Add a block with two transactions
	tx1 := &blockchain.Transaction{
		ID: [32]byte{1},
		Vout: []*blockchain.TxOutput{
			{Value: 10, ScriptPubKey: address},
		},
	}
	tx2 := &blockchain.Transaction{
		ID: [32]byte{2},
		Vout: []*blockchain.TxOutput{
			{Value: 20, ScriptPubKey: "other-address"},
		},
	}
	err = bc.MineBlock([]*blockchain.Transaction{tx1, tx2})
	require.NoError(t, err)

	// Add another block with a transaction spending tx1's output
	tx3 := &blockchain.Transaction{
		ID: [32]byte{3},
		Vin: []*blockchain.TxInput{
			{TxID: [32]byte{1}, Vout: 0, ScriptSig: address},
		},
		Vout: []*blockchain.TxOutput{
			{Value: 5, ScriptPubKey: address},
			{Value: 5, ScriptPubKey: "other-address"},
		},
	}
	err = bc.MineBlock([]*blockchain.Transaction{tx3})
	require.NoError(t, err)

	unspentTxOutputs := bc.FindUnspentTxOutputs(address)

	require.Len(t, unspentTxOutputs, 2)

	assert.Equal(t, int32(5), unspentTxOutputs[0].Value)
	assert.Equal(t, address, unspentTxOutputs[0].ScriptPubKey)
	assert.Equal(t, int32(10), unspentTxOutputs[1].Value)
	assert.Equal(t, "test-address", unspentTxOutputs[1].ScriptPubKey)
}
