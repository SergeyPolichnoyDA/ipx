package ipx

import (
	u128 "github.com/Pilatuz/bigx/v2/uint128"
)

type Uint128 = u128.Uint128

// load128 reads uint128 integer from raw 16 bytes.
// big endian format is assumed.
func load128(buf []byte) Uint128 {
	return u128.LoadBigEndian(buf)
}

// store128 writes uint128 integer into the raw 16 bytes.
// big endian format is assumed.
func store128(u Uint128, buf []byte) {
	u128.StoreBigEndian(buf, u)
}
