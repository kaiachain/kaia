// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2020 The klaytn Authors
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
// This file is derived from eth/downloader/queue_test.go (2020/07/24).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package downloader

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/faker"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	testdb  = database.NewMemoryDBManager()
	genesis = blockchain.GenesisBlockForTesting(testdb, testAddress, big.NewInt(1000000000))
)

// makeChain creates a chain of n blocks starting at and including parent.
// the returned hash chain is ordered head->parent. In addition, every 2nd block
// contains a transaction.
func makeChain(n int, seed byte, parent *types.Block, empty bool) ([]*types.Block, []types.Receipts) {
	blocks, receipts := blockchain.GenerateChain(params.TestChainConfig, parent, faker.NewFaker(), testdb, n, func(i int, block *blockchain.BlockGen) {
		block.SetRewardbase(common.Address{seed})
		// Add one tx to every second block
		if !empty && i%2 == 0 {
			signer := types.MakeSigner(params.TestChainConfig, block.Number())
			tx, err := types.SignTx(types.NewTransaction(block.TxNonce(testAddress), common.Address{seed}, big.NewInt(1000), params.TxGas, nil, nil), signer, testKey)
			if err != nil {
				panic(err)
			}
			block.AddTx(tx)
		}
	})
	return blocks, receipts
}

type chainData struct {
	blocks       []*types.Block
	stakingInfos []*staking.P2PStakingInfo
	offset       int
}

var (
	chain        *chainData
	emptyChain   *chainData
	testInterval uint64 = 4
)

func init() {
	// Create a chain of blocks to import. 128 blocks are created and a transaction is contained on every 2nd block
	targetBlocks := 128

	var stakingInfos []*staking.P2PStakingInfo
	for i := 4; i <= 128; i += 4 {
		stakingInfos = append(stakingInfos, &staking.P2PStakingInfo{BlockNum: uint64(i)})
	}

	blocks, _ := makeChain(targetBlocks, 0, genesis, false)
	chain = &chainData{blocks, stakingInfos, 0}

	blocks, _ = makeChain(targetBlocks, 0, genesis, true)
	emptyChain = &chainData{blocks, stakingInfos, 0}
}

func (chain *chainData) headers() []*types.Header {
	hdrs := make([]*types.Header, len(chain.blocks))
	for i, b := range chain.blocks {
		hdrs[i] = b.Header()
	}
	return hdrs
}

func (chain *chainData) Len() int {
	return len(chain.blocks)
}

func dummyPeer(id string) *peerConnection {
	p := &peerConnection{
		id:      id,
		lacking: make(map[common.Hash]struct{}),
	}
	return p
}

func newBodyDeliveryQueue(t *testing.T, id string) (*queue, *peerConnection, *fetchRequest) {
	t.Helper()

	q := newQueue(50, 50, uint64(istanbul.WeightedRandom), params.TestChainConfig)
	q.Prepare(1, FullSync)
	q.Schedule(chain.headers()[:12], 1)

	peer := dummyPeer(id)
	request, _, _ := q.ReserveBodies(peer, 2)
	if request == nil || len(request.Headers) != 2 {
		t.Fatalf("expected two body requests, got %#v", request)
	}
	return q, peer, request
}

func bodyTransactions(request *fetchRequest) [][]*types.Transaction {
	bodies := make([][]*types.Transaction, len(request.Headers))
	for i, header := range request.Headers {
		bodies[i] = chain.blocks[header.Number.Uint64()-1].Transactions()
	}
	return bodies
}

func newReceiptDeliveryQueue(t *testing.T, id string) (*queue, *peerConnection, *fetchRequest) {
	t.Helper()

	q := newQueue(50, 50, uint64(istanbul.WeightedRandom), params.TestChainConfig)
	q.Prepare(1, FastSync)
	headers := make([]*types.Header, 12)
	for i, header := range chain.headers()[:12] {
		headers[i] = types.CopyHeader(header)
		headers[i].ReceiptHash = types.DeriveReceiptsRoot(receiptsForHeader(headers[i]), headers[i].Number)
		if i > 0 {
			headers[i].ParentHash = headers[i-1].Hash()
		}
	}
	q.Schedule(headers, 1)

	peer := dummyPeer(id)
	request, _, _ := q.ReserveReceipts(peer, 2)
	if request == nil || len(request.Headers) != 2 {
		t.Fatalf("expected two receipt requests, got %#v", request)
	}
	return q, peer, request
}

func receiptsForHeader(header *types.Header) []*types.Receipt {
	return []*types.Receipt{types.NewReceipt(types.ReceiptStatusSuccessful, common.Hash{}, header.Number.Uint64())}
}

func receiptLists(request *fetchRequest) [][]*types.Receipt {
	receipts := make([][]*types.Receipt, len(request.Headers))
	for i, header := range request.Headers {
		receipts[i] = receiptsForHeader(header)
	}
	return receipts
}

type deliveryQueueTest struct {
	name            string
	item            string
	newQueue        func(*testing.T, string) (*queue, *peerConnection, *fetchRequest)
	reserve         func(*queue, *peerConnection, int) (*fetchRequest, bool, bool)
	expire          func(*queue, time.Duration) map[string]int
	deliver         func(*queue, string, *fetchRequest) (int, error)
	deliverInvalid  func(*queue, string) (int, error)
	deliverEmpty    func(*queue, string) (int, error)
	pendingRequest  func(*queue, string) *fetchRequest
	expiredRequests func(*queue, string) int
	invalidErr      error
}

func deliveryQueueTests() []deliveryQueueTest {
	return []deliveryQueueTest{
		{
			name:     "bodies",
			item:     "body",
			newQueue: newBodyDeliveryQueue,
			reserve:  (*queue).ReserveBodies,
			expire:   (*queue).ExpireBodies,
			deliver: func(q *queue, id string, request *fetchRequest) (int, error) {
				return q.DeliverBodies(id, bodyTransactions(request))
			},
			deliverInvalid: func(q *queue, id string) (int, error) {
				return q.DeliverBodies(id, [][]*types.Transaction{nil})
			},
			deliverEmpty:    func(q *queue, id string) (int, error) { return q.DeliverBodies(id, nil) },
			pendingRequest:  func(q *queue, id string) *fetchRequest { return q.blockPendPool[id] },
			expiredRequests: func(q *queue, id string) int { return len(q.blockExpired[id]) },
			invalidErr:      errInvalidBody,
		},
		{
			name:     "receipts",
			item:     "receipt",
			newQueue: newReceiptDeliveryQueue,
			reserve:  (*queue).ReserveReceipts,
			expire:   (*queue).ExpireReceipts,
			deliver: func(q *queue, id string, request *fetchRequest) (int, error) {
				return q.DeliverReceipts(id, receiptLists(request))
			},
			deliverInvalid: func(q *queue, id string) (int, error) {
				return q.DeliverReceipts(id, [][]*types.Receipt{nil})
			},
			deliverEmpty:    func(q *queue, id string) (int, error) { return q.DeliverReceipts(id, nil) },
			pendingRequest:  func(q *queue, id string) *fetchRequest { return q.receiptPendPool[id] },
			expiredRequests: func(q *queue, id string) int { return len(q.receiptExpired[id]) },
			invalidErr:      errInvalidReceipt,
		},
	}
}

func TestTrackExpiredRequestEvictsOldest(t *testing.T) {
	expired := make(map[string][][]*types.Header)
	headers := chain.headers()[:maxExpiredRequests+1]

	for _, header := range headers {
		trackExpiredRequest(expired, "peer-1", &fetchRequest{Headers: []*types.Header{header}})
	}
	if got := len(expired["peer-1"]); got != maxExpiredRequests {
		t.Fatalf("expected %d retained requests, got %d", maxExpiredRequests, got)
	}
	if got := expired["peer-1"][0][0]; got != headers[1] {
		t.Fatalf("expected oldest retained header %s, got %s", headers[1].Hash(), got.Hash())
	}
}

func TestDeliverBodiesRejectsInvalidDelivery(t *testing.T) {
	q, peer, request := newBodyDeliveryQueue(t, "peer-1")
	pending := q.PendingBlocks()

	first := request.Headers[0].Number.Uint64()
	validPrefix := chain.blocks[first-1].Transactions()
	accepted, err := q.DeliverBodies(peer.id, [][]*types.Transaction{validPrefix, nil})

	if accepted != 1 {
		t.Fatalf("expected one accepted body, got %d", accepted)
	}
	if !errors.Is(err, errInvalidBody) {
		t.Fatalf("expected %v, got %v", errInvalidBody, err)
	}
	if got, want := q.PendingBlocks(), pending+1; got != want {
		t.Fatalf("expected %d pending bodies after rejection, got %d", want, got)
	}
}

func TestDeliverBodiesPreservesInvalidBodyOverResultSlotError(t *testing.T) {
	q, peer, request := newBodyDeliveryQueue(t, "peer-1")

	// Make result slot lookup fail after validation has already identified the
	// invalid suffix.
	q.resultCache.lock.Lock()
	q.resultCache.items = nil
	q.resultCache.lock.Unlock()

	first := request.Headers[0].Number.Uint64()
	validPrefix := chain.blocks[first-1].Transactions()
	accepted, err := q.DeliverBodies(peer.id, [][]*types.Transaction{validPrefix, nil})

	if accepted != 1 {
		t.Fatalf("expected one accepted body, got %d", accepted)
	}
	if !errors.Is(err, errInvalidBody) {
		t.Fatalf("expected %v, got %v", errInvalidBody, err)
	}
	if errors.Is(err, errStaleDelivery) {
		t.Fatalf("stale slot error replaced invalid body error: %v", err)
	}
}

func TestDeliverBodiesEmptyResponseMarksBodiesLacking(t *testing.T) {
	q, peer, request := newBodyDeliveryQueue(t, "peer-1")
	pending := q.PendingBlocks()

	accepted, err := q.DeliverBodies(peer.id, nil)
	if accepted != 0 || err != nil {
		t.Fatalf("expected empty response to be accepted as missing, got accepted=%d err=%v", accepted, err)
	}
	if got, want := q.PendingBlocks(), pending+len(request.Headers); got != want {
		t.Fatalf("expected %d pending bodies after empty response, got %d", want, got)
	}
	for _, header := range request.Headers {
		if !peer.Lacks(header.Hash()) {
			t.Fatalf("expected peer to be marked lacking body %s", header.Hash())
		}
	}
}

func TestDeliveriesIdentifyLateResponse(t *testing.T) {
	for _, tt := range deliveryQueueTests() {
		t.Run(tt.name, func(t *testing.T) {
			q, peer, expired := tt.newQueue(t, "peer-1")
			expired.Time = time.Now().Add(-time.Hour)
			if got := tt.expire(q, time.Second)[peer.id]; got != len(expired.Headers) {
				t.Fatalf("expected %d expired %s requests, got %d", len(expired.Headers), tt.item, got)
			}

			// Assign the expired work elsewhere, then give the original peer a newer
			// request. The late response must not consume that newer request.
			if request, _, _ := tt.reserve(q, dummyPeer("peer-2"), 2); request == nil {
				t.Fatalf("expected expired %s requests to be reassigned", tt.item)
			}
			current, _, _ := tt.reserve(q, peer, 2)
			if current == nil {
				t.Fatalf("expected a newer %s request", tt.item)
			}

			if accepted, err := tt.deliver(q, peer.id, expired); accepted != 0 || !errors.Is(err, errLateDelivery) {
				t.Fatalf("expected late delivery, got accepted=%d err=%v", accepted, err)
			}
			if tt.pendingRequest(q, peer.id) != current {
				t.Fatalf("late delivery consumed the current %s request", tt.item)
			}
			if accepted, err := tt.deliver(q, peer.id, current); accepted != len(current.Headers) || err != nil {
				t.Fatalf("current %s delivery failed: accepted=%d err=%v", tt.item, accepted, err)
			}
		})
	}
}

func TestDeliveriesPreferCurrentResponse(t *testing.T) {
	for _, tt := range deliveryQueueTests() {
		t.Run(tt.name, func(t *testing.T) {
			q, peer, expired := tt.newQueue(t, "peer-1")
			expired.Time = time.Now().Add(-time.Hour)
			tt.expire(q, time.Second)

			tt.reserve(q, dummyPeer("peer-2"), 2)
			current, _, _ := tt.reserve(q, peer, 2)
			if current == nil {
				t.Fatalf("expected a current %s request", tt.item)
			}

			if accepted, err := tt.deliver(q, peer.id, current); accepted != len(current.Headers) || err != nil {
				t.Fatalf("current %s delivery failed: accepted=%d err=%v", tt.item, accepted, err)
			}
			if tt.expiredRequests(q, peer.id) != 1 {
				t.Fatalf("current %s delivery consumed expired request history", tt.item)
			}
			if accepted, err := tt.deliver(q, peer.id, expired); accepted != 0 || !errors.Is(err, errLateDelivery) {
				t.Fatalf("expected late %s delivery without a current request, got accepted=%d err=%v", tt.item, accepted, err)
			}
		})
	}
}

func TestDeliveriesRejectInvalidWithExpiredHistory(t *testing.T) {
	for _, tt := range deliveryQueueTests() {
		t.Run(tt.name, func(t *testing.T) {
			q, peer, expired := tt.newQueue(t, "peer-1")
			expired.Time = time.Now().Add(-time.Hour)
			tt.expire(q, time.Second)

			tt.reserve(q, dummyPeer("peer-2"), 2)
			if current, _, _ := tt.reserve(q, peer, 2); current == nil {
				t.Fatalf("expected a current %s request", tt.item)
			}
			if accepted, err := tt.deliverInvalid(q, peer.id); accepted != 0 || !errors.Is(err, tt.invalidErr) {
				t.Fatalf("expected %v, got accepted=%d err=%v", tt.invalidErr, accepted, err)
			}
		})
	}
}

func TestExpiredRequestsAcrossResetAndRevoke(t *testing.T) {
	for _, tt := range deliveryQueueTests() {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("reset retains pending request", func(t *testing.T) {
				q, peer, request := tt.newQueue(t, "peer-1")
				q.Reset(50, 50)

				if accepted, err := tt.deliver(q, peer.id, request); accepted != 0 || !errors.Is(err, errLateDelivery) {
					t.Fatalf("expected late %s delivery after reset, got accepted=%d err=%v", tt.item, accepted, err)
				}
			})

			t.Run("revoke clears history", func(t *testing.T) {
				q, peer, request := tt.newQueue(t, "peer-1")
				q.Reset(50, 50)
				q.Revoke(peer.id)

				if accepted, err := tt.deliver(q, peer.id, request); accepted != 0 || !errors.Is(err, errNoFetchesPending) {
					t.Fatalf("expected unrequested %s delivery after revoke, got accepted=%d err=%v", tt.item, accepted, err)
				}
			})
		})
	}
}

func TestEmptyDeliveriesDoNotConsumeExpiredHistory(t *testing.T) {
	for _, tt := range deliveryQueueTests() {
		t.Run(tt.name, func(t *testing.T) {
			q, peer, _ := tt.newQueue(t, "peer-1")
			q.Reset(50, 50)
			before := tt.expiredRequests(q, peer.id)

			if accepted, err := tt.deliverEmpty(q, peer.id); accepted != 0 || !errors.Is(err, errLateDelivery) {
				t.Fatalf("expected empty %s response to be late, got accepted=%d err=%v", tt.item, accepted, err)
			}
			if got := tt.expiredRequests(q, peer.id); got != before {
				t.Fatalf("empty %s response consumed expired history: have=%d want=%d", tt.item, got, before)
			}
		})
	}
}

func TestEmptyDeliveriesTrackCurrentRequest(t *testing.T) {
	for _, tt := range deliveryQueueTests() {
		t.Run(tt.name, func(t *testing.T) {
			q, peer, expired := tt.newQueue(t, "peer-1")
			expired.Time = time.Now().Add(-time.Hour)
			tt.expire(q, time.Second)

			tt.reserve(q, dummyPeer("peer-2"), 2)
			current, _, _ := tt.reserve(q, peer, 2)
			if current == nil {
				t.Fatalf("expected a current %s request", tt.item)
			}
			before := tt.expiredRequests(q, peer.id)

			if accepted, err := tt.deliverEmpty(q, peer.id); accepted != 0 || err != nil {
				t.Fatalf("expected empty current %s response to be accepted, got accepted=%d err=%v", tt.item, accepted, err)
			}
			if pending := tt.pendingRequest(q, peer.id); pending != nil {
				t.Fatalf("empty current %s response left request pending", tt.item)
			}
			if got, want := tt.expiredRequests(q, peer.id), before+1; got != want {
				t.Fatalf("empty current %s response was not tracked: have=%d want=%d", tt.item, got, want)
			}

			next, _, _ := tt.reserve(q, peer, 2)
			if next == nil {
				t.Fatalf("expected a newer %s request", tt.item)
			}
			if accepted, err := tt.deliver(q, peer.id, current); accepted != 0 || !errors.Is(err, errLateDelivery) {
				t.Fatalf("expected delayed current %s response to be late, got accepted=%d err=%v", tt.item, accepted, err)
			}
			if pending := tt.pendingRequest(q, peer.id); pending != next {
				t.Fatalf("late %s delivery consumed newer request", tt.item)
			}
		})
	}
}

func TestDeliverReceiptsRejectsInvalidDelivery(t *testing.T) {
	q, peer, request := newReceiptDeliveryQueue(t, "peer-1")
	pending := q.PendingReceipts()

	validPrefix := receiptLists(request)[0]
	accepted, err := q.DeliverReceipts(peer.id, [][]*types.Receipt{validPrefix, nil})

	if accepted != 1 {
		t.Fatalf("expected one accepted receipt list, got %d", accepted)
	}
	if !errors.Is(err, errInvalidReceipt) {
		t.Fatalf("expected %v, got %v", errInvalidReceipt, err)
	}
	if got, want := q.PendingReceipts(), pending+1; got != want {
		t.Fatalf("expected %d pending receipts after rejection, got %d", want, got)
	}
}

func TestBasics(t *testing.T) {
	numOfBlocks := len(chain.blocks)
	numOfReceipts := len(chain.blocks) / 2
	numOfStakingInfos := len(chain.stakingInfos)

	config := params.TestChainConfig.Copy()
	config.Governance.Reward.StakingUpdateInterval = testInterval
	config.KaiaCompatibleBlock = nil
	q := newQueue(10, 10, uint64(istanbul.WeightedRandom), config)
	if !q.Idle() {
		t.Errorf("new queue should be idle")
	}
	q.Prepare(1, FastSync)
	if res := q.Results(false); len(res) != 0 {
		t.Fatal("new queue should have 0 results")
	}

	// Schedule a batch of headers
	q.Schedule(chain.headers(), 1)
	if q.Idle() {
		t.Errorf("queue should not be idle")
	}
	if got, exp := q.PendingBlocks(), numOfBlocks; got != exp {
		t.Errorf("wrong pending block count, got %d, exp %d", got, exp)
	}
	// Only non-empty receipts get added to task-queue
	if got, exp := q.PendingReceipts(), numOfReceipts; got != exp {
		t.Errorf("wrong pending receipt count, got %d, exp %d", got, exp)
	}
	// staking info on every 4th block get added to task-queue
	if got, exp := q.PendingStakingInfos(), numOfStakingInfos; got != exp {
		t.Errorf("wrong pending receipt count, got %d, exp %d", got, exp)
	}
	// Items are now queued for downloading, next step is that we tell the
	// queue that a certain peer will deliver them for us
	{
		peer := dummyPeer("peer-1")
		fetchReq, _, throttle := q.ReserveBodies(peer, 50)
		if !throttle {
			// queue size is only 10, so throttling should occur
			t.Fatal("should throttle")
		}
		// But we should still get the first things to fetch
		if got, exp := len(fetchReq.Headers), 5; got != exp {
			t.Fatalf("expected %d requests, got %d", exp, got)
		}
		if got, exp := fetchReq.Headers[0].Number.Uint64(), uint64(1); got != exp {
			t.Fatalf("expected header %d, got %d", exp, got)
		}
	}
	if got, exp := q.blockTaskQueue.Size(), numOfBlocks-10; got != exp {
		t.Errorf("expected block task queue to be %d, got %d", exp, got)
	}
	if got, exp := q.receiptTaskQueue.Size(), numOfReceipts; got != exp {
		t.Errorf("expected receipt task queue to be %d, got %d", exp, got)
	}
	if got, exp := q.stakingInfoTaskQueue.Size(), numOfStakingInfos; got != exp {
		t.Errorf("expected staking info task queue to be %d, got %d", exp, got)
	}
	{
		peer := dummyPeer("peer-2")
		fetchReq, _, throttle := q.ReserveBodies(peer, 50)

		// The second peer should hit throttling
		if !throttle {
			t.Fatalf("should not throttle")
		}
		// And not get any fetches at all, since it was throttled to begin with
		if fetchReq != nil {
			t.Fatalf("should have no fetches, got %d", len(fetchReq.Headers))
		}
	}
	if got, exp := q.blockTaskQueue.Size(), numOfBlocks-10; got != exp {
		t.Errorf("expected block task queue to be %d, got %d", exp, got)
	}
	if got, exp := q.receiptTaskQueue.Size(), numOfReceipts; got != exp {
		t.Errorf("expected receipt task queue to be %d, got %d", exp, got)
	}
	if got, exp := q.stakingInfoTaskQueue.Size(), numOfStakingInfos; got != exp {
		t.Errorf("expected staking info task queue to be %d, got %d", exp, got)
	}
	{
		// The receipt delivering peer should not be affected
		// by the throttling of body deliveries
		peer := dummyPeer("peer-3")
		fetchReq, _, throttle := q.ReserveReceipts(peer, 50)
		if !throttle {
			// queue size is only 10, so throttling should occur
			t.Fatal("should throttle")
		}
		// But we should still get the first things to fetch
		if got, exp := len(fetchReq.Headers), 5; got != exp {
			t.Fatalf("expected %d requests, got %d", exp, got)
		}
		if got, exp := fetchReq.Headers[0].Number.Uint64(), uint64(1); got != exp {
			t.Fatalf("expected header %d, got %d", exp, got)
		}
	}
	if got, exp := q.blockTaskQueue.Size(), numOfBlocks-10; got != exp {
		t.Fatalf("expected block task queue size %d, got %d", exp, got)
	}
	if got, exp := q.receiptTaskQueue.Size(), numOfReceipts-5; got != exp {
		t.Fatalf("expected receipt task queue size %d, got %d", exp, got)
	}
	if got, exp := q.stakingInfoTaskQueue.Size(), numOfStakingInfos; got != exp {
		t.Fatalf("expected staking info task queue size %d, got %d", exp, got)
	}
	{
		// The staking info delivering peer should not be affected
		// by the throttling of body deliveries
		peer := dummyPeer("peer-4")
		fetchReq, _, throttle := q.ReserveStakingInfos(peer, 50)
		if !throttle {
			// queue size is only 10, so throttling should occur
			t.Fatal("should throttle")
		}
		// But we should still get the first things to fetch
		if got, exp := len(fetchReq.Headers), 2; got != exp {
			t.Fatalf("expected %d requests, got %d", exp, got)
		}
		if got, exp := fetchReq.Headers[0].Number.Uint64(), uint64(4); got != exp {
			t.Fatalf("expected header %d, got %d", exp, got)
		}
	}
	if got, exp := q.blockTaskQueue.Size(), numOfBlocks-10; got != exp {
		t.Fatalf("expected block task queue size %d, got %d", exp, got)
	}
	if got, exp := q.receiptTaskQueue.Size(), numOfReceipts-5; got != exp {
		t.Fatalf("expected receipt task queue size %d, got %d", exp, got)
	}
	if got, exp := q.stakingInfoTaskQueue.Size(), numOfStakingInfos-2; got != exp {
		t.Fatalf("expected staking info task queue size %d, got %d", exp, got)
	}
	if got, exp := q.resultCache.countCompleted(), 0; got != exp {
		t.Errorf("wrong processable count, got %d, exp %d", got, exp)
	}
}

func TestScheduleAfterKaia(t *testing.T) {
	config := params.TestChainConfig.Copy()
	config.Governance.Reward.StakingUpdateInterval = testInterval
	config.KaiaCompatibleBlock = big.NewInt(21)

	numOfStakingInfos := 5 // [4, 8, 12, 16, 20]; After kaia fork, it won't be scheduled.

	q := newQueue(50, 50, uint64(istanbul.WeightedRandom), config)
	if !q.Idle() {
		t.Errorf("new queue should be idle")
	}
	q.Prepare(1, FastSync)
	if res := q.Results(false); len(res) != 0 {
		t.Fatal("new queue should have 0 results")
	}
	// Schedule a batch of headers
	q.Schedule(chain.headers(), 1)
	if q.Idle() {
		t.Errorf("queue should not be idle")
	}
	// staking info on every 4th block get added to task-queue
	if got, exp := q.PendingStakingInfos(), numOfStakingInfos; got != exp {
		t.Errorf("wrong pending receipt count, got %d, exp %d", got, exp)
	}
	if got, exp := q.stakingInfoTaskQueue.Size(), numOfStakingInfos; got != exp {
		t.Errorf("expected staking info task queue to be %d, got %d", exp, got)
	}
	{
		peer := dummyPeer("peer-1")
		fetchReq, _, _ := q.ReserveStakingInfos(peer, 50)
		// But we should still get the first things to fetch
		if got, exp := len(fetchReq.Headers), 5; got != exp {
			t.Fatalf("expected %d requests, got %d", exp, got)
		}
		if got, exp := fetchReq.Headers[0].Number.Uint64(), uint64(4); got != exp {
			t.Fatalf("expected header %d, got %d", exp, got)
		}
	}
	if got, exp := q.stakingInfoTaskQueue.Size(), 0; got != exp {
		t.Fatalf("expected staking info task queue size %d, got %d", exp, got)
	}
}

func TestEmptyBlocks(t *testing.T) {
	numOfBlocks := len(emptyChain.blocks)
	numOfStakingInfos := len(emptyChain.stakingInfos)
	config := params.TestChainConfig.Copy()
	config.Governance.Reward.StakingUpdateInterval = testInterval
	config.KaiaCompatibleBlock = nil

	q := newQueue(10, 10, uint64(istanbul.WeightedRandom), config)

	q.Prepare(1, FastSync)
	// Schedule a batch of headers
	q.Schedule(emptyChain.headers(), 1)
	if q.Idle() {
		t.Errorf("queue should not be idle")
	}
	if got, exp := q.PendingBlocks(), numOfBlocks; got != exp {
		t.Errorf("wrong pending block count, got %d, exp %d", got, exp)
	}
	if got, exp := q.PendingReceipts(), 0; got != exp {
		t.Errorf("wrong pending receipt count, got %d, exp %d", got, exp)
	}
	if got, exp := q.PendingStakingInfos(), numOfStakingInfos; got != exp {
		t.Errorf("wrong pending staking infos count, got %d, exp %d", got, exp)
	}
	// They won't be processable, because the fetchresults haven't been
	// created yet
	if got, exp := q.resultCache.countCompleted(), 0; got != exp {
		t.Errorf("wrong processable count, got %d, exp %d", got, exp)
	}

	// Items are now queued for downloading, next step is that we tell the
	// queue that a certain peer will deliver them for us
	// That should trigger all of them to suddenly become 'done'
	{
		// Reserve blocks
		peer := dummyPeer("peer-1")
		fetchReq, _, _ := q.ReserveBodies(peer, 50)

		// there should be nothing to fetch, blocks are empty
		if fetchReq != nil {
			t.Fatal("there should be no body fetch tasks remaining")
		}
	}
	if q.blockTaskQueue.Size() != numOfBlocks-10 {
		t.Errorf("expected block task queue to be %d, got %d", numOfBlocks-10, q.blockTaskQueue.Size())
	}
	if q.receiptTaskQueue.Size() != 0 {
		t.Errorf("expected receipt task queue to be %d, got %d", 0, q.receiptTaskQueue.Size())
	}
	if got, exp := q.stakingInfoTaskQueue.Size(), numOfStakingInfos; got != exp {
		t.Fatalf("expected staking info task queue size %d, got %d", exp, got)
	}
	{
		peer := dummyPeer("peer-3")
		fetchReq, _, _ := q.ReserveReceipts(peer, 50)

		// there should be nothing to fetch, blocks are empty
		if fetchReq != nil {
			t.Fatal("there should be no body fetch tasks remaining")
		}
	}
	if q.blockTaskQueue.Size() != numOfBlocks-10 {
		t.Errorf("expected block task queue to be %d, got %d", numOfBlocks-10, q.blockTaskQueue.Size())
	}
	if q.receiptTaskQueue.Size() != 0 {
		t.Errorf("expected receipt task queue to be %d, got %d", 0, q.receiptTaskQueue.Size())
	}
	if got, exp := q.stakingInfoTaskQueue.Size(), numOfStakingInfos; got != exp {
		t.Fatalf("expected staking info task queue size %d, got %d", exp, got)
	}
	{
		// The staking info delivering peer should not be affected
		// by the throttling of body deliveries
		peer := dummyPeer("peer-4")
		fetchReq, _, throttle := q.ReserveStakingInfos(peer, 50)
		if !throttle {
			// queue size is only 10, so throttling should occur
			t.Fatal("should throttle")
		}
		// But we should still get the first things to fetch
		if got, exp := len(fetchReq.Headers), 2; got != exp {
			t.Fatalf("expected %d requests, got %d", exp, got)
		}
		if got, exp := fetchReq.Headers[0].Number.Uint64(), uint64(4); got != exp {
			t.Fatalf("expected header %d, got %d", exp, got)
		}
	}
	if q.blockTaskQueue.Size() != numOfBlocks-10 {
		t.Errorf("expected block task queue to be %d, got %d", numOfBlocks-10, q.blockTaskQueue.Size())
	}
	if q.receiptTaskQueue.Size() != 0 {
		t.Errorf("expected receipt task queue to be %d, got %d", 0, q.receiptTaskQueue.Size())
	}
	if got, exp := q.stakingInfoTaskQueue.Size(), numOfStakingInfos-2; got != exp {
		t.Fatalf("expected staking info task queue size %d, got %d", exp, got)
	}
	if got, exp := q.resultCache.countCompleted(), 3; got != exp {
		t.Errorf("wrong processable count, got %d, exp %d", got, exp)
	}
}

// XTestDelivery does some more extensive testing of events that happen,
// blocks that become known and peers that make reservations and deliveries.
// disabled since it's not really a unit-test, but can be executed to test
// some more advanced scenarios
func XTestDelivery(t *testing.T) {
	// the outside network, holding blocks
	blo, rec := makeChain(128, 0, genesis, false)
	world := newNetwork()
	world.receipts = rec
	world.chain = blo
	world.progress(10)
	if false {
		log.Root().SetHandler(log.StdoutHandler)
	}
	q := newQueue(10, 10, uint64(istanbul.WeightedRandom), nil)
	var wg sync.WaitGroup
	q.Prepare(1, FastSync)
	wg.Go(func() {
		// deliver headers
		c := 1
		for {
			// fmt.Printf("getting headers from %d\n", c)
			hdrs := world.headers(c)
			l := len(hdrs)
			// fmt.Printf("scheduling %d headers, first %d last %d\n",
			//	l, hdrs[0].Number.Uint64(), hdrs[len(hdrs)-1].Number.Uint64())
			q.Schedule(hdrs, uint64(c))
			c += l
		}
	})
	wg.Go(func() {
		// collect results
		tot := 0
		for {
			res := q.Results(true)
			tot += len(res)
			fmt.Printf("got %d results, %d tot\n", len(res), tot)
			// Now we can forget about these
			world.forget(res[len(res)-1].Header.Number.Uint64())

		}
	})
	wg.Go(func() {
		// reserve body fetch
		i := 4
		for {
			peer := dummyPeer(fmt.Sprintf("peer-%d", i))
			f, _, _ := q.ReserveBodies(peer, rand.Intn(30))
			if f != nil {
				var txs [][]*types.Transaction
				numToSkip := rand.Intn(len(f.Headers))
				for _, hdr := range f.Headers[0 : len(f.Headers)-numToSkip] {
					txs = append(txs, world.getTransactions(hdr.Number.Uint64()))
				}
				time.Sleep(100 * time.Millisecond)
				_, err := q.DeliverBodies(peer.id, txs)
				if err != nil {
					fmt.Printf("delivered %d bodies %v\n", len(txs), err)
				}
			} else {
				i++
				time.Sleep(200 * time.Millisecond)
			}
		}
	})
	go func() {
		defer wg.Done()
		// reserve receiptfetch
		peer := dummyPeer("peer-3")
		for {
			f, _, _ := q.ReserveReceipts(peer, rand.Intn(50))
			if f != nil {
				var rcs [][]*types.Receipt
				for _, hdr := range f.Headers {
					rcs = append(rcs, world.getReceipts(hdr.Number.Uint64()))
				}
				_, err := q.DeliverReceipts(peer.id, rcs)
				if err != nil {
					fmt.Printf("delivered %d receipts %v\n", len(rcs), err)
				}
				time.Sleep(100 * time.Millisecond)
			} else {
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
	wg.Go(func() {
		for range 50 {
			time.Sleep(300 * time.Millisecond)
			// world.tick()
			// fmt.Printf("trying to progress\n")
			world.progress(rand.Intn(100))
		}
		for range 50 {
			time.Sleep(2990 * time.Millisecond)
		}
	})
	wg.Go(func() {
		for {
			time.Sleep(990 * time.Millisecond)
			fmt.Printf("world block tip is %d\n",
				world.chain[len(world.chain)-1].Header().Number.Uint64())
			fmt.Println(q.Stats())
		}
	})
	wg.Wait()
}

func newNetwork() *network {
	var l sync.RWMutex
	return &network{
		cond:   sync.NewCond(&l),
		offset: 1, // block 1 is at blocks[0]
	}
}

// represents the network
type network struct {
	offset   int
	chain    []*types.Block
	receipts []types.Receipts
	lock     sync.RWMutex
	cond     *sync.Cond
}

func (n *network) getTransactions(blocknum uint64) types.Transactions {
	index := blocknum - uint64(n.offset)
	return n.chain[index].Transactions()
}

func (n *network) getReceipts(blocknum uint64) types.Receipts {
	index := blocknum - uint64(n.offset)
	if got := n.chain[index].Header().Number.Uint64(); got != blocknum {
		fmt.Printf("Err, got %d exp %d\n", got, blocknum)
		panic("sd")
	}
	return n.receipts[index]
}

func (n *network) forget(blocknum uint64) {
	index := blocknum - uint64(n.offset)
	n.chain = n.chain[index:]
	n.receipts = n.receipts[index:]
	n.offset = int(blocknum)
}

func (n *network) progress(numBlocks int) {
	n.lock.Lock()
	defer n.lock.Unlock()
	// fmt.Printf("progressing...\n")
	newBlocks, newR := makeChain(numBlocks, 0, n.chain[len(n.chain)-1], false)
	n.chain = append(n.chain, newBlocks...)
	n.receipts = append(n.receipts, newR...)
	n.cond.Broadcast()
}

func (n *network) headers(from int) []*types.Header {
	numHeaders := 128
	var hdrs []*types.Header
	index := from - n.offset

	for index >= len(n.chain) {
		// wait for progress
		n.cond.L.Lock()
		// fmt.Printf("header going into wait\n")
		n.cond.Wait()
		index = from - n.offset
		n.cond.L.Unlock()
	}
	n.lock.RLock()
	defer n.lock.RUnlock()
	for i, b := range n.chain[index:] {
		hdrs = append(hdrs, b.Header())
		if i >= numHeaders {
			break
		}
	}
	return hdrs
}

// newStakingInfoQueue schedules the test chain on a fresh queue and reserves the
// first two staking info fetches (blocks 4 and 8) for "peer-1".
func newStakingInfoQueue(t *testing.T) *queue {
	config := params.TestChainConfig.Copy()
	config.Governance.Reward.StakingUpdateInterval = testInterval
	config.KaiaCompatibleBlock = nil

	q := newQueue(50, 50, uint64(istanbul.WeightedRandom), config)
	q.Prepare(1, FastSync)
	q.Schedule(chain.headers(), 1)

	fetchReq, _, _ := q.ReserveStakingInfos(dummyPeer("peer-1"), 2)
	if got, exp := len(fetchReq.Headers), 2; got != exp {
		t.Fatalf("expected %d requests, got %d", exp, got)
	}
	for i, exp := range []uint64{4, 8} {
		if got := fetchReq.Headers[i].Number.Uint64(); got != exp {
			t.Fatalf("expected header %d at index %d, got %d", exp, i, got)
		}
	}
	return q
}

// Tests that a staking info whose block number does not match the requested
// header is rejected, the fetch is returned to the queue, and a matching
// redelivery is accepted.
func TestDeliverStakingInfosRejectsMismatchedBlockNum(t *testing.T) {
	q := newStakingInfoQueue(t)
	pending := q.PendingStakingInfos()

	// The first entry claims block 8 for the header of block 4.
	accepted, err := q.DeliverStakingInfos("peer-1", []*staking.P2PStakingInfo{{BlockNum: 8}, {BlockNum: 8}})
	if accepted != 0 {
		t.Fatalf("expected no accepted staking info, got %d", accepted)
	}
	if !errors.Is(err, errInvalidStakingInfo) {
		t.Fatalf("expected %v, got %v", errInvalidStakingInfo, err)
	}
	if got, exp := q.PendingStakingInfos(), pending+2; got != exp {
		t.Fatalf("expected %d pending staking infos after rejection, got %d", exp, got)
	}

	// The same blocks can be reserved again and a matching delivery is accepted.
	fetchReq, _, _ := q.ReserveStakingInfos(dummyPeer("peer-1"), 2)
	if got, exp := len(fetchReq.Headers), 2; got != exp {
		t.Fatalf("expected %d requests, got %d", exp, got)
	}
	accepted, err = q.DeliverStakingInfos("peer-1", []*staking.P2PStakingInfo{{BlockNum: 4}, {BlockNum: 8}})
	if accepted != 2 || err != nil {
		t.Fatalf("expected 2 accepted staking infos, got %d (err %v)", accepted, err)
	}
	if got, exp := q.PendingStakingInfos(), pending; got != exp {
		t.Fatalf("expected %d pending staking infos after delivery, got %d", exp, got)
	}
}

// Tests that a batch is accepted only up to the first mismatched entry and the
// rest is returned to the queue.
func TestDeliverStakingInfosAcceptsValidPrefix(t *testing.T) {
	q := newStakingInfoQueue(t)
	pending := q.PendingStakingInfos()

	accepted, err := q.DeliverStakingInfos("peer-1", []*staking.P2PStakingInfo{{BlockNum: 4}, {BlockNum: 12}})
	if accepted != 1 {
		t.Fatalf("expected 1 accepted staking info, got %d", accepted)
	}
	if err == nil {
		t.Fatal("expected a partial failure, got nil")
	}
	if got, exp := q.PendingStakingInfos(), pending+1; got != exp {
		t.Fatalf("expected %d pending staking infos, got %d", exp, got)
	}
	fetchReq, _, _ := q.ReserveStakingInfos(dummyPeer("peer-1"), 1)
	if got, exp := fetchReq.Headers[0].Number.Uint64(), uint64(8); got != exp {
		t.Fatalf("expected the rejected block %d to be requeued, got %d", exp, got)
	}
}
