package engine

import (
	"crypto/ecdsa"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	istanbulBackend "github.com/kaiachain/kaia/consensus/istanbul/backend"
	"github.com/kaiachain/kaia/consensus/kaiabft"
)

// EngineType identifies the consensus engine variant.
type EngineType string

const (
	Istanbul EngineType = "istanbul"
	KaiaBFT  EngineType = "kaiabft"
)

// Opts bundles the inputs needed to instantiate a consensus engine.
type Opts struct {
	EngineType     EngineType
	IstanbulConfig *istanbul.Config
	PrivateKey     *ecdsa.PrivateKey
	NodeType       common.ConnType
}

// NewEngine creates the consensus engine selected by the given configuration.
func NewEngine(opts *Opts) consensus.Engine {
	if opts == nil || opts.IstanbulConfig == nil || opts.PrivateKey == nil {
		panic("NewEngine receives invalid opts")
	}

	sealer := istanbul.NewSealerImpl(opts.PrivateKey)

	if opts.EngineType == KaiaBFT {
		return kaiabft.New(&kaiabft.Opts{
			Timeout:    opts.IstanbulConfig.Timeout,
			PrivateKey: opts.PrivateKey,
			NodeType:   opts.NodeType,
			Sealer:     sealer,
		})
	}

	return istanbulBackend.New(&istanbulBackend.BackendOpts{
		IstanbulConfig: opts.IstanbulConfig,
		PrivateKey:     opts.PrivateKey,
		NodeType:       opts.NodeType,
	})
}
