package cli

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Blockchain CLI",
	Long:  `Blockchain CLI is a command-line interface for managing a blockchain.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help() // Display help if no subcommand is provided
	},
}

func NewRootCmd(storage blockchain.Storage, powFactory blockchain.ProofOfWorkFactory) *cobra.Command {
	rootCmd.AddCommand(
		newCreateWalletCmd(storage),
		newCreateBlockchainCmd(storage, powFactory),
	)

	return rootCmd
}
