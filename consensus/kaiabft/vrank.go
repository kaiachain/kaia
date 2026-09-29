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

package kaiabft

import (
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/bft"
)

// prepreparedEvent is posted when this node accepts a PRE-PREPARE for a
// (block, view). The relay below forwards it to the VRank module.
type prepreparedEvent struct {
	Block *types.Block
	View  *bft.View
}

// postPrepreparedEvent notifies the VRank relay that this node accepted pp.
// The view is deep-copied so later round mutations don't alias the payload.
func (m *machine) postPrepreparedEvent(pp *bft.Preprepare) {
	block, ok := pp.Proposal.(*types.Block)
	if !ok {
		logger.Warn("Failed to post preprepared event due to unexpected proposal type")
		return
	}
	view := &bft.View{
		Round:    new(big.Int).Set(pp.View.Round),
		Sequence: new(big.Int).Set(pp.View.Sequence),
	}
	m.b.eventMux.Post(prepreparedEvent{Block: block, View: view})
}

// startPrepreparedRelay subscribes to prepreparedEvent and forwards each
// accepted (block, view) to the VRank module on dedicated goroutines, so the
// consensus loop is never blocked. No-op if no VRank module is registered.
func (b *backend) startPrepreparedRelay() {
	if b.vrankModule == nil || b.prepreparedSub != nil {
		return
	}
	const relayQueueSize = 32

	sub := b.eventMux.Subscribe(prepreparedEvent{})
	stopCh := make(chan struct{})
	queue := make(chan prepreparedEvent, relayQueueSize)

	b.prepreparedSub = sub
	b.prepreparedStopCh = stopCh

	b.prepreparedWg.Add(1)
	go func() {
		defer close(queue)
		defer b.prepreparedWg.Done()
		for {
			select {
			case ev, ok := <-sub.Chan():
				if !ok {
					return
				}
				preprepared, ok := ev.Data.(prepreparedEvent)
				if !ok || preprepared.Block == nil || preprepared.View == nil {
					continue
				}
				// Non-blocking enqueue keeps the TypeMux Post path free of backpressure.
				select {
				case queue <- preprepared:
				default:
					logger.Warn("Dropping preprepared event due to full relay queue",
						"blockNum", preprepared.Block.NumberU64(), "round", preprepared.View.Round.Uint64())
				}
			case <-stopCh:
				return
			}
		}
	}()

	b.prepreparedWg.Add(1)
	go func() {
		defer b.prepreparedWg.Done()
		for preprepared := range queue {
			b.vrankModule.HandleIstanbulPreprepare(preprepared.Block, preprepared.View)
		}
	}()
}

// stopPrepreparedRelay unsubscribes and waits for the relay goroutines to exit.
func (b *backend) stopPrepreparedRelay() {
	if b.prepreparedSub != nil {
		b.prepreparedSub.Unsubscribe()
		b.prepreparedSub = nil
	}
	if b.prepreparedStopCh != nil {
		close(b.prepreparedStopCh)
		b.prepreparedStopCh = nil
	}
	b.prepreparedWg.Wait()
}
