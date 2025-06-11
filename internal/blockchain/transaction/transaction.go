package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
)

const (
	subsidy = 10
)

type TxID [32]byte

// TX represents a transaction in the blockchain.
// It contains a list of inputs and outputs.
// Each transaction has a unique ID, which is a hash of the transaction data.
type TX struct {
	// ID is the unique identifier for the transaction.
	ID TxID
	// Vin is a slice of inputs for the transaction.
	// Each input references a previous output.
	Vin []TxInput
	// Vout is a slice of outputs for the transaction.
	// Each output specifies a value and a script that can unlock it.
	Vout []TxOutput
}

func NewCoinbaseTX(to, data string) (*TX, error) {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{
		TxID:      TxID{},
		Vout:      -1,
		Signature: nil,
		PubKey:    []byte(data),
	}
	txout := NewTxOutput(subsidy, to)
	tx := TX{
		ID:   TxID{},
		Vin:  []TxInput{txin},
		Vout: []TxOutput{txout},
	}
	tx.ID = tx.Hash()

	return &tx, nil
}

// IsCoinbase checks whether the transaction is coinbase.
func (tx *TX) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxID) == 0 && tx.Vin[0].Vout == -1
}

// Serialize serializes the transaction into a byte slice using gob encoding.
func (tx TX) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		// Error will only occur if the input contains unsupported types.
		// In this case, we panic because the block should only contain serializable types.
		panic(err)
	}

	return result.Bytes()
}

// Hash returns the hash of the transaction.
func (tx *TX) Hash() TxID {
	var hash TxID

	txCopy := *tx
	txCopy.ID = TxID{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash
}

// Sign signs the transaction using the provided private key.
func (tx *TX) Sign(privKey ecdsa.PrivateKey, prevTXs map[TxID]*TX) error {
	if tx.IsCoinbase() {
		return nil // Coinbase transactions do not require signing
	}

	for _, vin := range tx.Vin {
		if prevTXs[vin.TxID].ID == *new(TxID) {
			panic("invalid previous transaction is not correct")
		}
	}

	txCopy := tx.trimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[vin.TxID]
		txCopy.Vin[inID].Signature = nil                           // Clear the signature for signing
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash // Use the public key hash from the previous output
		txCopy.ID = txCopy.Hash()                                  // Recalculate the transaction ID
		txCopy.Vin[inID].PubKey = nil                              // Clear the public key for signing

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID[:])
		if err != nil {
			return err
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}

	return nil
}

// Verify checks the validity of the transaction against previous transactions.
func (tx *TX) Verify(prevTXs map[TxID]*TX) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[vin.TxID].ID == *new(TxID) {
			panic("invalid previous transaction is not correct")
		}
	}

	txCopy := tx.trimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[vin.TxID]
		txCopy.Vin[inID].Signature = nil                           // Clear the signature for verification
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash // Use the public key hash from the previous output
		txCopy.ID = txCopy.Hash()                                  // Recalculate the transaction ID
		txCopy.Vin[inID].PubKey = nil                              // Clear the public key for verification

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)]) //nolint:mnd // Read the first half of the signature
		s.SetBytes(vin.Signature[(sigLen / 2):]) //nolint:mnd // Read the second half of the signature

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)]) //nolint:mnd // Read the first half of the public key
		y.SetBytes(vin.PubKey[(keyLen / 2):]) //nolint:mnd // Read the second half of the public key

		rawPubKey := ecdsa.PublicKey{
			Curve: curve,
			X:     &x,
			Y:     &y,
		}

		if ecdsa.Verify(&rawPubKey, txCopy.ID[:], &r, &s) == false {
			return false
		}
	}

	return true
}

// trimmedCopy creates a copy of the transaction with inputs and outputs trimmed.
func (tx *TX) trimmedCopy() TX {
	var inputs []TxInput
	var outputs []TxOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TxInput{
			TxID:      vin.TxID,
			Vout:      vin.Vout,
			Signature: nil,
			PubKey:    nil,
		})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TxOutput{
			Value:      vout.Value,
			PubKeyHash: vout.PubKeyHash,
		})
	}

	txCopy := TX{tx.ID, inputs, outputs}

	return txCopy
}
