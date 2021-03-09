package ipx_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ns1/ipx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
