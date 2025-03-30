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

func (bs *badgerStorage) GetTip() (blockchain.BlockHash, error) {
	tip, err := bs.get([]byte(tipKey))
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

func (bs *badgerStorage) SetTip(hash blockchain.BlockHash) error {
	return bs.set([]byte(tipKey), hash[:])
}

func (bs *badgerStorage) GetBlock(hash blockchain.BlockHash) (*blockchain.Block, error) {
	blockData, err := bs.get(hash[:])
	if err != nil {
		return nil, err
	}

	block, err := blockchain.DeserializeBlock(blockData)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (bs *badgerStorage) AddBlock(block *blockchain.Block) error {
	blockData, err := block.Serialize()
	if err != nil {
		return err
	}

	hash := block.Hash
	err = bs.set(hash[:], blockData)
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
			return txn.Set(key, value)
		})
}
