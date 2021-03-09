package ipx

import (
	"encoding/binary"
)

// to32 reads uint32 integer from raw 4 bytes.
// big endian format is assumed.
func to32(buf []byte) uint32 {
	l := len(buf) // TODO: consider to accept only 4 bytes!
	return binary.BigEndian.Uint32(buf[l-4:])
}

// from32 writes uint32 integer into the raw 4 bytes.
// big endian format is assumed.
func from32(n uint32, buf []byte) {
	l := len(buf) // TODO: consider to accpent only 4 bytes!
	binary.BigEndian.PutUint32(buf[l-4:], n)
}
