package hashcash

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/jleipus/learn-blockchain/internal/blockchain"
	"github.com/jleipus/learn-blockchain/internal/utils"
)

const (
	targetBits uint32 = 18 // Target difficulty bits
	maxNonce   uint64 = math.MaxUint64
	sha256Bits uint32 = 256
	nonceBytes int    = 8
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
	target.Lsh(target, uint(sha256Bits-targetBits))

	pow := &hashCashPoW{
		target: target,
	}

	return pow
}

func (pow *hashCashPoW) Produce(block *blockchain.Block) (blockchain.BlockHash, []byte) {
	var hashInt big.Int
	var hash blockchain.BlockHash
	nonce := uint64(0)

	start := time.Now()

	var err error
	for nonce < maxNonce {
		hash, err = calculateHash(block, nonce)
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

	powData := make([]byte, nonceBytes)
	binary.BigEndian.PutUint64(powData, nonce)

	return hash, powData
}

func (pow *hashCashPoW) Validate(block *blockchain.Block) bool {
	var hashInt big.Int

	powData := block.PoW
	nonce := binary.BigEndian.Uint64(powData[:nonceBytes])

	hash, err := calculateHash(block, nonce)
	if err != nil {
		panic(err)
	}

	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

func calculateHash(block *blockchain.Block, nonce uint64) (blockchain.BlockHash, error) {
	var data []byte

	txHash := block.HashTransactions()

	timestampHex, err := utils.IntToHex(block.Timestamp)
	if err != nil {
		return blockchain.BlockHash{}, err
	}

	nonceHex, err := utils.IntToHex(nonce)
	if err != nil {
		return blockchain.BlockHash{}, err
	}

	data = append(data, block.PrevBlockHash[:]...)
	data = append(data, txHash[:]...)
	data = append(data, timestampHex...)
	data = append(data, nonceHex...)

	return sha256.Sum256(data), nil
}
