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

package impl

/*
// voteTestData is the testdata for vote tests.
// Not only votes, but also other functions such as VerifyHeader handles votes.
// To minimize the maintenance effort to manage the testdata, manage the expected result in voteTestData
type testNetworkInfo struct {
	validators []common.Address
	proposer   common.Address
}

type testHandleValidatorVoteExpectation struct {
	expectCList valset.AddressList
	expectError error
}

func TestValsetModule_HandleValidatorVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range []struct {
		name                      string
		networkInfo               testNetworkInfo
		voteData                  headergov.VoteData
		expectHandleValidatorVote testHandleValidatorVoteExpectation
	}{
		{
			name:                      "remove a validator",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", valset.AddressList([]common.Address{n[1]}).ToString()),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[2], n[3]}, nil},
		},
		{
			name:                      "add a validator",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", valset.AddressList([]common.Address{n[4]}).ToString()),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3], n[4]}, nil},
		},
		{
			name:                      "remove multiple validators",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", valset.AddressList([]common.Address{n[1], n[2], n[3]}).ToString()),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn}, nil},
		},
		{
			name:                      "add multiple validators",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", valset.AddressList([]common.Address{n[2], n[3]}).ToString()),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3]}, nil},
		},
		{
			name:                      "govnode cannot be removed",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", valset.AddressList([]common.Address{tgn}).ToString()),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3]}, nil},
		},
		{
			name:                      "cannot add existing validator",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", valset.AddressList([]common.Address{n[1]}).ToString()),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3]}, nil},
		},
		{
			name:                      "cannot remove non-exist validator",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", valset.AddressList([]common.Address{n[4]}).ToString()),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3]}, nil},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			vModule, tm, err := newTestVModule(ctrl)
			assert.NoError(t, err)

			tm.prepareMockExpectGovParam(2, testProposerPolicy, testSubGroupSize, tgn)

			err = writeCouncil(vModule.ChainKv, 1, tc.networkInfo.validators)
			assert.NoError(t, err)

			byteVote, err := tc.voteData.ToVoteBytes()
			assert.NoError(t, err)

			err = vModule.HandleValidatorVote(2, byteVote)
			assert.Equal(t, tc.expectHandleValidatorVote.expectError, err)

			cList, err := readCouncil(vModule.ChainKv, uint64(2))
			sort.Sort(tc.expectHandleValidatorVote.expectCList)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectHandleValidatorVote.expectCList, valset.AddressList(cList))
		})
	}
}
*/
