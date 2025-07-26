package cli

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/spf13/cobra"
)

func newGetBalanceCmd(storage blockchain.Storage, powFactory blockchain.ProofOfWorkFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "get-balance",
		Aliases: []string{"b"},
		Short:   "Get the balance of an address",
		Long: `Get the balance of a specific address in the blockchain.
This command retrieves the balance of the specified address by summing up all the transaction outputs that belong to it.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := wallet.ValidateAddress(args[0]); err != nil {
				cmd.PrintErrf("Invalid address %s: %v\n", args[0], err)
				return
			}

			wallets := wallet.NewCollection(storage)

			bc, err := blockchain.LoadBlockchain(storage, powFactory, wallets)
			if err != nil {
				cmd.PrintErrf("Error loading blockchain: %v\n", err)
				return
			}

			pubKeyHash, err := wallet.GetHashFromAddress([]byte(args[0]))
			if err != nil {
				cmd.PrintErrf("Error getting public key hash from address: %v\n", err)
				return
			}

			balance := 0
			for _, out := range bc.FindUnspentTxOutputs(pubKeyHash) {
				balance += int(out.Value)
			}

			cmd.Printf("%d\n", balance)
		},
	}
}
