// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "./Interfaces.sol";
import "./Deposit.sol";

contract Trade is Deposit {

    // Addresses
    address payable OWNER;
    address NATIVE_TOKEN;
	
     event InstaTraded(address indexed trader, address indexed _baseAsset, routeChain[] _routeData, uint256 _fromAmount, uint256 _gainedAmount);

    constructor(address native_token) {
        OWNER = payable(msg.sender);
		NATIVE_TOKEN = native_token;
    }
  
    function _tradeToken(         
		routeChain calldata routeTo,
        address assetTo,		
        uint256 amountIn,
        uint deadlineDeltaSec,
        uint256 initialBalance
    ) internal returns (uint256 tradeableAmount) {
		require(assetTo != routeTo.asset);
	
        _swapToken(routeTo, assetTo, amountIn, deadlineDeltaSec);
		
		uint256 afterBalance = 0;
		if (address(0x0) != assetTo) {
			afterBalance = IERC20(assetTo).balanceOf(address(this));
		} else {
			afterBalance = address(this).balance;
		}
		require(afterBalance > initialBalance, "Trade Reverted, No Profit Made");		
		return afterBalance - initialBalance;
    }    
	
	
    function _swapToken( routeChain calldata routeTo, address _tokenOut, uint256 _amountIn, uint deadlineDeltaSec) private {
		address _tokenIn = routeTo.asset;	
		if (address(0x0) != _tokenIn) {
			IERC20(_tokenIn).approve(routeTo.router, _amountIn);
		}
		if (routeTo.Itype == DexInterfaceType.IUniswapV4PoolManager) {
			require(address(0x0) != _tokenIn, "Direct ETH swap, not implemented here yet");
			require(address(0x0) != _tokenOut, "Direct ETH swap, not implemented here yet");			
			// routeTo.poolFee == 0x800000 // dynamic fee
			// experimental..
			IPoolManager.PoolKey memory pool = IPoolManager.PoolKey({
				currency0: /*Currency*/(_tokenIn < _tokenOut ? _tokenIn : _tokenOut),
				currency1: /*Currency*/(_tokenIn < _tokenOut ? _tokenOut : _tokenIn),
				fee: routeTo.poolFee,
				tickSpacing: routeTo.tickSpacing,
				hooks: /*IHooks*/(address(0))
			});			
			IPoolManager.SwapParams memory params = IPoolManager.SwapParams({
				zeroForOne: _tokenIn < _tokenOut,
				amountSpecified: int256(_amountIn),
				sqrtPriceLimitX96: _tokenIn < _tokenOut ? MIN_PRICE_LIMIT : MAX_PRICE_LIMIT // unlimited impact
			});
			bytes memory hookData = new bytes(0); // no hook data on the hookless pool
			//PoolSwapTest.TestSettings memory testSettings = PoolSwapTest.TestSettings({takeClaims: false, settleUsingBurn: false});
			//PoolSwapTest.TestSettings memory testSettings = PoolSwapTest.TestSettings({withdrawTokens: true, settleUsingTransfer: true});
			IPoolManager(routeTo.router).swap(pool, params, /*testSettings,*/ hookData);		
		} else if (routeTo.Itype == DexInterfaceType.IUniswapV3Router) {
		// V3
			require(address(0x0) != _tokenIn, "Router does not support direct ETH swap");
			require(address(0x0) != _tokenOut, "Router does not support direct ETH swap");
			ExactInputSingleParams memory params;
			params.tokenIn = _tokenIn;
			params.tokenOut = _tokenOut;
			params.fee = routeTo.poolFee;
			params.recipient = address(this);
			params.amountIn = _amountIn;
			params.amountOutMinimum = 0;
			params.sqrtPriceLimitX96 = MAX_PRICE_LIMIT;
            IUniswapV3Router(routeTo.router).exactInputSingle(params);
        } else { // DexInterfaceType.IUniswapV2Router
		// V2
			uint deadline = block.timestamp + deadlineDeltaSec;  
			address[] memory path;
			path = new address[](2);
			path[0] = _tokenIn;
			path[1] = _tokenOut;			
			if (address(0x0) == _tokenIn) {
				path[0] = NATIVE_TOKEN;			
				IUniswapV2Router(routeTo.router).swapExactETHForTokens{value: _amountIn}(0, path, address(this), deadline);
			} else if (address(0x0) == _tokenOut) {
				path[1] = NATIVE_TOKEN;			
				IUniswapV2Router(routeTo.router).swapExactTokensForETH(_amountIn, 0, path, address(this), deadline);
			} else {

				IUniswapV2Router(routeTo.router).swapExactTokensForTokens(_amountIn, 0, path, address(this), deadline);    
			}
        }
    }    

    function InstaTradeTokens(routeChain[] calldata _routedata, uint256 _startAmount, uint deadlineDeltaSec) payable public {
		require ( _routedata.length > 1, "Invalid param");
        
        uint256 startBalance = 0;
		if (address(0x0) == _routedata[0].asset) {
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
		if (address(0x0) == _routedata[0].asset) {
			depositEtherSucceded(msg.sender, gainedAmount);
		} else {
			depositTokenSucceded(msg.sender, _routedata[0].asset, gainedAmount);
		}			
        emit InstaTraded(msg.sender, _routedata[0].asset, _routedata, _startAmount, gainedAmount);
    }  
	
    function _instaTradeTokens(routeChain[] calldata _routedata, uint256 _amount, uint256 _startBalance, uint deadlineDeltaSec) internal returns (uint256 gainedAmount) {
		uint256[] memory balance = new uint256[](_routedata.length);
		for (uint b=1; b < _routedata.length; b++) {
			if (address(0x0) == _routedata[b].asset) {
				balance[b-1] = address(this).balance;
			} else {
				balance[b-1] = (IERC20(_routedata[b].asset).balanceOf(address(this)));
			}
		}
		balance[_routedata.length-1] = _startBalance;
		
		uint256 tradeableAmount = _amount;
		uint i = 0;
		for (i; i < _routedata.length-1; i++) {
			tradeableAmount = _tradeToken(_routedata[i], _routedata[i+1].asset, tradeableAmount, deadlineDeltaSec, balance[i]);
		}
		tradeableAmount = _tradeToken(_routedata[i], _routedata[0].asset, tradeableAmount, deadlineDeltaSec, balance[i]);
		return tradeableAmount;
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
