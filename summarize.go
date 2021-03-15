package ipx

import (
	"fmt"
	"math/bits"
	"net"

	u128 "github.com/Pilatuz/bigx/v2/uint128"
)

// SummarizeRange returns a series of networks which cover the range
// between the first and last addresses, inclusive.
func SummarizeRange(first, last net.IP) ([]*net.IPNet, error) {
	// first IPv4 or IPv6
	var firstV4, firstV6 net.IP
	switch len(first) {
	case net.IPv4len:
		firstV4 = first // first is IPv4

	case net.IPv6len:
		// note, even if IP address length is 128 it still can be IPv4!
		// need to do additional check converting it with `To4()`
		firstV4 = first.To4()
		if firstV4 == nil {
			firstV6 = first // first is IPv6
		}

	default:
		// invalid first IP address length
		return nil, fmt.Errorf("%w: first", ErrInvalidIP)
	}

	// last IPv4 or IPv6
	var lastV4, lastV6 net.IP
	switch len(last) {
	case net.IPv4len:
		lastV4 = last // last is IPv4

	case net.IPv6len:
		// note, even if IP address length is 128 it still can be IPv4!
		// need to do additional check converting it with `To4()`
		lastV4 = last.To4()
		if lastV4 == nil {
			lastV6 = last // last is IPv6
		}

	default:
		// invalid last IP address length
		return nil, fmt.Errorf("%w: last", ErrInvalidIP)
	}

	switch {
	case firstV4 != nil && lastV4 != nil:
		return summarizeRange4(to32(firstV4), to32(lastV4)), nil
	case firstV6 != nil && lastV6 != nil:
		return summarizeRange6(to128(firstV6), to128(lastV6)), nil
	}

	return nil, ErrVersionMismatch
}

// summarizeRange4 returns a series of IPv4 networks which cover the range
// between the first and last IPv4 addresses, inclusive.
func summarizeRange4(first, last uint32) (networks []*net.IPNet) {
	for first <= last {
		// the network will either be as long as all the trailing zeros of the first address OR the number of bits
		// necessary to cover the distance between first and last address -- whichever is smaller
		nBits := 32
		if z := bits.TrailingZeros32(first); z < nBits {
			nBits = z
		}

		if first != 0 || last != maxUint32 { // guard overflow; this would just be 32 anyway
			d := last - first + 1
			if z := 31 - bits.LeadingZeros32(d); z < nBits {
				nBits = z
			}
		}

		nwkMask := net.CIDRMask(32-nBits, 32)
		nwkIP := make(net.IP, net.IPv4len)
		from32(first, nwkIP)
		networks = append(networks,
			&net.IPNet{
				IP:   nwkIP,
				Mask: nwkMask,
			})

		first += 1 << nBits
		if first == 0 {
			break
		}
	}

	return
}

// summarizeRange6 returns a series of IPv6 networks which cover the range
// between the first and last IPv6 addresses, inclusive.
func summarizeRange6(first, last Uint128) (networks []*net.IPNet) {
	for first.Cmp(last) <= 0 { // first <= last
		// the network will either be as long as all the trailing zeros of the first address OR the number of bits
		// necessary to cover the distance between first and last address -- whichever is smaller
		nBits := 128
		if z := first.TrailingZeros(); z < nBits {
			nBits = z
		}

		// check extremes to make sure no overflow
		if !first.IsZero() || !last.Equals(u128.Max()) {
			d := last.Sub(first).Add64(1)
			if z := 127 - d.LeadingZeros(); z < nBits {
				nBits = z
			}
		}

		nwkMask := net.CIDRMask(128-nBits, 128)
		nwkIP := make(net.IP, net.IPv6len)
		from128(first, nwkIP)
		networks = append(networks,
			&net.IPNet{
				IP:   nwkIP,
				Mask: nwkMask,
			})

		first = first.Add(Uint128{Lo: 1}.Lsh(uint(nBits)))
		if first.IsZero() {
			break
		}
	}

	return
}
