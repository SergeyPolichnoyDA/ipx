package ipx

import (
	"net"
)

// Address represents vanilla IP address.
type Address = net.IP

// Network represents vanilla IP network.
// It's IP range defined by CIDR notation.
type Network = net.IPNet

// Networks is a slice of networks.
type Networks []*Network

// Strings get all networks as string representation.
func (nn Networks) Strings() []string {
	if nn == nil {
		return nil
	}

	out := make([]string, 0, len(nn))
	for _, n := range nn {
		out = append(out, n.String())
	}

	return out
}

// Supernet returns a supernet for the provided network with the specified prefix length.
// targetPrefixLen is the number of "one" bits in new supernet.
func Supernet(network *Network, targetPrefixLen int) *Network {
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
		ip := to32(v4)
		mask := ((uint32(1) << targetPrefixLen) - 1) << (bits - targetPrefixLen)

		outIP := make(Address, net.IPv4len)
		from32(ip&mask, outIP)
		return &Network{
			IP:   outIP,
			Mask: net.CIDRMask(targetPrefixLen, bits),
		}
	}

	// IPv6
	if v6 := network.IP.To16(); v6 != nil {
		ip := to128(v6)
		mask := uint128{0, 1}.
			Lsh(uint(targetPrefixLen)).
			Minus(uint128{0, 1}).
			Lsh(uint(bits - targetPrefixLen))

		outIP := make(Address, net.IPv6len)
		from128(ip.And(mask), outIP)
		return &Network{
			IP:   outIP,
			Mask: net.CIDRMask(targetPrefixLen, bits),
		}
	}

	return nil // bad input address length
}

// Broadcast returns the broadcast IP address for the provided network.
func Broadcast(network *Network) Address {
	if network == nil {
		return nil // no network, no address
	}

	ones, bits := network.Mask.Size()

	// IPv4
	if v4 := network.IP.To4(); v4 != nil {
		ip := to32(v4)
		mask := (uint32(1) << (bits - ones)) - 1

		out := make(Address, net.IPv4len)
		from32(ip|mask, out)
		return out
	}

	// IPv6
	if v6 := network.IP.To16(); v6 != nil {
		ip := to128(network.IP)
		mask := uint128{0, 1}.
			Lsh(uint(bits - ones)).
			Minus(uint128{0, 1})

		out := make(Address, net.IPv6len)
		from128(ip.Or(mask), out)
		return out
	}

	return nil // bad input address length
}

// IsSubnet returns whether b is a subnet of a.
func IsSubnet(a, b *Network) bool {
	if !a.Contains(b.IP) {
		return false
	}

	aOnes, aBits := a.Mask.Size()
	bOnes, bBits := b.Mask.Size()
	return aBits == bBits && aOnes <= bOnes
}

// IsSupernet returns whether b is a supernet of a.
func IsSupernet(a, b *Network) bool {
	return IsSubnet(b, a)
}
