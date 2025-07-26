package cli

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/spf13/cobra"
)

func newListAddressesCmd(storage blockchain.Storage) *cobra.Command {
	return &cobra.Command{
		Use:     "list-addresses",
		Aliases: []string{"ls"},
		Short:   "List all wallet addresses",
		Long:    `List all wallet addresses stored in the blockchain.`,
		Run: func(cmd *cobra.Command, args []string) {
			wallets := wallet.NewCollection(storage)
			addresses := wallets.GetAddresses()

			if len(addresses) == 0 {
				cmd.Println("No addresses found.")
				return
			}

			for _, address := range addresses {
				cmd.Println(address)
			}
		},
	}
}
