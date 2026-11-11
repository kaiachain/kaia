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
// This file is derived from p2p/util.go (2019/04/02).
// Modified and improved for the Kaia development.

package p2p

import (
	"container/heap"
	"net"

	"github.com/kaiachain/kaia/common/mclock"
)

// expHeap tracks strings and their expiry time.
type expHeap []expItem

// expItem is an entry in addrHistory.
type expItem struct {
	item string
	exp  mclock.AbsTime
}

// nextExpiry returns the next expiry time.
func (h *expHeap) nextExpiry() mclock.AbsTime {
	return (*h)[0].exp
}

// add adds an item and sets its expiry time.
func (h *expHeap) add(item string, exp mclock.AbsTime) {
	heap.Push(h, expItem{item, exp})
}

// count returns how many unexpired entries an item holds.
func (h expHeap) count(item string) int {
	n := 0
	for _, v := range h {
		if v.item == item {
			n++
		}
	}
	return n
}

// expire removes items with expiry time before 'now'.
func (h *expHeap) expire(now mclock.AbsTime, onExp func(string)) {
	for h.Len() > 0 && h.nextExpiry() < now {
		item := heap.Pop(h)
		if onExp != nil {
			onExp(item.(expItem).item)
		}
	}
}

// heap.Interface boilerplate
func (h expHeap) Len() int            { return len(h) }
func (h expHeap) Less(i, j int) bool  { return h[i].exp < h[j].exp }
func (h expHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *expHeap) Push(x interface{}) { *h = append(*h, x.(expItem)) }
func (h *expHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = expItem{}
	*h = old[0 : n-1]
	return x
}

// tcpAddrIP returns the IP of a *net.TCPAddr, or nil for any other address
// type (e.g. in-process test pipes that have no usable remote IP).
func tcpAddrIP(addr net.Addr) net.IP {
	if tcp, ok := addr.(*net.TCPAddr); ok {
		return tcp.IP
	}
	return nil
}
