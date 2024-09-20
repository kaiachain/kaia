package impl

import (
	"github.com/kaiachain/kaia/networks/rpc"
)

func (g *GovModule) APIs() []rpc.API {
	return g.hgm.APIs()
}
