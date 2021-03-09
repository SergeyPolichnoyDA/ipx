package ipx

import (
	"errors"
	"net"
)

// IncrIP returns the next IP
func IncrIP(ip net.IP, incr int) {
	if ip == nil {
		panic(errors.New("IP cannot be nil"))
	}
	if ip.To4() != nil {
		n := to32(ip)
		if incr >= 0 {
			n += uint32(incr)
		} else {
			n -= uint32(incr * -1)
		}
		from32(n, ip)
		return
	}

	// ipv6
	u := to128(ip)
	if incr >= 0 {
		u = u.Add(uint128{0, uint64(incr)})
	} else {
		u = u.Minus(uint128{0, uint64(incr * -1)})
	}
	from128(u, ip)
}

// IncrNet steps to the next net of the same mask
func IncrNet(ipNet *net.IPNet, incr int) {
	if ipNet.IP == nil {
		panic(errors.New("IP cannot be nil"))
	}
	if ipNet.Mask == nil {
		panic(errors.New("mask cannot be nil"))
	}
	if ipNet.IP.To4() != nil {
		n := to32(ipNet.IP)
		ones, bits := ipNet.Mask.Size()
		suffix := uint32(bits - ones)
		n >>= suffix
		if incr >= 0 {
			n += uint32(incr)
		} else {
			n -= uint32(incr * -1)
		}
		from32(n<<suffix, ipNet.IP)
		return
	}

	b := to128(ipNet.IP)

	ones, bits := ipNet.Mask.Size()
	suffix := uint(bits - ones)

	b = b.Rsh(suffix)
	if incr >= 0 {
		b = b.Add(uint128{0, uint64(incr)})
	} else {
		b = b.Minus(uint128{0, uint64(incr * -1)})
	}
	b = b.Lsh(suffix)

	from128(b, ipNet.IP)
}
