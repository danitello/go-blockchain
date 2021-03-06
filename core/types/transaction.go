package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/danitello/go-blockchain/common/byteutil"
	"github.com/danitello/go-blockchain/common/errutil"
)

// Transaction placed in Blocks
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// initTransaction initializes a new Tranaction
func initTransaction(inputs []TxInput, outputs []TxOutput) *Transaction {
	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	return &tx
}

// CreateTransaction creates a Transaction that will be added to a Block in the BlockChain -
// txoSum - sum of txos being spent
// utxos - map of txIDs and utxoIdxs
func CreateTransaction(from, to string, pubKeyHash []byte, amount, txoSum int, utxos map[string][]int) *Transaction {
	var newInputs []TxInput
	var newOutputs []TxOutput

	if txoSum < amount {
		pString := fmt.Sprintf("Error: Not enough funds in wallet address: %s", from)
		log.Panic(pString)
	}

	// New inputs for this Transaction
	for txID, utxoIdxs := range utxos {
		txID, err := hex.DecodeString(txID)
		errutil.Handle(err)

		for _, utxoIdx := range utxoIdxs {
			newInputs = append(newInputs, TxInput{txID, utxoIdx, nil, pubKeyHash}) // map outputs being spent by TxInputs
		}
	}

	// New outputs for this Transaction
	newOutputs = append(newOutputs, *InitTxOutput(amount, to))
	if txoSum > amount {
		newOutputs = append(newOutputs, *InitTxOutput(txoSum-amount, from)) // Keep left over
	}

	newTx := initTransaction(newInputs, newOutputs)
	return newTx

}

// Sign computes the signature for each txin in the tx with ecdsa -
// privKey - of signer
// prevTxs - containing the txos that will be referenced by new txins
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, txin := range tx.Inputs {
		if prevTxs[hex.EncodeToString(txin.TxID)].ID == nil {
			log.Panic("ERROR: tx.Sign cannot find previous txn with ID")
		}
	}

	txCopy := tx.TrimmedCopy()

	for txinID, txin := range txCopy.Inputs {
		prevTx := prevTxs[hex.EncodeToString(txin.TxID)]
		txCopy.Inputs[txinID].Signature = nil
		txCopy.Inputs[txinID].PubKey = prevTx.Outputs[txin.OutputIdx].PubKeyHash
		txCopy.ID = txCopy.Hash()

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		errutil.Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[txinID].Signature = signature // now update the actual tx
		txCopy.Inputs[txinID].PubKey = nil

	}
}

// Verify determines whether txins were signed correctly
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, txin := range tx.Inputs {
		if prevTxs[hex.EncodeToString(txin.TxID)].ID == nil {
			log.Panic("ERROR: tx.Verify cannot find previous txn with ID", prevTxs[hex.EncodeToString(txin.TxID)].ID)
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for txinID, txin := range tx.Inputs {
		// Get same information as signing flow
		prevTx := prevTxs[hex.EncodeToString(txin.TxID)]
		txCopy.Inputs[txinID].Signature = nil
		txCopy.Inputs[txinID].PubKey = prevTx.Outputs[txin.OutputIdx].PubKeyHash
		txCopy.ID = txCopy.Hash()

		// Signature information
		r := big.Int{}
		s := big.Int{}
		sigLen := len(txin.Signature)
		r.SetBytes(txin.Signature[:(sigLen / 2)])
		s.SetBytes(txin.Signature[(sigLen / 2):])

		// PubKey information
		x := big.Int{}
		y := big.Int{}
		keyLen := len(txin.PubKey)
		x.SetBytes(txin.PubKey[:(keyLen / 2)])
		y.SetBytes(txin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y} // reconstruct
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
		txCopy.Inputs[txinID].PubKey = nil
	}

	return true
}

// TrimmedCopy sets the Signature and PubKey fields of all txins to nil as these are unecessary for signing (btc spec)
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, txin := range tx.Inputs {
		inputs = append(inputs, TxInput{txin.TxID, txin.OutputIdx, nil, nil})
	}

	for _, txo := range tx.Outputs {
		outputs = append(outputs, TxOutput{txo.Amount, txo.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Hash computes the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(byteutil.Serialize(txCopy))

	return hash[:]
}

// CoinbaseTx is the transaction in each Block that rewards the miner
func CoinbaseTx(to string) *Transaction {
	amount := 100
	txin := TxInput{[]byte{}, -1, nil, []byte(fmt.Sprintf("CoinbaseTx: %d coins to %s", amount, to))} // referencing no output
	txout := InitTxOutput(amount, to)
	newTx := initTransaction([]TxInput{txin}, []TxOutput{*txout})
	return newTx
}

// IsCoinbase determines whether a Transaction is a coinbase tx
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0 && tx.Inputs[0].OutputIdx == -1
}

// String creates a string containing information to display about the tx
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, txin := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TxID:      %x", txin.TxID))
		lines = append(lines, fmt.Sprintf("       OutputIdx:       %d", txin.OutputIdx))
		lines = append(lines, fmt.Sprintf("       Signature: %x", txin.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", txin.PubKey))
	}
	for i, txo := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       amount:  %d", txo.Amount))
		lines = append(lines, fmt.Sprintf("       PubKeyHash: %x", txo.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
