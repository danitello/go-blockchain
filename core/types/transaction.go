package types

import (
	"crypto/sha256"
	"fmt"

	"github.com/danitello/go-blockchain/chaindb/dbutil"
)

/* Transaction placed in Blocks
@param ID - Transaction ID
@param TxInput - associated Transaction input
@param TxOutput - associated Transaction output
*/
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

/*TxInput is a reference to a previous TxOutput
@param ID - ID of Transaction that the TxOutput resides in
@param OutputIndex - index of the TxOutput in the Transaction
@param Sig - data used in TxOutput PubKey
*/
type TxInput struct {
	TxID        []byte
	OutputIndex int
	Sig         string
}

/*TxOutput specifies coin value made available to a user
@param Value - amount
@param PubKey - ID of user
*/
type TxOutput struct {
	Amount int
	PubKey string
}

/*CoinbaseTx is the transaction in each Block that rewards the miner
@param to - address of recipient
@return the created Transaction
*/
func CoinbaseTx(to string) *Transaction {
	value := 100
	txin := TxInput{[]byte{}, -1, fmt.Sprintf("%d coins to %s", value, to)} // referencing no output
	txout := TxOutput{value, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.setID()
	return &tx
}

/*setID computes the ID for a Transaction */
func (tx *Transaction) setID() {
	var hash [32]byte
	txEncoded := dbutil.Serialize(tx)

	hash = sha256.Sum256(txEncoded)
	tx.ID = hash[:]
}

/*IsCoinbase determines whether a Transaction is a coinbase tx
@return whether it's a coinbase tx
*/
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0 && tx.Inputs[0].OutputIndex == -1
}

/*CanUnlock determines whether the signature provided is the owner of the ouput referenced by txin
@param newSig - the signature in question
@return whether the signature is valid
*/
func (txin *TxInput) CanUnlock(newSig string) bool {
	return txin.Sig == newSig
}

/*CanBeUnlocked determines whether the PubKey is the owner of the output
@param newPubKey - the PubKey in question
@return whether the PubKey is valid
*/
func (txout *TxOutput) CanBeUnlocked(newPubKey string) bool {
	return txout.PubKey == newPubKey
}
