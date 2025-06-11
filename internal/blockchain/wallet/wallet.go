package wallet

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/jleipus/learn-blockchain/internal/utils"
	"golang.org/x/crypto/ripemd160"
)

const (
	version        = byte(0x00) // Version byte for the address
	ChecksumLength = 4          // Length of the checksum
)

// Wallet represents a cryptocurrency wallet containing a private key and a public key.
type Wallet struct {
	PrivateKey ecdh.PrivateKey
	PublicKey  []byte
}

// newWallet creates a newWallet Wallet with a randomly generated key pair.
func newWallet() (*Wallet, error) {
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

// newKeyPair generates a new ECDH key pair using the P-256 curve.
func newKeyPair() (ecdh.PrivateKey, []byte, error) {
	curve := ecdh.P256()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return ecdh.PrivateKey{}, nil, err
	}

	publicKey := privateKey.PublicKey().Bytes()

	return *privateKey, publicKey, nil
}

// getAddress generates a human-readable address from the wallet's public key.
func (w *Wallet) getAddress() (string, error) {
	pubKeyHash, err := HashPubKey(w.PublicKey)
	if err != nil {
		return "", err
	}

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := utils.Base58Encode(fullPayload)

	return hex.EncodeToString(address), nil
}

// Serialize serializes the Wallet into a byte slice.
func (w *Wallet) Serialize() []byte {
	var result bytes.Buffer
	result.Write(w.PublicKey)
	result.Write(w.PrivateKey.Bytes())
	return result.Bytes()
}

// Deserialize deserializes a byte slice into a Wallet.
func (w *Wallet) Deserialize(d []byte) error {
	const (
		publicKeyLen  = 65 // For P-256 uncompressed public key
		privateKeyLen = 32
	)

	if len(d) != publicKeyLen+privateKeyLen {
		return fmt.Errorf("invalid data length: expected %d bytes, got %d", publicKeyLen+privateKeyLen, len(d))
	}

	curve := ecdh.P256()

	publicKeyData := d[:publicKeyLen]
	privateKeyData := d[publicKeyLen:]

	privateKey, err := curve.NewPrivateKey(privateKeyData)
	if err != nil {
		return fmt.Errorf("failed to create private key: %w", err)
	}

	// Optional: Verify that public key matches
	generatedPublicKey := privateKey.PublicKey().Bytes()
	if !bytes.Equal(publicKeyData, generatedPublicKey) {
		return fmt.Errorf("public key does not match private key")
	}

	w.PrivateKey = *privateKey
	w.PublicKey = publicKeyData
	return nil
}

// HashPubKey hashes the public key using SHA-256 followed by RIPEMD-160.
func HashPubKey(pubKey []byte) ([]byte, error) {
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

	return hash[:ChecksumLength]
}
