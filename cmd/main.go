package main

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/badger"
	"github.com/jleipus/learn-blockchain/internal/blockchain/hashcash"
	"github.com/jleipus/learn-blockchain/internal/cli"
)

func main() {
	hashcash.SetVerbose()
	powFactory := hashcash.New(18)

	storage, err := badger.NewBadgerDB("blockchain.db")
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	bc, err := blockchain.NewBlockchain(storage, "start", powFactory)
	if err != nil {
		panic(err)
	}
	defer bc.Close()

	cli := cli.New(bc, powFactory)
	cli.Run()

	// bc.AddBlock([]*blockchain.Transaction{{}, {}, {}})
	// bc.AddBlock([]*blockchain.Transaction{{}, {}, {}})
	// bc.AddBlock([]*blockchain.Transaction{{}, {}, {}})

	// for i, block := range bc.Blocks() {
	// 	fmt.Printf("[Block %d]\n", i)
	// 	fmt.Printf("Hash:\t\t%x\n", block.Hash)
	// 	fmt.Printf("Prev. hash:\t%x\n", block.PrevBlockHash)
	// 	fmt.Printf("PoW:\t\t%v\n\n", powFactory.Validate(block))
	// }
}
