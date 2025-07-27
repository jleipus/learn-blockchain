package cli

import (
	"strconv"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
	"github.com/jleipus/learn-blockchain/internal/blockchain/wallet"
	"github.com/spf13/cobra"
)

func newSendCmd(storage blockchain.Storage, powFactory blockchain.ProofOfWorkFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "send",
		Short: "Send coins to an address",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if err := wallet.ValidateAddress(args[0]); err != nil {
				cmd.PrintErrf("Invalid sender address %s: %v\n", args[0], err)
				return
			}

			if err := wallet.ValidateAddress(args[1]); err != nil {
				cmd.PrintErrf("Invalid recipient address %s: %v\n", args[1], err)
				return
			}

			amount, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil || amount <= 0 {
				cmd.PrintErrf("Invalid amount %s: must be a positive integer\n", args[2])
				return
			}

			wallets := wallet.NewCollection(storage)

			bc, err := blockchain.LoadBlockchain(storage, powFactory, wallets)
			if err != nil {
				cmd.PrintErrf("Error loading blockchain: %v\n", err)
				return
			}

			tx, err := bc.NewUTXOTransaction(args[0], args[1], int32(amount))
			if err != nil {
				cmd.PrintErrf("Error creating transaction: %v\n", err)
				return
			}

			cbTx, err := transaction.NewCoinbaseTX(args[0], "")
			if err != nil {
				cmd.PrintErrf("Error creating coinbase transaction: %v\n", err)
				return
			}

			b, err := bc.MineBlock([]*transaction.Tx{tx, cbTx})
			if err != nil {
				cmd.PrintErrf("Error mining block: %v\n", err)
				return
			}

			if err := bc.Update(*b); err != nil {
				cmd.PrintErrf("Error updating UTXO set: %v\n", err)
				return
			}
		},
	}
}
