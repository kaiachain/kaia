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
package sc

import (
	"bytes"
	"strings"

	"github.com/kaiachain/kaia/v2/accounts/abi"
	"github.com/pkg/errors"
)

var ErrUnknownEvent = errors.New("Unknown event type")

var RequestValueTransferEncodeABIs = map[uint]string{
	2: `[{
			"anonymous":false,
			"inputs": [{
				"name": "uri",
				"type": "string"
			}],
			"name": "packedURI",
			"type": "event"
		}]`,
}

func UnpackEncodedData(ver uint8, packed []byte) map[string]interface{} {
	switch ver {
	case 2:
		encodedEvent := map[string]interface{}{}
		abi, err := abi.JSON(strings.NewReader(RequestValueTransferEncodeABIs[2]))
		if err != nil {
			logger.Error("Failed to ABI setup", "err", err)
			return nil
		}
		if err := abi.UnpackIntoMap(encodedEvent, "packedURI", packed); err != nil {
			logger.Error("Failed to unpack the values", "err", err)
			return nil
		}
		return encodedEvent
	default:
		logger.Error(ErrUnknownEvent.Error(), "encodingVer", ver)
		return nil
	}
}

func GetURI(ev IRequestValueTransferEvent) string {
	switch evType := ev.(type) {
	case RequestValueTransferEncodedEvent:
		decoded := UnpackEncodedData(evType.EncodingVer, evType.EncodedData)
		uri, ok := decoded["uri"].(string)
		if !ok {
			return ""
		}
		if len(uri) <= 64 {
			return ""
		}
		return string(bytes.Trim([]byte(uri[64:]), "\x00"))
	}
	return ""
}
