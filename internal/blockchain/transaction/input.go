package transaction

import (
	"bytes"
	"fmt"
)

// TxInput represents an input in a transaction.
type TxInput struct {
	// TxID is the ID of the transaction that created the output being spent.
	TxID [32]byte
	// Vout is the index of the output in the previous transaction.
	Vout int
	//
	Signature []byte
	// PubKey is the public key that can unlock the output being spent.
	PubKey []byte
}

// UsesKey checks if the input uses the specified public key hash to unlock the output.
func (in *TxInput) UsesKey(pubKeyHash []byte) (bool, error) {
	lockingHash, err := hashPubKey(in.PubKey)
	if err != nil {
		return false, fmt.Errorf("failed to hash public key: %w", err)
	}

	return bytes.Compare(lockingHash, pubKeyHash) == 0, nil
}
