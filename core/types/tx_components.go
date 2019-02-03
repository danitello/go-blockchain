package types

/*TxInput is a reference to a previous TxOutput
@param TxID - TxID of Transaction that the TxOutput resides in
@param OutputIndex - index of the TxOutput in the Transaction
@param Sig - data used in TxOutput PubKey
*/
type TxInput struct {
	TxID        []byte
	OutputIndex int
	Sig         string
}

/*TxOutput specifies coin value made available to a user
@param Amount - total
@param PubKey - ID of user
*/
type TxOutput struct {
	Amount int
	PubKey string
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
