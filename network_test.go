package ipx_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ns1/ipx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExampleSupernet is an example for Supernet
func ExampleSupernet() {
	_, nwk, _ := net.ParseCIDR("192.0.2.100/24")
	fmt.Println(ipx.Supernet(nwk, 16))
	// Output:
	// 192.0.0.0/16
}

// TestSupernet unit tests for Supernet
func TestSupernet(tt *testing.T) {
	tt.Run("bad", func(t *testing.T) {
		// no network, no supernet
		supernet := ipx.Supernet(nil, 0)
		assert.Nil(t, supernet)

		_, nwk, err := net.ParseCIDR("10.0.0.128/25")
		require.NoError(t, err, "failed to parse CIDR")
		supernet = ipx.Supernet(nwk, 26)
		assert.Nil(t, supernet)

		// bad address length
		nwk.IP = make(net.IP, 3)
		supernet = ipx.Supernet(nwk, 24)
		assert.Nil(t, supernet)
	})

	tt.Run("ipv4_one_level", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("10.0.0.128/25")
		require.NoError(t, err, "failed to parse CIDR")
		supernet := ipx.Supernet(nwk, 24)
		assert.Equal(t, "10.0.0.0/24", supernet.String())
	})

	tt.Run("ipv4_same_level", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("10.0.0.128/25")
		require.NoError(t, err, "failed to parse CIDR")
		supernet := ipx.Supernet(nwk, 25)
		assert.Equal(t, "10.0.0.128/25", supernet.String())
	})

	tt.Run("ipv4_8", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("10.0.0.128/25")
		require.NoError(t, err, "failed to parse CIDR")
		supernet := ipx.Supernet(nwk, 8)
		assert.Equal(t, "10.0.0.0/8", supernet.String())
	})

	tt.Run("ipv4_all", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("10.0.0.128/25")
		require.NoError(t, err, "failed to parse CIDR")
		supernet := ipx.Supernet(nwk, 0)
		assert.Equal(t, "0.0.0.0/0", supernet.String())
	})

	tt.Run("ipv6_one_level", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("29a2:241a:f62c::/64")
		require.NoError(t, err, "failed to parse CIDR")
		supernet := ipx.Supernet(nwk, 44)
		assert.Equal(t, "29a2:241a:f620::/44", supernet.String())
	})

	tt.Run("ipv6_same_level", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("29a2:241a:f62c::/64")
		require.NoError(t, err, "failed to parse CIDR")
		supernet := ipx.Supernet(nwk, 64)
		assert.Equal(t, "29a2:241a:f62c::/64", supernet.String())
	})

	tt.Run("ipv6_16", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("29a2:241a:f62c::/64")
		require.NoError(t, err, "failed to parse CIDR")
		supernet := ipx.Supernet(nwk, 16)
		assert.Equal(t, "29a2::/16", supernet.String())
	})

	tt.Run("ipv6_all", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("29a2:241a:f62c::/64")
		require.NoError(t, err, "failed to parse CIDR")
		supernet := ipx.Supernet(nwk, 0)
		assert.Equal(t, "::/0", supernet.String())
	})
}

// BenchmarkSupernet performance benchmarks for Supernet
func BenchmarkSupernet(bb *testing.B) {
	// helper function to run benchmark
	bench := func(network string, targetPrefixLen int) func(*testing.B) {
		return func(b *testing.B) {
			_, nwk, err := net.ParseCIDR(network)
			require.NoError(b, err, "failed to parse CIDR")

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				out := ipx.Supernet(nwk, targetPrefixLen)
				require.NotNil(b, out, "failed to get supernet")
			}
		}
	}

	// IPv4
	bb.Run("ipv4_20", bench("192.0.2.0/24", 20))
	bb.Run("ipv4_15", bench("192.0.2.0/24", 15))

	// IPv6
	bb.Run("ipv6_20", bench("::/24", 20))
	bb.Run("ipv6_15", bench("::/24", 15))
}

// ExampleBroadcast is an example for Broadcast
func ExampleBroadcast() {
	_, nwk, _ := net.ParseCIDR("10.0.1.0/24")
	fmt.Println(ipx.Broadcast(nwk))
	// Output:
	// 10.0.1.255
}

// TestBroadcast unit tests for Broadcast
func TestBroadcast(tt *testing.T) {
	tt.Run("bad", func(t *testing.T) {
		out := ipx.Broadcast(nil)
		assert.Nil(t, out)

		// bad address length
		out = ipx.Broadcast(&net.IPNet{IP: make(net.IP, 3)})
		assert.Nil(t, out)
	})

	tt.Run("ipv4", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("10.0.1.0/24")
		require.NoError(t, err, "failed to parse CIDR")

		out := ipx.Broadcast(nwk)
		assert.Equal(t, "10.0.1.255", out.String())
	})

	tt.Run("ipv6", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("29a2:241a:f620::/44")
		require.NoError(t, err, "failed to parse CIDR")

		out := ipx.Broadcast(nwk)
		assert.Equal(t, "29a2:241a:f62f:ffff:ffff:ffff:ffff:ffff", out.String())
	})
}

// BenchmarkBroadcast performance benchmarks for Broadcast
func BenchmarkBroadcast(bb *testing.B) {
	// helper function to run benchmark
	bench := func(network string) func(*testing.B) {
		return func(b *testing.B) {
			_, nwk, err := net.ParseCIDR(network)
			require.NoError(b, err, "failed to parse CIDR")

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				out := ipx.Broadcast(nwk)
				require.NotNil(b, out, "failed to get broadcast")
			}
		}
	}

	// IPv4
	bb.Run("ipv4_31", bench("10.0.1.0/31"))
	bb.Run("ipv4_16", bench("10.1.0.0/16"))
	bb.Run("ipv4_all", bench("0.0.0.0/0"))

	// IPv6
	bb.Run("ipv6_127", bench("::/127"))
	bb.Run("ipv6_64", bench("::/64"))
	bb.Run("ipv6_all", bench("::/0"))
}

// ExampleIsSubnet is an example of IsSubnet
func ExampleIsSubnet() {
	_, a, _ := net.ParseCIDR("10.0.0.0/16")
	_, b, _ := net.ParseCIDR("10.0.1.0/24")
	fmt.Println(ipx.IsSubnet(a, b))
	fmt.Println(ipx.IsSubnet(a, a))
	fmt.Println(ipx.IsSubnet(b, a))
	// Output:
	// true
	// true
	// false
}

// ExampleIsSupernet is an example for IsSupernet
func ExampleIsSupernet() {
	_, a, _ := net.ParseCIDR("10.0.0.0/16")
	_, b, _ := net.ParseCIDR("10.0.1.0/24")
	fmt.Println(ipx.IsSupernet(a, b))
	fmt.Println(ipx.IsSupernet(a, a))
	fmt.Println(ipx.IsSupernet(b, a))
	// Output:
	// false
	// true
	// true
}

func cidr(cidrS string) *net.IPNet {
	_, ipNet, _ := net.ParseCIDR(cidrS)
	return ipNet
}
