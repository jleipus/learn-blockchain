package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/jleipus/learn-blockchain/internal/utils"
	"golang.org/x/crypto/ripemd160"
)

const (
	version        = byte(0x00) // Version byte for the address
	checksumLength = 4          // Length of the checksum
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Walllets map[string]*Wallet
}

func NewWallet() (*Wallet, error) {
	privateKey, publicKey, err := newKeyPair()
	if err != nil {
		return nil, err
	}

	wallet := &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	return wallet, nil
}

func newKeyPair() (ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, nil, err
	}

	publicKey := append(privateKey.X.Bytes(), privateKey.Y.Bytes()...)

	return *privateKey, publicKey, nil
}

func (w *Wallet) GetAddress() ([]byte, error) {
	pubKeyHash, err := hashPubKey(w.PublicKey)
	if err != nil {
		return nil, err
	}

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := utils.Base58Encode(fullPayload)

	return address, nil
}

func hashPubKey(pubKey []byte) ([]byte, error) {
	pubKeySHA256 := sha256.Sum256(pubKey)

	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(pubKeySHA256[:])
	if err != nil {
		return nil, err
	}

	pubKeyRIPEMD160 := ripemd160Hasher.Sum(nil)

	return pubKeyRIPEMD160, nil
}

func checksum(payload []byte) []byte {
	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])

	return hash[:checksumLength]
}
