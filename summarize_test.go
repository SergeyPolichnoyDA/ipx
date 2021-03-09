package ipx_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ns1/ipx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExampleSummarizeRange_v4 is an example of SummarizeRange for IPv4
func ExampleSummarizeRange_v4() {
	networks, _ := ipx.SummarizeRange(
		net.ParseIP("192.0.2.0"),
		net.ParseIP("192.0.2.130"),
	)
	fmt.Println(networks)
	// Output:
	// [192.0.2.0/25 192.0.2.128/31 192.0.2.130/32]
}

// ExampleSummarizeRange_v6 is an example of SummarizeRange for IPv6
func ExampleSummarizeRange_v6() {
	networks, _ := ipx.SummarizeRange(
		net.ParseIP("::"),
		net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
	)
	fmt.Println(networks)
	// Output:
	// [::/0]
}

// BenchmarkSummarizeRange performance benchmarks for SummarizeRange
func BenchmarkSummarizeRange(bb *testing.B) {
	// helper function to run benchmark
	bench := func(first, last string, expectedLen int) func(*testing.B) {
		return func(b *testing.B) {
			ipFirst := net.ParseIP(first)
			ipLast := net.ParseIP(last)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				nwks, err := ipx.SummarizeRange(ipFirst, ipLast)
				require.NoError(b, err, "failed to summarize range")
				// note, we check only the output length.
				// all corner cases should be covered by
				// corresponding unit tests.
				require.Len(b, nwks, expectedLen)
			}
		}
	}

	// IPv4
	bb.Run("ipv4_all", bench("0.0.0.0", "255.255.255.255", 1))
	bb.Run("ipv4_24", bench("10.10.10.0", "10.10.10.255", 1))
	bb.Run("ipv4_32", bench("0.0.0.1", "255.255.255.255", 32))

	// IPv6
	bb.Run("ipv6_all", bench("::", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", 1))
	bb.Run("ipv6_128", bench("::1", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", 128))
}

// TestSummarizeRange unit tests for SummarizeRange
func TestSummarizeRange(tt *testing.T) {
	tt.Run("mismatched_versions", func(t *testing.T) {
		_, err := ipx.SummarizeRange(
			net.ParseIP("0.0.0.0"),
			net.ParseIP("::1"),
		)
		require.Error(t, err, "should not summarize range")
		assert.ErrorIs(t, err, ipx.ErrVersionMismatch)
		assert.Contains(t, err.Error(), "IP version mismatch")
	})

	tt.Run("bad_first", func(t *testing.T) {
		_, err := ipx.SummarizeRange(
			net.ParseIP("bad"), // nil
			net.ParseIP("192.168.2.200"),
		)
		require.Error(t, err, "should not summarize range")
		assert.ErrorIs(t, err, ipx.ErrInvalidIP)
		assert.Contains(t, err.Error(), "invalid IP address: first")
	})

	tt.Run("bad_last", func(t *testing.T) {
		_, err := ipx.SummarizeRange(
			net.ParseIP("192.168.2.100"),
			net.ParseIP("bad"), // nil
		)
		require.Error(t, err, "should not summarize range")
		assert.ErrorIs(t, err, ipx.ErrInvalidIP)
		assert.Contains(t, err.Error(), "invalid IP address: last")
	})

	tt.Run("bad_both", func(t *testing.T) {
		_, err := ipx.SummarizeRange(
			net.ParseIP("bad"), // nil
			net.ParseIP("bad"), // nil
		)
		require.Error(t, err, "should not summarize range")
		assert.ErrorIs(t, err, ipx.ErrInvalidIP)
		assert.Contains(t, err.Error(), "invalid IP address: first") // `first` is checked first
	})

	tt.Run("ipv4_no_overlap", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("192.168.2.200"),
			net.ParseIP("192.168.2.100"),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Empty(t, nwks.Strings())
	})

	tt.Run("ipv4_simple", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("192.0.2.0"),
			net.ParseIP("192.0.2.130"),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Equal(t, []string{
			"192.0.2.0/25",
			"192.0.2.128/31",
			"192.0.2.130/32",
		}, nwks.Strings())
	})

	tt.Run("ipv4_32", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("192.0.2.100").To4(),
			net.ParseIP("192.0.2.100").To4(),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Equal(t, []string{
			"192.0.2.100/32",
		}, nwks.Strings())
	})

	tt.Run("ipv4_16", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("192.168.0.0"),
			net.ParseIP("192.168.255.255"),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Equal(t, []string{
			"192.168.0.0/16",
		}, nwks.Strings())
	})

	tt.Run("ipv4_all", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("0.0.0.0"),
			net.ParseIP("255.255.255.255"),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Equal(t, []string{
			"0.0.0.0/0",
		}, nwks.Strings())
	})

	tt.Run("ipv4_odd_start", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("192.0.2.101"),
			net.ParseIP("192.0.2.130"),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Equal(t, []string{
			"192.0.2.101/32",
			"192.0.2.102/31",
			"192.0.2.104/29",
			"192.0.2.112/28",
			"192.0.2.128/31",
			"192.0.2.130/32",
		}, nwks.Strings())
	})

	tt.Run("ipv6_no_overlap", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("::200"),
			net.ParseIP("::100"),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Empty(t, nwks.Strings())
	})

	tt.Run("ipv6_128", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("::100").To16(),
			net.ParseIP("::100").To16(),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Equal(t, []string{
			"::100/128",
		}, nwks.Strings())
	})

	tt.Run("ipv6_16", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("1::"),
			net.ParseIP("1:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Equal(t, []string{
			"1::/16",
		}, nwks.Strings())
	})

	tt.Run("ipv6_all", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("::"),
			net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Equal(t, []string{
			"::/0",
		}, nwks.Strings())
	})

	tt.Run("ipv6_odd_start", func(t *testing.T) {
		nwks, err := ipx.SummarizeRange(
			net.ParseIP("1::1"),
			net.ParseIP("1::30"),
		)
		require.NoError(t, err, "failed to summarize range")
		assert.Equal(t, []string{
			"1::1/128",
			"1::2/127",
			"1::4/126",
			"1::8/125",
			"1::10/124",
			"1::20/124",
			"1::30/128",
		}, nwks.Strings())
	})
}
