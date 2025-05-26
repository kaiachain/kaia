import "../../../system_contracts/kip247/IWKAIA.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MaliciousWKAIA is IWKAIA {
    function deposit() external payable override {
        // Do nothing
    }

    function withdraw(uint256 _amount) external override {
        revert("Withdrawal failed");
    }

    function transfer(address _to, uint256 _value) external override returns (bool) {
        return true;
    }

    function balanceOf(address _account) external view override returns (uint256) {
        return 0;
    }

    function approve(address _spender, uint256 _value) external override returns (bool) {
        return true;
    }

    function transferFrom(address _from, address _to, uint256 _value) external override returns (bool) {
        return true;
    }
}
