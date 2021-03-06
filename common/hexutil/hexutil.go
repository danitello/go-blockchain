package hexutil

import (
	"bytes"
	"encoding/binary"
	"log"
)

// ToHex converts an integer into a []byte for compatibility (preceded by entries of 0 if necessary)
func ToHex(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}
