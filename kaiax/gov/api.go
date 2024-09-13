package gov

import (
	"github.com/kaiachain/kaia/networks/rpc"
)

func (g *govModule) APIs() []rpc.API {
	return g.hgm.APIs()
}
