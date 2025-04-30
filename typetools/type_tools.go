package typetools

import (
	"encoding/binary"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey/v2"
)

// Needs to be in little endian
func Uint16ToBytes(n uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, n)
	return b
}

func AccountIDToSS58(acc types.AccountID) string {
	recipientSS58 := subkey.SS58Encode(acc.ToBytes(), 42)
	return recipientSS58
}
