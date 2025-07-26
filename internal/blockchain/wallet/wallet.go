package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/jleipus/learn-blockchain/internal/utils"
	"golang.org/x/crypto/ripemd160"
)

const (
	version        = byte(0x00) // Version byte for the address
	VersionLength  = 1          // Length of the version field
	ChecksumLength = 4          // Length of the checksum
)

var (
	ErrAddressTooShort = errors.New("too short")
)

// Wallet represents a cryptocurrency wallet containing a private key and a public key.
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// New creates a New Wallet with a randomly generated key pair.
func New() (*Wallet, error) {
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
func newKeyPair() (ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, nil, err
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, publicKey, nil
}

// getAddress generates a human-readable address from the wallet's public key.
// The address consists of a version byte, the hashed public key, and a checksum.
// The full address is encoded in Base58 to make it human-readable.
func (w *Wallet) getAddress() ([]byte, error) {
	pubKeyHash, err := HashPubKey(w.PublicKey)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, 0, VersionLength+len(pubKeyHash)+ChecksumLength)

	payload = append(payload, version)       // Version
	payload = append(payload, pubKeyHash...) // Hash

	checksum := checksum(payload)
	payload = append(payload, checksum...) // Checksum

	address := utils.Base58Encode(payload)
	return address, nil
}

// Serialize serializes the Wallet into a byte slice to be stored.
func (w *Wallet) Serialize() ([]byte, error) {
	curve := w.PrivateKey.Curve
	byteLen := (curve.Params().BitSize + 7) / 8 //nolint:mnd // Magic number for byte length

	buf := new(bytes.Buffer)

	// Write D (private scalar)
	dBytes := w.PrivateKey.D.Bytes()
	paddedD := make([]byte, byteLen)
	copy(paddedD[byteLen-len(dBytes):], dBytes)
	buf.Write(paddedD)

	// Write X and Y (public point)
	xBytes := w.PrivateKey.PublicKey.X.Bytes()
	yBytes := w.PrivateKey.PublicKey.Y.Bytes()

	paddedX := make([]byte, byteLen)
	paddedY := make([]byte, byteLen)
	copy(paddedX[byteLen-len(xBytes):], xBytes)
	copy(paddedY[byteLen-len(yBytes):], yBytes)

	buf.Write(paddedX)
	buf.Write(paddedY)

	return buf.Bytes(), nil
}

// Deserialize deserializes a byte slice into a Wallet when reading from storage.
func (w *Wallet) Deserialize(data []byte) error {
	curve := elliptic.P256()
	byteLen := (curve.Params().BitSize + 7) / 8 //nolint:mnd // Magic number for byte length

	if len(data) != byteLen*3 {
		return errors.New("invalid wallet data length")
	}

	d := new(big.Int).SetBytes(data[:byteLen])
	x := new(big.Int).SetBytes(data[byteLen : 2*byteLen])
	y := new(big.Int).SetBytes(data[2*byteLen : 3*byteLen])

	priv := ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
	}

	pubKey := append(x.Bytes(), y.Bytes()...)

	w.PrivateKey = priv
	w.PublicKey = pubKey
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

func GetHashFromAddress(address []byte) ([]byte, error) {
	// Convert address to a public key hash
	pubKeyHash := utils.Base58Decode(address)

	minAddressLength := VersionLength + ChecksumLength
	if len(pubKeyHash) < minAddressLength {
		return nil, ErrAddressTooShort
	}

	// Remove version byte and checksum
	pubKeyHash = pubKeyHash[VersionLength : len(pubKeyHash)-ChecksumLength]

	return pubKeyHash, nil
}

func ValidateAddress(address string) error {
	addressPayload := utils.Base58Decode([]byte(address))

	foundChecksum := addressPayload[len(addressPayload)-ChecksumLength:]
	foundVersion := addressPayload[0]
	foundHash, err := GetHashFromAddress([]byte(address))
	if err != nil {
		return err
	}

	targetChecksum := checksum(append([]byte{foundVersion}, foundHash...))

	if ok := bytes.Equal(foundChecksum, targetChecksum); !ok {
		return errors.New("invalid checksum")
	}

	return nil
}

func checksum(payload []byte) []byte {
	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])

	return hash[:ChecksumLength]
}
