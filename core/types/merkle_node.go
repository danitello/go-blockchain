package types

import "crypto/sha256"

/*MerkleNode represents one node in a MerkleTree */
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

/*InitMerkleNode creates a new instance of a node
@param left - the node that will be it's left node
@param right - "" right
@param data - the txn data that will go in this node
@return the node */
func InitMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	node.Left = left
	node.Right = right

	return &node
}
