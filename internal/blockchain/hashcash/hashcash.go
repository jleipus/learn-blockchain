package hashcash

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/blockchain/block"
	"github.com/jleipus/learn-blockchain/internal/utils"
)

const (
	targetBits   uint32 = 18 // Target difficulty bits
	maxNonce     uint64 = math.MaxUint64
	sha256Length uint32 = 256
	nonceLength  int    = 8
)

var verbose = false //nolint:gochecknoglobals // Verbose flag for printing mining progress

func SetVerbose() {
	verbose = true
}

type hashCashPoW struct {
	target *big.Int
}

func New() blockchain.ProofOfWorkFactory {
	target := big.NewInt(1)
	target.Lsh(target, uint(sha256Length-targetBits))

	pow := &hashCashPoW{
		target: target,
	}

	return pow
}

var TimeNow = time.Now // Allow mocking time for testing

func (pow *hashCashPoW) Produce(b *block.Block) (block.Hash, []byte) {
	var hashInt big.Int
	var hash block.Hash
	nonce := uint64(0)

	start := TimeNow()

	var err error
	for nonce < maxNonce {
		hash, err = calculateHash(b, nonce)
		if err != nil {
			panic(err)
		}

		if verbose {
			//nolint:forbidigo // Print the hash and elapsed time
			fmt.Printf("\rMining: %x (%s)", hash, time.Since(start))
		}

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		}

		nonce++
	}

	if verbose {
		//nolint:forbidigo // Print the hash and elapsed time
		fmt.Printf("\rCompleted mining: %x (%s)\n", hash, time.Since(start))
	}

	powData := make([]byte, nonceLength)
	binary.BigEndian.PutUint64(powData, nonce)

	return hash, powData
}

func (pow *hashCashPoW) Validate(block *block.Block) bool {
	var hashInt big.Int

	powData := block.PoW

	if len(powData) < nonceLength {
		panic(fmt.Sprintf("Invalid PoW length: expected %d, got %d", nonceLength, len(powData)))
	}
	nonce := binary.BigEndian.Uint64(powData[:nonceLength])

	hash, err := calculateHash(block, nonce)
	if err != nil {
		panic(err)
	}

	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

func calculateHash(b *block.Block, nonce uint64) (block.Hash, error) {
	var data []byte

	txHash := b.HashTransactions()

	timestampHex, err := utils.IntToHex(b.Timestamp)
	if err != nil {
		return block.Hash{}, err
	}

	nonceHex, err := utils.IntToHex(nonce)
	if err != nil {
		return block.Hash{}, err
	}

	data = append(data, b.PrevBlockHash[:]...)
	data = append(data, txHash[:]...)
	data = append(data, timestampHex...)
	data = append(data, nonceHex...)

	return sha256.Sum256(data), nil
}
