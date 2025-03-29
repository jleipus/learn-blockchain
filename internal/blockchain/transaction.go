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

type Transaction struct {
	ID   [32]byte
	Vin  []*TXInput
	Vout []*TXOutput
}

type TXInput struct {
	TxID      []byte
	Vout      int32
	ScriptSig string
}

type TXOutput struct {
	Value        int32
	ScriptPubKey string
}

func NewCoinbaseTX(to, data string) (*Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{
		TxID:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := TXOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}
	tx := Transaction{
		Vin:  []*TXInput{&txin},
		Vout: []*TXOutput{&txout},
	}
	err := tx.setID()
	if err != nil {
		return nil, fmt.Errorf("failed to set transaction ID: %w", err)
	}

	return &tx, nil
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

func (in *TXInput) CanUnlockOutputWith(data string) bool {
	return in.ScriptSig == data
}

func (out *TXOutput) CanBeUnlockedWith(data string) bool {
	return out.ScriptPubKey == data
}
