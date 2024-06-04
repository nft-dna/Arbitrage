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

    struct routeChain {
        address router;
        address asset;
        uint24 poolFee;
    }		
	
	event DualDexTraded(address indexed trader, address indexed _fromToken, address indexed _toToken, address _fromDex, address _toDex, uint256 _fromAmount, uint256 _gainedAmount);
	event InstaTraded(address indexed trader, address indexed _baseAsset, routeChain[] _routeData, uint256 _fromAmount, uint256 _gainedAmount);

    address [] public tokens;
    address [] public stables;
    // Mapping to store pools by token pairs
    mapping(address => mapping(address => mapping(address => uint24))) public tokenV3PoolsFee;

    function sortTokens(address tokenA, address tokenB) internal pure returns (address token0, address token1) {
        require(tokenA != tokenB, 'IDENTICAL_ADDRESSES');
        (token0, token1) = tokenA < tokenB ? (tokenA, tokenB) : (tokenB, tokenA);
        require(token0 != address(0), 'ZERO_ADDRESS');
    }
	
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

    function AddTestV3PoolFee(address router, address token1, address token2, uint24 fee) external onlyOwner {
		(address tokenA, address tokenB) = sortTokens(token1, token2);
        tokenV3PoolsFee[router][tokenA][tokenB] = fee;
    }
	
    function getTestV3PoolFee(address router, address token1, address token2) internal view returns (uint24) {
		(address tokenA, address tokenB) = sortTokens(token1, token2);
        return tokenV3PoolsFee[router][tokenA][tokenB];
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

    function GetAmountOutMin(address router, uint24 poolfee, address _tokenIn, address _tokenOut, uint256 _amount) public view returns (uint256 ) {
        if (dexInterface[router] == DexInterfaceType.IUniswapV3Router || poolfee > 0) {
			require(NATIVE_TOKEN != _tokenIn, "Router does not support direct ETH swap");
			require(NATIVE_TOKEN != _tokenOut, "Router does not support direct ETH swap");
			if (poolfee >= 1000000) {
				poolfee = 0;
			}
            uint256 result = IQuoter(router).quoteExactInputSingle(_tokenIn, _tokenOut , poolfee, _amount, 0);
            return result;
        } else {
            uint256 result = 0;            
            address[] memory path;
            path = new address[](2);
            path[0] = _tokenIn;
            path[1] = _tokenOut;            
            try IUniswapV2Router(router).getAmountsOut(_amount, path) returns (uint256[] memory amountOutMins) {
                result = amountOutMins[path.length-1];
            } catch {
            }
            return result;      			
		//uint256[] memory amountOutMins = IUniswapV2Router(router).getAmountsOut(_amount, path);
		//return amountOutMins[path.length-1];         
        }
    }

    function EstimateDualDexTrade(address _fromToken, address _toToken, address _fromDex, uint24 _fromPoolFee, address _toDex, uint24 _toPoolFee, uint256 _fromAmount) external view returns (uint256) {
        uint256 amtBack1 = GetAmountOutMin(_fromDex, _fromPoolFee, _fromToken, _toToken, _fromAmount);
        uint256 amtBack2 = GetAmountOutMin(_toDex, _toPoolFee, _toToken, _fromToken, amtBack1);
        return amtBack2;
    }
  
    function DualDexTrade(address _fromToken, address _toToken, address _fromDex, uint24 _fromPoolFee, address _toDex, uint24 _toPoolFee, uint256 _fromAmount, uint deadlineDeltaSec) payable public {
        uint256 startBalance = 0;
		if (NATIVE_TOKEN == _fromToken) {
			if (msg.value >= 1000000) {
				depositEtherSucceded(msg.sender, msg.value);
			}
			startBalance = address(this).balance;
			require(startBalance >= _fromAmount, "Insufficient Ether balance");
        } else {
			startBalance = IERC20(_fromToken).balanceOf(address(this));
			require(startBalance >= _fromAmount, "Insufficient Token balance");		
        }
		uint256 tokenBalance = IERC20(_toToken).balanceOf(address(this));
		uint256 tradeableAmount = _tradeToken(_fromDex, _fromToken, _toToken, _fromPoolFee, _fromAmount, deadlineDeltaSec, tokenBalance);
		tradeableAmount = _tradeToken(_toDex, _toToken, _fromToken, _toPoolFee, tradeableAmount, deadlineDeltaSec, startBalance);
		require(tradeableAmount > 0, "End balance must exceed start balance.");
		
		if (NATIVE_TOKEN == _fromToken) {
			depositEtherSucceded(msg.sender, tradeableAmount);
		} else {
			depositTokenSucceded(msg.sender, _fromToken, tradeableAmount);
		}		
		emit DualDexTraded(msg.sender, _fromToken, _toToken, _fromDex, _toDex, _fromAmount, tradeableAmount);		
    }

    function _tradeToken(        
        address router,    
        address from,
        address to,
		uint24 poolFee,
        uint256 amount,
        uint deadlineDeltaSec,
        uint256 initialBalance
    ) internal returns (uint256 tradeableAmount) {
	
        _swapToken(router, poolFee, from, to, amount, deadlineDeltaSec);
		
		uint256 afterBalance = 0;
		if (NATIVE_TOKEN != to) {
			afterBalance = IERC20(to).balanceOf(address(this));
		} else {
			afterBalance = address(this).balance;
		}
		return afterBalance - initialBalance;
    }    
	
	
    function _swapToken(address router, uint24 _poolFee, address _tokenIn, address _tokenOut, uint256 _amount, uint deadlineDeltaSec) private {
		if (NATIVE_TOKEN != _tokenIn) {
			IERC20(_tokenIn).approve(router, _amount);
		}
        if (dexInterface[router] == DexInterfaceType.IUniswapV3Router || _poolFee > 0) {
			require(NATIVE_TOKEN != _tokenIn, "Router does not support direct ETH swap");
			require(NATIVE_TOKEN != _tokenOut, "Router does not support direct ETH swap");
			if (_poolFee == 1) {
				_poolFee = 0;
			}		
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
			uint deadline = block.timestamp + deadlineDeltaSec;            
			if (NATIVE_TOKEN == _tokenIn) {
				address[] memory path = new address[](1);
				path[0] = _tokenOut;			
				IUniswapV2Router(router).swapExactETHForTokens{value: _amount}(_amount, path, address(this), block.timestamp + deadlineDeltaSec);
			} else if (NATIVE_TOKEN == _tokenIn) {
				address[] memory path = new address[](1);
				path[0] = _tokenIn;			
				IUniswapV2Router(router).swapExactTokensForETH(_amount, 0, path, address(this), block.timestamp + deadlineDeltaSec);
			} else {
				address[] memory path;
				path = new address[](2);
				path[0] = _tokenIn;
				path[1] = _tokenOut;
				IUniswapV2Router(router).swapExactTokensForTokens(_amount, 0, path, address(this), deadline);    
			}
        }
    }    

	function InstaTradeTokens(routeChain[] calldata _routedata, uint256 _startAmount, uint deadlineDeltaSec) payable public {
		require ( _routedata.length > 1, "Invalid param");
        
        uint256 startBalance = 0;
		if (NATIVE_TOKEN == _routedata[0].asset) {
			if (msg.value > 0) {
				depositEtherSucceded(msg.sender, msg.value);
			}
			startBalance = address(this).balance;
			require(startBalance >= _startAmount, "Insufficient Ether balance");
        } else {
			startBalance = IERC20(_routedata[0].asset).balanceOf(address(this));
			require(startBalance >= _startAmount, "Insufficient Token balance");		
        }
	
		uint256 gainedAmount = _instaTradeTokens(_routedata, _startAmount, startBalance, deadlineDeltaSec);			
		if (NATIVE_TOKEN == _routedata[0].asset) {
			depositEtherSucceded(msg.sender, gainedAmount);
		} else {
			depositTokenSucceded(msg.sender, _routedata[0].asset, gainedAmount);
		}			
        emit InstaTraded(msg.sender, _routedata[0].asset, _routedata, _startAmount, gainedAmount);
    }  
	
    function _instaTradeTokens(routeChain[] calldata _routedata, uint256 _amount, uint256 _startBalance, uint deadlineDeltaSec) internal returns (uint256 gainedAmount) {
		uint256[] memory balance = new uint256[](_routedata.length);
		for (uint b=1; b < _routedata.length; b++) {
			if (NATIVE_TOKEN == _routedata[b].asset) {
				balance[b-1] = address(this).balance;
			} else {
				balance[b-1] = (IERC20(_routedata[b].asset).balanceOf(address(this)));
			}
		}
		balance[_routedata.length-1] = _startBalance;
		
		uint256 tradeableAmount = _amount;
		uint i = 0;
		for (i; i < _routedata.length-1; i++) {
			tradeableAmount = _tradeToken(_routedata[i].router, _routedata[i].asset, _routedata[i+1].asset, _routedata[i].poolFee, tradeableAmount, deadlineDeltaSec, balance[i]);
		}
		tradeableAmount = _tradeToken(_routedata[i].router, _routedata[i].asset, _routedata[0].asset, _routedata[i].poolFee, tradeableAmount, deadlineDeltaSec, balance[i]);
        require(tradeableAmount > 0, "Trade Reverted, No Profit Made");
		return tradeableAmount;
    }
		
    function AmountBack(
        address router,
        address baseAsset,
        uint256 amount,
        address token1,
        address token2,
        address token3
    ) internal view returns (uint256) {
        uint256 amtBack = GetAmountOutMin(router, getTestV3PoolFee(router, baseAsset, token1), baseAsset, token1, amount);
        amtBack = GetAmountOutMin(router, getTestV3PoolFee(router, token1, token2), token1, token2, amtBack);
        amtBack = GetAmountOutMin(router, getTestV3PoolFee(router, token2, token3), token2, token3, amtBack);
        amtBack = GetAmountOutMin(router, getTestV3PoolFee(router, token3, baseAsset), token3, baseAsset, amtBack);
        return amtBack;
    }

    // Base Asset > Altcoin > Stablecoin > Altcoin > Base Asset
    function CrossStableSearch(address/*[] calldata _routers*/_router, address _baseAsset, uint256 _amount) external view returns (uint256,address,address,address) {
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
