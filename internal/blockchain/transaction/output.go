package transaction

import (
	"bytes"
	"fmt"
)

// TxOutput represents an output in a transaction.
type TxOutput struct {
	// Value is the amount of cryptocurrency being transferred.
	Value int32
	// PubKeyHash is the hash of the public key that can unlock this output.
	PubKeyHash []byte
}

// Lock locks the output to a specific address.
func (out *TxOutput) Lock(address []byte) error {
	pubKeyHash, err := hashPubKey(address)
	if err != nil {
		return fmt.Errorf("failed to hash public key: %w", err)
	}

	out.PubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	return nil
}

// IsLockedWithKey checks if the output is locked with the specified public key hash.
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTxOutput create a new TxOutput.
func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}
