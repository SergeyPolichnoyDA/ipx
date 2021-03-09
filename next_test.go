package ipx_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ns1/ipx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExampleNextAddress_v4 is an example of NextAddress for IPv4.
func ExampleNextAddress_v4() {
	addr1 := net.ParseIP("0.0.0.1")
	addr2 := ipx.NextAddress(addr1, 257)
	fmt.Println(addr1, "+ 257 =", addr2)
	// Output:
	// 0.0.0.1 + 257 = 0.0.1.2
}

// ExampleNextAddress_v6 is an example of NextAddress for IPv6.
func ExampleNextAddress_v6() {
	addr1 := net.ParseIP("::1")
	addr2 := ipx.NextAddress(addr1, 65537)
	fmt.Println(addr1, "+ 65537 =", addr2)
	// Output:
	// ::1 + 65537 = ::1:2
}

// TestNextAddress unit tests for NextAddress
func TestNextAddress(t *testing.T) {
	// bad input
	assert.Nil(t, ipx.NextAddress(nil, 1))
	assert.Nil(t, ipx.NextAddress(make(net.IP, 3), 1))

	// step=0 means the same output
	assert.Equal(t, net.IPv4bcast, ipx.NextAddress(net.IPv4bcast, 0))
	assert.Equal(t, net.IPv6zero, ipx.NextAddress(net.IPv6zero, 0))

	// IPv4
	assert.Equal(t, "0.0.0.1", ipx.NextAddress(net.ParseIP("0.0.0.0"), +1).String())
	assert.Equal(t, "0.0.0.0", ipx.NextAddress(net.ParseIP("0.0.0.1"), -1).String())
	assert.Equal(t, "0.0.0.0", ipx.NextAddress(net.ParseIP("255.255.255.255"), +1).String())
	assert.Equal(t, "255.255.255.255", ipx.NextAddress(net.ParseIP("0.0.0.0"), -1).String())
	assert.Equal(t, "0.0.1.1", ipx.NextAddress(net.ParseIP("0.0.0.0"), +257).String())
	assert.Equal(t, "0.0.0.0", ipx.NextAddress(net.ParseIP("0.0.1.1"), -257).String())

	// IPv6
	assert.Equal(t, "::1", ipx.NextAddress(net.ParseIP("::"), +1).String())
	assert.Equal(t, "::", ipx.NextAddress(net.ParseIP("::1"), -1).String())
	assert.Equal(t, "::", ipx.NextAddress(net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"), +1).String())
	assert.Equal(t, "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", ipx.NextAddress(net.ParseIP("::"), -1).String())
	assert.Equal(t, "::1:0:0", ipx.NextAddress(net.ParseIP("::"), +(1<<32)).String())
	assert.Equal(t, "::", ipx.NextAddress(net.ParseIP("::1:0:0"), -(1<<32)).String())
}

// BenchmarkNextAddress performance benchmarks for NextAddress
func BenchmarkNextAddress(bb *testing.B) {
	// helper function to run benchmark
	bench := func(address string, step int) func(*testing.B) {
		return func(b *testing.B) {
			ip := net.ParseIP(address)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				out := ipx.NextAddress(ip, step)
				require.NotNil(b, out)
			}
		}
	}

	// IPv4
	bb.Run("ipv4_1", bench("10.0.0.0", 1))
	bb.Run("ipv4_2", bench("10.0.0.0", 2))

	// IPv6
	bb.Run("ipv6_1", bench("::", 1))
	bb.Run("ipv6_2", bench("::", 2))
}

// ExampleNextNetwork_v4 is an example of NextNetwork for IPv4
func ExampleNextNetwork_v4() {
	_, net1, _ := net.ParseCIDR("10.0.0.0/16")
	net2 := ipx.NextNetwork(net1, 2)
	fmt.Println(net1, "+ 2 =", net2)
	// Output:
	// 10.0.0.0/16 + 2 = 10.2.0.0/16
}

// ExampleNextNetwork_v6 is an example of NextNetwork for IPv6
func ExampleNextNetwork_v6() {
	_, net1, _ := net.ParseCIDR("0:1::/32")
	net2 := ipx.NextNetwork(net1, 2)
	fmt.Println(net1, "+ 2 =", net2)
	// Output:
	// 0:1::/32 + 2 = 0:3::/32
}

// TestNextNetwork unit tests for NextNetwork
func TestNextNetwork(t *testing.T) {
	// bad input
	assert.Nil(t, ipx.NextNetwork(nil, 1))
	assert.Nil(t, ipx.NextNetwork(&net.IPNet{IP: make(net.IP, 3)}, 1))

	// step=0 means the same output
	assert.Equal(t, &net.IPNet{}, ipx.NextNetwork(&net.IPNet{}, 0))

	// IPv4
	assert.Equal(t, "10.2.0.0/16", ipx.NextNetwork(cidr("10.0.0.0/16"), +2).String())
	assert.Equal(t, "10.0.0.0/16", ipx.NextNetwork(cidr("10.2.0.0/16"), -2).String())

	// IPv6
	assert.Equal(t, "0:2::/32", ipx.NextNetwork(cidr("::/32"), +2).String())
	assert.Equal(t, "::/32", ipx.NextNetwork(cidr("0:2::/32"), -2).String())
}

// BenchmarkNextNetwork performance benchmarks for NextNetwork
func BenchmarkNextNetwork(bb *testing.B) {
	// helper function to run benchmark
	bench := func(network string, step int) func(*testing.B) {
		return func(b *testing.B) {
			_, nwk, err := net.ParseCIDR(network)
			require.NoError(b, err, "failed to parse CIDR")

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				out := ipx.NextNetwork(nwk, step)
				require.NotNil(b, out)
			}
		}
	}

	// IPv4
	bb.Run("ipv4_1", bench("10.0.0.0/30", 1))
	bb.Run("ipv4_2", bench("10.0.0.0/30", 2))
	bb.Run("ipv4_3", bench("10.0.0.0/24", 1))

	// IPv6
	bb.Run("ipv6_1", bench("::/126", 1))
	bb.Run("ipv6_1", bench("::/32", 1))
}
