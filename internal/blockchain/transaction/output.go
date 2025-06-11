package transaction

import (
	"bytes"
	"fmt"

	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
)

// TxOutput represents an output in a transaction.
type TxOutput struct {
	// Value is the amount of cryptocurrency being transferred.
	Value int32
	// PubKeyHash is the hash of the public key that can unlock this output.
	PubKeyHash []byte
}

// NewTxOutput create a new TxOutput.
func NewTxOutput(value int32, address string) TxOutput {
	txo := TxOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	txo.Lock([]byte(address))

	return txo
}

// Lock locks the output to a specific address.
func (out *TxOutput) Lock(address []byte) error {
	pubKeyHash, err := wallet.HashPubKey(address)
	if err != nil {
		return fmt.Errorf("failed to hash public key: %w", err)
	}

	out.PubKeyHash = pubKeyHash[1 : len(pubKeyHash)-wallet.ChecksumLength]
	return nil
}

// IsLockedWithKey checks if the output is locked with the specified public key hash.
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}
