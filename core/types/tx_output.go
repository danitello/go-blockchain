package types

import (
	"bytes"
	"encoding/gob"

	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/danitello/go-blockchain/wallet"
	"github.com/danitello/go-blockchain/wallet/walletutil"
)

/*TxOutput specifies amount being made available in a block to a wallet
@param Amount - total
@param PubKeyHash - hash of pub key to unlock output
*/
type TxOutput struct {
	Amount     int
	PubKeyHash []byte
}

/*TxOutputs groups txos (for serialization) */
type TxOutputs struct {
	Outputs []TxOutput
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

/*DeserializeTxOutputs converts a []byte into []TxOutput
@param data - the []byte representation of []TxOutput
@returns a []TxOutput representation
*/
func DeserializeTxOutputs(data []byte) TxOutputs {
	var TXO TxOutputs

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&TXO)
	errutil.Handle(err)

	return TXO
}
