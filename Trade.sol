// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "./Interfaces.sol";
import "./Deposit.sol";

contract Trade is Deposit {

    // Addresses
    address payable OWNER;

    address[] public dexAddresses; // Array to store dex addresses
    mapping(address => DexInterfaceType) public dexInterface;    

    constructor() {
        OWNER = payable(msg.sender);
    }

    function isDexAdded(address _dex) internal view returns (bool) {
        for (uint256 i = 0; i < dexAddresses.length; i++) {
            if (dexAddresses[i] == _dex) {
                return true;
            }
        }
        return false;
    }

    function SetDex(address _dex, DexInterfaceType _interface) public onlyOwner {
        if (!isDexAdded(_dex)) {
            dexAddresses.push(_dex);
        }
        dexInterface[_dex] = _interface;
    }

    function getDexCount() public view returns (uint256) {
        return dexAddresses.length;
    }

    function arbEther(address _token, address _fromDex, address _toDex, uint256 _fromAmount) payable public {
        require(_fromAmount > 0, "Amount must be greater than 0");
        
        // Track original balance
        uint256 _startBalance = address(this).balance;
        require(_startBalance >= _fromAmount, "Insufficient token balance");
        
        // Perform the arb trade
        _tradeEther( _token, _fromDex, _toDex, _fromAmount, _startBalance);

        // Track result balance
        uint256 _endBalance = address(this).balance;

        // Require that arbitrage is profitable
        require(_endBalance > _startBalance, "End balance must exceed start balance.");

        depositEtherSucceded(msg.sender, _endBalance-_startBalance);
    }

    function arbToken(address _fromToken, address _toToken, address _fromDex, address _toDex, uint256 _fromAmount) /*payable*/ public {
        require(_fromAmount > 0, "Amount must be greater than 0");
        
        // Track original balance
        uint256 _startBalance = IERC20(_fromToken).balanceOf(address(this));
        require(_startBalance >= _fromAmount, "Insufficient token balance");
        
        // Perform the arb trade
        _tradeToken(_fromToken, _toToken, _fromDex, _toDex, _fromAmount, _startBalance);

        // Track result balance
        uint256 _endBalance = IERC20(_fromToken).balanceOf(address(this));

        // Require that arbitrage is profitable
        require(_endBalance > _startBalance, "End balance must exceed start balance.");

        depositTokenSucceded(msg.sender, _fromToken, _endBalance-_startBalance);
    }

    function _tradeEther(address _token, address _fromDex, address _toDex, uint256 _fromAmount, uint256 _startBalance) internal {
        require ( dexInterface[_fromDex] > DexInterfaceType.Unknown, "Unsupported from dex");
        require ( dexInterface[_toDex] > DexInterfaceType.Unknown, "Unsupported to dex");

        // Track the balance of the token RECEIVED from the trade
        //uint256 _startBalance = address(this).balance;

        address[] memory path = new address[](1);
        path[0] = _token;
        //path[1] = _toToken;
        /*uint256[] memory amounts_0 =*/ IUniswapV2Router(_fromDex).swapExactETHForTokens{value: _fromAmount}(_fromAmount, path, address(this), block.timestamp);
        // Get the amount of tokens received        
        //uint256 amountOut = amounts_0[1];
        // Calculate the how much of the token we received
        uint256 _afterBalance = address(this).balance;        
        uint256 _toAmount = _afterBalance - _startBalance;

        IERC20(_token).approve(address(_toDex), _toAmount); 
        path[0] = _token;
        //path[1] = _fromToken;
        /*uint256[] memory amounts_1 =*/ IUniswapV2Router(_toDex).swapExactTokensForETH(_toAmount,0, path, address(this), block.timestamp);
        // Read _toToken balance after swap        
        //uint256 amountOut_1 = amounts_1[1];
    }    

    function _tradeToken(address _fromToken, address _toToken, address _fromDex, address _toDex, uint256 _fromAmount, uint256 _startBalance) internal {
        require ( dexInterface[_fromDex] > DexInterfaceType.Unknown, "Unsupported from dex");
        require ( dexInterface[_toDex] > DexInterfaceType.Unknown, "Unsupported to dex");

        // Track the balance of the token RECEIVED from the trade
        //uint256 _startBalance = IERC20(_toToken).balanceOf(address(this));

        IERC20(_fromToken).approve(address(_fromDex) , _fromAmount); 
        address[] memory path = new address[](2);
        path[0] = _fromToken;
        path[1] = _toToken;
        /*uint256[] memory amounts_0 =*/ IUniswapV2Router(_fromDex).swapExactTokensForTokens(_fromAmount, 0, path, address(this), block.timestamp);
        // Get the amount of tokens received        
        //uint256 amountOut = amounts_0[1];
        // Calculate the how much of the token we received
        uint256 _afterBalance = IERC20(_toToken).balanceOf(address(this));        
        uint256 _toAmount = _afterBalance - _startBalance;

        IERC20(_toToken).approve(address(_toDex) ,_toAmount); 
        path[0] = _toToken;
        path[1] = _fromToken;
        /*uint256[] memory amounts_1 =*/ IUniswapV2Router(_toDex).swapExactTokensForTokens(_toAmount,0, path, address(this), block.timestamp);
        // Read _toToken balance after swap        
        //uint256 amountOut_1 = amounts_1[1];
    }

    // Allow the contract to receive Ether
    receive () external payable  {
    
    }    

    // Modifiers
    modifier onlyOwner() {
        require(msg.sender == OWNER, "caller is not the owner!");
        _;
    }

    // KEEP THIS FUNCTION IN CASE THE CONTRACT RECEIVES TOKENS!
    function withdrawToken(address _tokenAddress) public onlyOwner {
        uint256 balance = IERC20(_tokenAddress).balanceOf(address(this));
        IERC20(_tokenAddress).transfer(OWNER, balance);
    }

    // KEEP THIS FUNCTION IN CASE THE CONTRACT KEEPS LEFTOVER ETHER!
    function withdrawEther() public onlyOwner {
        address self = address(this); // workaround for a possible solidity bug
        uint256 balance = self.balance;
        // You need to mark the request.recipient as payable
        payable(address(OWNER)).transfer(balance);
    }
}