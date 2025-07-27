package transaction

import (
	"bytes"
	"encoding/gob"

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

// Serialize serializes the output into a byte slice using gob encoding.
func (out *TxOutput) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(out)
	if err != nil {
		// Error will only occur if the input contains unsupported types.
		// In this case, we panic because the output should only contain serializable types.
		panic(err)
	}

	return result.Bytes()
}

// Deserialize deserializes a byte slice into a TxOutput using gob encoding.
func (out *TxOutput) Deserialize(d []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(out)
	if err != nil {
		return err
	}

	return nil
}

func SerializeOutputs(outputs []TxOutput) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(outputs)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DeserializeOutputs(data []byte) ([]TxOutput, error) {
	var outputs []TxOutput
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outputs)
	if err != nil {
		return nil, err
	}
	return outputs, nil
}
