// Copyright 2024 The Kaia Authors
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

package supply

import (
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test the concurrency of catchup() and PostInsertBlock()

// Stop() must not block even if catchup() is not running.
func (s *SupplyTestSuite) TestStop() {
	doneCh := make(chan struct{})
	go func() {
		s.s.Stop() // immediately stop. there is no catchup() to receive <-quitCh.
		close(doneCh)
	}()

	timer := time.NewTimer(1 * time.Second)
	select {
	case <-timer.C:
		s.T().Fatal("timeout")
	case <-doneCh:
		return
	}
}

// Module started when head is 0. PostInsertBlock() will update the lastCheckpoint.
func (s *SupplyTestSuite) TestWithoutCatchup() {
	t := s.T()
	require.Nil(t, s.s.Start())
	defer s.s.Stop()
	s.insertBlocks()

	assert.Equal(t, uint64(400), s.s.lastNum)

	for _, tc := range s.testcases() {
		ts, err := s.s.GetTotalSupply(tc.number)
		require.NoError(t, err)
		assert.Equal(t, tc.expectTotalSupply, ts)
	}
}

// Module started when head is already 400. catchup() will update the lastCheckpoint.
func (s *SupplyTestSuite) TestWithCatchup() {
	t := s.T()
	require.Nil(t, s.s.Start())
	s.insertBlocks()
	s.s.Stop()

	// Though insertBlocks->InsertChain->PostInsertBlock() has already wrote to the database,
	// we pretend the database is empty so we can test catchup() `lastNum < currNum` path.
	WriteLastSupplyCheckpointNumber(s.s.ChainKv, 0)

	require.Nil(t, s.s.Start())
	s.waitCatchup()
	defer s.s.Stop()

	for _, tc := range s.testcases() {
		ts, err := s.s.GetTotalSupply(tc.number)
		require.NoError(t, err)
		assert.Equal(t, tc.expectTotalSupply, ts)
	}
}

func (s *SupplyTestSuite) waitCatchup() {
	for i := 0; i < 1000; i++ { // wait 10 seconds until catchup complete
		if s.s.lastNum == 400 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	s.T().Fatal("Catchup not finished in time")
}
