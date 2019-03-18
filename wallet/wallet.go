package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/wallet/walletutil"
	"golang.org/x/crypto/ripemd160"
)

const (
	// ChecksumLen is number of initial bytes to take from result of the sha256 hashes of the pub key hash
	ChecksumLen = 4
	// version of gen algo
	version = byte(0x00)
)

// Wallet is the entity for ownership on the chain
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// InitWallet initializes a new Wallet
func InitWallet() *Wallet {
	priv, pub := createKeyPair()
	return &Wallet{priv, pub}
}

// createKeyPair makes a new priv and pub key pair
func createKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	errutil.Handle(err)

	// Derive []byte representation of pub key
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	return *privKey, pubKey
}

// GetAddress derives the human readable address for a Wallet using pub key hash, version, and checksum (bitcoin spec)
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedHash := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedHash)
	fullHash := append(versionedHash, checksum...)

	return walletutil.Base58Encode(fullHash)
}

// ValidateAddress determines if a given address is correctly constructed
func ValidateAddress(address string) bool {
	decodedAddress := walletutil.Base58Decode([]byte(address))

	addressChecksum := decodedAddress[len(decodedAddress)-ChecksumLen:]
	targetChecksum := checksum(decodedAddress[0 : len(decodedAddress)-ChecksumLen])

	return bytes.Compare(addressChecksum, targetChecksum) == 0

}

// HashPubKey computes the pub key hash
func HashPubKey(pubKey []byte) []byte {
	shaPubKey := sha256.Sum256(pubKey)

	ripemdHasher := ripemd160.New()
	_, err := ripemdHasher.Write(shaPubKey[:])
	errutil.Handle(err)
	ripemdPubKey := ripemdHasher.Sum(nil)

	return ripemdPubKey

}

// checksum computes the checksum of a given payload
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:ChecksumLen]
}

// GetPubKeyHashFromAddress takes in an address and returns its pub key hash portion
func GetPubKeyHashFromAddress(address string) []byte {
	decodedAddress := walletutil.Base58Decode([]byte(address))
	pubKeyHash := decodedAddress[1 : len(decodedAddress)-ChecksumLen]
	return pubKeyHash
}
