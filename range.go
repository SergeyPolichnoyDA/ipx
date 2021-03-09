package ipx

// Range represents [first, last] IP range.
// Range can represent any IP range comparing to Network
// which can represent only ranges defined by CIDR notation.
// So Network is subset of Range.
//
// Note, first and last addresses are included into the IP range!
// If `first == last` IP range contains single address.
// If `first > last` IP range considered as empty.
type Range struct {
	First Address `json:"first,string"`
	Last  Address `json:"last,string"`
}

// NewRange is helper function to construct IP range.
func NewRange(first Address, last Address) Range {
	return Range{
		First: first,
		Last:  last,
	}
}

// Summarize returns a series of networks which cover the range.
func (r Range) Summarize() (Networks, error) {
	return SummarizeRange(r.First, r.Last)
}

// RangeFromNetwork returns the IP range for the given network.
// The first and last IP addresses are inclusive.
// Usually the first IP address is network address,
// while the last IP address is broadcast address.
func RangeFromNetwork(network *Network) (first Address, last Address) {
	if network == nil {
		return // no network, no IP range
	}

	n := len(network.IP)
	if n != len(network.Mask) {
		return // inconsistent length: IP address and IP mask
	}

	first = make(Address, n)
	last = make(Address, n)
	for i := 0; i < n; i++ {
		first[i] = network.IP[i] & network.Mask[i]
		last[i] = network.IP[i] | ^network.Mask[i]
	}

	return
}
