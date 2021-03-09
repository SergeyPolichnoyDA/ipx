package ipx_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ns1/ipx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExampleRangeFromNetwork_v4 is an example of RangeFromNetwork for IPv4
func ExampleRangeFromNetwork_v4() {
	_, nwk, _ := net.ParseCIDR("192.168.1.100/24")
	fmt.Println(ipx.RangeFromNetwork(nwk))
	// Output:
	// 192.168.1.0 192.168.1.255
}

// ExampleRangeFromNetwork_v6 is an example of RangeFromNetwork for IPv6
func ExampleRangeFromNetwork_v6() {
	_, nwk, _ := net.ParseCIDR("::100/96")
	fmt.Println(ipx.RangeFromNetwork(nwk))
	// Output:
	// :: ::ffff:ffff
}

// BenchmarkRangeFromNetwork performance benchmarks for RangeFromNetwork
func BenchmarkRangeFromNetwork(bb *testing.B) {
	// helper function to run benchmark
	bench := func(network string) func(*testing.B) {
		return func(b *testing.B) {
			_, nwk, err := net.ParseCIDR(network)
			require.NoError(b, err, "failed to parse CIDR")

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				first, last := ipx.RangeFromNetwork(nwk)
				require.NotNil(b, first)
				require.NotNil(b, last)
			}
		}
	}

	bb.Run("ipv4", bench("192.168.0.0/23"))
	bb.Run("ipv6", bench("2001:db8::8a2e:370:7334/107"))
}

// TestNewRange unit tests for NewRange
func TestNewRange(t *testing.T) {
	r := ipx.NewRange(net.IPv4zero, net.IPv4bcast)
	assert.Equal(t, net.IPv4zero, r.First)
	assert.Equal(t, net.IPv4bcast, r.Last)
}

// TestRangeSummarize unit tests for Range.Summarize
func TestRangeSummarize(t *testing.T) {
	r := ipx.NewRange(net.IPv4zero, net.IPv4bcast)
	nwks, err := r.Summarize()
	require.NoError(t, err, "failed to summarize IP range")
	assert.Equal(t, []string{"0.0.0.0/0"}, nwks.Strings())
}

// TestRangeFromNetwork unit tests for RangeFromNetwork
func TestRangeFromNetwork(tt *testing.T) {
	tt.Run("bad", func(t *testing.T) {
		first, last := ipx.RangeFromNetwork(nil)
		assert.Nil(t, first)
		assert.Nil(t, last)

		first, last = ipx.RangeFromNetwork(
			&net.IPNet{
				IP:   net.IPv4zero.To4(),    // address is IPv4
				Mask: net.CIDRMask(96, 128), // mask is IPv6
			},
		)
		assert.Nil(t, first)
		assert.Nil(t, last)
	})

	tt.Run("ipv4_29", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("192.168.0.10/29")
		require.NoError(t, err, "failed to parse CIDR")

		first, last := ipx.RangeFromNetwork(nwk)
		assert.Equal(t, "192.168.0.8", first.String())
		assert.Equal(t, "192.168.0.15", last.String())
	})

	tt.Run("ipv4_23", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("192.168.0.10/23")
		require.NoError(t, err, "failed to parse CIDR")

		first, last := ipx.RangeFromNetwork(nwk)
		assert.Equal(t, "192.168.0.0", first.String())
		assert.Equal(t, "192.168.1.255", last.String())
	})

	tt.Run("ipv4_29_mapped_to_ipv6_125", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("::ffff:192.168.0.10/125")
		require.NoError(t, err, "failed to parse CIDR")

		first, last := ipx.RangeFromNetwork(nwk)
		assert.Equal(t, "192.168.0.8", first.String())
		assert.Equal(t, "192.168.0.15", last.String())
	})

	tt.Run("ipv4_29_mapped_to_ipv6_125_hex", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("::ffff:c0a8:000a/125")
		require.NoError(t, err, "failed to parse CIDR")

		first, last := ipx.RangeFromNetwork(nwk)
		assert.Equal(t, "192.168.0.8", first.String())
		assert.Equal(t, "192.168.0.15", last.String())
	})

	tt.Run("ipv6_120", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("2001:db8::8a2e:370:7334/120")
		require.NoError(t, err, "failed to parse CIDR")

		first, last := ipx.RangeFromNetwork(nwk)
		assert.Equal(t, "2001:db8::8a2e:370:7300", first.String())
		assert.Equal(t, "2001:db8::8a2e:370:73ff", last.String())
	})

	tt.Run("ipv6_107", func(t *testing.T) {
		_, nwk, err := net.ParseCIDR("2001:db8::8a2e:370:7334/107")
		require.NoError(t, err, "failed to parse CIDR")

		first, last := ipx.RangeFromNetwork(nwk)
		assert.Equal(t, "2001:db8::8a2e:360:0", first.String())
		assert.Equal(t, "2001:db8::8a2e:37f:ffff", last.String())
	})
}
