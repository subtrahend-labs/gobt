package typetools

import "encoding/binary"

// Needs to be in little endian
func Uint16ToBytes(n uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, n)
	return b
}
