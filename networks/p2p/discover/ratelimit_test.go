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

// newPingTestUDP builds a minimal udp instance for exercising ping.preverify.
// It carries only the fields the rate-limit gate reads: the local network id,
// this node's type (BN gating), and the per-IP ping limiter.
func newPingTestUDP(nodeType NodeType, networkID uint64) *udp {
	return &udp{
		networkID:   networkID,
		ourEndpoint: rpcEndpoint{NType: nodeType},
		pingLimiter: newIPRateLimiter(pingRatePerIP, pingBurstPerIP, maxLimitedIPs, limiterIdleTTL),
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
		for i := 0; i < pingBurstPerIP; i++ {
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
		for i := 0; i < pingBurstPerIP*3; i++ {
			if err := req.preverify(u, from, fromID); err != nil {
				t.Fatalf("LAN ping %d must never be rate limited, got %v", i, err)
			}
		}
	})

	t.Run("non-BN node never rate-limits", func(t *testing.T) {
		u := newPingTestUDP(NodeTypeCN, networkID)
		from := &net.UDPAddr{IP: net.ParseIP("203.0.113.5"), Port: 30303}
		for i := 0; i < pingBurstPerIP*3; i++ {
			if err := req.preverify(u, from, fromID); err != nil {
				t.Fatalf("non-BN ping %d must never be rate limited, got %v", i, err)
			}
		}
	})
}
