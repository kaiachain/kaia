package bft

import (
	"crypto/ecdsa"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	istanbulBackend "github.com/kaiachain/kaia/consensus/istanbul/backend"
)

// Opts bundles the inputs needed to instantiate a consensus engine.
type Opts struct {
	IstanbulConfig *istanbul.Config
	PrivateKey     *ecdsa.PrivateKey
	NodeType       common.ConnType
}

// NewEngine creates the consensus engine selected by the given configuration.
func NewEngine(opts *Opts) consensus.Engine {
	return istanbulBackend.New(&istanbulBackend.BackendOpts{
		IstanbulConfig: opts.IstanbulConfig,
		PrivateKey:     opts.PrivateKey,
		NodeType:       opts.NodeType,
	})
}
