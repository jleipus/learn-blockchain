package blockchain_test

import (
	"testing"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/mock"
	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	mockStorage := mock.NewStorage()
	mockPowFactory := mock.NewPoWFactory()

	wallets := wallet.NewCollection(mockStorage)

	address1, err := wallets.AddWallet()
	require.NoError(t, err)
	wlt1, err := wallets.GetWallet(address1)
	require.NoError(t, err)
	address2, err := wallets.AddWallet()
	require.NoError(t, err)
	wlt2, err := wallets.GetWallet(address2)
	require.NoError(t, err)
	address3, err := wallets.AddWallet()
	require.NoError(t, err)
	wlt3, err := wallets.GetWallet(address3)
	require.NoError(t, err)

	err = blockchain.CreateBlockchain(mockStorage, mockPowFactory, address1)
	require.NoError(t, err)
	bc, err := blockchain.LoadBlockchain(mockStorage, mockPowFactory, wallets)
	require.NoError(t, err)

	tx1, err := bc.NewUTXOTransaction(address1, address2, 10)
	require.NoError(t, err)
	tx2, err := bc.NewUTXOTransaction(address2, address3, 5)
	require.NoError(t, err)

	err = bc.MineBlock([]*transaction.TX{tx1})
	require.NoError(t, err)
	err = bc.MineBlock([]*transaction.TX{tx2})
	require.NoError(t, err)

	amount1, outputs1 := bc.FindSpendableOutputs(wlt1.PublicKey, 0)
	require.Len(t, outputs1, 1)
	amount2, outputs2 := bc.FindSpendableOutputs(wlt2.PublicKey, 5)
	require.Len(t, outputs2, 1)
	amount3, outputs3 := bc.FindSpendableOutputs(wlt3.PublicKey, 5)
	require.Len(t, outputs3, 1)

	require.Equal(t, 10, amount1)
	require.Equal(t, 5, amount2)
	require.Equal(t, 0, amount3)
}
