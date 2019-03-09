package byteutil

import (
	"bytes"
	"encoding/gob"

	"github.com/danitello/go-blockchain/common/errutil"
)

/*Serialize converts a struct instance to []byte for database compatibility
@param b - the struct instance to serialize
@return the serialized data
*/
func Serialize(class interface{}) []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(class)
	errutil.Handle(err)

	return result.Bytes()
}
