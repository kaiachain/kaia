// Modifications Copyright 2026 The Kaia Authors
// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from p2p/util_test.go (2019/04/02).
// Modified and improved for the Kaia development.

package p2p

import (
	"net"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common/mclock"
	"github.com/kaiachain/kaia/networks/p2p/netutil"
)

func TestExpHeap(t *testing.T) {
	var h expHeap

	var (
		basetime = mclock.AbsTime(10)
		exptimeA = basetime + mclock.AbsTime(2*time.Second)
		exptimeB = basetime + mclock.AbsTime(3*time.Second)
		exptimeC = basetime + mclock.AbsTime(4*time.Second)
	)
	h.add("b", exptimeB)
	h.add("a", exptimeA)
	h.add("c", exptimeC)

	if h.nextExpiry() != exptimeA {
		t.Fatal("wrong nextExpiry")
	}
	if h.count("a") != 1 || h.count("b") != 1 || h.count("c") != 1 {
		t.Fatal("heap doesn't contain all live items")
	}

	h.expire(exptimeA+1, nil)
	if h.nextExpiry() != exptimeB {
		t.Fatal("wrong nextExpiry")
	}
	if h.count("a") != 0 {
		t.Fatal("heap contains a even though it has already expired")
	}
	if h.count("b") != 1 || h.count("c") != 1 {
		t.Fatal("heap doesn't contain all live items")
	}
}

// TestCheckInboundConnThrottle verifies that repeated attempts from the same
// non-LAN IP within inboundThrottleTime are rejected, and allowed again once
// the throttle window has elapsed.
func TestCheckInboundConnThrottle(t *testing.T) {
	srv := &BaseServer{}
	ip := net.ParseIP("203.0.113.7") // TEST-NET-3, treated as a public (non-LAN) IP
	now := mclock.AbsTime(0)

	if err := srv.checkInboundConn(ip, 1, now); err != nil {
		t.Fatalf("first attempt should be accepted, got %v", err)
	}
	if err := srv.checkInboundConn(ip, 1, now+mclock.AbsTime(inboundThrottleTime/2)); err != errTooManyInboundAttempts {
		t.Fatalf("attempt within throttle window should be rejected, got %v", err)
	}
	if err := srv.checkInboundConn(ip, 1, now+mclock.AbsTime(inboundThrottleTime)+1); err != nil {
		t.Fatalf("attempt after throttle window should be accepted, got %v", err)
	}
}

// TestCheckInboundConnLANExempt verifies that LAN/loopback IPs are never
// throttled, so co-located/internal nodes can reconnect freely.
func TestCheckInboundConnLANExempt(t *testing.T) {
	srv := &BaseServer{}
	ip := net.ParseIP("127.0.0.1")
	for i := 0; i < 5; i++ {
		if err := srv.checkInboundConn(ip, 1, mclock.AbsTime(0)); err != nil {
			t.Fatalf("LAN IP should never be throttled, got %v on attempt %d", err, i)
		}
	}
}

// TestCheckInboundConnNetRestrict verifies that NetRestrict is enforced before
// throttling.
func TestCheckInboundConnNetRestrict(t *testing.T) {
	list, err := netutil.ParseNetlist("192.0.2.0/24")
	if err != nil {
		t.Fatal(err)
	}
	srv := &BaseServer{}
	srv.NetRestrict = list

	if err := srv.checkInboundConn(net.ParseIP("203.0.113.7"), 1, 0); err != errNotWhitelisted {
		t.Fatalf("IP outside NetRestrict should be rejected, got %v", err)
	}
	if err := srv.checkInboundConn(net.ParseIP("192.0.2.5"), 1, 0); err != nil {
		t.Fatalf("IP within NetRestrict should be accepted, got %v", err)
	}
}
