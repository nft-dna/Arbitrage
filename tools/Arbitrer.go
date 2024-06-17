package main

import (
    "context"
    "fmt"
    //"log"
    "math/big"
    "strings"

    //"github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    //"github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

const (
	//zero address			= "0x0000000000000000000000000000000000000000"
	
	// Polygon
	NetworkRPC				= "https://polygon-rpc.com"
	UniswapV2RouterAddress	= "0xedf6066a2b290C185783862C7F4776A2C8077AD1"
	UniswapV3QuoterAddress	= "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6"
	UniswapV3RouterAddress	= "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	QuickswapV2RouterAddress= "0xa5E0829CaCEd8fFDD4De3c43696c57F7D7A678ff"
	QuickswapV3QuoterAddress= "0xa15F0D7377B2A0C0c10db057f641beD21028FC89"	
	QuickswapV3RouterAddress= "0xf5b509bB0909a69B1c207E495f687a596C168E12"
	WMATICAddress			= "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	GONEAddress				= "0x162539172b53e9a93b7d98fb6c41682de558a320"
	USDCAddress				= "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"	
	CHAINLAddress			= "0x53e0bca35ec356bd5dddfebbd1fc0fd03fabad39"
	AAVEAddress				= "0xd6df932a45c0f255f85145f286ea0b292b21c90b"
	RNDRAddress				= "0x61299774020dA444Af134c82fa83E3810b309991"
	//TradeAdddress			= "0x733641642aFDf6B5574f2af7969bfBe5730a8daB"
	TradeAdddress			= "0x9CC50ffCF5E76689E74AAC9B3353A0f4c62215E0"
	/*
	0xedf6066a2b290C185783862C7F4776A2C8077AD1 Uniswap V2Router Contract Address 
	0xE592427A0AEce92De3Edee1F18E0157C05861564 Uniswap V3Router
		0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6 Uniswap V3Quoter01
		
	0xa5E0829CaCEd8fFDD4De3c43696c57F7D7A678ff Quickswap V2 router address
	0xf5b509bB0909a69B1c207E495f687a596C168E12 Quickswap V3 Swap router
		0xa15F0D7377B2A0C0c10db057f641beD21028FC89 Quickswap V3 Quoter01 address	
	*/

	// Sepolia ETH
	/*
	NetworkRPC				= "https://ethereum-sepolia-rpc.publicnode.com"
	UniswapV2RouterAddress	= "0x425141165d3DE9FEC831896C016617a52363b687"
	UniswapV3QuoterAddress	= "0xEd1f6473345F45b75F8179591dd5bA1888cf2FB3" //QuoterV2 !!
	WETHAddress				= "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14" // "0xC778417E063141139Fce010982780140Aa0cD5Ab"
	WMCAddress				= "0x5fbad067f69ebbc276410d78ff52823be133ed48"
	*/
	
	/*
	// Mainnet ETH	
	NetworkRPC				= "https://ethereum-rpc.publicnode.com" // wss://ethereum-rpc.publicnode.com"	
    UniswapV2Factory		= "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
    UniswapV3Factory		= "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	UniswapV3RouterAddress	= "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	UniswapV3QuoterAddress	= "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6" // 0x61fFE014bA17989E743c5F6cB21bF9697530B21e QuoterV2
	UniswapV2RouterAddress	= "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"	
	WETHAddress				= "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
    WMCAddress				= "0x7fd4d7737597e7b4ee22acbf8d94362343ae0a79"
	*/

	
	quoteAmount				= 1000000000000000000	// 1 Token (in wei)
	
	UniswapV2RouterABI		="[{\"constant\":false,\"inputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\"},{\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"name\":\"path\",\"type\":\"address[]\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForTokens\",\"outputs\":[{\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\"},{\"name\":\"path\",\"type\":\"address[]\"}],\"name\":\"getAmountsOut\",\"outputs\":[{\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"amountA\",\"type\":\"uint256\"},{\"name\":\"reserveA\",\"type\":\"uint256\"},{\"name\":\"reserveB\",\"type\":\"uint256\"}],\"name\":\"quote\",\"outputs\":[{\"name\":\"amountB\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"
	UniswapV3RouterABI		="[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"}],\"internalType\":\"struct ISwapRouter.ExactInputParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"struct ISwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"}],\"internalType\":\"struct ISwapRouter.ExactOutputParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"struct ISwapRouter.ExactOutputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"refundETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowed\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowedIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sweepToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"sweepTokenWithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"unwrapWETH9\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"unwrapWETH9WithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"
	UniswapV3QuoterV1ABI	="[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"quoteExactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"name\":\"quoteExactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"quoteExactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"name\":\"quoteExactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	UniswapV3QuoterV2ABI	="[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"quoteExactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160[]\",\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\"},{\"internalType\":\"uint32[]\",\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"struct IQuoterV2.QuoteExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"quoteExactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160[]\",\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\"},{\"internalType\":\"uint32[]\",\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"struct IQuoterV2.QuoteExactOutputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	TradeABI				="[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"native_token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"DepositEther\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"DepositToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_baseAsset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enum DexInterfaceType\",\"name\":\"Itype\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"poolFee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"indexed\":false,\"internalType\":\"struct routeChain[]\",\"name\":\"_routeData\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_gainedAmount\",\"type\":\"uint256\"}],\"name\":\"InstaTraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"WithdrawEther\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"WithdrawToken\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enum DexInterfaceType\",\"name\":\"Itype\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"poolFee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"internalType\":\"struct routeChain\",\"name\":\"_routedata\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amountIn\",\"type\":\"uint256\"}],\"name\":\"GetAmountOutMin\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enum DexInterfaceType\",\"name\":\"Itype\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"poolFee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"internalType\":\"struct routeChain[]\",\"name\":\"_routedata\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"_startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadlineDeltaSec\",\"type\":\"uint256\"}],\"name\":\"InstaTradeTokens\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositEther\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"depositToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"etherBalances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEtherBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"getTokenBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalEtherBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"getTotalTokenBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"safeWithdrawEther\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"safeWithdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAddressesCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenBalances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"
		
)

func toEther(wei *big.Int) float64 {
    fWei := new(big.Float).SetInt(wei)
    ethValue := new(big.Float).Quo(fWei, big.NewFloat(1e18))
    value, _ := ethValue.Float64()
    return value
}

func createClient() (*ethclient.Client, error) {
    client, err := ethclient.Dial(NetworkRPC)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
    }
    return client, nil
}

func getTradeEtherBalance(client *ethclient.Client, trade common.Address) (*big.Int, error) {
	amount := new(big.Int)
    amount.SetInt64(0)	

	tradeABI, err := abi.JSON(strings.NewReader(TradeABI))
	if err != nil {
		fmt.Printf("Failed to read TradeABI: %v\n", err)
		return amount, err
	}

    Trade := bind.NewBoundContract(trade, tradeABI, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to the Trade Contract: %v\n", err)
		return amount, err
    }
	
    // Call the getAmountsOut function
	var result []interface{}
    err = Trade.Call(&bind.CallOpts{
        Context: context.Background(),
    }, &result, "getTotalEtherBalance")
    if err != nil {
        fmt.Printf("Failed to call contract function: %v\n", err)
		return amount, err
    }

	amount = *abi.ConvertType(result[0], new(*big.Int)).(**big.Int)
    return amount, nil
}

func getTradeTokenBalance(client *ethclient.Client, trade common.Address, token common.Address) (*big.Int, error) {
	amount := new(big.Int)
    amount.SetInt64(0)	

	tradeABI, err := abi.JSON(strings.NewReader(TradeABI))
	if err != nil {
		fmt.Printf("Failed to read TradeABI: %v\n", err)
		return amount, err
	}

    Trade := bind.NewBoundContract(trade, tradeABI, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to the Trade Contract: %v\n", err)
		return amount, err
    }
	
    // Call the getAmountsOut function
	var result []interface{}
    err = Trade.Call(&bind.CallOpts{
        Context: context.Background(),
    }, &result, "getTotalTokenBalance", token)
    if err != nil {
        fmt.Printf("Failed to call contract function: %v\n", err)
		return amount, err
    }

	amount = *abi.ConvertType(result[0], new(*big.Int)).(**big.Int)
    return amount, nil
}


func getAmountsOut(client *ethclient.Client, router common.Address, amountIn *big.Int, path []common.Address) (*big.Int, error) {
	quote := new(big.Int)
    quote.SetInt64(0)	

	uniswapV2RouterABI, err := abi.JSON(strings.NewReader(UniswapV2RouterABI))
	if err != nil {
		fmt.Printf("Failed to read UniswapV2RouterABI: %v\n", err)
		return quote, err
	}

    UniswapV2Router := bind.NewBoundContract(router, uniswapV2RouterABI, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to the Uniswap V2 Router contract: %v\n", err)
		return quote, err
    }
	
    // Call the getAmountsOut function
	var result []interface{}
    err = UniswapV2Router.Call(&bind.CallOpts{
        Context: context.Background(),
    }, &result, "getAmountsOut", amountIn, path)
    if err != nil {
        fmt.Printf("Failed to call contract function: %v\n", err)
		return quote, err
    }

	out0 := *abi.ConvertType(result[0], new([]*big.Int)).(*[]*big.Int)
	quote = out0[len(out0)-1]
    return quote, nil
}

func quoteExactInputSingle(client *ethclient.Client, router common.Address, useV2 bool, amountIn *big.Int, path []common.Address, fee *big.Int, sqrtPriceLimitX96 *big.Int) (*big.Int, error) {
	quote := new(big.Int)
    quote.SetInt64(0)
	
	var quoteabi = UniswapV3QuoterV1ABI;
	if (useV2) {
		quoteabi = UniswapV3QuoterV2ABI;
	}	
	uniswapV3QuoterABI, err := abi.JSON(strings.NewReader(quoteabi))
	if err != nil {
		fmt.Printf("Failed to read V3QuoterABI: %v\n", err)
		return quote, err
	}

	UniswapV3Quoter := bind.NewBoundContract(router, uniswapV3QuoterABI, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to the Uniswap V3 Quoter contract: %v\n", err)
		return quote, err
    }

    // Call the quoteExactInputSingle function
	var result []interface{}
	if (useV2) {
		// Set up the input parameters
		fee_32 := uint32(fee.Uint64())
		sqrtPriceLimitX96_160, _ := new(big.Int).SetString(sqrtPriceLimitX96.String(), 10)
		params := struct {
			TokenIn         common.Address
			TokenOut        common.Address
			Fee             uint32
			Amount			*big.Int
			SqrtPriceLimitX96 *big.Int
		}{
			TokenIn:         path[0],
			TokenOut:        path[1],
			Fee:             fee_32,
			Amount:			 amountIn,
			SqrtPriceLimitX96: sqrtPriceLimitX96_160,
		}
		input, err := uniswapV3QuoterABI.Pack("quoteExactOutputSingle", params)
		if err != nil {
			fmt.Printf("Failed to pack input: %v\n", err)
			return quote, err
		}		
		err = UniswapV3Quoter.Call(&bind.CallOpts{
			Context: context.Background(),
		}, &result, "quoteExactInputSingle", input)
		if err != nil {
			fmt.Printf("Failed to call contract function: %v\n", err)
			return quote, err
		}
	
	} else {
		err = UniswapV3Quoter.Call(&bind.CallOpts{
			Context: context.Background(),
		}, &result, "quoteExactInputSingle", path[0], path[1], fee, amountIn, sqrtPriceLimitX96)
		if err != nil {
			fmt.Printf("Failed to call contract function: %v\n", err)
			return quote, err
		}
	}

	quote = *abi.ConvertType(result[len(result)-1], new(*big.Int)).(**big.Int)
    return quote, nil
}

func main() {

	amountIn := big.NewInt(quoteAmount)
	//amountIn.Div(amountIn, big.NewInt(100))

    client, err := createClient()
    if err != nil {
        fmt.Printf("Error creating Ethereum client: %v\n", err)
		return
    }

    token1Address := WMATICAddress//USDCAddress
	token1Name := "WMATIC"//"USDC"
    token2Address := GONEAddress
	token2Name := "GONE"
	
	
	fmt.Printf("Token %s: %s - amountIn: %s (%.6f ETH)\n", token1Name, token1Address, amountIn.String(), toEther(amountIn))
	//fmt.Printf("Last verified value: 1 MATIC = 56.965,047 GONE\n")
	fmt.Printf("\n")
	ethAmount, err := getTradeEtherBalance(client, common.HexToAddress(TradeAdddress))
	fmt.Printf("Trade ETH - Balance: %s  (%.6f ETH)\n", ethAmount.String(), toEther(ethAmount))
	tknAmount, err := getTradeTokenBalance(client, common.HexToAddress(TradeAdddress), common.HexToAddress(WMATICAddress))
	fmt.Printf("Trade Token %s: %s - Balance: %s (%.6f ETH)\n", "WMATIC", WMATICAddress, tknAmount.String(), toEther(tknAmount))	
	//fmt.Printf("\n")
	tknAmount, err = getTradeTokenBalance(client, common.HexToAddress(TradeAdddress), common.HexToAddress(token1Address))
	fmt.Printf("Trade Token1 %s: %s - Balance: %s (%.6f ETH)\n", token1Name, token1Address, tknAmount.String(), toEther(tknAmount))
	tknAmount, err = getTradeTokenBalance(client, common.HexToAddress(TradeAdddress), common.HexToAddress(token2Address))
	fmt.Printf("Trade Token2 %s: %s - Balance: %s (%.6f ETH)\n", token2Name, token2Address, tknAmount.String(), toEther(tknAmount))	
	fmt.Printf("\n")

    path := []common.Address{common.HexToAddress(token1Address), common.HexToAddress(token2Address)}	
	
    // Call getAmountsOut
    amountUniswV2, err := getAmountsOut(client, common.HexToAddress(UniswapV2RouterAddress), amountIn, path)
    if err != nil {
        fmt.Printf("Failed to get Uniswap getAmountsOut: %v\n", err)
    } else {
		fmt.Printf("V2 Uniswap amounts Out: %s (%.6f %s)\n", amountUniswV2.String(), toEther(amountUniswV2), token2Name)		
	}
	
    // Call quoteExactInputSingle
    amountUniswV3, err := quoteExactInputSingle(client, common.HexToAddress(UniswapV3QuoterAddress), false, amountIn, path, big. NewInt(3000), big. NewInt(0))
    if err != nil {
        fmt.Printf("Failed to get Uniswap quoteExactInputSingle: %v\n", err)
    } else {
		fmt.Printf("V3 Uniswap amounts Out: %s (%.6f %s)\n", amountUniswV3.String(), toEther(amountUniswV3), token2Name)
	}
	
    // Call getAmountsOut
    amountQuickV2, err := getAmountsOut(client, common.HexToAddress(QuickswapV2RouterAddress), amountIn, path)
    if err != nil {
        fmt.Printf("Failed to get Quickswap getAmountsOut: %v\n", err)
    } else {
		fmt.Printf("V2 Quickswap amounts Out: %s (%.6f %s)\n", amountQuickV2.String(), toEther(amountQuickV2), token2Name)
	}	
	
	// Call quoteExactInputSingle
    amountQuickV3, err := quoteExactInputSingle(client, common.HexToAddress(QuickswapV3QuoterAddress), false, amountIn, path, big. NewInt(3000), big. NewInt(0))
    if err != nil {
        fmt.Printf("Failed to get Quickswap quoteExactInputSingle: %v\n", err)
    } else {
		fmt.Printf("V3 Quickswap amounts Out: %s (%.6f %s)\n", amountQuickV3.String(), toEther(amountQuickV3), token2Name)	
	}
	
	
	quoteBack := big.NewInt(0)
	revertpath := []common.Address{ common.HexToAddress(token2Address), common.HexToAddress(token1Address)}	
	if (amountUniswV3.Cmp(amountQuickV2) < 0) {
		quoteBack, err = quoteExactInputSingle(client, common.HexToAddress(UniswapV3QuoterAddress), false, amountQuickV2, revertpath, big. NewInt(3000), big. NewInt(0))
		fmt.Printf("Swap buyng from QuickV2 and selling to UniswV3 to get back: : %s (%.6f ETH)\n", quoteBack.String(), toEther(quoteBack))
		fmt.Printf("[[0,\"%s\",\"%s\",0,0],[1,\"%s\",\"%s\",3000,0]]\n%s - %.6f\n", QuickswapV2RouterAddress, token1Address, /*UniswapV3QuoterAddress*/UniswapV3RouterAddress, token2Address, amountIn.String(), toEther(amountIn))		
	} else {
		quoteBack, err = getAmountsOut(client, common.HexToAddress(QuickswapV2RouterAddress), amountUniswV3, revertpath)
		fmt.Printf("Swap buyng from UniswV3 and selling to QuickV2 to get back: : %s (%.6f ETH)\n", quoteBack.String(), toEther(quoteBack))
		fmt.Printf("[[1,\"%s\",\"%s\",3000,0],[0,\"%s\",\"%s\",0,0]]\n%s - %.6f\n", UniswapV3QuoterAddress, token1Address, QuickswapV2RouterAddress, token2Address, amountIn.String(), toEther(amountIn))
	}	
	
	/*
	rat := new(big.Rat)
	ratNum := new(big.Rat).SetInt(amountOut)
	ratDen := new(big.Rat).SetInt(quoteAmount)
	if (amountOut.Cmp(quoteAmount) > 0) {
		rat = new(big.Rat).Quo(ratNum, ratDen)
	} else {
		rat = new(big.Rat).Quo(ratDen, ratNum)
	}
	rat.Mul(rat, big.NewRat(100, 1))
	percentage, _ := rat.Float64()
	
	if (amountOut.Cmp(quoteAmount) == 0) {
		fmt.Printf("No trades available\n")
	} else if (amountOut.Cmp(quoteAmount) > 0) {
		fmt.Printf("Buy from: %s\nand sell to: %s\n", UniswapV2RouterAddress, UniswapV3QuoterAddress)
	} else {
		fmt.Printf("Buy from: %s\nand sell to: %s\n", UniswapV3QuoterAddress, UniswapV2RouterAddress)
	}
	fmt.Printf("to gain: %.2f%%\n", percentage)
	*/
}
