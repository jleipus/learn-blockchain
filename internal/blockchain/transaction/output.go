package transaction

import (
	"bytes"

	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
)

// TxOutput represents an output in a transaction.
type TxOutput struct {
	// Value is the amount of cryptocurrency being transferred.
	Value int32
	// PubKeyHash is the hash of the public key that locked and can unlock this output.
	PubKeyHash []byte
}

// NewTxOutput create a new TxOutput.
func NewTxOutput(value int32, address string) TxOutput {
	txo := TxOutput{
		Value:      value,
		PubKeyHash: nil, // Will be set when locking the output
	}
	if err := txo.lock([]byte(address)); err != nil {
		panic(err) // Expect address to be valid
	}

	return txo
}

// lock locks the output to a specific address by setting the PubKeyHash.
// It returns an error if the address is invalid.
func (out *TxOutput) lock(address []byte) error {
	pubKeyHash, err := wallet.GetHashFromAddress(address)
	if err != nil {
		return err
	}
	out.PubKeyHash = pubKeyHash
	return nil
}

// IsLockedWithKey checks if the output is locked with the specified public key hash.
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}
