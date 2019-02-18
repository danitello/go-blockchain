package types

import (
	"bytes"

	"github.com/danitello/go-blockchain/wallet"
)

/*TxInput is a reference to a previous TxOutput
@param TxID - ID of Transaction that the TxOutput resides in
@param OutputIdx - idx of the TxOutput in the Transaction
@param Signature - signs the txin as unlocking the txo
@param PubKey - the pub key used
*/
type TxInput struct {
	TxID      []byte
	OutputIdx int
	Signature []byte
	PubKey    []byte
}

/*UsesKey determines whether the pubKeyHash provided is the owner of the ouput referenced by txin
@param pubKeyHash - the pubKeyHash in question
@return whether it is valid
*/
func (txin *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(txin.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
