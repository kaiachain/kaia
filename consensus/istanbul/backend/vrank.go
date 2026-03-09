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

// RegisterVRankModule sets the VRank module of the Istanbul backend.
func (sb *backend) RegisterVRankModule(module vrank.VRankModule) {
	sb.vrankModule = module
}

func (sb *backend) startPrepreparedRelay() {
	if sb.vrankModule == nil || sb.prepreparedSub != nil {
		return
	}
	sub := sb.istanbulEventMux.Subscribe(istanbul.PrepreparedEvent{})
	stopCh := make(chan struct{})

	sb.prepreparedSub = sub
	sb.prepreparedCh = stopCh
	sb.prepreparedWg.Add(1)
	go func() {
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
				sb.vrankModule.HandleIstanbulPreprepare(preprepared.Block, preprepared.View)
			case <-stopCh:
				return
			}
		}
	}()
}

func (sb *backend) stopPrepreparedRelay() {
	if sb.prepreparedCh != nil {
		close(sb.prepreparedCh)
		sb.prepreparedCh = nil
	}
	if sb.prepreparedSub != nil {
		sb.prepreparedSub.Unsubscribe()
		sb.prepreparedSub = nil
	}
	sb.prepreparedWg.Wait()
}
