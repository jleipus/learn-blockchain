package blockchain

import (
	"fmt"
	"iter"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jleipus/learn-blockchain/internal/pow"
	"google.golang.org/protobuf/proto"
)

const (
	blocksBucket = "blocks"
	genesisData  = "Genesis Block"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
	pf  pow.ProofOfWorkFactory
}

func New(dbFile string, powFactory pow.ProofOfWorkFactory) (*Blockchain, error) {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := newGenesisBlock(powFactory)

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return fmt.Errorf("failed to create bucket: %w", err)
			}

			bytes, err := proto.Marshal(genesis)
			if err != nil {
				return fmt.Errorf("failed to serialize genesis block: %w", err)
			}

			err = b.Put(genesis.GetHash(), bytes)
			if err != nil {
				return fmt.Errorf("failed to put genesis block: %w", err)
			}

			err = b.Put([]byte("l"), genesis.GetHash())
			if err != nil {
				return fmt.Errorf("failed to put last hash: %w", err)
			}

			tip = genesis.GetHash()
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if tip == nil {
		return nil, fmt.Errorf("tip is nil")
	}

	return &Blockchain{
		tip: tip,
		db:  db,
		pf:  powFactory,
	}, nil
}

func (bc *Blockchain) Close() {
	bc.db.Close()
}

func newBlock(data string, prevBlockHash []byte, powFactory pow.ProofOfWorkFactory) pow.Block {
	block := &block.Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	nonce, hash := powFactory.Produce(block)

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func newGenesisBlock(powFactory pow.ProofOfWorkFactory) *block.Block {
	return newBlock(genesisData, []byte{}, powFactory)
}

func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte

	if err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	}); err != nil {
		return err
	}

	b := newBlock(data, lastHash, bc.pf)

	return bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		bytes, err := proto.Marshal(b)
		if err != nil {
			return fmt.Errorf("failed to serialize block: %w", err)
		}

		if err := bucket.Put(b.GetHash(), bytes); err != nil {
			return fmt.Errorf("failed to put block: %w", err)
		}

		if err := bucket.Put([]byte("l"), b.GetHash()); err != nil {
			return fmt.Errorf("failed to put last hash: %w", err)
		}

		bc.tip = b.GetHash()

		return nil
	})
}

func (bc *Blockchain) GetBlock(hash []byte) *block.Block {
	b := &block.Block{}

	if err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		encodedBlock := bucket.Get(hash)
		if encodedBlock == nil {
			return fmt.Errorf("block not found")
		}

		var err error
		err = proto.Unmarshal(encodedBlock, b)
		if err != nil {
			return fmt.Errorf("failed to deserialize block: %w", err)
		}

		return nil
	}); err != nil {
		return nil
	}

	return b
}

type blockchainIterator struct {
	currentIndex int
	currentHash  []byte
	bc           *Blockchain
}

func (bi *blockchainIterator) next() (int, *block.Block, bool) {
	currentBlock := bi.bc.GetBlock(bi.currentHash)
	if currentBlock == nil {
		return 0, nil, false
	}

	bi.currentHash = currentBlock.GetPrevBlockHash()
	defer func() { bi.currentIndex++ }()

	return bi.currentIndex, currentBlock, true
}

func (bc *Blockchain) Blocks() iter.Seq2[int, *block.Block] {
	return func(yield func(int, *block.Block) bool) {
		bi := &blockchainIterator{
			currentIndex: 0,
			currentHash:  bc.tip,
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
