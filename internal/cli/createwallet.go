package cli

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/spf13/cobra"
)

func newCreateWalletCmd(storage blockchain.Storage) *cobra.Command {
	return &cobra.Command{
		Use:   "create-wallet",
		Short: "Create a new wallet",
		Run: func(cmd *cobra.Command, args []string) {
			wallets := wallet.NewCollection(storage)
			address, err := wallets.AddWallet()
			if err != nil {
				cmd.PrintErrf("Error creating wallet: %v\n", err)
				return
			}
			cmd.Printf("%s\n", address)
		},
	}
}
