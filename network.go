package ipx

import (
	"net"
)

// Supernet returns a supernet for the provided network with the specified prefix length.
// targetPrefixLen is the number of "one" bits in new supernet mask.
func Supernet(network *net.IPNet, targetPrefixLen int) *net.IPNet {
	if network == nil {
		return nil // no network, no supernet
	}

	ones, bits := network.Mask.Size()
	if targetPrefixLen < 0 ||
		targetPrefixLen > ones ||
		targetPrefixLen > bits {
		return nil // invalid target prefix length
	}

	// IPv4
	if v4 := network.IP.To4(); v4 != nil {
		u := load32(v4)
		mask := ((uint32(1) << targetPrefixLen) - 1) << (bits - targetPrefixLen)

		out := make(net.IP, net.IPv4len)
		store32(u&mask, out)
		return &net.IPNet{
			IP:   out,
			Mask: net.CIDRMask(targetPrefixLen, bits),
		}
	}

	// IPv6
	if v6 := network.IP.To16(); v6 != nil {
		u := load128(v6)
		mask := Uint128{Lo: 1}.
			Lsh(uint(targetPrefixLen)).
			Sub64(1).
			Lsh(uint(bits - targetPrefixLen))

		out := make(net.IP, net.IPv6len)
		store128(u.And(mask), out)
		return &net.IPNet{
			IP:   out,
			Mask: net.CIDRMask(targetPrefixLen, bits),
		}
	}

	return nil // bad input address length
}

// Broadcast returns the broadcast IP address for the provided network.
func Broadcast(network *net.IPNet) net.IP {
	if network == nil {
		return nil // no network, no address
	}

	ones, bits := network.Mask.Size()

	// IPv4
	if v4 := network.IP.To4(); v4 != nil {
		u := load32(v4)
		bcmask := (uint32(1) << (bits - ones)) - 1

		out := make(net.IP, net.IPv4len)
		store32(u|bcmask, out)
		return out
	}

	// IPv6
	if v6 := network.IP.To16(); v6 != nil {
		u := load128(network.IP)
		bcmask := Uint128{Lo: 1}.
			Lsh(uint(bits - ones)).
			Sub64(1)

		out := make(net.IP, net.IPv6len)
		store128(u.Or(bcmask), out)
		return out
	}

	return nil // bad input address length
}

// IsSubnet returns whether b is a subnet of a.
func IsSubnet(a, b *net.IPNet) bool {
	if !a.Contains(b.IP) {
		return false
	}

	aOnes, aBits := a.Mask.Size()
	bOnes, bBits := b.Mask.Size()
	return aBits == bBits && aOnes <= bOnes
}

// IsSupernet returns whether b is a supernet of a.
func IsSupernet(a, b *net.IPNet) bool {
	return IsSubnet(b, a)
}

// NextNetwork returns the next network of the same mask.
// The step argument can be positive returning next network,
// or negative returning previous network.
func NextNetwork(network *net.IPNet, step int) *net.IPNet {
	if network == nil || step == 0 {
		return network // network is the same
	}

	// IPv4
	if v4 := network.IP.To4(); v4 != nil {
		ones, bits := network.Mask.Size()
		suffix := uint(bits - ones)

		u := load32(v4)
		if step > 0 {
			u += uint32(+step << suffix)
		} else {
			u -= uint32(-step << suffix)
		}

		out := make(net.IP, net.IPv4len)
		store32(u, out)
		return &net.IPNet{
			IP:   out,
			Mask: network.Mask, // shared
		}
	}

	// IPv6
	if v6 := network.IP.To16(); v6 != nil {
		ones, bits := network.Mask.Size()
		suffix := uint(bits - ones)

		u := load128(v6)
		if step > 0 {
			u = u.Add(Uint128{Lo: uint64(+step)}.Lsh(suffix))
		} else {
			u = u.Sub(Uint128{Lo: uint64(-step)}.Lsh(suffix))
		}

		out := make(net.IP, net.IPv6len)
		store128(u, out)
		return &net.IPNet{
			IP:   out,
			Mask: network.Mask, // shared
		}
	}

	return nil // bad input address length
}
