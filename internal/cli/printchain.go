package cli

import (
	"fmt"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/spf13/cobra"
)

func newPrintChainCmd(storage blockchain.Storage, powFactory blockchain.ProofOfWorkFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "print-chain",
		Short: "Print the blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			wallets := wallet.NewCollection(storage)

			bc, err := blockchain.LoadBlockchain(storage, powFactory, wallets)
			if err != nil {
				cmd.PrintErrf("Error loading blockchain: %v\n", err)
				return
			}

			for _, b := range bc.Blocks() {
				fmt.Printf("============ Block %x ============\n", b.Hash)
				fmt.Printf("Prev. block: %x\n", b.PrevBlockHash)
			}
		},
	}
}
