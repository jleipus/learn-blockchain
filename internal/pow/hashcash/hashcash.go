package hashcash

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/jleipus/learn-blockchain/internal/utils"
	"github.com/jleipus/learn-blockchain/proto/block"
)

const (
	maxNonce int64 = math.MaxInt64
)

var verbose bool = false

func SetVerbose() {
	verbose = true
}

type hashCashPoW struct {
	target *big.Int
}

func new(targetBits int64) *hashCashPoW {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &hashCashPoW{
		target: target,
	}

	return pow
}

func (pow *hashCashPoW) Produce(block *block.Block) (int64, [32]byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := int64(0)

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

	return nonce, hash
}

func (pow *hashCashPoW) Validate(block *block.Block) bool {
	var hashInt big.Int

	hash := calculateHash(block, block.GetNonce())
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

func calculateHash(block *block.Block, nonce int64) [32]byte {
	var data []byte

	data = append(data, block.PrevBlockHash...)
	data = append(data, block.Data...)
	data = append(data, utils.IntToHex(block.Timestamp)...)
	data = append(data, utils.IntToHex(nonce)...)

	return sha256.Sum256(data)
}
