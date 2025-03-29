package blockchain

import (
	"fmt"
	"iter"
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

func NewBlockchain(storage Storage, address string, powFactory ProofOfWorkFactory) (*Blockchain, error) {
	tip, err := storage.GetTip()
	if err != nil {
		return nil, fmt.Errorf("failed to get tip of blockchain: %w", err)
	}

	if tip == [32]byte{} { // If tip is empty
		cbtx, err := NewCoinbaseTX(address, genesisCoinbaseData)
		if err != nil {
			return nil, fmt.Errorf("failed to create coinbase transaction: %w", err)
		}

		genesis := newGenesisBlock(cbtx, powFactory)

		err = storage.AddBlock(genesis)
		if err != nil {
			return nil, fmt.Errorf("failed to add genesis block: %w", err)
		}

		err = storage.SetTip(genesis.Hash)
		if err != nil {
			return nil, fmt.Errorf("failed to set tip of blockchain: %w", err)
		}
	}

	return &Blockchain{
		storage:    storage,
		powFactory: powFactory,
	}, nil
}

func (bc *Blockchain) Close() {
	bc.storage.Close()
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

func (bc *Blockchain) AddBlock(transactions []*Transaction) error {
	tip, err := bc.storage.GetTip()
	if err != nil {
		return fmt.Errorf("failed to get tip of blockchain: %w", err)
	}

	if tip == [32]byte{} {
		return fmt.Errorf("tip is empty")
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

	bi.currentHash = BlockHash(currentBlock.PrevBlockHash)
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
