package wallet_test

import (
	"testing"

	"github.com/jleipus/learn-blockchain/internal/blockchain/mock"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWallet(t *testing.T) {
	wallets := wallet.NewCollection(mock.NewStorage())

	address, err := wallets.AddWallet()
	require.NoError(t, err)

	pubKeyHash, err := wallet.GetHashFromAddress([]byte(address))
	require.NoError(t, err)

	wlt, err := wallets.GetWallet(address)
	require.NoError(t, err)
	wltHash, err := wallet.HashPubKey(wlt.PublicKey)
	require.NoError(t, err)

	assert.Equal(t, pubKeyHash, wltHash)
}

func TestSerializeDeserialize(t *testing.T) {
	wlt, err := wallet.New()
	require.NoError(t, err)

	serialized, err := wlt.Serialize()
	require.NoError(t, err)
	require.NotEmpty(t, serialized)

	var deserialized wallet.Wallet
	err = deserialized.Deserialize(serialized)
	require.NoError(t, err)

	assert.Equal(t, wlt, &deserialized)
}
