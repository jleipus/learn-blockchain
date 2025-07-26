package badger

import (
	"errors"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/block"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
)

const (
	blocksPrefix  = "blocks_"
	walletsPrefix = "wallets_"
	tipKey        = "tip"
)

type badgerStorage struct {
	db *badger.DB
}

func NewStorage(path string) (blockchain.Storage, error) {
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

func (bs *badgerStorage) GetTip() (block.Hash, error) {
	tip, err := bs.blocksGet([]byte(tipKey))
	if errors.Is(err, badger.ErrKeyNotFound) {
		return block.Hash{}, nil
	}

	if err != nil {
		return block.Hash{}, err
	}

	var hash block.Hash
	copy(hash[:], tip)
	return hash, nil
}

func (bs *badgerStorage) SetTip(hash block.Hash) error {
	return bs.blocksSet([]byte(tipKey), hash[:])
}

func (bs *badgerStorage) GetBlock(hash block.Hash) (*block.Block, error) {
	blockData, err := bs.blocksGet(hash[:])
	if err != nil {
		return nil, err
	}

	block := &block.Block{}
	err = block.Deserialize(blockData)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (bs *badgerStorage) AddBlock(block block.Block) error {
	blockData := block.Serialize()
	hash := block.Hash
	return bs.blocksSet(hash[:], blockData)
}

func (bs *badgerStorage) AddWallet(address string, wallet wallet.Wallet) error {
	walletData, err := wallet.Serialize()
	if err != nil {
		return err
	}

	return bs.walletsSet([]byte(address), walletData)
}

func (bs *badgerStorage) GetAddresses() []string {
	addresses := make([]string, 0)

	err := bs.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(walletsPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()
			if len(key) > len(prefix) {
				address := string(key[len(prefix):])
				addresses = append(addresses, address)
			}
		}
		return nil
	})
	if err != nil {
		return nil
	}

	return addresses
}

func (bs *badgerStorage) GetWallet(address string) (*wallet.Wallet, error) {
	walletData, err := bs.walletsGet([]byte(address))
	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, errors.New("wallet not found")
	}

	if err != nil {
		return nil, err
	}

	w := &wallet.Wallet{}
	err = w.Deserialize(walletData)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (bs *badgerStorage) Close() error {
	return bs.db.Close()
}

func (bs *badgerStorage) blocksGet(key []byte) ([]byte, error) {
	return bs.get(append([]byte(blocksPrefix), key...))
}

func (bs *badgerStorage) blocksSet(key, value []byte) error {
	return bs.set(append([]byte(blocksPrefix), key...), value)
}

func (bs *badgerStorage) walletsGet(key []byte) ([]byte, error) {
	return bs.get(append([]byte(walletsPrefix), key...))
}

func (bs *badgerStorage) walletsSet(key, value []byte) error {
	return bs.set(append([]byte(walletsPrefix), key...), value)
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
