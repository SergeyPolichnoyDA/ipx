package ipx

import (
	"net"
)

// Network represents vanilla IP network.
// It's represented as CIDR notation.
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

// Supernet returns a supernet for the provided network with the specified prefix length
func Supernet(ipN *net.IPNet, newPrefix int) *net.IPNet {
	ones, bits := ipN.Mask.Size()
	if newPrefix < 0 || newPrefix > ones || newPrefix > bits {
		return nil
	}

	out := net.IPNet{
		IP:   make(net.IP, len(ipN.IP)),
		Mask: net.CIDRMask(newPrefix, bits),
	}
	if ipN.IP.To4() != nil {
		ip := to32(ipN.IP)
		ip &= (1<<newPrefix - 1) << (bits - newPrefix)
		from32(ip, out.IP)
		return &out
	}

	ip := to128(ipN.IP)

	mask := uint128{0, 1}.
		Lsh(uint(newPrefix)).
		Minus(uint128{0, 1}).
		Lsh(uint(bits - newPrefix))

	from128(ip.And(mask), out.IP)

	return &out
}

// Broadcast returns the broadcast address for the provided net.
func Broadcast(a *net.IPNet) net.IP {
	out := make(net.IP, len(a.IP))
	copy(out, a.IP)

	ones, bits := a.Mask.Size()

	if a.IP.To4() != nil {
		ip := to32(a.IP)

		ip |= 1<<(bits-ones) - 1

		from32(ip, out)

		return out
	}

	ip := to128(a.IP)

	hostMask := uint128{0, 1}.
		Lsh(uint(bits - ones)).
		Minus(uint128{0, 1})

	from128(ip.Or(hostMask), out)

	return out
}

// IsSubnet returns whether b is a subnet of a
func IsSubnet(a, b *net.IPNet) bool {
	return a.Contains(b.IP) && maskPrefix(a.Mask, b.Mask)
}

// IsSupernet returns whether b is a supernet of a
func IsSupernet(a, b *net.IPNet) bool {
	return IsSubnet(b, a)
}

func maskPrefix(a, b net.IPMask) bool {
	aOnes, aBits := a.Size()
	bOnes, bBits := b.Size()
	return aBits == bBits && aOnes <= bOnes
}
