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
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	// pingRatePerIP is the sustained discovery-ping rate allowed per source IP.
	// A node pings a bootnode at up to ~1/s while its dialer drives repeated
	// discovery refreshes, so this leaves headroom for a few nodes sharing one
	// (non-LAN) public IP.
	pingRatePerIP = rate.Limit(3) // 3 pings/sec sustained
	// pingBurstPerIP is the burst of pings tolerated per source IP.
	pingBurstPerIP = 10
	// maxLimitedIPs bounds the number of tracked source IPs so the limiter map
	// cannot grow without limit under source-IP rotation.
	maxLimitedIPs = 16384
	// limiterIdleTTL is how long an idle per-IP limiter is retained before it
	// becomes eligible for eviction.
	limiterIdleTTL = 10 * time.Minute
)

// ipRateLimiter applies a per-source-IP token-bucket rate limit. It is used by
// bootstrap nodes to throttle discovery pings, bounding amplification/DoS
// against the UDP endpoint.
//
// Source IPs in UDP are spoofable, so this is a partial, application-level
// defense intended to be complemented by network-layer (L3/L4) protections.
type ipRateLimiter struct {
	mu      sync.Mutex
	limit   rate.Limit
	burst   int
	maxIPs  int
	idleTTL time.Duration
	ips     map[string]*ipLimiterEntry
}

type ipLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func newIPRateLimiter(limit rate.Limit, burst, maxIPs int, idleTTL time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		limit:   limit,
		burst:   burst,
		maxIPs:  maxIPs,
		idleTTL: idleTTL,
		ips:     make(map[string]*ipLimiterEntry),
	}
}

// allow reports whether an event from ip is permitted at time now. now is passed
// in so the limiter is deterministic and testable.
func (l *ipRateLimiter) allow(ip net.IP, now time.Time) bool {
	key := ip.String()

	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.ips[key]
	if !ok {
		// Keep the map bounded: drop idle entries first, then evict the
		// least-recently-seen one if still at capacity.
		if len(l.ips) >= l.maxIPs {
			l.evictIdle(now)
		}
		if len(l.ips) >= l.maxIPs {
			l.evictOldest()
		}
		e = &ipLimiterEntry{limiter: rate.NewLimiter(l.limit, l.burst)}
		l.ips[key] = e
	}
	e.lastSeen = now
	return e.limiter.AllowN(now, 1)
}

// evictIdle removes entries idle for longer than idleTTL. An idle entry's token
// bucket is necessarily full, so dropping it loses no rate-limit state.
func (l *ipRateLimiter) evictIdle(now time.Time) {
	for k, e := range l.ips {
		if now.Sub(e.lastSeen) > l.idleTTL {
			delete(l.ips, k)
		}
	}
}

// evictOldest removes the single least-recently-seen entry to make room for a
// new IP when the map is at capacity.
func (l *ipRateLimiter) evictOldest() {
	var (
		oldestKey  string
		oldestTime time.Time
		found      bool
	)
	for k, e := range l.ips {
		if !found || e.lastSeen.Before(oldestTime) {
			oldestKey, oldestTime, found = k, e.lastSeen, true
		}
	}
	if found {
		delete(l.ips, oldestKey)
	}
}

// len returns the number of tracked IPs.
func (l *ipRateLimiter) len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.ips)
}
