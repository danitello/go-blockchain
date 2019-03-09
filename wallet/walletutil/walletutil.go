package walletutil

import (
	"github.com/danitello/go-blockchain/common/errutil"
	"github.com/mr-tron/base58"
)

/*Base58Encode encodes a byte array to base58
@param input - the byte array
@return encoded input
*/
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

/*Base58Decode decodes base58 encoded input
@param input - the byte array
@return encoded input
*/
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	errutil.Handle(err)

	return decode
}
