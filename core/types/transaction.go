package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/danitello/go-blockchain/chaindb/dbutil"
	"github.com/danitello/go-blockchain/common/errutil"
)

/*Transaction placed in Blocks
@param TxID - Transaction ID
@param TxInput - associated Transaction input
@param TxOutput - associated Transaction output
*/
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

/*initTransaction initializes a new Tranaction
@param TxID - Transaction ID
@param TxInput - associated Transaction input
@param TxOutput - associated Transaction output
@return the Transaction
*/
func initTransaction(inputs []TxInput, outputs []TxOutput) *Transaction {
	tx := Transaction{nil, inputs, outputs}
	tx.setID()
	return &tx
}

/*CreateTransaction creates a Transaction that will be added to a Block in the BlockChain
@param from - the sending address
@param to - the receiving address
@param amount - the amount being exchanged
@param txoSum - sum of txos being spent
@param utxos - map of txIDs and utxoIdxs
@return the new Transaction
*/
func CreateTransaction(from, to string, amount, txoSum int, utxos map[string][]int) *Transaction {
	var newInputs []TxInput
	var newOutputs []TxOutput

	if txoSum < amount {
		log.Panic("Error: Not enough funds")
	}

	// New inputs for this Transaction
	for txID, utxoIdxs := range utxos {
		txID, err := hex.DecodeString(txID)
		errutil.HandleErr(err)

		for _, utxoIdx := range utxoIdxs {
			newInputs = append(newInputs, TxInput{txID, utxoIdx, from}) // map outputs being spent to TxInputs
		}
	}

	// New outputs for this Transaction
	newOutputs = append(newOutputs, TxOutput{amount, to})
	if txoSum > amount {
		newOutputs = append(newOutputs, TxOutput{txoSum - amount, from}) // Keep left over
	}

	newTx := initTransaction(newInputs, newOutputs)
	return newTx

}

/*CoinbaseTx is the transaction in each Block that rewards the miner
@param to - address of recipient
@return the created Transaction
*/
func CoinbaseTx(to string) *Transaction {
	value := 100
	txin := TxInput{[]byte{}, -1, fmt.Sprintf("%d coins to %s", value, to)} // referencing no output
	txout := TxOutput{value, to}
	newTx := initTransaction([]TxInput{txin}, []TxOutput{txout})
	return newTx
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
