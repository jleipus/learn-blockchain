package main

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/cli"
	"github.com/jleipus/learn-blockchain/internal/pow/hashcash"
)

func main() {
	hashcash.SetVerbose()
	powFactory := hashcash.New(18)

	bc, err := blockchain.New("blockchain.db", powFactory)
	if err != nil {
		panic(err)
	}
	defer bc.Close()

	cli := cli.New(bc, powFactory)
	cli.Run()
}
