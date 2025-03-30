package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const (
	subsidy = 10
)

// Transaction represents a transaction in the blockchain.
// It contains a list of inputs and outputs.
// Each transaction has a unique ID, which is a hash of the transaction data.
type Transaction struct {
	// ID is the unique identifier for the transaction.
	ID [32]byte
	// Vin is a slice of inputs for the transaction.
	// Each input references a previous output.
	Vin []*TxInput
	// Vout is a slice of outputs for the transaction.
	// Each output specifies a value and a script that can unlock it.
	Vout []*TxOutput
}

// TxInput represents an input in a transaction.
type TxInput struct {
	// TxID is the ID of the transaction that created the output being spent.
	TxID [32]byte
	// Vout is the index of the output in the previous transaction.
	Vout int
	// ScriptSig is the data that is used in an outputs script to unlock the output.
	ScriptSig string
}

// TxOutput represents an output in a transaction.
type TxOutput struct {
	// Value is the amount of cryptocurrency being transferred.
	Value int32
	// ScriptPubKey is the script that locks the output.
	ScriptPubKey string
}

func NewCoinbaseTX(to, data string) (*Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{
		TxID:      [32]byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := TxOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}
	tx := Transaction{
		Vin:  []*TxInput{&txin},
		Vout: []*TxOutput{&txout},
	}
	err := tx.setID()
	if err != nil {
		return nil, fmt.Errorf("failed to set transaction ID: %w", err)
	}

	return &tx, nil
}

func NewTransaction(inputs []*TxInput, outputs []*TxOutput) (*Transaction, error) {
	tx := &Transaction{
		Vin:  inputs,
		Vout: outputs,
	}
	err := tx.setID()
	if err != nil {
		return nil, fmt.Errorf("failed to set transaction ID: %w", err)
	}

	return tx, nil
}

// IsCoinbase checks whether the transaction is coinbase.
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxID) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) setID() error {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		return err
	}

	tx.ID = sha256.Sum256(encoded.Bytes())
	return nil
}

func (in *TxInput) CanUnlockOutputWith(data string) bool {
	return in.ScriptSig == data
}

func (out *TxOutput) CanBeUnlockedWith(data string) bool {
	return out.ScriptPubKey == data
}
