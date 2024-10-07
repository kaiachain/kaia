package governance

import "github.com/kaiachain/kaia/common"

type GovernanceAPI struct {
	governance Engine // Node interfaced by this API
}

func NewGovernanceAPI(gov Engine) *GovernanceAPI {
	return &GovernanceAPI{governance: gov}
}

func (api *GovernanceAPI) NodeAddress() common.Address {
	return api.governance.NodeAddress()
}
