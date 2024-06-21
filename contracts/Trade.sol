// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "./Interfaces.sol";
import "./Deposit.sol";

contract Trade is Deposit {

    // Addresses
    address payable OWNER;
    address payable NATIVE_TOKEN;
	
     event InstaTraded(address indexed trader, address indexed _baseAsset, routeChain[] _routeData, uint256 _fromAmount, uint256 _gainedAmount);

    constructor(address native_token) {
        OWNER = payable(msg.sender);
		NATIVE_TOKEN = payable(native_token);
    }
  
    function _tradeToken(routeChain calldata routeTo, address _tokenOut, uint256 amountIn, uint deadlineDeltaSec, uint256 initialBalance/*, bool checkProfit*/) internal returns (uint256 tradeableAmount) {
		address _tokenIn = (address(0x0) == routeTo.asset) ? NATIVE_TOKEN : routeTo.asset;
		if (address(0x0) == _tokenOut) {
			_tokenOut = NATIVE_TOKEN;
		} 
		require(_tokenOut != _tokenIn);
	
        _swapToken(routeTo, _tokenOut, amountIn, deadlineDeltaSec);
		
		uint256 afterBalance = IERC20(_tokenOut).balanceOf(address(this));
		//if (checkProfit) {
			require(afterBalance > initialBalance, "Trade Reverted, No Profit Made");		
		//} else if (afterBalance < initialBalance) {
		//	tokenBalances[_tokenOut][msg.sender] = tokenBalances[_tokenOut][msg.sender] - (initialBalance - afterBalance);
		//	return 0;
		//}	
		return afterBalance - initialBalance;
    }    
	

    function _swapToken( routeChain calldata routeTo, address _tokenOut, uint256 _amountIn, uint deadlineDeltaSec) private {
		address _tokenIn = (address(0x0) == routeTo.asset) ? NATIVE_TOKEN : routeTo.asset;
		if (address(0x0) == _tokenOut) {
			_tokenOut = NATIVE_TOKEN;
		} 

		IERC20(_tokenIn).approve(routeTo.router, _amountIn);
		uint deadline = block.timestamp + deadlineDeltaSec; 		
		if (routeTo.Itype == DexInterfaceType.IUniswapV4PoolManager) {	
			// routeTo.poolFee == 0x800000 // dynamic fee
			// experimental..
			IUniswapV4PoolManager.PoolKey memory pool = IUniswapV4PoolManager.PoolKey({
				currency0: /*Currency*/(_tokenIn < _tokenOut ? _tokenIn : _tokenOut),
				currency1: /*Currency*/(_tokenIn < _tokenOut ? _tokenOut : _tokenIn),
				fee: routeTo.poolFee,
				tickSpacing: routeTo.tickSpacing,
				hooks: /*IHooks*/(address(0))
			});			
			IUniswapV4PoolManager.SwapParams memory params = IUniswapV4PoolManager.SwapParams({
				zeroForOne: _tokenIn < _tokenOut,
				amountSpecified: int256(_amountIn),
				sqrtPriceLimitX96: _tokenIn < _tokenOut ? MIN_PRICE_LIMIT : MAX_PRICE_LIMIT // unlimited impact
			});
			bytes memory hookData = new bytes(0); // no hook data on the hookless pool
			//PoolSwapTest.TestSettings memory testSettings = PoolSwapTest.TestSettings({takeClaims: false, settleUsingBurn: false});
			//PoolSwapTest.TestSettings memory testSettings = PoolSwapTest.TestSettings({withdrawTokens: true, settleUsingTransfer: true});
			IUniswapV4PoolManager(routeTo.router).swap(pool, params, /*testSettings,*/ hookData);		
		} else if (routeTo.Itype == DexInterfaceType.IUniswapV3RouterQuoter02) {
		// V3
			//V3SwapRouter
			V3SwapRouter.ExactInputSingleParams memory params;
			params.tokenIn = _tokenIn;
			params.tokenOut = _tokenOut;
			params.fee = routeTo.poolFee;
			params.recipient = address(this);
			params.amountIn = _amountIn;
			params.amountOutMinimum = 0;
			params.sqrtPriceLimitX96 = 0; //MAX_PRICE_LIMIT;
            V3SwapRouter(routeTo.router).exactInputSingle(params);			
		} else if (routeTo.Itype == DexInterfaceType.IUniswapV3RouterQuoter01) {
		// V3
			//IUniswapV3Router
			ISwapRouter.ExactInputSingleParams memory params;
			params.tokenIn = _tokenIn;
			params.tokenOut = _tokenOut;
			params.fee = routeTo.poolFee;
			params.recipient = address(this);
			params.deadline = deadline;
			params.amountIn = _amountIn;
			params.amountOutMinimum = 0;
			params.sqrtPriceLimitX96 = 0; //MAX_PRICE_LIMIT;
            ISwapRouter(routeTo.router).exactInputSingle(params);	
        } else if (routeTo.Itype == DexInterfaceType.IUniswapV2Router) {
		// V2 
			address[] memory path;
			path = new address[](2);
			path[0] = _tokenIn;
			path[1] = _tokenOut;			
			IUniswapV2Router0102(routeTo.router).swapExactTokensForTokens(_amountIn, 0, path, address(this), deadline);    
        } else { //  if (routeTo.Itype == DexInterfaceType.IQuickswapV3RouterQuoter) {
			IQuickswapV3Router(routeTo.router).exactInputSingle(_tokenIn, _tokenOut, address(this), deadline, _amountIn, 0, 0);
		}
    }    

	function InstaTradeTokens(routeChain[] calldata _routedata, uint256 _startAmount, uint deadlineDeltaSec) public payable {
	//	InstaTradeTokensChecked(_routedata, _startAmount, deadlineDeltaSec, true);
	//}	
    //function InstaTradeTokensChecked(routeChain[] calldata _routedata, uint256 _startAmount, uint deadlineDeltaSec, bool checkProfit) public payable {
		require ( _routedata.length > 1, "Invalid param");
		address tokenIn = (address(0x0) == _routedata[0].asset) ? NATIVE_TOKEN : _routedata[0].asset;
		if (NATIVE_TOKEN == tokenIn) {
			if (msg.value > 0) {
				INativeToken(NATIVE_TOKEN).deposit{ value: msg.value }();
				depositTokenSucceded(msg.sender, NATIVE_TOKEN, msg.value);				
			}
        }
        uint256 startBalance = IERC20(tokenIn).balanceOf(address(this));
		require(startBalance >= _startAmount, "Insufficient Token balance");
		//if (!checkProfit) {
		//	require(_startAmount <= tokenBalances[tokenIn][msg.sender], "Insufficient Token balance");
		//}
		
		uint256 gainedAmount = _instaTradeTokens(_routedata, _startAmount, startBalance, deadlineDeltaSec/*, checkProfit*/);			
		depositTokenSucceded(msg.sender, tokenIn, gainedAmount);
		if ((address(0x0) == _routedata[0].asset) && (msg.value == _startAmount)) {
			uint amount = msg.value + gainedAmount;
			INativeToken(NATIVE_TOKEN).withdraw(amount);
			tokenBalances[NATIVE_TOKEN][msg.sender] = tokenBalances[NATIVE_TOKEN][msg.sender] - amount;
			payable(msg.sender).transfer(amount);
			emit WithdrawToken(NATIVE_TOKEN, msg.sender, amount);
		} 	
        emit InstaTraded(msg.sender, tokenIn, _routedata, _startAmount, gainedAmount);
    }  

    function InstaSwapTokens(routeChain calldata _routedata, uint256 _startAmount, address _tokenOut, uint deadlineDeltaSec) public payable {
		require (_tokenOut != _routedata.asset, "Invalid param");
		address tokenIn = (address(0x0) == _routedata.asset) ? NATIVE_TOKEN : _routedata.asset;
		address tokenOut = (address(0x0) == _tokenOut) ? NATIVE_TOKEN : _tokenOut;
		if (NATIVE_TOKEN == tokenIn) {
			if (msg.value > 0) {
				INativeToken(NATIVE_TOKEN).deposit{ value: msg.value }();
				depositTokenSucceded(msg.sender, NATIVE_TOKEN, msg.value);				
			}
        }
        uint256 startAmount = getTokenBalance(tokenIn, msg.sender);
		require(startAmount >= _startAmount, "Insufficient Token balance");

		uint256 startBalance = IERC20(tokenOut).balanceOf(address(this));	
		_swapToken( _routedata, tokenOut, _startAmount, deadlineDeltaSec);
		uint256 endBalance = IERC20(tokenOut).balanceOf(address(this));

		tokenBalances[tokenIn][msg.sender] = tokenBalances[tokenIn][msg.sender] - _startAmount;
		if (endBalance > startBalance) {
			uint256 amount = endBalance - startBalance;
			depositTokenSucceded(msg.sender, tokenOut, amount);
			if (address(0x0) == _tokenOut) {
				INativeToken(NATIVE_TOKEN).withdraw(amount);
				tokenBalances[NATIVE_TOKEN][msg.sender] = tokenBalances[NATIVE_TOKEN][msg.sender] - amount;
				payable(msg.sender).transfer(amount);
				emit WithdrawToken(NATIVE_TOKEN, msg.sender, amount);
			}
		}
    } 	
	
    function _instaTradeTokens(routeChain[] calldata _routedata, uint256 _amount, uint256 _startBalance, uint deadlineDeltaSec/*, bool checkProfit*/) internal returns (uint256 gainedAmount) {
		uint256[] memory balance = new uint256[](_routedata.length);
		for (uint b=1; b < _routedata.length; b++) {
			if (address(0x0) == _routedata[b].asset) {
				balance[b-1] = IERC20(NATIVE_TOKEN).balanceOf(address(this));
			} else {
				balance[b-1] = IERC20(_routedata[b].asset).balanceOf(address(this));
			}
		}
		balance[_routedata.length-1] = _startBalance;
		
		uint256 tradeableAmount = _amount;
		uint i = 0;
		for (i; i < _routedata.length-1; i++) {
			tradeableAmount = _tradeToken(_routedata[i], _routedata[i+1].asset, tradeableAmount, deadlineDeltaSec, balance[i]/*, checkProfit*/);
		}
		tradeableAmount = _tradeToken(_routedata[i], _routedata[0].asset, tradeableAmount, deadlineDeltaSec, balance[i]/*, checkProfit*/);
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

    /*
    // KEEP THIS FUNCTION IN CASE THE CONTRACT RECEIVES TOKENS!
    function safeWithdrawToken(address _tokenAddress, uint256 amount) public onlyOwner {
        uint256 balance = IERC20(_tokenAddress).balanceOf(address(this));
        require(amount <= balance, "Insufficient Token balance");
        IERC20(_tokenAddress).transfer(OWNER, amount);
    }
    */

    // KEEP THIS FUNCTION IN CASE THE CONTRACT RECEIVS OR KEEPS LEFTOVER ETHER!
    function safeWithdrawEther(uint256 amount) public onlyOwner {
        address self = address(this); // workaround for a possible solidity bug
        uint256 balance = self.balance;
        require(amount <= balance, "Insufficient Ether balance");
        // You need to mark the request.recipient as payable
        payable(address(OWNER)).transfer(amount);
    }
}
