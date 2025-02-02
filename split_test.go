package ipx_test

import (
	"fmt"
	"github.com/ns1/ipx/v2"
	"net"
	"testing"
)

func TestSplit(t *testing.T) {
	for _, c := range []struct {
		name, net string
		newPrefix int
		expected  []string
	}{
		{
			"no-op",
			"10.0.0.0/24",
			24,
			[]string{"10.0.0.0/24"},
		},
		{
			"invalid prefix",
			"10.0.0.0/24",
			23,
			[]string{},
		},
		{
			"ipv4",
			"10.0.0.0/24",
			26,
			[]string{"10.0.0.0/26", "10.0.0.64/26", "10.0.0.128/26", "10.0.0.192/26"},
		},
		{
			"ipv6",
			"::/24",
			26,
			[]string{"::/26", "0:40::/26", "0:80::/26", "0:c0::/26"},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			var nets []string
			splitter := ipx.Split(cidr(c.net), c.newPrefix)
			for splitter.Next() {
				nets = append(nets, splitter.Net().String())
			}

			if len(nets) != len(c.expected) {
				t.Fatalf("expected %v nets but got %v nets: %v", len(c.expected), len(nets), nets)
			}

			for i := range nets {
				if c.expected[i] != nets[i] {
					t.Errorf("expected %v at position %v but got %v", c.expected[i], i, nets[i])
				}
			}
		})
	}
}

func ExampleSplit() {
	c := cidr("10.0.0.0/24")
	split := ipx.Split(c, 26)
	for split.Next() {
		fmt.Println(split.Net())
	}
	// Output:
	// 10.0.0.0/26
	// 10.0.0.64/26
	// 10.0.0.128/26
	// 10.0.0.192/26
}

func ExampleSplit_IP6() {
	c := cidr("::/24")
	split := ipx.Split(c, 26)
	for split.Next() {
		fmt.Println(split.Net())
	}
	// Output:
	// ::/26
	// 0:40::/26
	// 0:80::/26
	// 0:c0::/26
}

func BenchmarkSplit(b *testing.B) {
	type bench struct {
		cidr      string
		newPrefix int
	}
	for _, g := range []struct {
		name    string
		benches []bench
	}{
		{
			"ipv4",
			[]bench{
				{"192.0.2.0/24", 30},
				{"192.0.2.0/24", 28},
			},
		},
		{
			"ipv6",
			[]bench{
				{"::/24", 30},
				{"::/24", 28},
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.benches {
				ipN := cidr(c.cidr)
				ones, _ := ipN.Mask.Size()
				b.Run(fmt.Sprintf("%v-%v", ones, c.newPrefix), func(b *testing.B) {
					b.ReportAllocs()

					for i := 0; i < b.N; i++ {
						for split := ipx.Split(ipN, c.newPrefix); split.Next(); {
						}
					}
				})
			}
		})
	}
}

func ExampleAddresses() {
	c := cidr("10.0.0.0/30")
	addrs := ipx.Addresses(c)
	for addrs.Next() {
		fmt.Println(addrs.IP())
	}
	// Output:
	// 10.0.0.0
	// 10.0.0.1
	// 10.0.0.2
	// 10.0.0.3
}

func TestAddresses(t *testing.T) {
	for _, c := range []struct {
		name, net string
		expected  []string
	}{
		{"ipv4 30", "10.0.0.0/30", []string{"10.0.0.0", "10.0.0.1", "10.0.0.2", "10.0.0.3"}},
		{"ipv4 32", "10.0.0.3/32", []string{"10.0.0.3"}},
		{
			"ipv6 126",
			"1bc1:6d67:4ec8::/126",
			[]string{
				"1bc1:6d67:4ec8::",
				"1bc1:6d67:4ec8::1",
				"1bc1:6d67:4ec8::2",
				"1bc1:6d67:4ec8::3",
			},
		},
		{"ipv6 128", "1bc1:6d67:4ec8::3/128", []string{"1bc1:6d67:4ec8::3"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, ipN, _ := net.ParseCIDR(c.net)

			var ips []string
			iter := ipx.Addresses(ipN)
			for iter.Next() {
				ips = append(ips, iter.IP().String())
			}

			if len(c.expected) != len(ips) {
				t.Fatalf("expected %v addresses but got %v: %v", len(c.expected), len(ips), ips)
			}
			for i := range ips {
				if ips[i] != c.expected[i] {
					t.Errorf("expected %v at position %d but got %v", c.expected[i], i, ips[i])
				}
			}
		})
	}
}

func BenchmarkAddresses(b *testing.B) {
	for _, g := range []struct {
		name  string
		cidrs []string
	}{
		{
			"ipv4",
			[]string{
				"10.0.0.0/30", // 4
				"10.0.0.0/24", // 256
			},
		},
		{
			"ipv6",
			[]string{
				"::/126", // 4
				"::/120", // 256
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.cidrs {
				ipN := cidr(c)
				ones, _ := ipN.Mask.Size()
				b.Run(fmt.Sprint(ones), func(b *testing.B) {
					b.ReportAllocs()

					for i := 0; i < b.N; i++ {
						for hosts := ipx.Addresses(ipN); hosts.Next(); {
						}
					}
				})
			}
		})
	}
}

func ExampleHosts() {
	c := cidr("10.0.0.0/29")
	hosts := ipx.Hosts(c)
	for hosts.Next() {
		fmt.Println(hosts.IP())
	}
	// Output:
	// 10.0.0.1
	// 10.0.0.2
	// 10.0.0.3
	// 10.0.0.4
	// 10.0.0.5
	// 10.0.0.6
}

func ExampleHosts_IP6() {
	c := cidr("::/125")
	hosts := ipx.Hosts(c)
	for hosts.Next() {
		fmt.Println(hosts.IP())
	}
	// Output:
	// ::1
	// ::2
	// ::3
	// ::4
	// ::5
	// ::6
}

func BenchmarkHosts(b *testing.B) {
	for _, g := range []struct {
		name  string
		cidrs []string
	}{
		{
			"ipv4",
			[]string{
				"10.0.0.0/30", // 4-2
				"10.0.0.0/24", // 256-2
			},
		},
		{
			"ipv6",
			[]string{
				"::/126", // 4-2
				"::/120", // 256-2
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.cidrs {
				ipN := cidr(c)
				ones, _ := ipN.Mask.Size()
				b.Run(fmt.Sprint(ones), func(b *testing.B) {
					b.ReportAllocs()

					for i := 0; i < b.N; i++ {
						for hosts := ipx.Hosts(ipN); hosts.Next(); {
						}
					}
				})
			}
		})
	}
}
