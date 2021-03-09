package ipx

import (
	"net"
)

// NextAddress returns the next IP address.
// The step argument can be positive returning next address,
// or negative returning previous address.
// No network mask applied.
func NextAddress(address net.IP, step int) net.IP {
	if len(address) == 0 || step == 0 {
		return address // address is the same
	}

	// IPv4
	if v4 := address.To4(); v4 != nil {
		u := to32(v4)
		if step > 0 {
			u += uint32(+step)
		} else {
			u -= uint32(-step)
		}

		out := make(net.IP, net.IPv4len)
		from32(u, out)
		return out
	}

	// IPv6
	if v6 := address.To16(); v6 != nil {
		u := to128(v6)
		if step > 0 {
			u = u.Add(uint128{0, uint64(+step)})
		} else {
			u = u.Sub(uint128{0, uint64(-step)})
		}

		out := make(net.IP, net.IPv6len)
		from128(u, out)
		return out
	}

	return nil // bad input address length
}

// NextNetwork returns the next network of the same mask.
// The step argument can be positive returning next network,
// or negative returning previous network.
func NextNetwork(network *Network, step int) *Network {
	if network == nil || step == 0 {
		return network // network is the same
	}

	// IPv4
	if v4 := network.IP.To4(); v4 != nil {
		ones, bits := network.Mask.Size()
		suffix := uint(bits - ones)

		u := to32(v4)
		if step > 0 {
			u += uint32(+step << suffix)
		} else {
			u -= uint32(-step << suffix)
		}

		outIP := make(net.IP, net.IPv4len)
		from32(u, outIP)
		return &Network{
			IP:   outIP,
			Mask: network.Mask,
		}
	}

	// IPv6
	if v6 := network.IP.To16(); v6 != nil {
		ones, bits := network.Mask.Size()
		suffix := uint(bits - ones)

		u := to128(v6)
		if step > 0 {
			u = u.Add(uint128{0, uint64(+step)}.Lsh(suffix))
		} else {
			u = u.Sub(uint128{0, uint64(-step)}.Lsh(suffix))
		}

		outIP := make(net.IP, net.IPv6len)
		from128(u, outIP)
		return &Network{
			IP:   outIP,
			Mask: network.Mask,
		}
	}

	return nil // bad input address length
}
