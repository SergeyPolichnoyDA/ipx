package ipx

import (
	"net"
)

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

		u := to32(v4)
		if step > 0 {
			u += uint32(+step << suffix)
		} else {
			u -= uint32(-step << suffix)
		}

		outIP := make(net.IP, net.IPv4len)
		from32(u, outIP)
		return &net.IPNet{
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
		return &net.IPNet{
			IP:   outIP,
			Mask: network.Mask,
		}
	}

	return nil // bad input address length
}
