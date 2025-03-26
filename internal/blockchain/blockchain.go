package blockchain

import (
	"fmt"
	"iter"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jleipus/learn-blockchain/internal/transaction"
	pb "github.com/jleipus/learn-blockchain/proto"
	"google.golang.org/protobuf/proto"
)

const (
	blocksBucket        = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

type Blockchain struct {
	tip BlockHash
	db  *bolt.DB
	pf  ProofOfWorkFactory
}

func New(dbFile string, address string, powFactory ProofOfWorkFactory) (*Blockchain, error) {
	var tip BlockHash

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			cbtx := transaction.NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := newGenesisBlock(cbtx, powFactory)

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

			tip = BlockHash(genesis.GetHash())
		} else {
			tip = BlockHash(b.Get([]byte("l")))
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if tip == [32]byte{} {
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

func newBlock(transactions []*pb.Transaction, prevBlockHash BlockHash, powFactory ProofOfWorkFactory) *pb.BlockEntity {
	block := &pb.BlockEntity{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash[:],
	}

	hash, powData := powFactory.Produce(block)

	block.Hash = hash[:]
	block.PoW = powData

	return block
}

func newGenesisBlock(coinbase *pb.Transaction, powFactory ProofOfWorkFactory) *pb.BlockEntity {
	return newBlock([]*pb.Transaction{coinbase}, BlockHash{}, powFactory)
}

func (bc *Blockchain) AddBlock(transactions []*pb.Transaction) error {
	var lastHash BlockHash

	if err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = BlockHash(b.Get([]byte("l")))

		return nil
	}); err != nil {
		return err
	}

	b := newBlock(transactions, lastHash, bc.pf)

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

		bc.tip = BlockHash(b.GetHash())

		return nil
	})
}

func (bc *Blockchain) GetBlock(hash BlockHash) *pb.BlockEntity {
	b := &pb.BlockEntity{}

	if err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		encodedBlock := bucket.Get(hash[:])
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
	currentHash  BlockHash
	bc           *Blockchain
}

func (bi *blockchainIterator) next() (int, *pb.BlockEntity, bool) {
	currentBlock := bi.bc.GetBlock(bi.currentHash)
	if currentBlock == nil {
		return 0, nil, false
	}

	bi.currentHash = BlockHash(currentBlock.GetPrevBlockHash())
	defer func() { bi.currentIndex++ }()

	return bi.currentIndex, currentBlock, true
}

func (bc *Blockchain) Blocks() iter.Seq2[int, *pb.BlockEntity] {
	return func(yield func(int, *pb.BlockEntity) bool) {
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
