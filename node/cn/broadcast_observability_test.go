package cn

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/internal/poolobserve"
	"github.com/prometheus/client_golang/prometheus"
)

// Exercise real peer selection: a connected peer may already know the entire
// batch, so connection count alone cannot explain an absent send attempt.
func TestBroadcastNoRecipients(t *testing.T) {
	for _, role := range []common.ConnType{common.CONSENSUSNODE, common.PROXYNODE, common.ENDPOINTNODE} {
		for _, scenario := range []string{"disconnected", "known", "eligible", "empty"} {
			t.Run(role.String()+"/"+scenario, func(t *testing.T) {
				defer poolobserve.EnableForTesting()()
				reg := prometheus.NewRegistry()
				poolobserve.Register(reg)
				peers := newPeerSet()
				pm := &ProtocolManager{nodetype: role, peers: peers}
				tx := types.NewTransaction(0, common.Address{}, big.NewInt(1), 21000, big.NewInt(25000000000), nil)
				txs := types.Transactions{tx}
				if scenario == "empty" {
					txs = nil
				}
				if scenario == "known" || scenario == "eligible" {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					peer := NewMockPeer(ctrl)
					peers.cnpeers[common.Address{1}] = peer
					peers.peers["test"] = peer
					peer.EXPECT().ConnType().Return(common.CONSENSUSNODE).AnyTimes()
					peer.EXPECT().KnowsTx(tx.Hash()).Return(scenario == "known")
					if scenario == "eligible" {
						if role == common.CONSENSUSNODE {
							peer.EXPECT().AsyncSendTransactions(gomock.Any())
						} else {
							peer.EXPECT().SendTransactions(gomock.Any()).Return(nil)
						}
					}
				}
				pm.BroadcastTxs(txs)
				want := float64(0)
				if scenario == "known" || scenario == "disconnected" {
					want = 1
				}
				if got := propagationMetric(t, reg, "kaia_poolobs_events_total", "broadcast_no_recipients"); got != want {
					t.Fatalf("got %g want %g", got, want)
				}
			})
		}
	}
}

func TestRebroadcastNoRecipients(t *testing.T) {
	for _, role := range []common.ConnType{common.CONSENSUSNODE, common.PROXYNODE, common.ENDPOINTNODE} {
		for _, scenario := range []string{"disconnected", "connected", "empty"} {
			t.Run(role.String()+"/"+scenario, func(t *testing.T) {
				defer poolobserve.EnableForTesting()()
				reg := prometheus.NewRegistry()
				poolobserve.Register(reg)
				peers := newPeerSet()
				pm := &ProtocolManager{nodetype: role, peers: peers}
				txs := types.Transactions{types.NewTransaction(0, common.Address{}, big.NewInt(1), 21000, big.NewInt(25000000000), nil)}
				if scenario == "empty" {
					txs = nil
				}
				if scenario == "connected" {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()
					peer := NewMockPeer(ctrl)
					peers.cnpeers[common.Address{1}] = peer
					peers.peers["test"] = peer
					peer.EXPECT().ConnType().Return(common.CONSENSUSNODE).AnyTimes()
					// Resend deliberately ignores the known-transaction filter.
					if role != common.CONSENSUSNODE {
						peer.EXPECT().SendTransactions(gomock.Any()).Return(nil)
					}
				}
				pm.ReBroadcastTxs(txs)
				want := float64(0)
				if scenario == "disconnected" && role != common.CONSENSUSNODE {
					want = 1
				}
				if got := propagationMetric(t, reg, "kaia_poolobs_events_total", "rebroadcast_no_recipients"); got != want {
					t.Fatalf("got %g want %g", got, want)
				}
			})
		}
	}
}
