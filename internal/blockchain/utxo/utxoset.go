package utxo

import (
	"fmt"

	"github.com/jleipus/learn-blockchain/internal/blockchain/block"
	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
)

type Storage interface {
	GetUTXOs() (map[transaction.TxID][]transaction.TxOutput, error)
	SetUTXOs(txID transaction.TxID, outputs []transaction.TxOutput) error
}

type UTXOSet struct {
	storage Storage
}

// NewUTXOSet creates a new UTXOSet with the provided storage.
func NewUTXOSet(storage Storage) *UTXOSet {
	return &UTXOSet{storage: storage}
}

func (u *UTXOSet) Set(utxos map[transaction.TxID][]transaction.TxOutput) error {
	for txID, outputs := range utxos {
		if err := u.storage.SetUTXOs(txID, outputs); err != nil {
			return fmt.Errorf("failed to set UTXOs for transaction %s: %w", txID, err)
		}
	}
	return nil
}

// FindSpendableOutputIndexes finds and returns a map of trasnsaction IDs to their unspent output indexes
// that can be used to spend the specified amount.
func (u *UTXOSet) FindSpendableOutputIndexes(
	pubKeyHash []byte,
	amount int32,
) (int32, map[transaction.TxID][]int, error) {
	utxos, err := u.storage.GetUTXOs()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get UTXOs: %w", err)
	}

	var accumulated int32
	unspentOutputs := make(map[transaction.TxID][]int)
	for txID, txos := range utxos {
		for outIDx, out := range txos {
			if accumulated >= amount {
				return accumulated, unspentOutputs, nil
			}

			if out.IsLockedWithKey(pubKeyHash) {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIDx)
			}
		}
	}

	return accumulated, unspentOutputs, nil
}

// FindUnspentTxOutputs finds and returns all unspent transaction outputs.
func (u *UTXOSet) FindUnspentTxOutputs(pubKeyHash []byte) ([]transaction.TxOutput, error) {
	utxos, err := u.storage.GetUTXOs()
	if err != nil {
		return nil, fmt.Errorf("failed to get UTXOs: %w", err)
	}

	var unspentTxOs []transaction.TxOutput
	for _, txos := range utxos {
		for _, out := range txos {
			if out.IsLockedWithKey(pubKeyHash) {
				unspentTxOs = append(unspentTxOs, out)
			}
		}
	}

	return unspentTxOs, nil
}

func (u *UTXOSet) Update(b block.Block) error {
	utxos, err := u.storage.GetUTXOs()
	if err != nil {
		return fmt.Errorf("failed to get UTXOs: %w", err)
	}

	for _, tx := range b.Transactions {
		// Remove spent outputs
		for _, in := range tx.Vin {
			updatedOutputs := make([]transaction.TxOutput, 0)

			outputs, ok := utxos[in.TxID]
			if !ok {
				continue // No outputs found for this transaction ID
			}

			for outIDx, out := range outputs {
				if outIDx != in.Vout {
					updatedOutputs = append(updatedOutputs, out) // Keep unspent outputs
				}
			}

			err := u.storage.SetUTXOs(in.TxID, updatedOutputs)
			if err != nil {
				return fmt.Errorf("failed to update UTXOs for transaction %s: %w", in.TxID, err)
			}
		}

		// Add new outputs
		newOutputs := make([]transaction.TxOutput, 0)
		for _, out := range tx.Vout {
			newOutputs = append(newOutputs, out)
		}

		err := u.storage.SetUTXOs(tx.ID, newOutputs)
		if err != nil {
			return fmt.Errorf("failed to set UTXOs for transaction %s: %w", tx.ID, err)
		}
	}

	return nil
}
