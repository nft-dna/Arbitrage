// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;


interface IERC20 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address recipient, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

enum DexInterfaceType {
	IUniswapV2Router,
	IUniswapV3RouterQuoter01,
	IUniswapV3RouterQuoter02,
	IUniswapV4PoolManager
}

struct routeChain {
	DexInterfaceType Itype;
	address router;
	address asset;
	uint24 poolFee;
	int24 tickSpacing;
}		


// UniswapV2

interface IUniswapV2Router0102 {
	function WETH() external pure returns (address);
    function getAmountsOut(uint amountIn, address[] memory path) external view returns (uint[] memory amounts);
    function swapExactTokensForETH(uint amountIn,uint amountOutMin,address[] calldata path,address to,uint deadline) external returns (uint[] memory amounts);
    function swapExactTokensForTokens(uint amountIn, uint amountOutMin, address[] calldata path, address to, uint deadline) external returns (uint[] memory amounts);
    function swapExactETHForTokens(uint amountOutMin,address[] calldata path,address to,uint deadline) external payable returns (uint[] memory amounts);
    //function swapETHForExactTokens(uint amountOut,address[] calldata path,address to,uint deadline) external payable returns (uint[] memory amounts);	
}

interface IUniswapV2Pair {
  function token0() external view returns (address);
  function token1() external view returns (address);
  function swap(uint256 amount0Out,	uint256 amount1Out,	address to,	bytes calldata data) external;
}

// UniswapV3

interface IUniswapV3Factory {
    function getPool(address tokenA, address tokenB, uint24 fee) external view returns (address pool);
}


interface IUniswapV3Router {
    struct ExactInputSingleParams {
        address tokenIn;
        address tokenOut;
        uint24 fee;
        address recipient;
        uint256 deadline;
        uint256 amountIn;
        uint256 amountOutMinimum;
        uint160 sqrtPriceLimitX96;
    }
	function WETH9() external pure returns (address);	
	function exactInputSingle(IUniswapV3Router.ExactInputSingleParams calldata params) external returns (uint256 amountOut);
}

interface IUniswapV3Quoter01 {
    function quoteExactInputSingle(address tokenIn, address tokenOut, uint24 fee, uint256 amountIn, uint160 sqrtPriceLimitX96) external returns (uint256 amountOut);
}

interface IUniswapV3Quoter02 {
	struct QuoteExactInputSingleParams {
		address tokenIn;
		address tokenOut;
		uint24 fee;
		uint256 amountIn;
		uint160 sqrtPriceLimitX96;
	}	
	function quoteExactInputSingle(IUniswapV3Quoter02.QuoteExactInputSingleParams calldata params) external returns (uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate);
}

// UniswapV4

//type Currency is address;
//type BalanceDelta is int256;
//type IHooks is address;
/// @dev The minimum value that can be returned from #getSqrtPriceAtTick. Equivalent to getSqrtPriceAtTick(MIN_TICK)
uint160 constant MIN_SQRT_PRICE = 4295128739;
/// @dev The maximum value that can be returned from #getSqrtPriceAtTick. Equivalent to getSqrtPriceAtTick(MAX_TICK)
uint160 constant MAX_SQRT_PRICE = 1461446703485210103287273052203988822378723970342;
// slippage tolerance to allow for unlimited price impact
uint160 constant MIN_PRICE_LIMIT = MIN_SQRT_PRICE + 1;
uint160 constant MAX_PRICE_LIMIT = MAX_SQRT_PRICE - 1;

interface IUniswapV4PoolManager {

	struct PoolKey {
		/// @notice The lower currency of the pool, sorted numerically
		/*Currency*/address currency0;
		/// @notice The higher currency of the pool, sorted numerically
		/*Currency*/address currency1;
		/// @notice The pool swap fee, capped at 1_000_000. If the first bit is 1, the pool has a dynamic fee and must be exactly equal to 0x800000
		uint24 fee;
		/// @notice Ticks that involve positions must be a multiple of tick spacing
		int24 tickSpacing;
		/// @notice The hooks of the pool
		/*IHooks*/address hooks;
	}

	struct SwapParams {
		bool zeroForOne;
		int256 amountSpecified;
		uint160 sqrtPriceLimitX96;
	}

    function swap(PoolKey memory key, SwapParams memory params, bytes calldata hookData) external returns (/*BalanceDelta*/int256);
}

interface IUniswapV4QuoterV4 {
	struct QuoteExactSingleParams {
		IUniswapV4PoolManager.PoolKey poolKey;
		bool zeroForOne;
		address recipient;
		uint128 exactAmount;
		uint160 sqrtPriceLimitX96;
		bytes hookData;
	}		
	function quoteExactInputSingle(QuoteExactSingleParams calldata params) external returns (int128[] memory deltaAmounts, uint160 sqrtPriceX96After, uint32 initializedTicksLoaded);	
}