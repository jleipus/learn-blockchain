package hashcash

import (
	"time"

	"github.com/jleipus/learn-blockchain/internal/pow"
	pb "github.com/jleipus/learn-blockchain/proto"
)

type hashcashBlock struct {
	*pb.BlockEntity

	nonce int64
	hash  [32]byte
}

func New(transactions []*pb.Transaction, prevBlockHash []byte) pow.Block {
	return &hashcashBlock{
		BlockEntity: &pb.BlockEntity{
			Timestamp:     time.Now().Unix(),
			Transactions:  transactions,
			PrevBlockHash: prevBlockHash,
		},
	}
}

func (b *hashcashBlock) Marshal() ([]byte, error) {

}

func (b *hashcashBlock) Unmarshal(data []byte) error {
	return b.BlockEntity.Unmarshal(data)
}
