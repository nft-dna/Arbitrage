// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "./../Interfaces.sol";

contract MockDEX {
    mapping(address => mapping(address => uint256)) public prices; // Mock prices for token pairs
    mapping(address => mapping(address => uint24)) public fees;    // Mock fees for token pairs
    address payable OWNER;

    constructor() {
    OWNER = payable(msg.sender);
    }    

    modifier onlyOwner() {
        require(msg.sender == OWNER, "caller is not the owner!");
        _;
    }    

    // Allow the contract to receive Ether
    receive () external payable  {    
    }   

    // Set mock price for a token pair
    function setPrice(address tokenIn, address tokenOut, uint256 price) external onlyOwner {
        prices[tokenIn][tokenOut] = price;
    }

    // Set mock fee for a token pair
    function setFee(address tokenIn, address tokenOut, uint24 fee) external onlyOwner {
        fees[tokenIn][tokenOut] = fee;
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

        amounts = new uint[](path.length);
        amounts[0] = amountIn;
        
        for (uint i = 0; i < path.length - 1; i++) {
            address tokenIn = path[i];
            address tokenOut = path[i + 1];
            uint256 price = prices[tokenIn][tokenOut];
            require(price > 0, "Price not set");

            uint256 fee = fees[tokenIn][tokenOut];
            uint256 amountOut = (amounts[i] * price * (10000 - fee)) / 10000; // Apply price and fee
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

        amounts = new uint[](path.length);
        amounts[0] = msg.value;

        for (uint i = 0; i < path.length - 1; i++) {
            address tokenIn = path[i];
            address tokenOut = path[i + 1];
            uint256 price = prices[tokenIn][tokenOut];
            require(price > 0, "Price not set");

            uint256 fee = fees[tokenIn][tokenOut];
            uint256 amountOut = (amounts[i] * price * (10000 - fee)) / 10000; // Apply price and fee
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

        amounts = new uint[](path.length);
        amounts[path.length - 1] = amountOut;

        for (uint i = path.length - 1; i > 0; i--) {
            address tokenIn = path[i - 1];
            address tokenOut = path[i];
            uint256 price = prices[tokenIn][tokenOut];
            require(price > 0, "Price not set");

            uint256 fee = fees[tokenIn][tokenOut];
            uint256 amountIn = (amounts[i] * 10000) / (price * (10000 - fee)); // Apply price and fee
            require(msg.value >= amountIn, "Insufficient input amount");

            amounts[i - 1] = amountIn;

            // Transfer tokens
            IERC20(tokenOut).transfer(to, amountOut);
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

        uint256 price = prices[tokenIn][tokenOut];
        require(price > 0, "Price not set");

        uint256 amountOut = (amountIn * price * (10000 - fee)) / 10000; // Apply price and fee
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
        uint256 price = prices[tokenIn][tokenOut];
        require(price > 0, "Price not set");

        amountOut = (amountIn * price * (10000 - fee)) / 10000; // Apply price and fee
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
    
