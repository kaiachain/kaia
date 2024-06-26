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

import "../../libs/openzeppelin-contracts-v2/contracts/token/ERC20/IERC20.sol";
import "../../libs/openzeppelin-contracts-v2/contracts/token/ERC20/ERC20Mintable.sol";
import "../../libs/openzeppelin-contracts-v2/contracts/math/SafeMath.sol";
import "../../libs/openzeppelin-contracts-v2/contracts/ownership/Ownable.sol";


contract BridgeOperator is Ownable {
    struct VotesData {
        address[] voters;   // voter list for deleting voted map
        mapping(address => bytes32) voted; // <operator, sha3(type, args, nonce)>

        bytes32[] voteKeys; // voteKey list for deleting voteCounts map
        mapping(bytes32 => uint8) voteCounts; // <sha3(type, args, nonce), uint8>
    }

    mapping(uint8 => mapping (uint64 => VotesData)) private votes; // <voteType, <nonce, VotesData>
    mapping(uint64 => bool) public closedValueTransferVotes; // <nonce, bool>

    uint64 public constant MAX_OPERATOR = 12;
    mapping(address => bool) public operators;
    address[] public operatorList;

    mapping(uint8 => uint8) public operatorThresholds; // <vote type, uint8>
    uint64 public configurationNonce;

    enum VoteType {
        ValueTransfer,
        Configuration,
        Max
    }

    event OperatorRegistered(address indexed operator);
    event OperatorDeregistered(address indexed operator);
    event OperatorThresholdChanged(VoteType indexed voteType, uint8 threshold);

    constructor() internal {
        for (uint8 i = 0; i < uint8(VoteType.Max); i++) {
            operatorThresholds[uint8(i)] = 1;
        }

        operators[msg.sender] = true;
        operatorList.push(msg.sender);
    }

    modifier onlyOperators()
    {
        require(operators[msg.sender], "msg.sender is not an operator");
        _;
    }

    function getOperatorList() external view returns(address[] memory) {
        return operatorList;
    }

    // _voteCommon handles common functionality for voting.
    function _voteCommon(VoteType _voteType, uint64 _nonce, bytes32 _voteKey)
        private
        returns(bool)
    {
        VotesData storage vote = votes[uint8(_voteType)][_nonce];

        // If the same voter voted again, revoke previous vote.
        bytes32 oldVoteKeyOfVoter = vote.voted[msg.sender];
        if (oldVoteKeyOfVoter == bytes32(0)) {
            vote.voters.push(msg.sender);
        } else {
            vote.voteCounts[oldVoteKeyOfVoter]--;
        }

        // Either the current voter has voted before or not, update the vote data.
        vote.voted[msg.sender] = _voteKey;

        if (vote.voteCounts[_voteKey] == 0) {
            vote.voteKeys.push(_voteKey);
        }
        vote.voteCounts[_voteKey]++;

        if (vote.voteCounts[_voteKey] >= operatorThresholds[uint8(_voteType)]) {
            return true;
        }
        return false;
    }

    // _voteValueTransfer votes value transfer transaction with the operator.
    function _voteValueTransfer(uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(!closedValueTransferVotes[_requestNonce], "closed vote");

        bytes32 voteKey = keccak256(msg.data);
        if (_voteCommon(VoteType.ValueTransfer, _requestNonce, voteKey)) {
            closedValueTransferVotes[_requestNonce] = true;
            return true;
        }

        return false;
    }

    // _voteConfiguration votes contract configuration transaction with the operator.
    function _voteConfiguration(uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(configurationNonce == _requestNonce, "nonce mismatch");

        bytes32 voteKey = keccak256(msg.data);
        if (_voteCommon(VoteType.Configuration, _requestNonce, voteKey)) {
            configurationNonce++;
            return true;
        }

        return false;
    }

    // registerOperator registers a new operator.
    function registerOperator(address _operator)
    external
    onlyOwner
    {
        require(operatorList.length < MAX_OPERATOR, "max operator limit");
        require(!operators[_operator], "exist operator");
        operators[_operator] = true;
        operatorList.push(_operator);
        emit OperatorRegistered(_operator);
    }

    // deregisterOperator deregisters the operator.
    //
    // Note that outstanding votes by the deregistered operator are not revoked.
    // This enables a subtle counterintuitive scenario.
    //
    // Suppose there are two operators A, B and C with threshold 2.
    // 1. Operator A votes on nonce N
    // 2. Owner deregisters A
    // 3. Operator B votes on nonce N, thereby executing the request N.
    // In this case the request was executed with A's vote after A is deregistered.
    //
    // The Owner shall recognize this issue and expect that operator deregistration
    // takes some time to be fully effective.
    function deregisterOperator(address _operator)
    external
    onlyOwner
    {
        require(operators[_operator]);
        delete operators[_operator];

        for (uint i = 0; i < operatorList.length; i++) {
           if (operatorList[i] == _operator) {
               operatorList[i] = operatorList[operatorList.length-1];
               operatorList.length--;
               break;
           }
        }
        emit OperatorDeregistered(_operator);
    }

    // setOperatorThreshold sets the operator threshold.
    function setOperatorThreshold(VoteType _voteType, uint8 _threshold)
    external
    onlyOwner
    {
        require(_threshold > 0, "zero threshold");
        require(operatorList.length >= _threshold, "bigger than num of operators");
        operatorThresholds[uint8(_voteType)] = _threshold;
        emit OperatorThresholdChanged(_voteType, _threshold);
    }
}
