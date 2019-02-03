package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/wallet/walletutil"
	"golang.org/x/crypto/ripemd160"
)

const (
	// number of inital bytes to take from result of the shas256 hashes of the pub key hash
	checksumLen = 4
	// version of gen algo
	version = byte(0x00)
)

/*Wallet is the entity for ownership on the chain
@param PrivateKey - unique identifier
@param PublicKey - to derive public address
*/
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

/*InitWallet initializes a new Wallet
@return the Wallet
*/
func InitWallet() *Wallet {
	priv, pub := createKeyPair()
	return &Wallet{priv, pub}
}

/*createKeyPair makes a new priv and pub key pair
@return ecdsa.PrivateKey - ecdsa priv key
@return []byte - derived pub key
*/
func createKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	errutil.HandleErr(err)

	// Derive []byte representation of pub key
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	return *privKey, pubKey
}

/*GetAddress derives the human readable address for a Wallet using pub key hash, version, and checksum (bitcoin spec)
@param w - Wallet in question
@return the address
*/
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedHash := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedHash)
	fullHash := append(versionedHash, checksum...)

	return walletutil.Base58Encode(fullHash)
}

/*HashPubKey computes the pub key hash
@param pubKey - pub key
@return the pub key hash
*/
func HashPubKey(pubKey []byte) []byte {
	shaPubKey := sha256.Sum256(pubKey)

	ripemdHasher := ripemd160.New()
	_, err := ripemdHasher.Write(shaPubKey[:])
	errutil.HandleErr(err)
	ripemdPubKey := ripemdHasher.Sum(nil)

	return ripemdPubKey

}

/*checksum computes the checksum of a given payload
@param payload - to checksum
@return the checksum
*/
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:checksumLen]
}
