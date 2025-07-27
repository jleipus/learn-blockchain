package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"iter"
	"slices"
	"time"

	"github.com/jleipus/learn-blockchain/internal/blockchain/block"
	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
	"github.com/jleipus/learn-blockchain/internal/blockchain/utxo"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
)

const (
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

// Blockchain represents a blockchain structure that holds blocks and manages transactions.
type Blockchain struct {
	storage    block.Storage
	powFactory ProofOfWorkFactory
	wallets    *wallet.Collection
	utxoSet    *utxo.UTXOSet
}

// CreateBlockchain initializes a new blockchain with a genesis block.
// It requires a storage implementation to persist the blockchain data and a proof-of-work factory to create the genesis block.
func CreateBlockchain(storage block.Storage, powFactory ProofOfWorkFactory, address string) error {
	tip, err := storage.GetTip()
	if err != nil {
		return fmt.Errorf("failed to get tip of blockchain: %w", err)
	}

	if tip != *new(block.Hash) {
		return errors.New("blockchain already exists")
	}

	cbtx, err := transaction.NewCoinbaseTX(address, genesisCoinbaseData)
	if err != nil {
		return fmt.Errorf("failed to create coinbase transaction: %w", err)
	}

	genesis := newGenesisBlock(cbtx, powFactory)

	err = storage.AddBlock(*genesis)
	if err != nil {
		return fmt.Errorf("failed to add genesis block: %w", err)
	}

	err = storage.SetTip(genesis.Hash)
	if err != nil {
		return fmt.Errorf("failed to set tip of blockchain: %w", err)
	}

	return nil
}

// LoadBlockchain loads an existing blockchain from storage.
func LoadBlockchain(
	storage Storage,
	powFactory ProofOfWorkFactory,
	wallets *wallet.Collection,
) (*Blockchain, error) {
	tip, err := storage.GetTip()
	if err != nil {
		return nil, fmt.Errorf("failed to get tip of blockchain: %w", err)
	}

	if tip == *new(block.Hash) {
		return nil, errors.New("blockchain does not exist")
	}

	return &Blockchain{
		storage:    storage,
		powFactory: powFactory,
		wallets:    wallets,
		utxoSet:    utxo.NewUTXOSet(storage),
	}, nil
}

// newBlock creates a new block with the given transactions and previous block hash.
func newBlock(
	transactions []*transaction.Tx,
	prevBlockHash block.Hash,
	powFactory ProofOfWorkFactory,
) *block.Block {
	block := &block.Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}

	hash, powData := powFactory.Produce(block)

	block.Hash = hash
	block.PoW = powData

	return block
}

// newGenesisBlock creates a new genesis block with the given coinbase transaction.
func newGenesisBlock(coinbase *transaction.Tx, powFactory ProofOfWorkFactory) *block.Block {
	return newBlock([]*transaction.Tx{coinbase}, block.Hash{}, powFactory)
}

// GetBlock retrieves a block by its hash from the blockchain storage.
func (bc *Blockchain) GetBlock(hash block.Hash) (*block.Block, error) {
	return bc.storage.GetBlock(hash)
}

// MineBlock mines a new block with the provided transactions and adds it to the blockchain.
func (bc *Blockchain) MineBlock(transactions []*transaction.Tx) (*block.Block, error) {
	for _, tx := range transactions {
		ok, err := bc.verifyTransaction(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to verify transaction %x: %w", tx.ID, err)
		}
		if !ok {
			return nil, fmt.Errorf("invalid transaction: %x", tx.ID)
		}
	}

	tip, err := bc.storage.GetTip()
	if err != nil {
		return nil, fmt.Errorf("failed to get tip of blockchain: %w", err)
	}

	if tip == *new(block.Hash) {
		return nil, errors.New("tip is empty")
	}

	b := newBlock(transactions, tip, bc.powFactory)

	err = bc.storage.AddBlock(*b)
	if err != nil {
		return nil, fmt.Errorf("failed to add block: %w", err)
	}

	err = bc.storage.SetTip(b.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to set tip of blockchain: %w", err)
	}

	return b, nil
}

// findTransaction finds a transaction by its ID.
func (bc *Blockchain) findTransaction(id transaction.TxID) (*transaction.Tx, error) {
	for _, block := range bc.Blocks() {
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID[:], id[:]) {
				return tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return nil, errors.New("transaction not found")
}

// signTransaction signs inputs of a Transaction.
func (bc *Blockchain) signTransaction(tx *transaction.Tx, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[transaction.TxID]*transaction.Tx)

	for _, vin := range tx.Vin {
		prevTX, err := bc.findTransaction(vin.TxID)
		if err != nil {
			panic(err)
		}
		prevTXs[prevTX.ID] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// verifyTransaction verifies transaction input signatures.
func (bc *Blockchain) verifyTransaction(tx *transaction.Tx) (bool, error) {
	if tx.IsCoinbase() {
		return true, nil // Coinbase transactions are always valid
	}

	prevTXs := make(map[transaction.TxID]*transaction.Tx)
	for _, vin := range tx.Vin {
		prevTX, err := bc.findTransaction(vin.TxID)
		if err != nil {
			return false, fmt.Errorf("failed to find previous transaction %x: %w", vin.TxID, err)
		}
		prevTXs[prevTX.ID] = prevTX
	}

	return tx.Verify(prevTXs), nil
}

// NewUTXOTransaction creates a new transaction with unspent transaction outputs (UTXO).
func (bc *Blockchain) NewUTXOTransaction(fromAddress, toAddress string, amount int32) (*transaction.Tx, error) {
	wlt, err := bc.wallets.GetWallet(fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet for address %s: %w", fromAddress, err)
	}

	pubKeyHash, err := wallet.HashPubKey(wlt.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to hash public key: %w", err)
	}

	acc, validOutputs, err := bc.utxoSet.FindSpendableOutputIndexes(pubKeyHash, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to find spendable outputs: %w", err)
	}
	if acc < amount {
		return nil, fmt.Errorf("not enough funds: %d < %d", acc, amount)
	}

	// Build a list of inputs
	var inputs []transaction.TxInput
	for txID, outs := range validOutputs {
		for _, out := range outs {
			input := transaction.TxInput{
				TxID:      txID,
				Vout:      out,
				Signature: nil, // This will be filled later with the signature
				PubKey:    wlt.PublicKey,
			}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	var outputs []transaction.TxOutput
	outputs = append(outputs, transaction.NewTxOutput(amount, toAddress))
	if acc > amount {
		outputs = append(outputs, transaction.NewTxOutput(acc-amount, fromAddress)) // The change
	}

	tx := transaction.Tx{
		ID:   transaction.TxID{}, // This will be filled later with the hash
		Vin:  inputs,
		Vout: outputs,
	}
	tx.ID = tx.Hash()

	bc.signTransaction(&tx, wlt.PrivateKey)

	return &tx, nil
}

// findUnspentTxOutputs returns a list of transaction outputs.
func (bc *Blockchain) findUnspentTxOutputs() map[transaction.TxID][]transaction.TxOutput {
	unspentTxOs := make(map[transaction.TxID][]transaction.TxOutput)
	spentTxOs := make(map[transaction.TxID][]int) // Map of transaction ID to slice of output indexes

	for _, b := range bc.Blocks() {
		for _, tx := range b.Transactions {
			for outID, out := range tx.Vout {
				// Was the output spent?
				if outputIDs, ok := spentTxOs[tx.ID]; ok {
					if slices.Contains(outputIDs, outID) {
						continue // Skip this output if it was spent
					}
				}

				unspentTxOs[tx.ID] = append(unspentTxOs[tx.ID], out)
			}

			if tx.IsCoinbase() {
				continue
			}

			for _, in := range tx.Vin {
				spentTxOs[in.TxID] = append(spentTxOs[in.TxID], in.Vout)
			}
		}
	}

	return unspentTxOs
}

func (bc *Blockchain) ReindexUTXOSet() error {
	utxoSet := bc.findUnspentTxOutputs()
	if err := bc.utxoSet.Set(utxoSet); err != nil {
		return fmt.Errorf("failed to set UTXO set: %w", err)
	}
	return nil
}

func (bc *Blockchain) Update(b block.Block) error {
	return bc.utxoSet.Update(b)
}

func (bc *Blockchain) FindUnspentTxOutputs(pubKeyHash []byte) ([]transaction.TxOutput, error) {
	return bc.utxoSet.FindUnspentTxOutputs(pubKeyHash)
}

type blockchainIterator struct {
	currentIndex int
	currentHash  block.Hash
	bc           *Blockchain
}

func (bi *blockchainIterator) next() (int, *block.Block, bool) {
	currentBlock, err := bi.bc.GetBlock(bi.currentHash)
	if err != nil || currentBlock == nil {
		return 0, nil, false
	}

	bi.currentHash = currentBlock.PrevBlockHash
	defer func() { bi.currentIndex++ }()

	return bi.currentIndex, currentBlock, true
}

func (bc *Blockchain) Blocks() iter.Seq2[int, *block.Block] {
	return func(yield func(int, *block.Block) bool) {
		tip, err := bc.storage.GetTip()
		if err != nil {
			panic(fmt.Errorf("failed to get tip of blockchain: %w", err))
		}

		if tip == *new(block.Hash) {
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
