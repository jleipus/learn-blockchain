package cli

import (
	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/spf13/cobra"
)

func newListAddressesCmd(storage blockchain.Storage) *cobra.Command {
	return &cobra.Command{
		Use:   "list-addresses",
		Short: "List all wallet addresses",
		Run: func(cmd *cobra.Command, args []string) {
			wallets := wallet.NewCollection(storage)
			addresses, err := wallets.GetAddresses()
			if err != nil {
				cmd.Println("Error retrieving addresses:", err)
				return
			}

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
