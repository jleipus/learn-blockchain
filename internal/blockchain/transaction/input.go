package transaction

import (
	"bytes"
	"fmt"

	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
)

// TxInput represents an input of a transaction.
type TxInput struct {
	// TxID is the ID of the transaction that created the output being spent.
	TxID TxID
	// Vout is the index of the output in the previous transaction.
	Vout int
	// Signature is the signature that unlocks the output being spent.
	Signature []byte
	// PubKey is the public key that can unlock the output being spent.
	PubKey []byte
}

// UsesKey checks if the input uses the specified public key hash to unlock the output.
func (in *TxInput) UsesKey(pubKeyHash []byte) (bool, error) {
	lockingHash, err := wallet.HashPubKey(in.PubKey)
	if err != nil {
		return false, fmt.Errorf("failed to hash public key: %w", err)
	}

	return bytes.Equal(lockingHash, pubKeyHash), nil
}
