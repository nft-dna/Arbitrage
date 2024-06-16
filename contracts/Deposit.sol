// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "./Interfaces.sol";

contract Deposit {
    //mapping(address => uint256) public etherBalances;
    mapping(address => mapping(address => uint256)) public tokenBalances;
    address[] public tokenAddresses; // Array to store token addresses
    uint256 public tokenAddressesCount = 0;

    mapping(address => bool) private tokenAddressesPresent;


    //event DepositEther(address indexed _from, uint256 _value);
    event DepositToken(address indexed _token, address indexed _from, uint256 _value);
    //event WithdrawEther(address indexed _to, uint256 _value);
    event WithdrawToken(address indexed _token, address indexed _to, uint256 _value);

    //function depositEther() public payable {
    //    depositEtherSucceded(msg.sender, msg.value);
    //}

    //function depositEtherSucceded(address issuer, uint256 _amount) internal {
    //    etherBalances[issuer] += _amount;
    //    emit DepositEther(issuer, _amount);
    //}

    function depositToken(address _token, uint256 _amount) public {
        require(_amount > 0, "Deposit amount must be greater than 0");
        require(IERC20(_token).allowance(msg.sender, address(this)) >= _amount, "Allowance too low");

        require(IERC20(_token).transferFrom(msg.sender, address(this), _amount), "Token transfer failed");
        depositTokenSucceded(msg.sender, _token, _amount);
    }

    function depositTokenSucceded(address issuer, address _token, uint256 _amount) internal {
        tokenBalances[_token][issuer] = tokenBalances[_token][issuer] + _amount;        
        if (tokenAddressesPresent[_token] == false) {
            tokenAddresses.push(_token);
            tokenAddressesCount = tokenAddressesCount + 1;
            tokenAddressesPresent[_token] = true;
        }
        emit DepositToken(_token, issuer, _amount);
    }

    //function getEtherBalance() public view returns (uint256) {
    //    return etherBalances[msg.sender];
    //}

    function getTokenBalance(address _token) public view returns (uint256) {
        return tokenBalances[_token][msg.sender];
    }

    //function getTotalEtherBalance() public view returns (uint256) {
    //    return address(this).balance;
    //}

    function getTotalTokenBalance(address _token) public view returns (uint256) {
        return IERC20(_token).balanceOf(address(this));
    }    

    //function withdrawEther(uint256 _amount) public {
    //    require(etherBalances[msg.sender] >= _amount, "Insufficient Ether balance");
    //    etherBalances[msg.sender] -= _amount;
    //    payable(msg.sender).transfer(_amount);
    //    
    //    emit WithdrawEther(msg.sender, _amount);
    //}

    function withdrawToken(address _token, uint256 _amount) public {
        require(tokenBalances[_token][msg.sender] >= _amount, "Insufficient token balance");
        tokenBalances[_token][msg.sender] = tokenBalances[_token][msg.sender] - _amount;
        require(IERC20(_token).transfer(msg.sender, _amount), "Token transfer failed");
        
        emit WithdrawToken(_token, msg.sender, _amount);
    }
}
