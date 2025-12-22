// Copyright 2025 The kaia Authors
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

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.24;

import "./IHolderVerifier.sol";
import "./FnsaVerify.sol";
import "./IOperator.sol";
import "./IBridge.sol";
import "openzeppelin-contracts-5.0/access/Ownable.sol";
import {ReentrancyGuard} from "openzeppelin-contracts-5.0/utils/ReentrancyGuard.sol";

contract HolderVerifier is Ownable, ReentrancyGuard, IHolderVerifier {
    /* ========== CONSTANTS ========== */

    // 148.079656 KAIA/FNSA = 148079656000000 kei/cony ~ 148e12 kei/cony
    // e.g. 1 FNSA = 1e6 cony ~ (1e6 cony) * (148e12 kei/cony) = 148e18 kei = 148 KAIA
    uint256 public constant CONV_RATE = 148079656000000; // 148.079656 KAIA per FNSA, in kei

    /* ========== STATE VARIABLES ========== */

    mapping(string => uint256) public conyBalances; // fnsaAddr => conyBalance
    mapping(string => uint64) public provisionSeq; // fnsaAddr => provision seq (0 means not provisioned)
    string[] public fnsaAddrs;

    uint256 public override allConyBalances; // sum of all records in this contract. Includes provisioned amounts.
    uint256 public override provisionedConyBalances; // sum of all provisioned amounts in this contract.
    uint256 public override provisionedAccounts; // number of accounts that have been provisioned.

    address public operator;

    event OperatorChanged(address indexed previousOperator, address indexed newOperator);

    /* ========== CONSTRUCTOR ========== */

    constructor(address operatorAddr) Ownable(msg.sender) {
        _setOperator(operatorAddr);
    }

    /* ========== ADMIN FUNCTIONS ========== */

    function addRecord(string calldata fnsaAddr, uint256 conyBalance) public override onlyOwner {
        require(bytes(fnsaAddr).length > 0, "HolderVerifier: empty fnsaAddr");
        require(conyBalance > 0, "HolderVerifier: zero conyBalance");
        require(!isProvisioned(fnsaAddr), "HolderVerifier: cannot overwrite after provisioned");

        uint256 previousBalance = conyBalances[fnsaAddr];
        if (previousBalance == 0) {
            fnsaAddrs.push(fnsaAddr);
        } else {
            allConyBalances -= previousBalance;
        }

        conyBalances[fnsaAddr] = conyBalance;
        allConyBalances += conyBalance;

        emit RecordAdded(fnsaAddr, conyBalance);
    }

    function addRecords(
        string[] calldata fnsaAddrList,
        uint256[] calldata conyBalanceList
    ) external override onlyOwner {
        uint256 length = fnsaAddrList.length;
        require(length == conyBalanceList.length, "HolderVerifier: array length mismatch");
        require(length > 0, "HolderVerifier: empty arrays");

        for (uint256 i = 0; i < length; i++) {
            string calldata fnsaAddr = fnsaAddrList[i];
            uint256 conyBalance = conyBalanceList[i];

            addRecord(fnsaAddr, conyBalance);
        }

        emit RecordsAdded(length);
    }

    function changeOperator(address newOperator) external onlyOwner {
        _setOperator(newOperator);
    }

    /* ========== PUBLIC FUNCTIONS ========== */

    function requestProvision(
        bytes calldata publicKey,
        string calldata fnsaAddress,
        bytes32 messageHash,
        bytes calldata signature
    ) external override nonReentrant {
        require(bytes(fnsaAddress).length > 0, "HolderVerifier: empty fnsaAddress");

        uint256 conyBalance = conyBalances[fnsaAddress];
        require(conyBalance > 0, "HolderVerifier: no claimable balance");
        require(!isProvisioned(fnsaAddress), "HolderVerifier: already provisioned");

        address holderAddr = FnsaVerify.verify(publicKey, fnsaAddress, messageHash, signature);

        uint256 kaiaAmount = conyBalance * CONV_RATE;
        require(kaiaAmount > 0, "HolderVerifier: invalid amount");

        require(operator != address(0), "HolderVerifier: operator not set");
        address currentBridge = IOperator(operator).bridge();
        require(currentBridge != address(0), "HolderVerifier: bridge not set");

        // NOTE: Using nextProvisionSeq + 1 may be seen counter-intuitive,
        // but Bridge.sol expects nextProvisionSeq + 1 in the new provision.
        // Since HolderVerifier can execute the provision by itself,
        // nextProvisionSeq should increment as soon as submitTransaction is finished
        uint64 bridgeSeq = IBridge(currentBridge).nextProvisionSeq();
        uint64 seq = bridgeSeq + 1;

        IBridge.ProvisionData memory provisionData = IBridge.ProvisionData({
            seq: seq,
            sender: fnsaAddress,
            receiver: holderAddr,
            amount: kaiaAmount
        });

        bytes memory payload = abi.encodeWithSelector(IBridge.provision.selector, provisionData);

        IOperator(operator).submitTransaction(currentBridge, payload, uint256(seq));

        provisionSeq[fnsaAddress] = seq;
        provisionedConyBalances += conyBalance;
        provisionedAccounts += 1;

        emit ProvisionRequested(fnsaAddress, holderAddr, conyBalance, kaiaAmount, seq);
    }

    /* ========== VIEW FUNCTIONS ========== */

    function getRecord(string calldata fnsaAddr) external view override returns (uint256 conyBalance, uint64 seq) {
        return (conyBalances[fnsaAddr], provisionSeq[fnsaAddr]);
    }

    function getRecords(
        uint256 startIdx,
        uint256 maxCount
    )
        external
        view
        override
        returns (string[] memory fnsaAddrsResult, uint256[] memory conyBalancesResult, uint64[] memory seqs)
    {
        uint256 length = fnsaAddrs.length;
        require(startIdx < length, "HolderVerifier: startIdx out of bounds");

        // Return [startIdx, endIdx)
        uint256 endIdx = startIdx + maxCount;
        if (endIdx > length) {
            endIdx = length;
        }

        uint256 resultLength = endIdx - startIdx;
        fnsaAddrsResult = new string[](resultLength);
        conyBalancesResult = new uint256[](resultLength);
        seqs = new uint64[](resultLength);

        for (uint256 i = 0; i < resultLength; i++) {
            string memory addr = fnsaAddrs[startIdx + i];
            fnsaAddrsResult[i] = addr;
            conyBalancesResult[i] = conyBalances[addr];
            seqs[i] = provisionSeq[addr];
        }
    }

    function getRecordCount() external view override returns (uint256) {
        return fnsaAddrs.length;
    }

    function isProvisioned(string memory fnsaAddr) public view override returns (bool) {
        return provisionSeq[fnsaAddr] != 0;
    }

    /* ========== INTERNAL FUNCTIONS ========== */

    function _setOperator(address newOperator) internal {
        require(newOperator != address(0), "HolderVerifier: operator zero address");
        address previousOperator = operator;
        operator = newOperator;
        emit OperatorChanged(previousOperator, newOperator);
    }
}
