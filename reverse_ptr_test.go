package ipx_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ns1/ipx/v2"
	"github.com/stretchr/testify/assert"
)

// ExampleReversePTR_v4 is an example of ReversePTR for IPv4.
func ExampleReversePTR_v4() {
	addr := net.ParseIP("192.168.0.10")
	fmt.Println(ipx.ReversePTR(addr))
	// Output:
	// 10.0.168.192.in-addr.arpa
}

// ExampleReversePTR_v6 is an example of ReversePTR for IPv6.
func ExampleReversePTR_v6() {
	addr := net.ParseIP("2001:db8::1")
	fmt.Println(ipx.ReversePTR(addr))
	// Output:
	// 1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa
}

// TestReversePTR unit tests for ReversePTR
func TestReversePTR(t *testing.T) {
	// IPv4
	assert.Equal(t, "10.0.168.192.in-addr.arpa",
		ipx.ReversePTR(net.ParseIP("192.168.0.10")))
	assert.Equal(t, "0.0.0.0.in-addr.arpa",
		ipx.ReversePTR(net.ParseIP("0.0.0.0")))
	assert.Equal(t, "255.255.255.255.in-addr.arpa",
		ipx.ReversePTR(net.ParseIP("255.255.255.255")))

	// IPv6
	assert.Equal(t, "4.3.3.7.0.7.3.0.e.2.a.8.0.0.0.0.0.0.0.0.3.a.5.8.8.b.d.0.1.0.0.2.ip6.arpa",
		ipx.ReversePTR(net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")))
	assert.Equal(t, "1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa",
		ipx.ReversePTR(net.ParseIP("2001:db8::1")))
	assert.Equal(t, "0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		ipx.ReversePTR(net.ParseIP("::")))
	assert.Equal(t, "f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.ip6.arpa",
		ipx.ReversePTR(net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")))

	// bad
	assert.Equal(t, "", ipx.ReversePTR(net.ParseIP("bad")))
	assert.Equal(t, "", ipx.ReversePTR(make(net.IP, 3)))
}
