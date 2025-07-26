package cli

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/spf13/cobra"
)

func newCreateBlockchainCmd(storage blockchain.Storage, powFactory blockchain.ProofOfWorkFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "create-blockchain",
		Short: "Create a new blockchain",
		Long: `Create a new blockchain with a genesis block.
This command initializes a new blockchain with a genesis block and sets the initial state.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := wallet.ValidateAddress(args[0]); err != nil {
				cmd.PrintErrf("Invalid address %s: %v\n", args[0], err)
				return
			}

			if err := blockchain.CreateBlockchain(storage, powFactory, args[0]); err != nil {
				cmd.PrintErrf("Error creating blockchain: %v\n", err)
				return
			}
			cmd.Printf("Blockchain created with genesis block for address: %s\n", args[0])
		},
	}
}
