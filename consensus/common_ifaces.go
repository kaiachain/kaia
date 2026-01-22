package consensus

import (
	"github.com/kaiachain/kaia/blockchain/types"
)

type Verifier interface {
	VerifyHeader(chain ChainReader, header *types.Header) error
	VerifySeals(chain ChainReader, header *types.Header) error
}
