package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
)

type CLI struct {
	storage    blockchain.Storage
	powFactory blockchain.ProofOfWorkFactory
}

func New(storage blockchain.Storage, powFactory blockchain.ProofOfWorkFactory) *CLI {
	return &CLI{
		storage:    storage,
		powFactory: powFactory,
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

func (cli *CLI) printUsage() {
	printOutput("Usage:")
	printOutput("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	printOutput("  printchain - Print all the blocks of the blockchain")
	printOutput("  getbalance -address ADDRESS - Get balance of ADDRESS")
	printOutput("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 { //nolint:mnd // Fixed number of arguments
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createBlockchain(address string) {
	err := blockchain.CreateBlockchain(cli.storage, cli.powFactory, address)
	if err != nil {
		printOutput("failed to create blockchain: %v", err)
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	bc, err := blockchain.LoadBlockchain(cli.storage, cli.powFactory)
	if err != nil {
		printOutput("failed to load blockchain: %v", err)
		os.Exit(1)
	}

	for i, block := range bc.Blocks() {
		printOutput("[Block %d]", i)
		printOutput("Hash:\t\t%x", block.Hash)
		printOutput("Prev. hash:\t%x", block.PrevBlockHash)
		printOutput("PoW:\t\t%v\n\n", cli.powFactory.Validate(block))
	}
}

func (cli *CLI) getBalance(address string) {
	bc, err := blockchain.LoadBlockchain(cli.storage, cli.powFactory)
	if err != nil {
		printOutput("failed to load blockchain: %v", err)
		os.Exit(1)
	}

	var balance int32
	utxos := bc.FindUnspentTxOutputs(address)

	for _, out := range utxos {
		balance += out.Value
	}

	printOutput("Balance of '%s': %d", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	bc, err := blockchain.LoadBlockchain(cli.storage, cli.powFactory)
	if err != nil {
		printOutput("failed to load blockchain: %v", err)
		os.Exit(1)
	}

	tx, err := bc.NewUnspentTx(from, to, int32(amount)) //nolint:gosec // Ignore integer overflow
	if err != nil {
		printOutput("failed to create transaction: %v", err)
		os.Exit(1)
	}

	err = bc.MineBlock([]*blockchain.Transaction{tx})
	if err != nil {
		printOutput("failed to mine block: %v", err)
		os.Exit(1)
	}
}

func printOutput(format string, args ...any) {
	if format == "" {
		return
	}

	if format[len(format)-1] != '\n' {
		format += "\n"
	}

	fmt.Printf(format, args...) //nolint:forbidigo // Using fmt.Printf for formatted output
}
