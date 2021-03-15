package ipx

import (
	"net"
)

// CompareIP compares two IP addresses.
// Returns:
//  - `0` if a == b,
//  - `-1` if a < b,
//  - `+1` if a > b.
func CompareIP(a, b net.IP) (int, error) {
	// assume both are IPv4
	if a4, b4 := a.To4(), b.To4(); a4 != nil && b4 != nil {
		a, b := load32(a4), load32(b4)
		if a < b {
			return -1, nil // a < b
		} else if a > b {
			return +1, nil // a > b
		}
		return 0, nil // a == b
	}

	// assume both or at least one is IPv6
	if a6, b6 := a.To16(), b.To16(); a6 != nil && b6 != nil {
		a, b := load128(a6), load128(b6)
		return a.Cmp(b), nil
	}

	return 0, ErrInvalidIP // either a or b or both
}

// MustCompareIP compares two IP addresses.
// The same as CompareIP() but panics in case of bad input.
func MustCompareIP(a, b net.IP) int {
	cmp, err := CompareIP(a, b)
	if err != nil {
		panic(err)
	}
	return cmp
}

// NextIP returns the next IP address.
// The step argument can be positive returning next address,
// or negative returning previous address.
// No network mask applied so it worth to call Mask()
// after to fit resulting address into desired network.
// Returns nil on bad input.
func NextIP(addr net.IP, step int) net.IP {
	if step == 0 {
		return addr // the same address
	}

	// IPv4
	if v4 := addr.To4(); v4 != nil {
		u := load32(v4)
		if step > 0 {
			u += uint32(+step)
		} else {
			u -= uint32(-step)
		}

		out := make(net.IP, net.IPv4len)
		store32(u, out)
		return out
	}

	// IPv6
	if v6 := addr.To16(); v6 != nil {
		u := load128(v6)
		if step > 0 {
			u = u.Add64(uint64(+step))
		} else {
			u = u.Sub64(uint64(-step))
		}

		out := make(net.IP, net.IPv6len)
		store128(u, out)
		return out
	}

	return nil // bad input address length
}
