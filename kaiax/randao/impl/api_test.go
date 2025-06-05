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
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/v2/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/v2/common/hexutil"
	"github.com/kaiachain/kaia/v2/datasync/downloader"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/networks/rpc"
	"github.com/kaiachain/kaia/v2/storage/database"
	"github.com/stretchr/testify/assert"
)

func TestAPI_getBlsInfos(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		db     = database.NewMemoryDBManager()
		alloc  = testAllocStorage()
		config = testRandaoForkChainConfig(big.NewInt(0))

		pubKey1 = hexutil.MustDecode("0x" + "b716443d8d1b3c1230d1d186b1db0db80f79f72805646ba8135b98242df276bdbfb5dea0201c0258d6b60f30724f28e3")
		pop1    = hexutil.MustDecode("0x" + "85ffe933f8bdf4d86ddbb7060355987838acf84f39f45eea309f0a7e4cc2f63afb7a57682f75b8f44e68b64cc12299a701b7acbc5a650c7bc9cbac98a93e76c06c0607a567cbac14eb02e2596ae2b48d11a36bda4c7166dea4ba8b28db8d7d63")
		pubKey2 = (hexutil.MustDecode("0x" + "a5b6d96a1bb2bd8ec5480d112dc6bbad46ec08937b9320187928c0ed27339791186f581397c5a9679e49f6ac459d5a48"))
		pop2    = hexutil.MustDecode("0x" + "9658933c3a8618765618afddf6031465f4ab550a11c47b9edcfa205c01a5b498e02f2584821a8aa1ee4d5db9f9db9ef10337b129e814f09266447dca4a2eb8643c18ba797fdb699b14e6fbb68d49f775f3981bd2bdec58bcb53ae9eeeab45165")
		pubKey3 = hexutil.MustDecode("0x" + "a2093da481a55e7e374de2fa19a8d9acbf055a52048d697d87de864fab9d334bbd4d838c68d53022f355c06fb4cd6722")
		pop3    = hexutil.MustDecode("0x" + "b0e820e4eef472f45853b704c36d1b536d10235c9ab1517a66795d24908d1d81cc23de6714f339771c54e32361276a8d0efecf154983b4386d03d36a011f847e7b5e657637e506fd5298f97b64026b790ca9b2e23eeddab9c2ec0adb3c4157d5")

		expected = map[string]interface{}{
			"0x0000000000000000000000000000000000000001": map[string]interface{}{
				"publicKey": hex.EncodeToString(pubKey1),
				"pop":       hex.EncodeToString(pop1),
				"verifyErr": nil,
			},
			"0x0000000000000000000000000000000000000002": map[string]interface{}{
				"publicKey": hex.EncodeToString(pubKey2),
				"pop":       hex.EncodeToString(pop2),
				"verifyErr": nil,
			},
			"0x0000000000000000000000000000000000000003": map[string]interface{}{
				"publicKey": hex.EncodeToString(pubKey3),
				"pop":       hex.EncodeToString(pop3),
				"verifyErr": nil,
			},
		}
	)

	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	mRandao := NewRandaoModule()
	fakeDownloader := &downloader.FakeDownloader{}
	mRandao.Init(&InitOpts{
		ChainConfig: config,
		Chain:       backend.BlockChain(),
		Downloader:  fakeDownloader,
	})

	kaiaAPI := newRandaoAPI(mRandao)

	actual, err := kaiaAPI.GetBlsInfos(rpc.BlockNumber(1))
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
