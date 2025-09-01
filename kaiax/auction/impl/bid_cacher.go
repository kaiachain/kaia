// Copyright 2025 The Kaia Authors
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

package impl

import (
	"math"
	"math/big"
	"runtime"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/auction"
)

// bidSigCacher is a concurrent bid recoverer and cacher.
var bidSigCacher = newBidCacher(calcNumBidCachers())

func calcNumBidCachers() int {
	numWorkers := math.Ceil(float64(runtime.NumCPU()) * 2.0 / 3.0)
	return int(numWorkers)
}

// bidCacherRequest is a request for recovering bid with a
// specific signature scheme and caching it into the bid itself.
type bidCacherRequest struct {
	bid *auction.Bid

	chainID           *big.Int
	verifyingContract common.Address
	auctioneer        common.Address
}

// bidCacher is a helper structure to concurrently ecrecover transaction
// bids from digital signatures on background threads.
type bidCacher struct {
	threads int
	taskCh  chan *bidCacherRequest
}

// newBidCacher creates a new bid background cacher and starts
// as many processing goroutines as allowed by the GOMAXPROCS on construction.
func newBidCacher(threads int) *bidCacher {
	cacher := &bidCacher{
		taskCh:  make(chan *bidCacherRequest, threads),
		threads: threads,
	}
	for range threads {
		go cacher.cache()
	}
	return cacher
}

// cache is an infinite loop, caching bids from various forms of
// data structures.
func (cacher *bidCacher) cache() {
	for task := range cacher.taskCh {
		task.bid.ValidateSig(task.chainID, task.verifyingContract, task.auctioneer)
	}
}

// recover verifies the signatures and caches the bid.
func (cacher *bidCacher) recover(bid *auction.Bid, chainID *big.Int, verifyingContract common.Address, auctioneer common.Address) {
	// Early return if bid is nil.
	if bid == nil {
		return
	}

	cacher.taskCh <- &bidCacherRequest{
		bid:               bid,
		chainID:           chainID,
		verifyingContract: verifyingContract,
		auctioneer:        auctioneer,
	}
}
