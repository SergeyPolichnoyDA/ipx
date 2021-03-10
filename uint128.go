package ipx

import (
	"encoding/binary"
	"math/bits"
)

// largely cribbed from https://github.com/davidminor/uint128 and https://github.com/lukechampine/uint128
type uint128 struct {
	H, L uint64
}

func (u uint128) And(other uint128) uint128 {
	u.H &= other.H
	u.L &= other.L
	return u
}

func (u uint128) Or(other uint128) uint128 {
	u.H |= other.H
	u.L |= other.L
	return u
}

func (u uint128) Cmp(other uint128) int {
	switch {
	case u.H > other.H:
		return 1
	case u.H < other.H,
		u.L < other.L:
		return -1
	case u.L > other.L:
		return 1
	default:
		return 0
	}
}

// Equal checks 128-bit values are equal.
// Equivalent of `u.Cmp(other) == 0` but a bit faster.
func (u uint128) Equal(other uint128) bool {
	return (u.H == other.H) &&
		(u.L == other.L)
}

// Add does addition of two uint128 integers.
func (u uint128) Add(v uint128) uint128 {
	oldL := u.L
	u.H += v.H
	u.L += v.L
	if u.L < oldL { // wrapped
		u.H++
	}
	return u
}

// Add64 does addition of uint128 and uint64 integers.
func (u uint128) Add64(v uint64) uint128 {
	oldL := u.L
	// u.H += 0
	u.L += v
	if u.L < oldL { // wrapped
		u.H++
	}
	return u
}

// Sub does subtraction of two uint128 integers.
func (u uint128) Sub(v uint128) uint128 {
	oldL := u.L
	u.H -= v.H
	u.L -= v.L
	if u.L > oldL { // wrapped
		u.H--
	}
	return u
}

// Sub64 does subtraction of uint128 and uint64 integers.
func (u uint128) Sub64(v uint64) uint128 {
	oldL := u.L
	// u.H -= 0
	u.L -= v
	if u.L > oldL { // wrapped
		u.H--
	}
	return u
}

func (u uint128) Lsh(bits uint) uint128 {
	switch {
	case bits >= 128:
		u.H, u.L = 0, 0
	case bits >= 64:
		u.H, u.L = u.L<<(bits-64), 0
	default:
		u.H <<= bits
		u.H |= u.L >> (64 - bits) // set top with prefix that cross from bottom
		u.L <<= bits
	}
	return u
}

func (u uint128) Rsh(bits uint) uint128 {
	switch {
	case bits >= 128:
		u.H, u.L = 0, 0
	case bits >= 64:
		u.H, u.L = 0, u.H>>(bits-64)
	default:
		u.L >>= bits
		u.L |= u.H << (64 - bits) // set bottom with prefix that cross from top
		u.H >>= bits
	}
	return u
}

func (u uint128) Not() uint128 {
	return uint128{^u.H, ^u.L}
}

// TrailingZeros counts the number of trailing zeros.
// Just as standard `bits.TrailingZerosXX` does.
func (u uint128) TrailingZeros() int {
	z := bits.TrailingZeros64(u.L)
	if z == 64 {
		z += bits.TrailingZeros64(u.H)
	}
	return z
}

// LeadingZeros counts the number of leading zeros.
// Just as standard `bits.LeadingZerosXX` does.
func (u uint128) LeadingZeros() int {
	z := bits.LeadingZeros64(u.H)
	if z == 64 {
		z += bits.LeadingZeros64(u.L)
	}
	return z
}

// to128 reads uint128 integer from raw 16 bytes.
// big endian format is assumed.
func to128(buf []byte) uint128 {
	return uint128{
		H: binary.BigEndian.Uint64(buf[:8]),
		L: binary.BigEndian.Uint64(buf[8:]),
	}
}

// from128 writes uint128 integer into the raw 16 bytes.
// big endian format is assumed.
func from128(u uint128, buf []byte) {
	binary.BigEndian.PutUint64(buf[:8], u.H)
	binary.BigEndian.PutUint64(buf[8:], u.L)
}
