package ipx_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ns1/ipx/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExampleCompareIP_v4 an example of CompareIP for IPv4
func ExampleCompareIP_v4() {
	a := net.ParseIP("192.168.0.10")
	b := net.ParseIP("192.168.0.20")
	fmt.Println(ipx.MustCompareIP(a, b))
	fmt.Println(ipx.MustCompareIP(a, a))
	fmt.Println(ipx.MustCompareIP(b, a))
	// Output:
	// -1
	// 0
	// 1
}

// ExampleCompareIP_v6 an example of CompareIP for IPv6
func ExampleCompareIP_v6() {
	a := net.ParseIP("2001:db8::10")
	b := net.ParseIP("2001:db8::20")
	fmt.Println(ipx.MustCompareIP(a, b))
	fmt.Println(ipx.MustCompareIP(a, a))
	fmt.Println(ipx.MustCompareIP(b, a))
	// Output:
	// -1
	// 0
	// 1
}

// TestCompareIP unit tests for CompareIP
func TestCompareIP(t *testing.T) {
	// bad input
	_, err := ipx.CompareIP(make(net.IP, 3), net.ParseIP("bad"))
	require.Error(t, err, "should not compare IP addresses")
	assert.ErrorIs(t, err, ipx.ErrInvalidIP)
	assert.Panics(t, func() {
		ipx.MustCompareIP(make(net.IP, 3), net.ParseIP("bad"))
	})

	// IPv4
	assert.EqualValues(t, -1, // less than /24
		ipx.MustCompareIP(
			net.ParseIP("192.168.0.10"),
			net.ParseIP("192.168.0.20"),
		))
	assert.EqualValues(t, +1, // greater than /24
		ipx.MustCompareIP(
			net.ParseIP("192.168.0.20"),
			net.ParseIP("192.168.0.10"),
		))
	assert.EqualValues(t, 0, // equal /24
		ipx.MustCompareIP(
			net.ParseIP("192.168.0.10"),
			net.ParseIP("192.168.0.10"),
		))
	assert.EqualValues(t, -1, // less than /16
		ipx.MustCompareIP(
			net.ParseIP("192.168.10.20"),
			net.ParseIP("192.168.20.10"),
		))
	assert.EqualValues(t, +1, // greater than /16
		ipx.MustCompareIP(
			net.ParseIP("192.168.20.10"),
			net.ParseIP("192.168.10.20"),
		))
	assert.EqualValues(t, 0, // equal /16
		ipx.MustCompareIP(
			net.ParseIP("192.168.10.0"),
			net.ParseIP("192.168.10.0"),
		))

	// IPv4 in IPv6
	assert.EqualValues(t, -1, // less than /24
		ipx.MustCompareIP(
			net.ParseIP("::192.168.0.10"),
			net.ParseIP("::192.168.0.20"),
		))
	assert.EqualValues(t, +1, // greater than /24
		ipx.MustCompareIP(
			net.ParseIP("::192.168.0.20"),
			net.ParseIP("::192.168.0.10"),
		))
	assert.EqualValues(t, 0, // equal /24
		ipx.MustCompareIP(
			net.ParseIP("::192.168.0.10"),
			net.ParseIP("::192.168.0.10"),
		))
	assert.EqualValues(t, -1, // less than /16
		ipx.MustCompareIP(
			net.ParseIP("::192.168.10.20"),
			net.ParseIP("::192.168.20.10"),
		))
	assert.EqualValues(t, +1, // greater than /16
		ipx.MustCompareIP(
			net.ParseIP("::192.168.20.10"),
			net.ParseIP("::192.168.10.20"),
		))
	assert.EqualValues(t, 0, // equal /16
		ipx.MustCompareIP(
			net.ParseIP("::192.168.10.0"),
			net.ParseIP("::192.168.10.0"),
		))

	// IPv6
	assert.EqualValues(t, -1, // less than /16
		ipx.MustCompareIP(
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:6000"),
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7000"),
		))
	assert.EqualValues(t, +1, // greater than /16
		ipx.MustCompareIP(
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7000"),
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:6000"),
		))
	assert.EqualValues(t, 0, // equal /16
		ipx.MustCompareIP(
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:6000"),
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:6000"),
		))
	assert.EqualValues(t, -1, // less than /32
		ipx.MustCompareIP(
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:6000:7000"),
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:7000:6000"),
		))
	assert.EqualValues(t, +1, // greater than /32
		ipx.MustCompareIP(
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:7000:6000"),
			net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:6000:7000"),
		))
	assert.EqualValues(t, -1, // less than /128
		ipx.MustCompareIP(
			net.ParseIP("2001:eeee:eeee:eeee:eeee:eeee:eeee:eeee"),
			net.ParseIP("3001:eeee:eeee:eeee:eeee:eeee:eeee:eeee"),
		))
	assert.EqualValues(t, +1, // greater than /128
		ipx.MustCompareIP(
			net.ParseIP("3001:eeee:eeee:eeee:eeee:eeee:eeee:eeee"),
			net.ParseIP("2001:eeee:eeee:eeee:eeee:eeee:eeee:eeee"),
		))
}

// ExampleNextIP_v4 is an example of NextIP for IPv4.
func ExampleNextIP_v4() {
	addr1 := net.ParseIP("0.0.0.1")
	addr2 := ipx.NextIP(addr1, 257)
	fmt.Println(addr1, "+ 257 =", addr2)
	// Output:
	// 0.0.0.1 + 257 = 0.0.1.2
}

// ExampleNextIP_v6 is an example of NextIP for IPv6.
func ExampleNextIP_v6() {
	addr1 := net.ParseIP("::1")
	addr2 := ipx.NextIP(addr1, 65537)
	fmt.Println(addr1, "+ 65537 =", addr2)
	// Output:
	// ::1 + 65537 = ::1:2
}

// TestNextIP unit tests for NextIP
func TestNextIP(t *testing.T) {
	// bad input
	assert.Nil(t, ipx.NextIP(nil, 1))
	assert.Nil(t, ipx.NextIP(make(net.IP, 3), 1))

	// step=0 means the same output
	assert.Equal(t, net.IPv4bcast, ipx.NextIP(net.IPv4bcast, 0))
	assert.Equal(t, net.IPv6zero, ipx.NextIP(net.IPv6zero, 0))

	// IPv4
	assert.Equal(t, "0.0.0.1", ipx.NextIP(net.ParseIP("0.0.0.0"), +1).String())
	assert.Equal(t, "0.0.0.0", ipx.NextIP(net.ParseIP("0.0.0.1"), -1).String())
	assert.Equal(t, "0.0.0.0", ipx.NextIP(net.ParseIP("255.255.255.255"), +1).String())
	assert.Equal(t, "255.255.255.255", ipx.NextIP(net.ParseIP("0.0.0.0"), -1).String())
	assert.Equal(t, "0.0.1.1", ipx.NextIP(net.ParseIP("0.0.0.0"), +257).String())
	assert.Equal(t, "0.0.0.0", ipx.NextIP(net.ParseIP("0.0.1.1"), -257).String())

	// IPv6
	assert.Equal(t, "::1", ipx.NextIP(net.ParseIP("::"), +1).String())
	assert.Equal(t, "::", ipx.NextIP(net.ParseIP("::1"), -1).String())
	assert.Equal(t, "::", ipx.NextIP(net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"), +1).String())
	assert.Equal(t, "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", ipx.NextIP(net.ParseIP("::"), -1).String())
	assert.Equal(t, "::1:0:0", ipx.NextIP(net.ParseIP("::"), +(1<<32)).String())
	assert.Equal(t, "::", ipx.NextIP(net.ParseIP("::1:0:0"), -(1<<32)).String())
}

// BenchmarkNextIP performance benchmarks for NextIP
func BenchmarkNextIP(bb *testing.B) {
	// helper function to run benchmark
	bench := func(address string, step int) func(*testing.B) {
		return func(b *testing.B) {
			ip := net.ParseIP(address)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				out := ipx.NextIP(ip, step)
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
