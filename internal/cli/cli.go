package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
)

type CLI struct {
	bc *blockchain.Blockchain
	pf blockchain.ProofOfWorkFactory
}

func New(bc *blockchain.Blockchain, powFactory blockchain.ProofOfWorkFactory) *CLI {
	return &CLI{
		bc: bc,
		pf: powFactory,
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	var err error
	switch os.Args[1] {
	case "addblock":
		err = addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		err = printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	for i, block := range cli.bc.Blocks() {
		fmt.Printf("[Block %d]\n", i)
		fmt.Printf("Data:\t\t%s\n", block.GetTransactions())
		fmt.Printf("Hash:\t\t%x\n", block.GetHash())
		fmt.Printf("Prev. hash:\t%x\n", block.GetPrevBlockHash())
		fmt.Printf("PoW:\t\t%v\n\n", cli.pf.Validate(block))
	}
}
