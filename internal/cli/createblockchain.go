package cli

import (
	"github.com/spf13/cobra"
)

var createBlockchainCmd = &cobra.Command{
	Use:     "create-blockchain",
	Aliases: []string{"cb"},
	Short:   "Create a new blockchain",
	Long: `Create a new blockchain with a genesis block.
This command initializes a new blockchain with a genesis block and sets the initial state.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}
