// Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package discover

import (
	"net"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

// TestIPRateLimiterBurstThenThrottle verifies a single IP gets its burst, is
// then throttled, and recovers as tokens replenish.
func TestIPRateLimiterBurstThenThrottle(t *testing.T) {
	burst := 5
	l := newIPRateLimiter(rate.Limit(1), burst, 1024, time.Minute)
	ip := net.ParseIP("203.0.113.7")
	now := time.Unix(0, 0)

	for i := 0; i < burst; i++ {
		if !l.allow(ip, now) {
			t.Fatalf("burst ping %d should be allowed", i)
		}
	}
	if l.allow(ip, now) {
		t.Fatal("ping beyond burst should be throttled")
	}
	// rate = 1/s, so one token is back after a second.
	if !l.allow(ip, now.Add(time.Second)) {
		t.Fatal("ping after 1s should be allowed again")
	}
}

// TestIPRateLimiterPerIP verifies budgets are independent per source IP.
func TestIPRateLimiterPerIP(t *testing.T) {
	l := newIPRateLimiter(rate.Limit(1), 1, 1024, time.Minute)
	now := time.Unix(0, 0)
	a := net.ParseIP("203.0.113.1")
	b := net.ParseIP("203.0.113.2")

	if !l.allow(a, now) || !l.allow(b, now) {
		t.Fatal("first ping from each IP should be allowed")
	}
	if l.allow(a, now) {
		t.Fatal("second ping from A within the same instant should be throttled")
	}
	if !l.allow(b, now.Add(time.Second)) {
		t.Fatal("B's budget should be independent of A")
	}
}

// TestIPRateLimiterBounded verifies the map stays within its cap under
// source-IP rotation (eviction works).
func TestIPRateLimiterBounded(t *testing.T) {
	const max = 64
	l := newIPRateLimiter(rate.Limit(1), 1, max, time.Hour)
	now := time.Unix(0, 0)
	for i := 0; i < max*4; i++ {
		ip := net.IPv4(10, 0, byte(i>>8), byte(i))
		l.allow(ip, now)
		now = now.Add(time.Millisecond)
	}
	if got := l.len(); got > max {
		t.Fatalf("limiter map exceeded cap: got %d, want <= %d", got, max)
	}
}

// TestIPRateLimiterEvictsLeastRecentlySeen verifies that, at capacity, the key
// seen least recently is evicted and a recently touched key survives.
func TestIPRateLimiterEvictsLeastRecentlySeen(t *testing.T) {
	const max = 4
	l := newIPRateLimiter(rate.Limit(1), 1, max, time.Hour)
	now := time.Unix(0, 0)
	ips := make([]net.IP, max)
	for i := range ips {
		ips[i] = net.IPv4(10, 0, 0, byte(i+1))
		l.allow(ips[i], now)
		now = now.Add(time.Millisecond)
	}
	// Touch the oldest key so it becomes the most recent.
	l.allow(ips[0], now)
	now = now.Add(time.Millisecond)

	// Inserting a new key at capacity evicts ips[1], now the least recent.
	l.allow(net.IPv4(10, 0, 0, 100), now)
	if got := l.len(); got != max {
		t.Fatalf("tracked keys = %d, want %d", got, max)
	}
	if !l.has(ips[0]) {
		t.Fatal("recently touched key must not be evicted")
	}
	if l.has(ips[1]) {
		t.Fatal("least recently seen key should have been evicted")
	}
	for _, ip := range ips[2:] {
		if !l.has(ip) {
			t.Fatalf("key %s should still be tracked", ip)
		}
	}
}

// TestIPRateLimiterEvictIdle verifies that, at capacity, every idle key is
// dropped before any active key is considered for eviction.
func TestIPRateLimiterEvictIdle(t *testing.T) {
	const max = 4
	const ttl = time.Minute
	l := newIPRateLimiter(rate.Limit(1), 1, max, ttl)
	now := time.Unix(0, 0)
	idle := []net.IP{net.IPv4(10, 0, 0, 1), net.IPv4(10, 0, 0, 2)}
	active := []net.IP{net.IPv4(10, 0, 0, 3), net.IPv4(10, 0, 0, 4)}
	for _, ip := range idle {
		l.allow(ip, now)
	}
	now = now.Add(ttl + time.Second)
	for _, ip := range active {
		l.allow(ip, now)
	}

	l.allow(net.IPv4(10, 0, 0, 100), now)
	for _, ip := range idle {
		if l.has(ip) {
			t.Fatalf("idle key %s should have been evicted", ip)
		}
	}
	for _, ip := range active {
		if !l.has(ip) {
			t.Fatalf("active key %s must survive idle eviction", ip)
		}
	}
	if got := l.len(); got != len(active)+1 {
		t.Fatalf("tracked keys = %d, want %d", got, len(active)+1)
	}
}

// BenchmarkIPRateLimiterRotation measures the per-packet cost when every packet
// comes from a new source and the map is at capacity, which is the eviction
// path an attacker rotating source IPs drives.
func BenchmarkIPRateLimiterRotation(b *testing.B) {
	l := newIPRateLimiter(rate.Limit(1), 1, maxLimitedIPs, time.Hour)
	now := time.Unix(0, 0)
	ip := make(net.IP, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Vary the address so each packet is a new source IP.
		ip[0], ip[1], ip[2], ip[3] = byte(i>>24), byte(i>>16), byte(i>>8), byte(i)
		l.allow(ip, now)
		now = now.Add(time.Microsecond)
	}
}

const (
	testPingRate  = 3  // pings/sec for ping.preverify tests
	testPingBurst = 10 // token-bucket burst for ping.preverify tests
)

// newPingTestUDP builds a minimal udp instance for exercising ping.preverify.
// It carries only the fields the rate-limit gate reads: the local network id,
// this node's type (BN gating), and the per-IP ping limiter.
func newPingTestUDP(nodeType NodeType, networkID uint64) *udp {
	return &udp{
		networkID:   networkID,
		ourEndpoint: rpcEndpoint{NType: nodeType},
		pingLimiter: newPingLimiter(testPingRate, testPingBurst),
	}
}

// TestPingPreverifyRateLimit exercises the full ping rate-limit gate in
// ping.preverify end to end: a bootstrap node throttles pings from a non-LAN
// source IP once the burst is spent, exempts LAN/loopback sources, and does not
// rate-limit at all when the node is not a bootstrap node.
func TestPingPreverifyRateLimit(t *testing.T) {
	const networkID = 12345
	req := &ping{NetworkID: networkID, Expiration: futureExp}
	fromID := NodeID{}

	t.Run("BN throttles non-LAN flood after burst", func(t *testing.T) {
		u := newPingTestUDP(NodeTypeBN, networkID)
		from := &net.UDPAddr{IP: net.ParseIP("203.0.113.5"), Port: 30303}
		// The burst (token bucket capacity) is allowed in one flood.
		for i := 0; i < testPingBurst; i++ {
			if err := req.preverify(u, from, fromID); err != nil {
				t.Fatalf("ping %d within burst should pass, got %v", i, err)
			}
		}
		// The very next ping (no time for tokens to refill) is throttled.
		if err := req.preverify(u, from, fromID); err != errPingRateLimited {
			t.Fatalf("ping beyond burst should be rate limited, got %v", err)
		}
	})

	t.Run("LAN source is exempt", func(t *testing.T) {
		u := newPingTestUDP(NodeTypeBN, networkID)
		from := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 30303}
		for i := 0; i < testPingBurst*3; i++ {
			if err := req.preverify(u, from, fromID); err != nil {
				t.Fatalf("LAN ping %d must never be rate limited, got %v", i, err)
			}
		}
	})

	t.Run("non-BN node never rate-limits", func(t *testing.T) {
		u := newPingTestUDP(NodeTypeCN, networkID)
		from := &net.UDPAddr{IP: net.ParseIP("203.0.113.5"), Port: 30303}
		for i := 0; i < testPingBurst*3; i++ {
			if err := req.preverify(u, from, fromID); err != nil {
				t.Fatalf("non-BN ping %d must never be rate limited, got %v", i, err)
			}
		}
	})
}
