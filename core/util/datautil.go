package util

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/danitello/go-blockchain/core/types"
)

/*SerializeBlock converts a Block to []byte for database compatibility
@param b - the Block to serialize
@return the serialized data
*/
func SerializeBlock(b *types.Block) []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	HandleErr(err)

	return result.Bytes()
}

/*DeserializeBlock converts a []byte into a Block for database compatibility
@param data - the []byte representation of a Block
@returns a types.Block representation of a Block
*/
func DeserializeBlock(data []byte) *types.Block {
	var block types.Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	HandleErr(err)

	return &block
}

/*HandleErr displays errors to terminal
@param err - the error in question
*/
func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
