package main

import (
	"fmt"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/hashcash"
	pb "github.com/jleipus/learn-blockchain/proto"
)

func main() {
	hashcash.SetVerbose()
	powFactory := hashcash.New(18)

	bc, err := blockchain.New("blockchain.db", "start", powFactory)
	if err != nil {
		panic(err)
	}
	defer bc.Close()

	// cli := cli.New(bc, powFactory)
	// cli.Run()

	bc.AddBlock([]*pb.Transaction{{}, {}, {}})
	bc.AddBlock([]*pb.Transaction{{}, {}, {}})
	bc.AddBlock([]*pb.Transaction{{}, {}, {}})

	for i, block := range bc.Blocks() {
		fmt.Printf("[Block %d]\n", i)
		// fmt.Printf("Data:\t\t%s\n", block.GetTransactions())
		fmt.Printf("Hash:\t\t%x\n", block.GetHash())
		fmt.Printf("Prev. hash:\t%x\n", block.GetPrevBlockHash())
		fmt.Printf("PoW:\t\t%v\n\n", powFactory.Validate(block))
	}
}
