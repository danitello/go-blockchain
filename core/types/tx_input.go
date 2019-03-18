package types

import (
	"bytes"

	"github.com/danitello/go-blockchain/wallet"
)

// TxInput spends (references) a previous TxOutput -
// TxID - ID of Transaction that the TxOutput resides in
// OutputIdx - idx of the TxOutput in the Transaction
// Signature - signs the txin as unlocking the txo
// PubKey - the pub key used
type TxInput struct {
	TxID      []byte
	OutputIdx int
	Signature []byte
	PubKey    []byte
}

// UsesKey determines whether the pubKeyHash provided is the owner of the output referenced by txin
func (txin *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(txin.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
