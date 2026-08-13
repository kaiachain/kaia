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

package backend

import (
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// RegisterVRankModule sets the VRank module of the Istanbul backend so that
// accepted PRE-PREPAREs are relayed to it for consensus-participation scoring.
func (sb *backend) RegisterVRankModule(module vrank.VRankModule) {
	sb.vrankModule = module
}

// startPrepreparedRelay subscribes to PrepreparedEvent and forwards each accepted
// (block, view) to the VRank module on dedicated goroutines, so the consensus hot
// path is never blocked. No-op if no VRank module is registered.
func (sb *backend) startPrepreparedRelay() {
	if sb.vrankModule == nil || sb.prepreparedSub != nil {
		return
	}
	const relayQueueSize = 32

	sub := sb.istanbulEventMux.Subscribe(istanbul.PrepreparedEvent{})
	stopCh := make(chan struct{})
	queue := make(chan istanbul.PrepreparedEvent, relayQueueSize)

	sb.prepreparedSub = sub
	sb.prepreparedStopCh = stopCh

	sb.prepreparedWg.Add(1)
	go func() {
		defer close(queue)
		defer sb.prepreparedWg.Done()
		for {
			select {
			case ev, ok := <-sub.Chan():
				if !ok {
					return
				}
				preprepared, ok := ev.Data.(istanbul.PrepreparedEvent)
				if !ok || preprepared.Block == nil || preprepared.View == nil {
					continue
				}
				// Non-blocking enqueue keeps the TypeMux Post path free of backpressure.
				select {
				case queue <- preprepared:
				default:
					logger.Warn("Dropping PrepreparedEvent due to full relay queue",
						"blockNum", preprepared.Block.NumberU64(), "round", preprepared.View.Round.Uint64())
				}
			case <-stopCh:
				return
			}
		}
	}()

	sb.prepreparedWg.Add(1)
	go func() {
		defer sb.prepreparedWg.Done()
		for preprepared := range queue {
			sb.vrankModule.HandleIstanbulPreprepare(preprepared.Block, preprepared.View)
		}
	}()
}

// stopPrepreparedRelay unsubscribes and waits for the relay goroutines to exit.
func (sb *backend) stopPrepreparedRelay() {
	if sb.prepreparedSub != nil {
		sb.prepreparedSub.Unsubscribe()
		sb.prepreparedSub = nil
	}
	if sb.prepreparedStopCh != nil {
		close(sb.prepreparedStopCh)
		sb.prepreparedStopCh = nil
	}
	sb.prepreparedWg.Wait()
}
