// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from quorum/consensus/istanbul/backend/snapshot.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package backend

import (
	"bytes"
	"encoding/json"
	"math/big"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/consensus/istanbul/validator"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

const (
	dbKeySnapshotPrefix = "istanbul-snapshot"
)

// Snapshot is the state of the authorization voting at a given point in time.
type Snapshot struct {
	Epoch         uint64                // The number of blocks after which to checkpoint and reset the pending votes
	Number        uint64                // Block number where the snapshot was created
	Hash          common.Hash           // Block hash where the snapshot was created
	ValSet        istanbul.ValidatorSet // Set of authorized validators at this moment
	Policy        uint64
	CommitteeSize uint64
	Votes         []governance.GovernanceVote      // List of votes cast in chronological order
	Tally         []governance.GovernanceTallyItem // Current vote tally to avoid recalculating
}

func effectiveParams(g gov.GovModule, number uint64) (epoch uint64, policy uint64, committeeSize uint64) {
	pset := g.EffectiveParamSet(number)
	epoch = pset.Epoch
	policy = pset.ProposerPolicy
	committeeSize = pset.CommitteeSize

	return
}

// newSnapshot create a new snapshot with the specified startup parameters. This
// method does not initialize the set of recent validators, so only ever use if for
// the genesis block.
func newSnapshot(g gov.GovModule, number uint64, hash common.Hash, valSet istanbul.ValidatorSet, chainConfig *params.ChainConfig) *Snapshot {
	epoch, policy, committeeSize := effectiveParams(g, number+1)

	snap := &Snapshot{
		Epoch:         epoch,
		Number:        number,
		Hash:          hash,
		ValSet:        valSet,
		Policy:        policy,
		CommitteeSize: committeeSize,
		Votes:         make([]governance.GovernanceVote, 0),
		Tally:         make([]governance.GovernanceTallyItem, 0),
	}
	return snap
}

// loadSnapshot loads an existing snapshot from the database.
func loadSnapshot(db database.DBManager, hash common.Hash) (*Snapshot, error) {
	blob, err := db.ReadIstanbulSnapshot(hash)
	if err != nil {
		return nil, err
	}
	snap := new(Snapshot)
	if err := json.Unmarshal(blob, snap); err != nil {
		return nil, err
	}
	return snap, nil
}

// store inserts the snapshot into the database.
func (s *Snapshot) store(db database.DBManager) error {
	blob, err := json.Marshal(s)
	if err != nil {
		return err
	}

	db.WriteIstanbulSnapshot(s.Hash, blob)
	return nil
}

// copy creates a deep copy of the snapshot, though not the individual votes.
func (s *Snapshot) copy() *Snapshot {
	cpy := &Snapshot{
		Epoch:         s.Epoch,
		Number:        s.Number,
		Hash:          s.Hash,
		ValSet:        s.ValSet.Copy(),
		Policy:        s.Policy,
		CommitteeSize: s.CommitteeSize,
		Votes:         make([]governance.GovernanceVote, len(s.Votes)),
		Tally:         make([]governance.GovernanceTallyItem, len(s.Tally)),
	}

	copy(cpy.Votes, s.Votes)
	copy(cpy.Tally, s.Tally)

	return cpy
}

// checkVote return whether it's a valid vote
func (s *Snapshot) checkVote(address common.Address, authorize bool) bool {
	_, validator := s.ValSet.GetByAddress(address)
	return (validator != nil && !authorize) || (validator == nil && authorize)
}

// apply creates a new authorization snapshot by applying the given headers to
// the original one.
func (s *Snapshot) apply(headers []*types.Header, gov governance.Engine, govModule gov.GovModule, addr common.Address, policy uint64, chain consensus.ChainReader, stakingModule staking.StakingModule, writable bool) (*Snapshot, error) {
	// Allow passing in no headers for cleaner code
	if len(headers) == 0 {
		return s, nil
	}
	// Sanity check that the headers can be applied
	for i := 0; i < len(headers)-1; i++ {
		if headers[i+1].Number.Uint64() != headers[i].Number.Uint64()+1 {
			return nil, errInvalidVotingChain
		}
	}
	if headers[0].Number.Uint64() != s.Number+1 {
		return nil, errInvalidVotingChain
	}

	// Iterate through the headers and create a new snapshot
	snap := s.copy()

	// Copy values which might be changed by governance vote
	snap.Epoch, snap.Policy, snap.CommitteeSize = effectiveParams(govModule, snap.Number+1)

	for _, header := range headers {
		// Remove any votes on checkpoint blocks
		number := header.Number.Uint64()

		// Resolve the authorization key and check against validators
		validator, err := ecrecover(header)
		if err != nil {
			return nil, err
		}
		if _, v := snap.ValSet.GetByAddress(validator); v == nil {
			return nil, errUnauthorized
		}

		if number%snap.Epoch == 0 {
			if writable {
				gov.UpdateCurrentSet(number)
				if len(header.Governance) > 0 {
					gov.WriteGovernanceForNextEpoch(number, header.Governance)
				}
				gov.ClearVotes(number)
			}
			snap.Votes = make([]governance.GovernanceVote, 0)
			snap.Tally = make([]governance.GovernanceTallyItem, 0)
		}

		// Reload governance values
		snap.Epoch, snap.Policy, snap.CommitteeSize = effectiveParams(govModule, number+1)

		snap.ValSet, snap.Votes, snap.Tally = gov.HandleGovernanceVote(snap.ValSet, snap.Votes, snap.Tally, header, validator, addr, writable)
		if policy == uint64(params.WeightedRandom) {
			// Snapshot of block N (Snapshot_N) should contain proposers for N+1 and following blocks.
			// Validators for Block N+1 can be calculated based on the staking information from the previous stakingUpdateInterval block.
			// If the governance mode is single, the governing node is added to validator all the time.
			//
			// Proposers for Block N+1 can be calculated from the nearest previous proposersUpdateInterval block.
			// Refresh proposers in Snapshot_N using previous proposersUpdateInterval block for N+1, if not updated yet.

			// because snapshot(num)'s ValSet = validators for num+1
			pset := govModule.EffectiveParamSet(number + 1)

			isSingle := (pset.GovernanceMode == "single")
			govNode := pset.GoverningNode
			minStaking := pset.MinimumStake.Uint64()

			if err := snap.ValSet.RefreshValSet(number+1, chain.Config(), isSingle, govNode, minStaking, stakingModule); err != nil {
				logger.Trace("Skip refreshing validators while creating snapshot", "snap.Number", snap.Number, "err", err)
			}

			// Do not refresh proposers from the kaia fork block.
			// The proposer is calculated every block in `CalcProposer` function after the randao fork.
			if !chain.Config().IsKaiaForkEnabled(big.NewInt(int64(number + 1))) {
				pHeader := chain.GetHeaderByNumber(params.CalcProposerBlockNumber(number + 1))
				if pHeader != nil {
					if err := snap.ValSet.RefreshProposers(pHeader.Hash(), pHeader.Number.Uint64(), chain.Config()); err != nil {
						// There are three error cases and they just don't refresh proposers
						// (1) no validator at all
						// (2) invalid formatted hash
						// (3) no staking info available
						logger.Trace("Skip refreshing proposers while creating snapshot", "snap.Number", snap.Number, "pHeader.Number", pHeader.Number.Uint64(), "err", err)
					}
				} else {
					logger.Trace("Can't refreshing proposers while creating snapshot due to lack of required header", "snap.Number", snap.Number)
				}
			}
		}
	}
	snap.Number += uint64(len(headers))
	snap.Hash = headers[len(headers)-1].Hash()

	if snap.ValSet.Policy() == istanbul.WeightedRandom {
		snap.ValSet.SetBlockNum(snap.Number)

		bigNum := new(big.Int).SetUint64(snap.Number)
		if chain.Config().IsRandaoForkBlockParent(bigNum) {
			// The ForkBlock must select proposers using MixHash but (ForkBlock - 1) has no MixHash. Using ZeroMixHash instead.
			snap.ValSet.SetMixHash(params.ZeroMixHash)
		} else if chain.Config().IsRandaoForkEnabled(bigNum) {
			// Feed parent MixHash
			snap.ValSet.SetMixHash(headers[len(headers)-1].MixHash)
		}
	}
	snap.ValSet.SetSubGroupSize(snap.CommitteeSize)

	if writable {
		gov.SetTotalVotingPower(snap.ValSet.TotalVotingPower())
		gov.SetMyVotingPower(snap.getMyVotingPower(addr))
	}

	return snap, nil
}

func (s *Snapshot) getMyVotingPower(addr common.Address) uint64 {
	for _, a := range s.ValSet.List() {
		if a.Address() == addr {
			return a.VotingPower()
		}
	}
	return 0
}

// validators retrieves the list of authorized validators in ascending order.
func (s *Snapshot) validators() []common.Address {
	validators := make([]common.Address, 0, s.ValSet.Size())
	for _, validator := range s.ValSet.List() {
		validators = append(validators, validator.Address())
	}
	return sortValidatorArray(validators)
}

// demotedValidators retrieves the list of authorized, but demoted validators in ascending order.
func (s *Snapshot) demotedValidators() []common.Address {
	demotedValidators := make([]common.Address, 0, len(s.ValSet.DemotedList()))
	for _, demotedValidator := range s.ValSet.DemotedList() {
		demotedValidators = append(demotedValidators, demotedValidator.Address())
	}
	return sortValidatorArray(demotedValidators)
}

func (s *Snapshot) committee(prevHash common.Hash, view *istanbul.View) []common.Address {
	committeeList := s.ValSet.SubList(prevHash, view)

	committee := make([]common.Address, 0, len(committeeList))
	for _, v := range committeeList {
		committee = append(committee, v.Address())
	}
	return committee
}

func sortValidatorArray(validators []common.Address) []common.Address {
	for i := 0; i < len(validators); i++ {
		for j := i + 1; j < len(validators); j++ {
			if bytes.Compare(validators[i][:], validators[j][:]) > 0 {
				validators[i], validators[j] = validators[j], validators[i]
			}
		}
	}
	return validators
}

type snapshotJSON struct {
	Epoch  uint64                           `json:"epoch"`
	Number uint64                           `json:"number"`
	Hash   common.Hash                      `json:"hash"`
	Votes  []governance.GovernanceVote      `json:"votes"`
	Tally  []governance.GovernanceTallyItem `json:"tally"`

	// for validator set
	Validators   []common.Address        `json:"validators"`
	Policy       istanbul.ProposerPolicy `json:"policy"`
	SubGroupSize uint64                  `json:"subgroupsize"`

	// for weighted validator
	RewardAddrs       []common.Address `json:"rewardAddrs"`
	VotingPowers      []uint64         `json:"votingPower"`
	Weights           []uint64         `json:"weight"`
	Proposers         []common.Address `json:"proposers"`
	ProposersBlockNum uint64           `json:"proposersBlockNum"`
	DemotedValidators []common.Address `json:"demotedValidators"`
	MixHash           []byte           `json:"mixHash,omitempty"`
}

func (s *Snapshot) toJSONStruct() *snapshotJSON {
	var rewardAddrs []common.Address
	var votingPowers []uint64
	var weights []uint64
	var proposers []common.Address
	var proposersBlockNum uint64
	var validators []common.Address
	var demotedValidators []common.Address
	var mixHash []byte

	if s.ValSet.Policy() == istanbul.WeightedRandom {
		validators, demotedValidators, rewardAddrs, votingPowers, weights, proposers, proposersBlockNum, mixHash = validator.GetWeightedCouncilData(s.ValSet)
	} else {
		validators = s.validators()
	}

	return &snapshotJSON{
		Epoch:             s.Epoch,
		Number:            s.Number,
		Hash:              s.Hash,
		Votes:             s.Votes,
		Tally:             s.Tally,
		Validators:        validators,
		Policy:            istanbul.ProposerPolicy(s.Policy),
		SubGroupSize:      s.CommitteeSize,
		RewardAddrs:       rewardAddrs,
		VotingPowers:      votingPowers,
		Weights:           weights,
		Proposers:         proposers,
		ProposersBlockNum: proposersBlockNum,
		DemotedValidators: demotedValidators,
		MixHash:           mixHash,
	}
}

// Unmarshal from a json byte array
func (s *Snapshot) UnmarshalJSON(b []byte) error {
	var j snapshotJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	s.Epoch = j.Epoch
	s.Number = j.Number
	s.Hash = j.Hash
	s.Votes = j.Votes
	s.Tally = j.Tally

	if j.Policy == istanbul.WeightedRandom {
		s.ValSet = validator.NewWeightedCouncil(j.Validators, j.DemotedValidators, j.RewardAddrs, j.VotingPowers, j.Weights, j.Policy, j.SubGroupSize, j.Number, j.ProposersBlockNum, nil)
		validator.RecoverWeightedCouncilProposer(s.ValSet, j.Proposers)
		s.ValSet.SetMixHash(j.MixHash)
	} else {
		s.ValSet = validator.NewSubSet(j.Validators, j.Policy, j.SubGroupSize)
	}
	return nil
}

// Marshal to a json byte array
func (s *Snapshot) MarshalJSON() ([]byte, error) {
	j := s.toJSONStruct()
	return json.Marshal(j)
}

// prepareSnapshotApply is a helper function to prepare snapshot and headers for the given block number and hash.
// It returns the snapshot, headers, and error if any.
func (sb *backend) prepareSnapshotApply(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) (*Snapshot, []*types.Header, error) {
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *Snapshot
	)

	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := sb.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			break
		}
		// If an on-disk checkpoint snapshot can be found, use that
		if params.IsCheckpointInterval(number) {
			if s, err := loadSnapshot(sb.db, hash); err == nil {
				logger.Trace("Loaded voting snapshot form disk", "number", number, "hash", hash)
				snap = s
				break
			}
		}
		// If we're at block zero, make a snapshot
		if number == 0 {
			var err error
			if snap, err = sb.initSnapshot(chain); err != nil {
				return nil, nil, err
			}
			break
		}
		// No snapshot for this header, gather the header and move backward
		if header := getPrevHeaderAndUpdateParents(chain, number, hash, &parents); header == nil {
			return nil, nil, consensus.ErrUnknownAncestor
		} else {
			headers = append(headers, header)
			number, hash = number-1, header.ParentHash
		}
	}
	// Previous snapshot found, apply any pending headers on top of it
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}

	return snap, headers, nil
}

// GetKaiaHeadersForSnapshotApply returns the headers need to be applied to create snapshot for the given block number.
// Note that it only returns headers for kaia fork enabled blocks.
func (sb *backend) GetKaiaHeadersForSnapshotApply(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) ([]*types.Header, error) {
	_, headers, err := sb.prepareSnapshotApply(chain, number, hash, parents)
	if err != nil {
		return nil, err
	}

	kaiaHeaders := []*types.Header{}
	for i := 0; i < len(headers); i++ {
		if chain.Config().IsKaiaForkEnabled(new(big.Int).Add(headers[i].Number, big.NewInt(1))) {
			kaiaHeaders = headers[i:]
			break
		}
	}

	return kaiaHeaders, nil
}

// snapshot retrieves the state of the authorization voting at a given point in time.
// There's in-memory snapshot and on-disk snapshot. On-disk snapshot is stored every checkpointInterval blocks.
// Moreover, if the block has no in-memory or on-disk snapshot, before generating snapshot, it gathers the header and apply the vote in it.
func (sb *backend) snapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header, writable bool) (*Snapshot, error) {
	snap, headers, err := sb.prepareSnapshotApply(chain, number, hash, parents)
	if err != nil {
		return nil, err
	}

	pset := sb.govModule.EffectiveParamSet(snap.Number)
	snap, err = snap.apply(headers, sb.governance, sb.govModule, sb.address, pset.ProposerPolicy, chain, sb.stakingModule, writable)
	if err != nil {
		return nil, err
	}

	// If we've generated a new checkpoint snapshot, save to disk
	if writable && params.IsCheckpointInterval(snap.Number) && len(headers) > 0 {
		if sb.governance.CanWriteGovernanceState(snap.Number) {
			sb.governance.WriteGovernanceState(snap.Number, true)
		}
		if err = snap.store(sb.db); err != nil {
			return nil, err
		}
		logger.Trace("Stored voting snapshot to disk", "number", snap.Number, "hash", snap.Hash)
	}

	sb.regen(chain, headers)

	sb.recents.Add(snap.Hash, snap)
	return snap, err
}

// regen commits snapshot data to database
// regen is triggered if there is any checkpoint block in the `headers`.
// For each checkpoint block, this function verifies the existence of its snapshot in DB and stores one if missing.
/*
 Triggered:
 |   ^                          ^                          ^                          ^  ...|
     SI                 SI*(last snapshot)                 SI                         SI
       			   | header1, .. headerN |
 Not triggered: (Guaranteed SI* was committed before )
 |   ^                          ^                          ^                          ^  ...|
     SI                 SI*(last snapshot)                 SI                         SI
	                            | header1, .. headerN |
*/
func (sb *backend) regen(chain consensus.ChainReader, headers []*types.Header) {
	// Prevent nested call. Ignore header length one
	// because it was handled before the `regen` called.
	if !sb.isRestoringSnapshots.CompareAndSwap(false, true) || len(headers) <= 1 {
		return
	}
	defer func() {
		sb.isRestoringSnapshots.Store(false)
	}()

	var (
		from        = headers[0].Number.Uint64()
		to          = headers[len(headers)-1].Number.Uint64()
		start       = time.Now()
		commitTried = false
	)

	// Shortcut: No missing snapshot data to be processed.
	if to-(to%uint64(params.CheckpointInterval)) < from {
		return
	}

	for _, header := range headers {
		var (
			hn = header.Number.Uint64()
			hh = header.Hash()
		)
		if params.IsCheckpointInterval(hn) {
			// Store snapshot data if it was not committed before
			if loadSnap, _ := sb.db.ReadIstanbulSnapshot(hh); loadSnap != nil {
				continue
			}
			snap, err := sb.snapshot(chain, hn, hh, nil, false)
			if err != nil {
				logger.Warn("[Snapshot] Snapshot restoring failed", "len(headers)", len(headers), "from", from, "to", to, "headerNumber", hn)
				continue
			}
			if err = snap.store(sb.db); err != nil {
				logger.Warn("[Snapshot] Snapshot restoring failed", "len(headers)", len(headers), "from", from, "to", to, "headerNumber", hn)
			}
			commitTried = true
		}
	}
	if commitTried { // This prevents pushing too many logs by potential DoS attack
		logger.Trace("[Snapshot] Snapshot restoring completed", "len(headers)", len(headers), "from", from, "to", to, "elapsed", time.Since(start))
	}
}
