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
	"container/list"
	"net"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	// defaultPingBurst is the token-bucket burst used when ping rate limiting is
	// enabled without an explicit burst. LAN/loopback sources are exempt, so the
	// limiter only sees public IPs — where several nodes can sit behind one IP
	// via NAT and share its bucket. Each node pings a bootnode at ~1/s (its
	// dialer drives repeated discovery refreshes), so a burst of 10 absorbs
	// short spikes from a handful of such nodes without throttling them.
	defaultPingBurst = 10
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
//
// Entries are kept in a recency list so that eviction under source-IP rotation
// costs O(1) per new IP rather than a scan of the whole map while holding mu.
type ipRateLimiter struct {
	mu      sync.Mutex
	limit   rate.Limit
	burst   int
	maxIPs  int
	idleTTL time.Duration
	// entries maps a source IP (as a string) to its element in lru. The
	// element value is the *ipLimiterEntry.
	entries map[string]*list.Element
	// lru orders entries by lastSeen: front is the most recently seen IP,
	// back is the least recently seen one.
	lru *list.List
}

type ipLimiterEntry struct {
	key      string
	limiter  *rate.Limiter
	lastSeen time.Time
}

func newIPRateLimiter(limit rate.Limit, burst, maxIPs int, idleTTL time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		limit:   limit,
		burst:   burst,
		maxIPs:  maxIPs,
		idleTTL: idleTTL,
		entries: make(map[string]*list.Element),
		lru:     list.New(),
	}
}

// newPingLimiter builds the per-IP discovery-ping limiter for a bootstrap node
// from operator-provided values, applying defaultPingBurst when burst <= 0.
func newPingLimiter(ratePerSec float64, burst int) *ipRateLimiter {
	if burst <= 0 {
		burst = defaultPingBurst
	}
	return newIPRateLimiter(rate.Limit(ratePerSec), burst, maxLimitedIPs, limiterIdleTTL)
}

// allow reports whether an event from ip is permitted at time now. now is passed
// in so the limiter is deterministic and testable.
func (l *ipRateLimiter) allow(ip net.IP, now time.Time) bool {
	key := ip.String()

	l.mu.Lock()
	defer l.mu.Unlock()

	if el, ok := l.entries[key]; ok {
		e := el.Value.(*ipLimiterEntry)
		e.lastSeen = now
		l.lru.MoveToFront(el)
		return e.limiter.AllowN(now, 1)
	}

	// Keep the map bounded: drop idle entries first, then evict the
	// least-recently-seen one if still at capacity.
	if l.lru.Len() >= l.maxIPs {
		l.evictIdle(now)
	}
	if l.lru.Len() >= l.maxIPs {
		l.evictOldest()
	}
	e := &ipLimiterEntry{key: key, limiter: rate.NewLimiter(l.limit, l.burst), lastSeen: now}
	l.entries[key] = l.lru.PushFront(e)
	return e.limiter.AllowN(now, 1)
}

// evictIdle removes entries idle for longer than idleTTL. An idle entry's token
// bucket is necessarily full, so dropping it loses no rate-limit state.
// Entries are visited from the least recently seen end and the walk stops at
// the first active one, so the cost is proportional to the number evicted.
// The caller must hold l.mu.
func (l *ipRateLimiter) evictIdle(now time.Time) {
	for el := l.lru.Back(); el != nil; el = l.lru.Back() {
		e := el.Value.(*ipLimiterEntry)
		if now.Sub(e.lastSeen) <= l.idleTTL {
			return
		}
		l.remove(el)
	}
}

// evictOldest removes the single least-recently-seen entry to make room for a
// new IP when the map is at capacity. The caller must hold l.mu.
func (l *ipRateLimiter) evictOldest() {
	if el := l.lru.Back(); el != nil {
		l.remove(el)
	}
}

// remove deletes el from both the map and the recency list. The caller must
// hold l.mu.
func (l *ipRateLimiter) remove(el *list.Element) {
	e := el.Value.(*ipLimiterEntry)
	delete(l.entries, e.key)
	l.lru.Remove(el)
}

// len returns the number of tracked IPs.
func (l *ipRateLimiter) len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}

// has reports whether ip is currently tracked.
func (l *ipRateLimiter) has(ip net.IP) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.entries[ip.String()]
	return ok
}
