package types

/*MerkleTree holds the root node of the representation */
type MerkleTree struct {
	Root *MerkleNode
}

/*InitMerkleTree creates an instance of a MerkleTree
@param data - the txns to be merked
@return the MerkleTree
*/
func InitMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	// Create nodes for each tx
	for _, hash := range data {
		node := InitMerkleNode(nil, nil, hash)
		nodes = append(nodes, *node)
	}

	// Create tree structure
	for i := 0; i < len(data)/2; i++ {
		var level []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := InitMerkleNode(&nodes[j], &nodes[j+1], nil)
			level = append(level, *node)
		}

		nodes = level
	}

	return &MerkleTree{&nodes[0]}
}
