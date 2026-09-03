// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
//
// The kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the kaia library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTxInternalDataUnmarshalJSONEmptySignatures ensures that decoding a JSON
// transaction with a missing, empty, or null "signatures" field returns an
// error instead of panicking. Such payloads can be produced by a malicious or
// tampered JSON-RPC endpoint, and the decoders must not index into an empty
// TxSignatures slice.
func TestTxInternalDataUnmarshalJSONEmptySignatures(t *testing.T) {
	types := []struct {
		name string
		gen  func() TxInternalData
	}{
		{"Legacy", genLegacyTransaction},
		{"EthereumAccessList", genAccessListTransaction},
		{"EthereumDynamicFee", genDynamicFeeTransaction},
		{"EthereumSetCode", genSetCodeTransaction},
		{"EthereumBlob", genBlobTransaction},
	}

	// Each malformed "signatures" value that must be rejected gracefully.
	malformed := map[string]interface{}{
		"empty":   []interface{}{},    // "signatures": []
		"null":    []interface{}{nil}, // "signatures": [null]
		"missing": nil,                // "signatures" key absent
	}

	for _, tt := range types {
		for sigName, sigValue := range malformed {
			t.Run(tt.name+"/"+sigName, func(t *testing.T) {
				// Marshal a well-formed transaction, then tamper the signatures field.
				raw, err := json.Marshal(tt.gen())
				require.NoError(t, err)

				var m map[string]interface{}
				require.NoError(t, json.Unmarshal(raw, &m))
				if sigName == "missing" {
					delete(m, "signatures")
				} else {
					m["signatures"] = sigValue
				}
				tampered, err := json.Marshal(m)
				require.NoError(t, err)

				// Decoding must not panic and must be rejected specifically
				// by the empty-signatures guard.
				dec := newTxInternalDataSerializer()
				require.NotPanics(t, func() {
					err = json.Unmarshal(tampered, dec)
				})
				require.ErrorIs(t, err, errEmptyTxSignatures)
			})
		}
	}
}
