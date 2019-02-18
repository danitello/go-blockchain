package types

import (
	"bytes"

	"github.com/danitello/go-blockchain/wallet"
	"github.com/danitello/go-blockchain/wallet/walletutil"
)

/*TxOutput specifies amount being made available to a wallet
@param Amount - total
@param PubKeyHash - hash of pub key to unlock output
*/
type TxOutput struct {
	Amount     int
	PubKeyHash []byte
}

/*InitTxOutput creates a new txo and locks it using a given address
@param amount - the amout in the txo
@param address - the locking address
@return the new txo
*/
func InitTxOutput(amount int, address string) *TxOutput {
	txo := &TxOutput{amount, nil}
	txo.Lock([]byte(address))

	return txo
}

/*Lock signs the TxOutput with a given address
@param the address
*/
func (txo *TxOutput) Lock(address []byte) {
	pubKeyHash := walletutil.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-wallet.ChecksumLen]
	txo.PubKeyHash = pubKeyHash
}

/*IsLockedWithKey determines whether a given pubKeyHash is the one used to lock the txo
@param pubKeyHash - the hash in question
@return whether it is valid
*/
func (txo *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(txo.PubKeyHash, pubKeyHash) == 0
}
