package transaction

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	pb "github.com/jleipus/learn-blockchain/proto"
	"google.golang.org/protobuf/proto"
)

const (
	subsidy = 10
)

func NewCoinbaseTX(to, data string) *pb.Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := pb.TXInput{
		TxID:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := pb.TXOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}
	tx := pb.Transaction{
		ID:   nil,
		Vin:  []*pb.TXInput{&txin},
		Vout: []*pb.TXOutput{&txout},
	}

	setTransactionID(&tx)

	return &tx
}

func setTransactionID(tx *pb.Transaction) {
	bytes, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(bytes)
	tx.ID = hash[:]
}

func HashTransactions(txs []*pb.Transaction) []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range txs {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
