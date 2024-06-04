// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "./Interfaces.sol";
import "./Trade.sol";

contract Trader is Trade {

    address[] public dexAddresses; // Array to store dex addresses
    mapping(address => DexInterfaceType) public dexInterface;    

    address [] public tokens;
    address [] public stables;
    // Mapping to store pools by token pairs
    mapping(address => mapping(address => mapping(address => uint24))) public tokenV3PoolsFee;
	
    function sortTokens(address tokenA, address tokenB) internal pure returns (address token0, address token1) {
        require(tokenA != tokenB, 'IDENTICAL_ADDRESSES');
        (token0, token1) = tokenA < tokenB ? (tokenA, tokenB) : (tokenB, tokenA);
        //require(token0 != address(0), 'ZERO_ADDRESS');
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

    function GetAmountOutMin(address _router, uint24 _poolfee, address _tokenIn, address _tokenOut, uint256 _amount) public view returns (uint256 ) {
        if (dexInterface[_router] == DexInterfaceType.IUniswapV3Router || _poolfee > 0) {
			require(NATIVE_TOKEN != _tokenIn, "Router does not support direct ETH swap");
			require(NATIVE_TOKEN != _tokenOut, "Router does not support direct ETH swap");
			if (_poolfee >= 100000) {
				_poolfee = 0;
			}
            uint256 result = IQuoter(_router).quoteExactInputSingle(_tokenIn, _tokenOut , _poolfee, _amount, 0);
            return result;
        } else {
            //uint256 result = 0;            
            address[] memory path;
            path = new address[](2);
            path[0] = _tokenIn;
            path[1] = _tokenOut;            
            //try IUniswapV2Router(_router).getAmountsOut(_amount, path) returns (uint256[] memory amountOutMins) {
            //    result = amountOutMins[path.length-1];
            //} catch {
            //}
            //return result;      			
			uint256[] memory amountOutMins = IUniswapV2Router(_router).getAmountsOut(_amount, path);
			return amountOutMins[path.length-1];      
        }
    }

    function EstimateDualDexTrade(address _fromToken, address _toToken, address _fromDex, uint24 _fromPoolFee, address _toDex, uint24 _toPoolFee, uint256 _fromAmount) external view returns (uint256) {
        uint256 amtBack1 = GetAmountOutMin(_fromDex, _fromPoolFee, _fromToken, _toToken, _fromAmount);
        uint256 amtBack2 = GetAmountOutMin(_toDex, _toPoolFee, _toToken, _fromToken, amtBack1);
        return amtBack2;
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
}
