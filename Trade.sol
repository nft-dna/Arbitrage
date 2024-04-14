// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "./Interfaces.sol";
import "./Deposit.sol";

contract Trade is Deposit {

    // Addresses
    address payable OWNER;
    address NATIVE_TOKEN = address(0x0);

    address[] public dexAddresses; // Array to store dex addresses
    mapping(address => DexInterfaceType) public dexInterface;    

    event DualDexEtherTraded(address indexed trader, address indexed _toToken, address _fromDex, address _toDex, uint256 _fromAmount, uint256 _gainedAmount);
    event DualDexTokenTraded(address indexed trader, address indexed _fromToken, address indexed _toToken, address _fromDex, address _toDex, uint256 _fromAmount, uint256 _gainedAmount);
    event InstaTraded(address indexed trader, address indexed _baseAsset, address _token2, address _token3, address _token4, address _fromDex, uint256 _fromAmount, uint256 _gainedAmount);

    address [] public tokens;
    address [] public stables;

    function addTestTokens(address[] calldata _tokens) external onlyOwner {
        for (uint i=0; i<_tokens.length; i++) {
            tokens.push(_tokens[i]);
        }
    }

    function addTestStables(address[] calldata _stables) external onlyOwner {
        for (uint i=0; i<_stables.length; i++) {
            stables.push(_stables[i]);
        }
    }    

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

    function SetDex(address[] calldata  _dex, DexInterfaceType[] calldata  _interface) public onlyOwner {
        require ( _dex.length == _interface.length, "Invalid param");
        for (uint i=0; i<_dex.length; i++) {
            if (!isDexAdded(_dex[i])) {
                dexAddresses.push(_dex[i]);
            }
            dexInterface[_dex[i]] = _interface[i];
        }
    }

    function getDexCount() public view returns (uint256) {
        return dexAddresses.length;
    }

    function getAmountOutMin(address router, address _tokenIn, address _tokenOut, uint256 _amount) public view returns (uint256 ) {
        address[] memory path;
        path = new address[](2);
        path[0] = _tokenIn;
        path[1] = _tokenOut;
        uint256 result = 0;
        try IUniswapV2Router(router).getAmountsOut(_amount, path) returns (uint256[] memory amountOutMins) {
            result = amountOutMins[path.length -1];
        } catch {
        }
        return result;
    }

    function estimateDualDexTrade(address _fromToken, address _toToken, address _fromDex, address _toDex, uint256 _fromAmount) external view returns (uint256) {
        uint256 amtBack1 = getAmountOutMin(_fromDex, _fromToken, _toToken, _fromAmount);
        uint256 amtBack2 = getAmountOutMin(_toDex, _toToken, _fromToken, amtBack1);
        return amtBack2;
    }
  
    function DualDexTrade(address _fromToken, address _toToken, address _fromDex, address _toDex, uint256 _fromAmount, uint deadlineDeltaSec) /*payable*/ public {
        if (NATIVE_TOKEN == _fromToken) {
            _arbEther(_toToken, _fromDex, _toDex, _fromAmount, deadlineDeltaSec);
        } else {
            _arbToken(_fromToken, _toToken, _fromDex, _toDex, _fromAmount, deadlineDeltaSec);
        }
    }

    function _arbEther(address _token, address _fromDex, address _toDex, uint256 _fromAmount, uint deadlineDeltaSec) /*payable*/ internal {
        require(_fromAmount > 0, "Amount must be greater than 0");
        
        // Track original balance
        uint256 _startBalance = address(this).balance;
        require(_startBalance >= _fromAmount, "Insufficient token balance");
        
        // Perform the arb trade
        _tradeEther( _token, _fromDex, _toDex, _fromAmount, deadlineDeltaSec, _startBalance);

        // Track result balance
        uint256 _endBalance = address(this).balance;

        // Require that arbitrage is profitable
        require(_endBalance > _startBalance, "End balance must exceed start balance.");

        depositEtherSucceded(msg.sender, _endBalance-_startBalance);
        emit DualDexEtherTraded(msg.sender, _token, _fromDex, _toDex, _fromAmount, _endBalance-_startBalance);
    }

    function _arbToken(address _fromToken, address _toToken, address _fromDex, address _toDex, uint256 _fromAmount, uint deadlineDeltaSec) /*payable*/ internal {
        require(_fromAmount > 0, "Amount must be greater than 0");
        
        // Track original balance
        uint256 _startBalance = IERC20(_fromToken).balanceOf(address(this));
        require(_startBalance >= _fromAmount, "Insufficient token balance");
        
        // Perform the arb trade
        _tradeToken(_fromToken, _toToken, _fromDex, _toDex, _fromAmount, deadlineDeltaSec, _startBalance);

        // Track result balance
        uint256 _endBalance = IERC20(_fromToken).balanceOf(address(this));

        // Require that arbitrage is profitable
        require(_endBalance > _startBalance, "End balance must exceed start balance.");

        depositTokenSucceded(msg.sender, _fromToken, _endBalance-_startBalance);
        emit DualDexTokenTraded(msg.sender, _fromToken, _toToken, _fromDex, _toDex, _fromAmount, _endBalance-_startBalance);
    }

    function _tradeEther(address _token, address _fromDex, address _toDex, uint256 _fromAmount, uint deadlineDeltaSec, uint256 _startBalance) internal {
        //require ( dexInterface[_fromDex] > DexInterfaceType.Unknown, "Unsupported from dex");
        //require ( dexInterface[_toDex] > DexInterfaceType.Unknown, "Unsupported to dex");

        // Track the balance of the token RECEIVED from the trade
        //uint256 _startBalance = address(this).balance;

        address[] memory path = new address[](1);
        path[0] = _token;
        //path[1] = _toToken;
        /*uint256[] memory amounts_0 =*/ IUniswapV2Router(_fromDex).swapExactETHForTokens{value: _fromAmount}(_fromAmount, path, address(this), block.timestamp + deadlineDeltaSec);
        // Get the amount of tokens received        
        //uint256 amountOut = amounts_0[1];
        // Calculate the how much of the token we received
        uint256 _afterBalance = address(this).balance;        
        uint256 _toAmount = _afterBalance - _startBalance;

        IERC20(_token).approve(address(_toDex), _toAmount); 
        path[0] = _token;
        //path[1] = _fromToken;
        /*uint256[] memory amounts_1 =*/ IUniswapV2Router(_toDex).swapExactTokensForETH(_toAmount,0, path, address(this), block.timestamp + deadlineDeltaSec);
        // Read _toToken balance after swap        
        //uint256 amountOut_1 = amounts_1[1];
    }    

    function _tradeToken(address _fromToken, address _toToken, address _fromDex, address _toDex, uint256 _fromAmount, uint deadlineDeltaSec, uint256 _startBalance) internal {
        //require ( dexInterface[_fromDex] > DexInterfaceType.Unknown, "Unsupported from dex");
        //require ( dexInterface[_toDex] > DexInterfaceType.Unknown, "Unsupported to dex");

        // Track the balance of the token RECEIVED from the trade
        //uint256 _startBalance = IERC20(_toToken).balanceOf(address(this));

        _swapToken(_fromDex, _fromToken, _toToken, _fromAmount, deadlineDeltaSec);
        /*
        IERC20(_fromToken).approve(address(_fromDex) , _fromAmount); 
        address[] memory path = new address[](2);
        path[0] = _fromToken;
        path[1] = _toToken;
        /uint256[] memory amounts_0 =/ IUniswapV2Router(_fromDex).swapExactTokensForTokens(_fromAmount, 0, path, address(this), block.timestamp);
        // Get the amount of tokens received        
        //uint256 amountOut = amounts_0[1];
        // Calculate the how much of the token we received
        */
        uint256 _afterBalance = IERC20(_toToken).balanceOf(address(this));        
        uint256 _toAmount = _afterBalance - _startBalance;

        _swapToken(_toDex, _toToken, _fromToken, _toAmount, deadlineDeltaSec);
        /*
        IERC20(_toToken).approve(address(_toDex) ,_toAmount); 
        path[0] = _toToken;
        path[1] = _fromToken;
        /uint256[] memory amounts_1 =/ IUniswapV2Router(_toDex).swapExactTokensForTokens(_toAmount,0, path, address(this), block.timestamp);
        // Read _toToken balance after swap        
        //uint256 amountOut_1 = amounts_1[1];
        */
    }

    function _swapToken(address router, address _tokenIn, address _tokenOut, uint256 _amount, uint deadlineDeltaSec) private {
        IERC20(_tokenIn).approve(router, _amount);
        address[] memory path;
        path = new address[](2);
        path[0] = _tokenIn;
        path[1] = _tokenOut;
        uint deadline = block.timestamp + deadlineDeltaSec;
        IUniswapV2Router(router).swapExactTokensForTokens(_amount, 1, path, address(this), deadline);
    }    

    /*
    Base Asset > Altcoin > Stablecoin > Altcoin > Base Asset
    */ 
    function InstaSearch(address/*[] calldata _routers*/_router, address _baseAsset, uint256 _amount) external view returns (uint256,address,address,address) {
        uint256 amtBack;
        address token1;
        address token2;
        address token3;
        //for (uint i0=0; i0<_routers.length; i0++) {
            for (uint i1=0; i1<tokens.length; i1++) {
                for (uint i2=0; i2<stables.length; i2++) {
                    for (uint i3=0; i3<tokens.length; i3++) {
                        amtBack = getAmountOutMin(_router, _baseAsset, tokens[i1], _amount);
                        amtBack = getAmountOutMin(_router, tokens[i1], stables[i2], amtBack);
                        amtBack = getAmountOutMin(_router, stables[i2], tokens[i3], amtBack);
                        amtBack = getAmountOutMin(_router, tokens[i3], _baseAsset, amtBack);
                        if (amtBack > _amount) {
                        token1 = tokens[i1];
                        token2 = tokens[i2];
                        token3 = tokens[i3];
                        break;
                        }
                    }
                }
            }
        //}
        return (amtBack,token1,token2,token3);
    }   

    function InstaTrade(address _router1, address _baseAsset, address _token2, address _token3, address _token4, uint256 _amount, uint deadlineDeltaSec) external {
        uint256 _startBalance = IERC20(_baseAsset).balanceOf(address(this));
        require(_startBalance > 0, "StartBalance must be greater than 0");
        uint256 token2InitialBalance = IERC20(_token2).balanceOf(address(this));
        uint256 token3InitialBalance = IERC20(_token3).balanceOf(address(this));
        uint256 token4InitialBalance = IERC20(_token4).balanceOf(address(this));
        _swapToken(_router1,_baseAsset, _token2, _amount, deadlineDeltaSec);
        uint256 tradeableAmount2 = IERC20(_token2).balanceOf(address(this)) - token2InitialBalance;
        _swapToken(_router1,_token2, _token3, tradeableAmount2, deadlineDeltaSec);
        uint256 tradeableAmount3 = IERC20(_token3).balanceOf(address(this)) - token3InitialBalance;
        _swapToken(_router1,_token3, _token4, tradeableAmount3, deadlineDeltaSec);
        uint256 tradeableAmount4 = IERC20(_token4).balanceOf(address(this)) - token4InitialBalance;
        _swapToken(_router1,_token4, _baseAsset, tradeableAmount4, deadlineDeltaSec);
        uint256 _endBalance = IERC20(_baseAsset).balanceOf(address(this));
        require(_endBalance > _startBalance, "Trade Reverted, No Profit Made");
        depositTokenSucceded(msg.sender, _baseAsset, _endBalance-_startBalance);
        //emit InstaTraded(msg.sender, _baseAsset, _token2, _token3, _token4, _router1, _amount, _endBalance-_startBalance);
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
    function safeWithdrawToken(address _tokenAddress) public onlyOwner {
        uint256 balance = IERC20(_tokenAddress).balanceOf(address(this));
        IERC20(_tokenAddress).transfer(OWNER, balance);
    }

    // KEEP THIS FUNCTION IN CASE THE CONTRACT KEEPS LEFTOVER ETHER!
    function safeWithdrawEther() public onlyOwner {
        address self = address(this); // workaround for a possible solidity bug
        uint256 balance = self.balance;
        // You need to mark the request.recipient as payable
        payable(address(OWNER)).transfer(balance);
    }
}
