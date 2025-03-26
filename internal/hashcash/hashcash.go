package hashcash

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/transaction"
	"github.com/jleipus/learn-blockchain/internal/utils"
	pb "github.com/jleipus/learn-blockchain/proto"
)

const (
	maxNonce uint64 = math.MaxUint64
)

var verbose bool = false

func SetVerbose() {
	verbose = true
}

type hashCashPoW struct {
	target *big.Int
}

func New(targetBits int64) blockchain.ProofOfWorkFactory {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &hashCashPoW{
		target: target,
	}

	return pow
}

func (pow *hashCashPoW) Produce(block *pb.BlockEntity) (blockchain.BlockHash, []byte) {
	var hashInt big.Int
	var hash blockchain.BlockHash
	nonce := uint64(0)

	start := time.Now()
	for nonce < maxNonce {
		hash = calculateHash(block, nonce)

		if verbose {
			fmt.Printf("\rMining: %x (%s)", hash, time.Since(start))
		}

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	if verbose {
		fmt.Printf("\rCompleted mining: %x (%s)\n", hash, time.Since(start))
	}

	powData := make([]byte, 8)
	binary.BigEndian.PutUint64(powData, nonce)

	return hash, powData
}

func (pow *hashCashPoW) Validate(block *pb.BlockEntity) bool {
	var hashInt big.Int

	powData := block.GetPoW()
	nonce := binary.BigEndian.Uint64(powData[:8])

	hash := calculateHash(block, nonce)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

func calculateHash(block *pb.BlockEntity, nonce uint64) [32]byte {
	var data []byte

	data = append(data, block.GetPrevBlockHash()...)
	data = append(data, transaction.HashTransactions(block.GetTransactions())...)
	data = append(data, utils.IntToHex(block.Timestamp)...)
	data = append(data, utils.IntToHex(nonce)...)

	return sha256.Sum256(data)
}
