package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Blockchain CLI",
	Long: `Blockchain CLI is a command-line interface for managing a blockchain.
It allows you to initialize a new blockchain, add blocks, and manage transactions.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help() // Display help if no subcommand is provided
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
