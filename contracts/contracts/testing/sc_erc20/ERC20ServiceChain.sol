// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

pragma solidity 0.5.6;

import "../../libs/openzeppelin-contracts-v2/contracts/token/ERC20/ERC20.sol";
import "../../libs/openzeppelin-contracts-v2/contracts/utils/Address.sol";
import "../../libs/openzeppelin-contracts-v2/contracts/ownership/Ownable.sol";
import "../../service_chain/sc_erc20/IERC20BridgeReceiver.sol";


/**
 * @title ERC20ServiceChain
 * @dev ERC20 service chain value transfer logic for 1-step transfer.
 */
contract ERC20ServiceChain is ERC20, Ownable {
    using Address for address;
    address public bridge;

    constructor(address _bridge) internal {
        setBridge(_bridge);
    }

    function setBridge(address _bridge) public onlyOwner {
        if (!_bridge.isContract()) {
            revert("bridge is not a contract");
        }
        bridge = _bridge;
    }

    function requestValueTransfer(uint256 _amount, address _to, uint256 _feeLimit, bytes calldata _extraData) external {
        require(transfer(bridge, _amount.add(_feeLimit)), "requestValueTransfer: transfer failed");
        IERC20BridgeReceiver(bridge).onERC20Received(msg.sender, _to, _amount, _feeLimit, _extraData);
    }
}
