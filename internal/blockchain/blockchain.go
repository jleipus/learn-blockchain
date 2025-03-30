package blockchain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"iter"
	"slices"
	"time"
)

const (
	blocksBucket        = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

type Blockchain struct {
	storage    Storage
	powFactory ProofOfWorkFactory
}

func CreateBlockchain(storage Storage, powFactory ProofOfWorkFactory, address string) error {
	tip, err := storage.GetTip()
	if err != nil {
		return fmt.Errorf("failed to get tip of blockchain: %w", err)
	}

	if tip != [32]byte{} {
		return errors.New("blockchain already exists")
	}

	cbtx, err := NewCoinbaseTX(address, genesisCoinbaseData)
	if err != nil {
		return fmt.Errorf("failed to create coinbase transaction: %w", err)
	}

	genesis := newGenesisBlock(cbtx, powFactory)

	err = storage.AddBlock(genesis)
	if err != nil {
		return fmt.Errorf("failed to add genesis block: %w", err)
	}

	err = storage.SetTip(genesis.Hash)
	if err != nil {
		return fmt.Errorf("failed to set tip of blockchain: %w", err)
	}

	return nil
}

func LoadBlockchain(storage Storage, powFactory ProofOfWorkFactory) (*Blockchain, error) {
	tip, err := storage.GetTip()
	if err != nil {
		return nil, fmt.Errorf("failed to get tip of blockchain: %w", err)
	}

	if tip == [32]byte{} {
		return nil, errors.New("blockchain does not exist")
	}

	return &Blockchain{
		storage:    storage,
		powFactory: powFactory,
	}, nil
}

func newBlock(transactions []*Transaction, prevBlockHash BlockHash, powFactory ProofOfWorkFactory) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}

	hash, powData := powFactory.Produce(block)

	block.Hash = hash
	block.PoW = powData

	return block
}

func newGenesisBlock(coinbase *Transaction, powFactory ProofOfWorkFactory) *Block {
	return newBlock([]*Transaction{coinbase}, BlockHash{}, powFactory)
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) error {
	tip, err := bc.storage.GetTip()
	if err != nil {
		return fmt.Errorf("failed to get tip of blockchain: %w", err)
	}

	if tip == [32]byte{} {
		return errors.New("tip is empty")
	}

	b := newBlock(transactions, tip, bc.powFactory)

	err = bc.storage.AddBlock(b)
	if err != nil {
		return fmt.Errorf("failed to add block: %w", err)
	}

	err = bc.storage.SetTip(b.Hash)
	if err != nil {
		return fmt.Errorf("failed to set tip of blockchain: %w", err)
	}

	return nil
}

func (bc *Blockchain) GetBlock(hash BlockHash) (*Block, error) {
	return bc.storage.GetBlock(hash)
}

func (bc *Blockchain) findUnspentTransactions(address string) []*Transaction {
	var unspentTxs []*Transaction
	spentTxOs := make(map[string][]int) // Map of transaction ID to slice of output indexes

	for _, b := range bc.Blocks() {
		for _, tx := range b.Transactions {
			txID := hex.EncodeToString(tx.ID[:])

			for outID, out := range tx.Vout {
				// Was the output spent?
				if outputIDs, ok := spentTxOs[txID]; ok {
					if slices.Contains(outputIDs, outID) {
						continue // Skip this output if it was spent
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTxs = append(unspentTxs, tx)
				}
			}

			if tx.IsCoinbase() {
				continue
			}

			for _, in := range tx.Vin {
				if in.CanUnlockOutputWith(address) {
					inTxID := hex.EncodeToString(in.TxID[:])
					spentTxOs[inTxID] = append(spentTxOs[inTxID], in.Vout)
				}
			}
		}
	}

	return unspentTxs
}

func (bc *Blockchain) FindUnspentTxOutputs(address string) []*TxOutput {
	utxos := []*TxOutput{}
	unspentTxs := bc.findUnspentTransactions(address)

	for _, tx := range unspentTxs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				utxos = append(utxos, out)
			}
		}
	}

	return utxos
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int32) (int32, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTxs := bc.findUnspentTransactions(address)
	var accumulated int32

	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID[:])

		for outIdx, out := range tx.Vout {
			if accumulated >= amount {
				return accumulated, unspentOutputs
			}

			if out.CanBeUnlockedWith(address) {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
		}
	}

	return accumulated, unspentOutputs
}

func (bc *Blockchain) NewUnspentTx(from, to string, amount int32) (*Transaction, error) {
	accumulated, validOutputs := bc.FindSpendableOutputs(from, amount)

	if accumulated < amount {
		return nil, fmt.Errorf("not enough funds: %d < %d", accumulated, amount)
	}

	// Build a list of inputs
	inputs := []*TxInput{}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			return nil, fmt.Errorf("failed to decode txid %s: %w", txid, err)
		}

		for _, out := range outs {
			inputs = append(inputs, &TxInput{
				TxID:      [32]byte(txID),
				Vout:      out,
				ScriptSig: from,
			})
		}
	}

	// Build a list of outputs
	outputs := []*TxOutput{}
	outputs = append(outputs, &TxOutput{Value: amount, ScriptPubKey: to}) // Send to recipient
	if accumulated > amount {
		outputs = append(outputs, &TxOutput{Value: accumulated - amount, ScriptPubKey: from}) // Send change back to sender
	}

	return NewTransaction(inputs, outputs)
}

type blockchainIterator struct {
	currentIndex int
	currentHash  BlockHash
	bc           *Blockchain
}

func (bi *blockchainIterator) next() (int, *Block, bool) {
	currentBlock, err := bi.bc.GetBlock(bi.currentHash)
	if err != nil || currentBlock == nil {
		return 0, nil, false
	}

	bi.currentHash = currentBlock.PrevBlockHash
	defer func() { bi.currentIndex++ }()

	return bi.currentIndex, currentBlock, true
}

func (bc *Blockchain) Blocks() iter.Seq2[int, *Block] {
	return func(yield func(int, *Block) bool) {
		tip, err := bc.storage.GetTip()
		if err != nil {
			panic(fmt.Errorf("failed to get tip of blockchain: %w", err))
		}

		if tip == [32]byte{} {
			return
		}

		bi := &blockchainIterator{
			currentIndex: 0,
			currentHash:  tip,
			bc:           bc,
		}

		for {
			index, b, ok := bi.next()
			if !ok {
				break
			}

			if !yield(index, b) {
				break
			}
		}
	}
}
