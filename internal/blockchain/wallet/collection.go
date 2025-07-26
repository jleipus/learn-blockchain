package wallet

import (
	"fmt"
)

// Storage is an interface for a storage system that can store and retrieve wallets.
type Storage interface {
	// AddWallet adds a new wallet to the storage and returns its address.
	AddWallet(address string, wallet Wallet) error
	// GetAddresses returns a slice of all wallet addresses in the storage.
	GetAddresses() []string
	// GetWallet retrieves a wallet by its address.
	GetWallet(address string) (*Wallet, error)
}

// Collection stores a collection of wallets.
type Collection struct {
	storage Storage
}

// NewCollection creates new Collection.
func NewCollection(storage Storage) *Collection {
	return &Collection{
		storage: storage,
	}
}

// AddWallet adds a Wallet to Collection and returns its address.
func (c *Collection) AddWallet() (string, error) {
	wallet, err := New()
	if err != nil {
		return "", err
	}

	address, err := wallet.getAddress()
	if err != nil {
		return "", err
	}

	addressStr := fmt.Sprintf("%s", address)
	err = c.storage.AddWallet(addressStr, *wallet)
	if err != nil {
		return "", err
	}

	return addressStr, nil
}

// GetAddresses returns an array of addresses stored in the Collection.
func (c *Collection) GetAddresses() []string {
	return c.storage.GetAddresses()
}

// GetWallet returns a Wallet by its address.
func (c Collection) GetWallet(address string) (*Wallet, error) {
	return c.storage.GetWallet(address)
}
