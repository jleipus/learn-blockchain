package main

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain/badger"
	"github.com/jleipus/learn-blockchain/internal/blockchain/hashcash"
	"github.com/jleipus/learn-blockchain/internal/cli"
)

func main() {
	hashcash.SetVerbose()
	// powFactory := hashcash.New()

	storage, err := badger.NewStorage("blockchain.db")
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	cli.Execute()
}
