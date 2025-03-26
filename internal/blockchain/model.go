package blockchain

import (
	pb "github.com/jleipus/learn-blockchain/proto"
)

type BlockHash [32]byte

type ProofOfWorkFactory interface {
	Produce(block *pb.BlockEntity) (hash BlockHash, powData []byte)
	Validate(block *pb.BlockEntity) bool
}
