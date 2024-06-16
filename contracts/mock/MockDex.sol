// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "./../Interfaces.sol";

contract MockDEX {

    address payable OWNER;
	address NATIVE_TOKEN;

    constructor(address weth) {
    OWNER = payable(msg.sender);
	NATIVE_TOKEN = weth;
    }    

    modifier onlyOwner() {
        require(msg.sender == OWNER, "caller is not the owner!");
        _;
    }    
	
    // Allow the contract to receive Ether
    receive () external payable  {    
    }   	
	
    struct PairInfo {
        uint24 price; // Mock price (amount of tokenOut per 1 unit of tokenIn in percentage points out of 100)
						// i.e 100 = 1 tokenIn costs as 1 tokenOut
						// i.e  50 = 1 tokenIn costs as half tokenOut		
        uint24 fee;    // Mock fee in basis 100000 points (e.g., 300 for 0.3%)
    }

    function sortTokens(address tokenA, address tokenB) internal pure returns (address token0, address token1) {
        require(tokenA != tokenB, 'IDENTICAL_ADDRESSES');
        (token0, token1) = tokenA < tokenB ? (tokenA, tokenB) : (tokenB, tokenA);
        require(token0 != address(0), 'ZERO_ADDRESS');
    }
	
    mapping(address => mapping(address => PairInfo)) public pairs; // Stores price and fee info for token pairs

    // Set mock price and fee for a token pair
    function setPairInfo(address tokenIn, address tokenOut, uint24 price, uint24 fee) external {
		(address token0, address token1) = sortTokens(tokenIn, tokenOut);
        pairs[token0][token1] = PairInfo(price, fee);
    }
	
    function getAmountsOut(uint amountIn, address[] memory path) external view returns (uint[] memory amounts) {
		require(path.length >= 2, "Invalid path");
        amounts = new uint[](path.length);
        amounts[0] = amountIn;
        
        for (uint i = 0; i < path.length - 1; i++) {
            address tokenIn = path[i];
			require(i + 1 < path.length, "Invalid path");		
            address tokenOut = path[i + 1];
            (address token0, address token1) = sortTokens(tokenIn, tokenOut);
			PairInfo memory pinfo = pairs[token0][token1];
            require(pinfo.price > 0, "Price not set");

            uint256 amountOut = (((amounts[i] * pinfo.price) / 100) * (100000 - pinfo.fee)) / 100000; // Apply price and fee
 
            amounts[i + 1] = amountOut;

		}	
	}

    // Mock Uniswap V2 swap
    function swapExactTokensForTokens(
        uint amountIn,
        uint amountOutMin,
        address[] calldata path,
        address to,
        uint deadline
    ) external returns (uint[] memory amounts) {
        require(block.timestamp <= deadline, "Transaction expired");
		require(path.length >= 2, "Invalid path");

        amounts = new uint[](path.length);
        amounts[0] = amountIn;
        
        for (uint i = 0; i < path.length - 1; i++) {
            address tokenIn = path[i];
			require(i + 1 < path.length, "Invalid path");			
            address tokenOut = path[i + 1];
            (address token0, address token1) = sortTokens(tokenIn, tokenOut);
			PairInfo memory pinfo = pairs[token0][token1];
            require(pinfo.price > 0, "Price not set");

            uint256 amountOut = (((amounts[i] * pinfo.price) / 100) * (100000 - pinfo.fee)) / 100000; // Apply price and fee			
            amounts[i + 1] = amountOut;
        }
        // Transfer tokens
		require(amounts[path.length - 1] >= amountOutMin, "Insufficient output amount");
        IERC20(path[0]).transferFrom(msg.sender, address(this), amounts[0]);
        IERC20(path[path.length - 1]).transfer(to, amounts[path.length - 1]);		
    }
    // Mock Uniswap V2 swapExactETHForTokens
    function swapExactETHForTokens(
        uint amountOutMin,
        address[] calldata path,
        address to,
        uint deadline
    ) external payable returns (uint[] memory amounts) {
        require(block.timestamp <= deadline, "Transaction expired");
        require(path.length >= 2, "Invalid path");
		require(path[0] == NATIVE_TOKEN, "Invalid token");
		
        amounts = new uint[](path.length);
        amounts[0] = msg.value;

        for (uint i = 0; i < path.length - 1; i++) {
            address tokenIn = path[i];
            address tokenOut = path[i + 1];
            (address token0, address token1) = sortTokens(tokenIn, tokenOut);
			PairInfo memory pinfo = pairs[token0][token1];
            require(pinfo.price > 0, "Price not set");

            uint256 amountOut = (((amounts[i] * pinfo.price) / 100) * (100000 - pinfo.fee)) / 100000; // Apply price and fee	
            amounts[i + 1] = amountOut;
        }

		require(amounts[path.length - 1] >= amountOutMin, "Insufficient output amount");
        IERC20(path[path.length - 1]).transfer(to, amounts[path.length - 1]);		
    }

    // Mock Uniswap V2 swapExactTokensForETH
    function swapExactTokensForETH(
        uint amountIn,
		uint amountOutMin,
        address[] calldata path,
        address to,
        uint deadline
    ) external returns (uint[] memory amounts) {
        require(block.timestamp <= deadline, "Transaction expired");
        require(path.length >= 2, "Invalid path");
		require(path[path.length - 1] == NATIVE_TOKEN, "Invalid token");		

        amounts = new uint[](path.length);
        amounts[0] = amountIn;

        for (uint i = 0; i < path.length - 1; i++) {
            address tokenIn = path[i];
            address tokenOut = path[i + 1];
            (address token0, address token1) = sortTokens(tokenIn, tokenOut);
			PairInfo memory pinfo = pairs[token0][token1];
            require(pinfo.price > 0, "Price not set");

            uint256 amountOut = (((amounts[i] * pinfo.price) / 100) * (100000 - pinfo.fee)) / 100000; // Apply price and fee
            amounts[i + 1] = amountOut;
        }
		
		require(amounts[path.length - 1] >= amountOutMin, "Insufficient output amount");
		IERC20(path[0]).transferFrom(msg.sender, address(this), amounts[0]);
        payable(to).transfer(amounts[path.length - 1]);	
    }    	


    // Mock Uniswap V3 swap
	// IUniswapV3Router02 mock
    function exactInputSingle(IUniswapV3Router02.ExactInputSingleParams calldata params) external payable returns (uint256 amountOut) {	

		(address token0, address token1) = sortTokens(params.tokenIn, params.tokenOut);
		PairInfo memory pinfo = pairs[token0][token1];
		require(pinfo.price > 0, "Price not set");

        amountOut = (((params.amountIn * pinfo.price) / 100) * (100000 - params.fee)) / 100000; // Apply price and fee	
        require(amountOut >= params.amountOutMinimum, "Insufficient output amount");		

        // Transfer tokens
        IERC20(params.tokenIn).transferFrom(msg.sender, address(this), params.amountIn);
        IERC20(params.tokenOut).transfer(params.recipient, amountOut);
    }
	
	// IUniswapV3Router01 mock
    function exactInputSingle(IUniswapV3Router01.ExactInputSingleParams calldata params) external payable returns (uint256 amountOut) {	

		(address token0, address token1) = sortTokens(params.tokenIn, params.tokenOut);
		PairInfo memory pinfo = pairs[token0][token1];
		require(pinfo.price > 0, "Price not set");

        amountOut = (((params.amountIn * pinfo.price) / 100) * (100000 - /*pinfo.*/params.fee)) / 100000; // Apply price and fee	
        require(amountOut >= params.amountOutMinimum, "Insufficient output amount");		

        // Transfer tokens
        IERC20(params.tokenIn).transferFrom(msg.sender, address(this), params.amountIn);
        IERC20(params.tokenOut).transfer(params.recipient, amountOut);

    }	
	
    // Mock Uniswap V3 quoteExactInputSingle
    function quoteExactInputSingle(
        address tokenIn,
        address tokenOut,
        uint24 fee,
        uint256 amountIn,
		uint160 sqrtPriceLimitX96
    ) external view returns (uint256 amountOut) {
		(address token0, address token1) = sortTokens(tokenIn, tokenOut);
		PairInfo memory pinfo = pairs[token0][token1];
		require(pinfo.price > 0, "Price not set");

        amountOut = (((amountIn * pinfo.price) / 100) * (100000 - /*pinfo.*/fee)) / 100000; // Apply price and fee	
    }	
	
	// Mock Uniswap V4
    function swap(IUniswapV4PoolManager.PoolKey memory key, IUniswapV4PoolManager.SwapParams memory params, bytes calldata hookData) external returns (/*BalanceDelta*/int256) {
		PairInfo memory pinfo = pairs[key.currency0][key.currency1];
		require(pinfo.price > 0, "Price not set");

        int256 amountOut = int256((((uint256(params.amountSpecified) * pinfo.price) / 100) * (100000 - key.fee)) / 100000); // Apply price and fee		
		
        IERC20(params.zeroForOne ? key.currency0 : key.currency1).transferFrom(msg.sender, address(this), uint256(params.amountSpecified));
        IERC20(params.zeroForOne ? key.currency1 : key.currency0).transfer(msg.sender, uint256(amountOut));
		
		return amountOut;
	}
	

    // Allow the contract to receive tokens
    function depositToken(address token, uint256 amount) external onlyOwner {
        IERC20(token).transferFrom(msg.sender, address(this), amount);
    }

    // Allow the contract owner to withdraw tokens
    function withdrawToken(address token, uint256 amount) external onlyOwner {
        IERC20(token).transfer(msg.sender, amount);
    }
}
    
