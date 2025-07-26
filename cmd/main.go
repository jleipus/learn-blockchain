package main

import (
	"log"

	"github.com/jleipus/learn-blockchain/internal/blockchain/badger"
	"github.com/jleipus/learn-blockchain/internal/blockchain/hashcash"
	"github.com/jleipus/learn-blockchain/internal/cli"
)

func main() {
	hashcash.SetVerbose()
	powFactory := hashcash.New()

	storage, err := badger.NewStorage("blockchain.db")
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	rootCmd := cli.NewRootCmd(storage, powFactory)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
