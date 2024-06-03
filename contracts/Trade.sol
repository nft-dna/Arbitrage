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
    // Mapping to store pools by token pairs
    mapping(address => mapping(address => uint24)) public tokenV3PoolsFee;

    function AddTestTokens(address[] calldata _tokens) external onlyOwner {
        for (uint i=0; i<_tokens.length; i++) {
            tokens.push(_tokens[i]);
        }
    }

    function AddTestStables(address[] calldata _stables) external onlyOwner {
        for (uint i=0; i<_stables.length; i++) {
            stables.push(_stables[i]);
        }
    }    

    function AddTestV3PoolFee(address token1, address token2, uint24 fee) external onlyOwner {
        tokenV3PoolsFee[token1][token2] = fee;
        tokenV3PoolsFee[token2][token1] = fee;
    }      

    constructor() {
        OWNER = payable(msg.sender);
    }

    function IsDexAdded(address _dex) internal view returns (bool) {
        for (uint256 i = 0; i < dexAddresses.length; i++) {
            if (dexAddresses[i] == _dex) {
                return true;
            }
        }
        return false;
    }

    function AddDex(address[] calldata  _dex, DexInterfaceType[] calldata  _interface) public onlyOwner {
        require ( _dex.length == _interface.length, "Invalid param");
        for (uint i=0; i<_dex.length; i++) {
            if (!IsDexAdded(_dex[i])) {
                dexAddresses.push(_dex[i]);
            }
            dexInterface[_dex[i]] = _interface[i];
        }
    }

    function getDexCount() public view returns (uint256) {
        return dexAddresses.length;
    }

    function getAmountOutMin(address router, uint24 poolfee, address _tokenIn, address _tokenOut, uint256 _amount) public view returns (uint256 ) {
        if (dexInterface[router] == DexInterfaceType.IUniswapV3Router || poolfee > 0) {
            uint256 result = IQuoter(router).quoteExactInputSingle(_tokenIn, _tokenOut , poolfee, _amount, 0);
            return result;
        } else {
            uint256 result = 0;            
            address[] memory path;
            path = new address[](2);
            path[0] = _tokenIn;
            path[1] = _tokenOut;            
            try IUniswapV2Router(router).getAmountsOut(_amount, path) returns (uint256[] memory amountOutMins) {
                result = amountOutMins[path.length -1];
            } catch {
            }
            return result;            
        }
    }

    function EstimateDualDexTrade(address _fromToken, address _toToken, address _fromDex, uint24 _fromPoolFee, address _toDex, uint24 _toPoolFee, uint256 _fromAmount) external view returns (uint256) {
        uint256 amtBack1 = getAmountOutMin(_fromDex, _fromPoolFee, _fromToken, _toToken, _fromAmount);
        uint256 amtBack2 = getAmountOutMin(_toDex, _toPoolFee, _toToken, _fromToken, amtBack1);
        return amtBack2;
    }
  
    function DualDexTrade(address _fromToken, address _toToken, address _fromDex, uint24 _fromPoolFee, address _toDex, uint24 _toPoolFee, uint256 _fromAmount, uint deadlineDeltaSec) /*payable*/ public {
        if (NATIVE_TOKEN == _fromToken) {
            _arbEther(_toToken, _fromDex, _toDex, _fromAmount, deadlineDeltaSec);
        } else {
            _arbToken(_fromToken, _toToken, _fromDex, _fromPoolFee, _toDex, _toPoolFee, _fromAmount, deadlineDeltaSec);
        }
    }

    function _arbEther(address _token, address _fromDex, address _toDex, uint256 _fromAmount, uint deadlineDeltaSec) /*payable*/ internal {
        require(_fromAmount > 0, "Amount must be greater than 0");
        require(dexInterface[_fromDex] == DexInterfaceType.IUniswapV2Router, "fromDex does not support direct ETH swap");
        require(dexInterface[_toDex] == DexInterfaceType.IUniswapV2Router, "toDex does not support direct ETH swap");
        
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

    function _arbToken(address _fromToken, address _toToken, address _fromDex, uint24 _fromPoolFee, address _toDex, uint24 _toPoolFee, uint256 _fromAmount, uint deadlineDeltaSec) /*payable*/ internal {
        require(_fromAmount > 0, "Amount must be greater than 0");
        
        // Track original balance
        uint256 _startBalance = IERC20(_fromToken).balanceOf(address(this));
        require(_startBalance >= _fromAmount, "Insufficient token balance");
        
        // Perform the arb trade
        _tradeToken(_fromToken, _toToken, _fromDex, _fromPoolFee, _toDex, _toPoolFee, _fromAmount, deadlineDeltaSec, _startBalance);

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

    function _tradeToken(address _fromToken, address _toToken, address _fromDex, uint24 _fromPoolFee, address _toDex, uint24 _toPoolFee, uint256 _fromAmount, uint deadlineDeltaSec, uint256 _startBalance) internal {
        //require ( dexInterface[_fromDex] > DexInterfaceType.Unknown, "Unsupported from dex");
        //require ( dexInterface[_toDex] > DexInterfaceType.Unknown, "Unsupported to dex");

        // Track the balance of the token RECEIVED from the trade
        //uint256 _startBalance = IERC20(_toToken).balanceOf(address(this));

        _swapToken(_fromDex, _fromPoolFee, _fromToken, _toToken, _fromAmount, deadlineDeltaSec);
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

        _swapToken(_toDex, _toPoolFee, _toToken, _fromToken, _toAmount, deadlineDeltaSec);
        /*
        IERC20(_toToken).approve(address(_toDex) ,_toAmount); 
        path[0] = _toToken;
        path[1] = _fromToken;
        /uint256[] memory amounts_1 =/ IUniswapV2Router(_toDex).swapExactTokensForTokens(_toAmount,0, path, address(this), block.timestamp);
        // Read _toToken balance after swap        
        //uint256 amountOut_1 = amounts_1[1];
        */
    }

    function _swapToken(address router, uint24 _poolFee, address _tokenIn, address _tokenOut, uint256 _amount, uint deadlineDeltaSec) private {
        IERC20(_tokenIn).approve(router, _amount);
        if (dexInterface[router] == DexInterfaceType.IUniswapV3Router || _poolFee > 0) {
            bytes memory params = abi.encode(
                _tokenIn,
                _tokenOut,
                _poolFee,
                address(this),
                _amount,
                0,
                0
            );           
            IUniswapV3Router(router).exactInputSingle(params);
        } else {
            address[] memory path;
            path = new address[](2);
            path[0] = _tokenIn;
            path[1] = _tokenOut;
            uint deadline = block.timestamp + deadlineDeltaSec;            
            IUniswapV2Router(router).swapExactTokensForTokens(_amount, 0, path, address(this), deadline);    
        }
    }    

    function AmountBack(
        address router,
        address baseAsset,
        uint256 amount,
        address token1,
        address stable,
        address token3
    ) internal view returns (uint256) {
        uint256 amtBack = getAmountOutMin(router, tokenV3PoolsFee[baseAsset][token1], baseAsset, token1, amount);
        amtBack = getAmountOutMin(router, tokenV3PoolsFee[token1][stable], token1, stable, amtBack);
        amtBack = getAmountOutMin(router, tokenV3PoolsFee[stable][token3], stable, token3, amtBack);
        amtBack = getAmountOutMin(router, tokenV3PoolsFee[token3][baseAsset], token3, baseAsset, amtBack);
        return amtBack;
    }

    /*
    Base Asset > Altcoin > Stablecoin > Altcoin > Base Asset
    */ 
    function InstaSearch(address/*[] calldata _routers*/_router, address _baseAsset, uint256 _amount) external view returns (uint256,address,address,address) {
        uint256 maxAmtBack = 0;
        address token1;
        address token2;
        address token3;
        //for (uint i0=0; i0<dexAddresses.length; i0++) {
            for (uint i1=0; i1<tokens.length; i1++) {
				if (_baseAsset != tokens[i1]) {
					for (uint i2=0; i2<stables.length; i2++) {
						for (uint i3=0; i3<tokens.length; i3++) {
							if (_baseAsset != tokens[i3]) {
								uint256 amtBack = AmountBack(_router, _baseAsset, _amount, tokens[i1], stables[i2], tokens[i3]);
								if (amtBack > _amount && amtBack > maxAmtBack) {
									maxAmtBack = amtBack;
									token1 = tokens[i1];
									token2 = tokens[i2];
									token3 = tokens[i3];
								}
							}
						}
					}
				}
            }
        //}
        return (maxAmtBack,token1,token2,token3);
    } 

    function InstaTradeTokens(address _router1, address _baseAsset, address _token2, address _token3, address _token4, uint256 _amount, uint deadlineDeltaSec) external {
        uint256 _startBalance = IERC20(_baseAsset).balanceOf(address(this));
        require(_startBalance > 0, "StartBalance must be greater than 0");
        _instaTradeTokens(_router1, _baseAsset, _token2, _token3, _token4, _amount, _startBalance, deadlineDeltaSec);
        uint256 _endBalance = IERC20(_baseAsset).balanceOf(address(this));
        depositTokenSucceded(msg.sender, _baseAsset, _endBalance-_startBalance);
        emit InstaTraded(msg.sender, _baseAsset, _token2, _token3, _token4, _router1, _amount, _endBalance-_startBalance);
    }  

    function _instaTradeTokens(address _router, address _baseAsset, address _token2, address _token3, address _token4, uint256 _amount, uint256 _startBalance, uint deadlineDeltaSec) internal {        
        uint256 token3InitialBalance = IERC20(_token3).balanceOf(address(this));
        uint256 token4InitialBalance = IERC20(_token4).balanceOf(address(this));
        uint256 tradeableAmount2 = _instaTradeSwap(_router, _baseAsset, _token2, _amount, deadlineDeltaSec, IERC20(_token2).balanceOf(address(this)));
        uint256 tradeableAmount3 = _instaTradeSwap(_router, _token2, _token3, tradeableAmount2, deadlineDeltaSec, token3InitialBalance);
        uint256 tradeableAmount4 = _instaTradeSwap(_router, _token3, _token4, tradeableAmount3, deadlineDeltaSec, token4InitialBalance);
        require(_instaTradeSwap(_router, _token4, _baseAsset, tradeableAmount4, deadlineDeltaSec, _startBalance) > 0, "Trade Reverted, No Profit Made");
    }    

    function _instaTradeSwap(        
        address router,    
        address from,
        address to,
        uint256 amount,
        uint deadlineDeltaSec,
        uint256 initialBalance
    ) internal returns (uint256 tradeableAmount) {
        _swapToken(router, tokenV3PoolsFee[from][to], from, to, amount, deadlineDeltaSec);
        return IERC20(to).balanceOf(address(this)) - initialBalance;
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
    function safeWithdrawToken(address _tokenAddress, uint256 amount) public onlyOwner {
        uint256 balance = IERC20(_tokenAddress).balanceOf(address(this));
        require(amount <= balance, "Insufficient Token balance");
        IERC20(_tokenAddress).transfer(OWNER, amount);
    }

    // KEEP THIS FUNCTION IN CASE THE CONTRACT KEEPS LEFTOVER ETHER!
    function safeWithdrawEther(uint256 amount) public onlyOwner {
        address self = address(this); // workaround for a possible solidity bug
        uint256 balance = self.balance;
        require(amount <= balance, "Insufficient Ether balance");
        // You need to mark the request.recipient as payable
        payable(address(OWNER)).transfer(amount);
    }
}
