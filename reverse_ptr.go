package ipx

import (
	"bytes"
	"net"
	"strconv"
)

var (
	// is used to convert 4 bits integer into a rune
	nibbles = [...]rune{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
	}
)

// ReversePTR returns the name of the reverse DNS PTR record for the given IP address.
func ReversePTR(address net.IP) string {
	var buf bytes.Buffer

	// IPv4
	if v4 := address.To4(); v4 != nil {
		// write IP address bytes (in reverse order)
		for i := len(v4) - 1; i >= 0; i-- {
			buf.WriteString(strconv.Itoa(int(v4[i])))
			buf.WriteRune('.')
		}

		buf.WriteString("in-addr.arpa")
		return buf.String()
	}

	// IPv6
	if v6 := address.To16(); v6 != nil {
		// write IP address (in reverse order)
		for i := len(v6) - 1; i >= 0; i-- {
			buf.WriteRune(nibbles[v6[i]&0x0F])
			buf.WriteRune('.')
			buf.WriteRune(nibbles[v6[i]>>4])
			buf.WriteRune('.')
		}

		buf.WriteString("ip6.arpa")
		return buf.String()
	}

	return "" // bad address length
}
