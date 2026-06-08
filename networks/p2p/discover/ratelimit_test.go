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
