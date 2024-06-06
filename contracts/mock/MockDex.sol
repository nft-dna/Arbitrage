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
            address tokenOut = path[i + 1];
            (address token0, address token1) = sortTokens(tokenIn, tokenOut);
			PairInfo memory pinfo = pairs[token0][token1];
            require(pinfo.price > 0, "Price not set");

            uint256 amountOut = (((amounts[i] * pinfo.price) / 100) * (100000 - pinfo.fee)) / 100000; // Apply price and fee
            require(amountOut >= amountOutMin, "Insufficient output amount");

            amounts[i + 1] = amountOut;

            // Transfer tokens
            IERC20(tokenIn).transferFrom(msg.sender, address(this), amounts[i]);
            IERC20(tokenOut).transfer(to, amountOut);
        }
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
            require(amountOut >= amountOutMin, "Insufficient output amount");

            amounts[i + 1] = amountOut;

            // Transfer tokens
            IERC20(tokenOut).transfer(to, amountOut);
        }

        // Refund excess ETH
        if (msg.value > amounts[path.length - 1]) {
            payable(msg.sender).transfer(msg.value - amounts[path.length - 1]);
        }
    }

    // Mock Uniswap V2 swapETHForExactTokens
    function swapETHForExactTokens(
        uint amountOut,
        address[] calldata path,
        address to,
        uint deadline
    ) external payable returns (uint[] memory amounts) {
        require(block.timestamp <= deadline, "Transaction expired");
        require(path.length >= 2, "Invalid path");
		require(path[path.length - 1] == NATIVE_TOKEN, "Invalid token");		

        amounts = new uint[](path.length);
        amounts[path.length - 1] = amountOut;

        for (uint i = path.length - 1; i > 0; i--) {
            address tokenIn = path[i - 1];
            address tokenOut = path[i];
            (address token0, address token1) = sortTokens(tokenIn, tokenOut);
			PairInfo memory pinfo = pairs[token0][token1];
            require(pinfo.price > 0, "Price not set");

            uint256 amountIn = (((amounts[i] * pinfo.price) / 100) * (100000 - pinfo.fee)) / 100000; // Apply price and fee
            require(msg.value >= amountIn, "Insufficient input amount");

            amounts[i - 1] = amountIn;

            // Transfer eth
            //IERC20(tokenOut).transfer(to, amountOut);
			payable(to).transfer(amountOut);
        }

        // Refund excess ETH
        if (msg.value > amounts[0]) {
            payable(msg.sender).transfer(msg.value - amounts[0]);
        }
    }    

    // Mock Uniswap V2 swapExactTokensForETH
    function swapExactTokensForETH(
        uint amountIn,
		uint amountOutMin,
        address[] calldata path,
        address to,
        uint deadline
    ) external payable returns (uint[] memory amounts) {
        require(block.timestamp <= deadline, "Transaction expired");
        require(path.length >= 2, "Invalid path");
		require(path[path.length - 1] == NATIVE_TOKEN, "Invalid token");		

        amounts = new uint[](path.length);
        amounts[0] = msg.value;

        for (uint i = path.length - 1; i > 0; i--) {
            address tokenIn = path[i - 1];
            address tokenOut = path[i];
            (address token0, address token1) = sortTokens(tokenIn, tokenOut);
			PairInfo memory pinfo = pairs[token0][token1];
            require(pinfo.price > 0, "Price not set");

            uint256 amountOut = (((amounts[i] * pinfo.price) / 100) * (100000 - pinfo.fee)) / 100000; // Apply price and fee			
            require(amountOut >= amountOutMin, "Insufficient output amount");

            amounts[i - 1] = amountIn;

            // Transfer eth
            //IERC20(tokenOut).transfer(to, amountOut);
			payable(to).transfer(amountOut);
        }

        // Refund excess ETH
        if (msg.value > amounts[0]) {
            payable(msg.sender).transfer(msg.value - amounts[0]);
        }
    }    	

    // Mock Uniswap V3 swap
    function exactInputSingle(
        uint256 amountIn,
        uint256 amountOutMin,
        address tokenIn,
        address tokenOut,
        uint24 fee,
        address recipient,
        uint256 deadline
    ) external returns (uint256 ) {
        require(block.timestamp <= deadline, "Transaction expired");

		(address token0, address token1) = sortTokens(tokenIn, tokenOut);
		PairInfo memory pinfo = pairs[token0][token1];
		require(pinfo.price > 0, "Price not set");

        uint256 amountOut = (((amountIn * pinfo.price) / 100) * (100000 - /*pinfo.*/fee)) / 100000; // Apply price and fee	
        require(amountOut >= amountOutMin, "Insufficient output amount");		

        // Transfer tokens
        IERC20(tokenIn).transferFrom(msg.sender, address(this), amountIn);
        IERC20(tokenOut).transfer(recipient, amountOut);

        return amountOut;
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

    // Allow the contract to receive tokens
    function depositToken(address token, uint256 amount) external onlyOwner {
        IERC20(token).transferFrom(msg.sender, address(this), amount);
    }

    // Allow the contract owner to withdraw tokens
    function withdrawToken(address token, uint256 amount) external onlyOwner {
        IERC20(token).transfer(msg.sender, amount);
    }
}
    
