package badger

import (
	"errors"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/jleipus/learn-blockchain/internal/blockchain"
)

const tipKey = "tip"

type badgerStorage struct {
	db *badger.DB
}

func NewBadgerDB(path string) (blockchain.Storage, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &badgerStorage{
		db: db,
	}, nil
}

func (b *badgerStorage) GetTip() (blockchain.BlockHash, error) {
	tip, err := b.get([]byte(tipKey))
	if errors.Is(err, badger.ErrKeyNotFound) {
		return blockchain.BlockHash{}, nil
	}

	if err != nil {
		return blockchain.BlockHash{}, err
	}

	var hash blockchain.BlockHash
	copy(hash[:], tip)
	return hash, nil
}

func (b *badgerStorage) SetTip(hash blockchain.BlockHash) error {
	return b.set([]byte(tipKey), hash[:])
}

func (b *badgerStorage) GetBlock(hash blockchain.BlockHash) (*blockchain.Block, error) {
	blockData, err := b.get(hash[:])
	if err != nil {
		return nil, err
	}

	block, err := blockchain.DeserializeBlock(blockData)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (b *badgerStorage) AddBlock(block *blockchain.Block) error {
	blockData, err := block.Serialize()
	if err != nil {
		return err
	}

	hash := block.Hash
	err = b.set(hash[:], blockData)
	if err != nil {
		return err
	}

	return nil
}

func (bs *badgerStorage) Close() error {
	return bs.db.Close()
}

func (bs *badgerStorage) get(key []byte) ([]byte, error) {
	value := make([]byte, 0)
	return value, bs.db.View(
		func(tx *badger.Txn) error {
			item, err := tx.Get(key)
			if err != nil {
				return err
			}

			valueCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			value = valueCopy
			return nil
		})
}

func (bs *badgerStorage) set(key, value []byte) error {
	return bs.db.Update(
		func(txn *badger.Txn) error {
			return txn.Set([]byte(key), []byte(value))
		})
}
