package main

import (
    "context"
    "fmt"
	"errors"
	"sort"
    //"log"
    "math/big"
    "strings"
	//"encoding/hex"
    //"github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    //"github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

const (
	//zero address				= "0x0000000000000000000000000000000000000000"
	
	// Polygon
	NetworkRPC					= "https://polygon-rpc.com"
	UniswapV2RouterAddress		= "0xedf6066a2b290C185783862C7F4776A2C8077AD1"
	//UniswapV2FactoryAddress		= "0x9e5A52f57b3038F1B8EeE45F28b3C1967e22799C" // get it as 'factory()' from UniswapV3RouterAddress !
	UniswapV3QuoterAddress		= "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6"
	UniswapV3RouterAddress		= "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	//UniswapV3FactoryAddress		= "0x1F98431c8aD98523631AE4a59f267346ea31F984" // get it as 'factory()' from UniswapV3RouterAddress !
	QuickswapV2RouterAddress	= "0xa5E0829CaCEd8fFDD4De3c43696c57F7D7A678ff"
	//QuickswapV2FactoryAddress	= "0x5757371414417b8C6CAad45bAeF941aBc7d3Ab32"
	QuickswapV3QuoterAddress	= "0xa15F0D7377B2A0C0c10db057f641beD21028FC89"	
	QuickswapV3RouterAddress	= "0xf5b509bB0909a69B1c207E495f687a596C168E12"
	//QuickswapV3FactoryAddress	= "0x411b0fAcC3489691f28ad58c47006AF5E3Ab3A28"
	WMATICAddress				= "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	GONEAddress					= "0x162539172b53e9a93b7d98fb6c41682de558a320"
	USDCAddress					= "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"	
	CHAINLAddress				= "0x53e0bca35ec356bd5dddfebbd1fc0fd03fabad39"
	AAVEAddress					= "0xd6df932a45c0f255f85145f286ea0b292b21c90b"
	RNDRAddress					= "0x61299774020dA444Af134c82fa83E3810b309991"
	//TradeAdddress				= "0x733641642aFDf6B5574f2af7969bfBe5730a8daB"
	TradeAdddress				= "0x9CC50ffCF5E76689E74AAC9B3353A0f4c62215E0"
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
	
	UniswapV2RouterABI		= "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountADesired\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBDesired\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"addLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenDesired\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"addLiquidityETH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveOut\",\"type\":\"uint256\"}],\"name\":\"getAmountIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveOut\",\"type\":\"uint256\"}],\"name\":\"getAmountOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"}],\"name\":\"getAmountsIn\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"}],\"name\":\"getAmountsOut\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveB\",\"type\":\"uint256\"}],\"name\":\"quote\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"removeLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"removeLiquidityETH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"removeLiquidityETHSupportingFeeOnTransferTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"approveMax\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"removeLiquidityETHWithPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"approveMax\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"removeLiquidityETHWithPermitSupportingFeeOnTransferTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"approveMax\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"removeLiquidityWithPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapETHForExactTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactETHForTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactETHForTokensSupportingFeeOnTransferTokens\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForETH\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForETHSupportingFeeOnTransferTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForTokensSupportingFeeOnTransferTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMax\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapTokensForExactETH\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMax\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapTokensForExactTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"
	UniswapV2FactoryABI		= "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"}],\"name\":\"getPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	UniswapV2PairABI		= "[{\"constant\":true,\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"name\":\"_reserve0\",\"type\":\"uint112\"},{\"name\":\"_reserve1\",\"type\":\"uint112\"},{\"name\":\"_blockTimestampLast\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
	UniswapV3RouterABI		= "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"}],\"internalType\":\"struct ISwapRouter.ExactInputParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"struct ISwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"}],\"internalType\":\"struct ISwapRouter.ExactOutputParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"struct ISwapRouter.ExactOutputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"refundETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowed\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowedIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sweepToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"sweepTokenWithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"unwrapWETH9\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"unwrapWETH9WithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"	
	UniswapV3QuoterV1ABI	= "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"quoteExactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"name\":\"quoteExactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"quoteExactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"name\":\"quoteExactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	UniswapV3QuoterV2ABI	= "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"quoteExactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160[]\",\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\"},{\"internalType\":\"uint32[]\",\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"struct IQuoterV2.QuoteExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"quoteExactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160[]\",\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\"},{\"internalType\":\"uint32[]\",\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"struct IQuoterV2.QuoteExactOutputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]"	
	UniswapV3FactoryABI		= "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"}],\"name\":\"getPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	UniswapV3PoolABI		= "[{\"inputs\":[],\"name\":\"slot0\",\"outputs\":[{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"},{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"},{\"internalType\":\"uint16\",\"name\":\"observationIndex\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"observationCardinality\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"observationCardinalityNext\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"feeProtocol\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"unlocked\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"liquidity\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	
	// Quickswap variants..
	QuickswapV3QuoterV0ABI	= "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"quoteExactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"name\":\"quoteExactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"quoteExactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"name\":\"quoteExactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	QuickswapV3FactoryABI	= "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"}],\"name\":\"poolByPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	QuickswapAlgebraPoolABI	= "[{\"inputs\":[],\"name\":\"globalState\",\"outputs\":[{\"internalType\":\"uint160\",\"name\":\"price\",\"type\":\"uint160\"},{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"},{\"internalType\":\"uint16\",\"name\":\"fee\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"timepointIndex\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"communityFeeToken0\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"communityFeeToken1\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"unlocked\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"liquidity\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	
	TradeABI				= "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"native_token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"DepositToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_baseAsset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enum DexInterfaceType\",\"name\":\"Itype\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"poolFee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"indexed\":false,\"internalType\":\"struct routeChain[]\",\"name\":\"_routeData\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_gainedAmount\",\"type\":\"uint256\"}],\"name\":\"InstaTraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"WithdrawToken\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enum DexInterfaceType\",\"name\":\"Itype\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"poolFee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"internalType\":\"struct routeChain[]\",\"name\":\"_routedata\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"_startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadlineDeltaSec\",\"type\":\"uint256\"}],\"name\":\"InstaTradeTokens\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enum DexInterfaceType\",\"name\":\"Itype\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"poolFee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"internalType\":\"struct routeChain[]\",\"name\":\"_routedata\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"_startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadlineDeltaSec\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"checkProfit\",\"type\":\"bool\"}],\"name\":\"InstaTradeTokensChecked\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"depositToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"getTokenBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"getTotalTokenBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"safeWithdrawEther\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"safeWithdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAddressesCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenBalances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"
		
)

func SortAddresses(addr1, addr2 common.Address) (common.Address, common.Address) {
    addresses := []common.Address{addr1, addr2}
    sort.Slice(addresses, func(i, j int) bool {
        return addresses[i].Hex() < addresses[j].Hex()
    })
    return addresses[0], addresses[1]
}

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

	if (fee == nil) {
		return quoteExactInputSingle_NoFee(client, router, amountIn, path, sqrtPriceLimitX96)
	}

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

// quickswap variant..
func quoteExactInputSingle_NoFee(client *ethclient.Client, router common.Address, amountIn *big.Int, path []common.Address, sqrtPriceLimitX96 *big.Int) (*big.Int, error) {
	quote := new(big.Int)
    quote.SetInt64(0)
	
	var quoteabi = QuickswapV3QuoterV0ABI;

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
	err = UniswapV3Quoter.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &result, "quoteExactInputSingle", path[0], path[1], amountIn, sqrtPriceLimitX96)
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return quote, err
	}

	quote = *abi.ConvertType(result[len(result)-1], new(*big.Int)).(**big.Int)
    return quote, nil
}

func calculateV2PriceImpact(amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Float, *big.Float) {
	//fmt.Printf("calculateV2PriceImpact\n")

    amountInWithFee := new(big.Int).Mul(amountIn, big.NewInt(997))
    numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
    denominator := new(big.Int).Add(new(big.Int).Mul(reserveIn, big.NewInt(1000)), amountInWithFee)
    amountOut := new(big.Int).Div(numerator, denominator)

    initialPrice := new(big.Float).Quo(new(big.Float).SetInt(reserveOut), new(big.Float).SetInt(reserveIn))
    finalPrice := new(big.Float).Quo(new(big.Float).SetInt(reserveOut.Sub(reserveOut, amountOut)), new(big.Float).SetInt(reserveIn.Add(reserveIn, amountIn)))

    priceImpact := new(big.Float).Quo(finalPrice, initialPrice)
    priceImpact.Sub(priceImpact, big.NewFloat(1))
    priceImpact.Mul(priceImpact, big.NewFloat(100))

    return finalPrice, priceImpact
}

func getV2PairAddress(client *ethclient.Client, routerAddress common.Address, tokenA common.Address, tokenB common.Address) (common.Address, error) {
	//fmt.Printf("getV2PairAddress\n")
	
    uniswapV2Router, err := abi.JSON(strings.NewReader(UniswapV2RouterABI))
    if err != nil {
        return common.Address{}, fmt.Errorf("failed to parse ABI: %v\n", err)	
    }
	
	UniswapV2Router := bind.NewBoundContract(routerAddress, uniswapV2Router, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to the Contract: %v\n", err)
		return common.Address{}, err
    }	

	var fresult []interface{}
	err = UniswapV2Router.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &fresult, "factory")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return common.Address{}, err
	}
	
	factoryAddress := *abi.ConvertType(fresult[0], new(common.Address)).(*common.Address)
	//fmt.Printf("V2 factoryAddress %s\n", factoryAddress.String())	
	
    factoryContract, err := abi.JSON(strings.NewReader(UniswapV2FactoryABI))
    if err != nil {
        return common.Address{}, fmt.Errorf("failed to parse ABI: %v\n", err)
    }

	FactoryContract := bind.NewBoundContract(factoryAddress, factoryContract, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind contract: %v\n", err)
		return common.Address{}, err
    }	
	
	sortedAddress1, sortedAddress2 := SortAddresses(tokenA, tokenB)
	
	var result []interface{}
	err = FactoryContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &result, "getPair", sortedAddress1, sortedAddress2)
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return common.Address{}, err
	}	

	pairAddress := *abi.ConvertType(result[0], new(common.Address)).(*common.Address)
    return pairAddress, nil
}

func getV2Reserves(client *ethclient.Client, pairAddress common.Address) (*big.Int, *big.Int, error) {
	//fmt.Printf("getV2Reserves\n")
	
    pairContract, err := abi.JSON(strings.NewReader(UniswapV2PairABI))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to parse ABI: %v\n", err)	
    }
	
	PairContract := bind.NewBoundContract(pairAddress, pairContract, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to the Contract: %v\n", err)
		return nil, nil, err
    }	

	var result []interface{}
	err = PairContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &result, "getReserves")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return nil, nil, err
	}	

 	reserve0 := *abi.ConvertType(result[0], new(*big.Int)).(**big.Int)
	reserve1 := *abi.ConvertType(result[1], new(*big.Int)).(**big.Int)	
    return reserve0, reserve1, nil
}

func getV2PriceImpact(client *ethclient.Client, routerAddress common.Address, amountIn *big.Int, tokenA common.Address, tokenB common.Address) (*big.Float, *big.Float, error) {
  
	//fmt.Printf("getV2PriceImpact\n") 
  
	sortedAddress1, sortedAddress2 := SortAddresses(tokenA, tokenB)
  
    pairAddress, err := getV2PairAddress(client, routerAddress, sortedAddress1, sortedAddress2)
    if err != nil {
        fmt.Errorf("Failed to get pair address: %vv", err)
		return nil, nil, err
    }
	
    reserve0, reserve1, err := getV2Reserves(client, pairAddress)
    if err != nil {
        fmt.Errorf("Failed to get reserves: %v\n", err)
		return nil, nil, err
    }
	fmt.Printf("V2 Reserves: %s - %s\n", reserve0.String(), reserve1.String())

    price, priceImpact := calculateV2PriceImpact(amountIn, reserve0, reserve1)
	return price, priceImpact, nil		
}

func calculateV3PriceImpact(amountIn *big.Int, sqrtPriceX96 *big.Int, liquidity *big.Int) (*big.Float, *big.Float) {
	
	//fmt.Printf("calculateV3PriceImpact\n")
	
    // Perform the calculations based on Uniswap v3 math.
    // Convert sqrtPriceX96 to price
    price := new(big.Float).Quo(
        new(big.Float).SetInt(new(big.Int).Mul(sqrtPriceX96, sqrtPriceX96)),
        new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 192)),
    )

    // Calculate the new price after the trade
    amountInFloat := new(big.Float).SetInt(amountIn)
    liquidityFloat := new(big.Float).SetInt(liquidity)
    priceImpact := new(big.Float).Quo(amountInFloat, liquidityFloat)

    return price, priceImpact
}

func findV3PoolAddress(client *ethclient.Client, quoterAddress common.Address, tokenA common.Address, tokenB common.Address) (common.Address, *big.Int, error) {

	//fmt.Printf("findV3PoolAddress\n")
	
	addr, err := getV3PoolAddress(client, quoterAddress, tokenA, tokenB, big.NewInt(3000))
	if (err == nil) {
		return addr, big.NewInt(3000), nil		
	}
		
	addr, err = getV3PoolAddress(client, quoterAddress, tokenA, tokenB, big.NewInt(500))
	if (err == nil) {
		return addr, big.NewInt(500), nil
	}

	addr, err = getV3PoolAddress(client, quoterAddress, tokenA, tokenB, big.NewInt(1000))
	if (err == nil) {
		return addr, big.NewInt(1000), nil	
	}
		
	addr, err = getV3PoolAddress(client, quoterAddress, tokenA, tokenB, big.NewInt(10000))
	if (err == nil) {
		return addr, big.NewInt(10000), nil			
	}
	
	addr, err = getV3PoolAddress(client, quoterAddress, tokenA, tokenB, nil)
	if (err == nil) {
		return addr, nil, nil		
	}	
		
	//fmt.Printf("failed to find V3PoolAddress\n")
	
	return common.Address{}, nil, errors.New("failed to find V3PoolAddress")
}
	
		
func getV3PoolAddress(client *ethclient.Client, quoterAddress common.Address, tokenA common.Address, tokenB common.Address, fee *big.Int) (common.Address, error) {
	
	//fmt.Printf("getV3PoolAddress\n")
	
	// quickswap variant..
	if (fee == nil) {
		return getV3PoolAddress_NoFee(client, quoterAddress, tokenA, tokenB)
	}
	
    uniswapV3Quoter, err := abi.JSON(strings.NewReader(UniswapV3QuoterV1ABI))
    if err != nil {
        return common.Address{}, fmt.Errorf("failed to parse ABI: %v\n", err)	
    }
	
	UniswapV3Quoter := bind.NewBoundContract(quoterAddress, uniswapV3Quoter, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to the Contract: %v\n", err)
		return common.Address{}, err
    }	

	var fresult []interface{}
	err = UniswapV3Quoter.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &fresult, "factory")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return common.Address{}, err
	}
	
	factoryAddress := *abi.ConvertType(fresult[0], new(common.Address)).(*common.Address)
	//fmt.Printf("V3 factoryAddress %s\n", factoryAddress.String())		
	
    factoryContract, err := abi.JSON(strings.NewReader(UniswapV3FactoryABI))
    if err != nil {
        return common.Address{}, fmt.Errorf("failed to parse ABI: %v\n", err)
    }
	
	FactoryContract := bind.NewBoundContract(factoryAddress, factoryContract, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to Contract: %v\n", err)
		return common.Address{}, err
    }
    
	sortedAddress1, sortedAddress2 := SortAddresses(tokenA, tokenB)
	
	var result []interface{}
	err = FactoryContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &result, "getPool", sortedAddress1, sortedAddress2, fee)
	if err != nil {
		//fmt.Printf("Failed to call contract function: %v\n", err)
		return common.Address{}, err
	}	

	poolAddress := *abi.ConvertType(result[0], new(common.Address)).(*common.Address)
	
	//fmt.Printf("getV3PoolAddress %s\n", poolAddress.String())
	
    return poolAddress, nil
}

func getV3PoolAddress_NoFee(client *ethclient.Client, quoterAddress common.Address, tokenA common.Address, tokenB common.Address) (common.Address, error) {
	
	//fmt.Printf("getV3PoolAddress_NoFee\n")
	

    uniswapV3Quoter, err := abi.JSON(strings.NewReader(UniswapV3QuoterV1ABI))
    if err != nil {
        return common.Address{}, fmt.Errorf("failed to parse ABI: %v\n", err)	
    }
	
	UniswapV3Quoter := bind.NewBoundContract(quoterAddress, uniswapV3Quoter, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to the Contract: %v\n", err)
		return common.Address{}, err
    }	

	var fresult []interface{}
	err = UniswapV3Quoter.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &fresult, "factory")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return common.Address{}, err
	}
	
	factoryAddress := *abi.ConvertType(fresult[0], new(common.Address)).(*common.Address)
	//fmt.Printf("V3 factoryAddress %s\n", factoryAddress.String())		
	
    factoryContract, err := abi.JSON(strings.NewReader(QuickswapV3FactoryABI))
    if err != nil {
        return common.Address{}, fmt.Errorf("failed to parse ABI: %v\n", err)
    }
	
	FactoryContract := bind.NewBoundContract(factoryAddress, factoryContract, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to Contract: %v\n", err)
		return common.Address{}, err
    }
    
	var result []interface{}

	sortedAddress1, sortedAddress2 := SortAddresses(tokenA, tokenB)
	
	err = FactoryContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &result, "poolByPair", sortedAddress1, sortedAddress2)
	if err != nil {		
	
		fmt.Printf("Failed to call contract function: %v\n", err)
		return common.Address{}, err
	}

	poolAddress := *abi.ConvertType(result[0], new(common.Address)).(*common.Address)
	
	//fmt.Printf("getV3PoolAddress_NoFee %s\n", poolAddress.String())
	
    return poolAddress, nil
}

func getV3PoolSlot0(client *ethclient.Client, poolAddress common.Address) (*big.Int, *big.Int, bool, error) {

	//fmt.Printf("getV3PoolSlot0\n")
	
    poolContract, err := abi.JSON(strings.NewReader(UniswapV3PoolABI))
    if err != nil {
        return nil, nil, false, fmt.Errorf("failed to parse ABI: %v\n", err)
    }
	
	PoolContract := bind.NewBoundContract(poolAddress, poolContract, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to Contract: %v\n", err)
		return nil, nil, false, err
    }

	var result []interface{}
	err = PoolContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &result, "slot0")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return nil, nil, false, err
	}
	
    /*
    var (
        sqrtPriceX96             *big.Int
        tick                     int32
        observationIndex         uint16
        observationCardinality   uint16
        observationCardinalityNext uint16
        feeProtocol              uint8
        unlocked                 bool
    )	
	*/
	
	sqrtPriceX96 := *abi.ConvertType(result[0], new(*big.Int)).(**big.Int)
	tick := *abi.ConvertType(result[1], new(*big.Int)).(**big.Int)
	unlocked := *abi.ConvertType(result[6], new(bool)).(*bool)
		
    return sqrtPriceX96, tick, unlocked, nil	
}

func getV3PoolLiquidity(client *ethclient.Client, poolAddress common.Address) (*big.Int, error) {
	
	//fmt.Printf("getV3PoolLiquidity\n")
	
    // Fetch the liquidity of the pool
    poolContract, err := abi.JSON(strings.NewReader(UniswapV3PoolABI))
    if err != nil {
        return nil, fmt.Errorf("failed to parse ABI: %v", err)
    }

	PoolContract := bind.NewBoundContract(poolAddress, poolContract, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to Contract: %v\n", err)
		return nil, err
    }
	
	var result []interface{}
	err = PoolContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &result, "liquidity")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return nil, err
	}

	liquidity := *abi.ConvertType(result[0], new(*big.Int)).(**big.Int)
    return liquidity, nil
}

func calculateV3Reserves(sqrtPriceX96 *big.Int, liquidity *big.Int) (*big.Int, *big.Int) {
    Q96 := new(big.Int).Lsh(big.NewInt(1), 96)
    price := new(big.Int).Mul(sqrtPriceX96, sqrtPriceX96)
    price.Div(price, Q96)

    reserve0 := new(big.Int).Div(new(big.Int).Mul(liquidity, Q96), sqrtPriceX96)
    reserve1 := new(big.Int).Mul(liquidity, sqrtPriceX96)
    reserve1.Div(reserve1, Q96)

    return reserve0, reserve1
}

func getV3AlgebraPoolGlobalStateAndLiquidity(client *ethclient.Client, poolAddress common.Address) (*big.Int, *big.Int, *big.Int, error) {
    poolABI, err := abi.JSON(strings.NewReader(QuickswapAlgebraPoolABI))
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to parse ABI: %v", err)
    }
	
	PoolContract := bind.NewBoundContract(poolAddress, poolABI, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to Contract: %v\n", err)
		return nil, nil, nil, err
    }
	
	var resultl []interface{}
	err = PoolContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &resultl, "liquidity")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return nil, nil, nil, err
	}	
	
    liquidity := *abi.ConvertType(resultl[0], new(*big.Int)).(**big.Int)	
	
	var resultg []interface{}
	err = PoolContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &resultg, "globalState")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return nil, nil, nil, err
	}
	
	price := *abi.ConvertType(resultg[0], new(*big.Int)).(**big.Int)	
	tick := *abi.ConvertType(resultg[1], new(*big.Int)).(**big.Int)	

    return price, tick, liquidity, nil
}

func calculateV3AlgebraPoolReserves(price *big.Int, liquidity *big.Int) (*big.Int, *big.Int) {
    Q96 := new(big.Int).Lsh(big.NewInt(1), 96)
    sqrtPriceX96 := new(big.Int).Sqrt(price)
    
    reserve0 := new(big.Int).Div(new(big.Int).Mul(liquidity, Q96), sqrtPriceX96)
    reserve1 := new(big.Int).Mul(liquidity, sqrtPriceX96)
    reserve1.Div(reserve1, Q96)

    return reserve0, reserve1
}

func calculateV3AlgebraPoolPrice(price *big.Int) *big.Float {
    Q96 := new(big.Int).Lsh(big.NewInt(1), 96)
    //return new(big.Int).Div(price, Q96)
	return new(big.Float).SetInt(new(big.Int).Quo(price, Q96))
}

func calculateV3AlgebraPoolPriceImpact(tradeAmount *big.Int, reserve0 *big.Int, reserve1 *big.Int, isToken0ToToken1 bool) *big.Float {
    var newReserve0, newReserve1 *big.Int
    if isToken0ToToken1 {
        newReserve0 = new(big.Int).Add(reserve0, tradeAmount)
        //newReserve1 = new(big.Int).Sub(reserve1, new(big.Int).Mul(reserve1, tradeAmount).Div(newReserve0))
		sub1 := new(big.Int).Mul(reserve1, tradeAmount)
		sub1 = sub1.Div(sub1, newReserve0)
		newReserve1 = new(big.Int).Sub(reserve1, sub1)
    } else {
        newReserve1 = new(big.Int).Add(reserve1, tradeAmount)
        //newReserve0 = new(big.Int).Sub(reserve0, new(big.Int).Mul(reserve0, tradeAmount).Div(newReserve1))
		sub0 := new(big.Int).Mul(reserve0, tradeAmount)
		sub0 = sub0.Div(sub0, newReserve1)
		newReserve0 = new(big.Int).Sub(reserve0, sub0)
    }

    newPrice := new(big.Float).Quo(new(big.Float).SetInt(newReserve1), new(big.Float).SetInt(newReserve0))
    currentPrice := new(big.Float).Quo(new(big.Float).SetInt(reserve1), new(big.Float).SetInt(reserve0))

    priceImpact := new(big.Float).Quo(new(big.Float).Sub(newPrice, currentPrice), currentPrice)
    return priceImpact
}

func getV3PriceImpact(client *ethclient.Client, quoterAddress common.Address, amountIn *big.Int, tokenA common.Address, tokenB common.Address, fee *big.Int) (*big.Float, *big.Float, error) {

	//fmt.Printf("getV3PriceImpact\n")
	
	uniswapV3PoolAddress, err := getV3PoolAddress(client, quoterAddress, tokenA, tokenB, fee)
    if err != nil {
        fmt.Errorf("Failed to get pool address: %v\n", err)
		return nil, nil, err
    }
	
	if (fee == nil) {
	
		gsprice, _/*tick*/, liquidity, err := getV3AlgebraPoolGlobalStateAndLiquidity(client, uniswapV3PoolAddress)
		if err != nil {
			fmt.Errorf("Failed to get AlgebraPoolGlobalStateAndLiquidity: %v\n", err)
			return nil, nil, err
		}
		fmt.Printf("Liquidity: %s\n", liquidity.String())

		reserve0, reserve1 := calculateV3AlgebraPoolReserves(gsprice, liquidity)
		if err != nil {
			fmt.Errorf("Failed to calculate V3AlgebraPoolReserves: %v\n", err)
		} else {
			fmt.Printf("AlgebraPool Reserves: %s - %s\n", reserve0.String(), reserve1.String())			
		}
		
		price := calculateV3AlgebraPoolPrice(gsprice)		
		priceImpact := calculateV3AlgebraPoolPriceImpact(amountIn, reserve0, reserve1, true)		
		return price, priceImpact, nil
		
	} else {
	
		sqrtPriceX96, tick, unlocked, err := getV3PoolSlot0(client, uniswapV3PoolAddress)
		if err != nil {
			fmt.Errorf("Failed to get slot0: %v\n", err)
			return nil, nil, err
		}
		fmt.Printf("SqrtPriceX96: %s - tick: %s - unlocked: %t\n", sqrtPriceX96.String(), tick.String(), unlocked)

		liquidity, err := getV3PoolLiquidity(client, uniswapV3PoolAddress)
		if err != nil {
			fmt.Errorf("Failed to get liquidity: %v\n", err)
			return nil, nil, err
		}
		fmt.Printf("Liquidity: %s\n", liquidity.String())
		
		reserve0, reserve1 := calculateV3Reserves(sqrtPriceX96, liquidity)
		fmt.Printf("V3 Reserves: %s - %s\n", reserve0.String(), reserve1.String())	

		price, priceImpact := calculateV3PriceImpact(amountIn, sqrtPriceX96, liquidity)
		return price, priceImpact, nil
	}
}


func main() {

	quoteAmount := big.NewInt(quoteAmount)
	quoteAmount.Div(quoteAmount, big.NewInt(1000))
	//quoteAmount.Mul(quoteAmount, big.NewInt(100))
	
	amountIn := quoteAmount	

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
	tknAmount, err := getTradeTokenBalance(client, common.HexToAddress(TradeAdddress), common.HexToAddress(WMATICAddress))
	fmt.Printf("Trade Token %s: %s - Balance: %s (%.6f ETH)\n", "WMATIC", WMATICAddress, tknAmount.String(), toEther(tknAmount))	
	//fmt.Printf("\n")
	tknAmount, err = getTradeTokenBalance(client, common.HexToAddress(TradeAdddress), common.HexToAddress(token1Address))
	fmt.Printf("Trade Token1 %s: %s - Balance: %s (%.6f ETH)\n", token1Name, token1Address, tknAmount.String(), toEther(tknAmount))
	tknAmount, err = getTradeTokenBalance(client, common.HexToAddress(TradeAdddress), common.HexToAddress(token2Address))
	fmt.Printf("Trade Token2 %s: %s - Balance: %s (%.6f ETH)\n", token2Name, token2Address, tknAmount.String(), toEther(tknAmount))	

	
    path := []common.Address{common.HexToAddress(token1Address), common.HexToAddress(token2Address)}	
	
	amountUniswV2 := big.NewInt(0)
	amountUniswV3 := big.NewInt(0)	
	amountUniswV3Fee := big.NewInt(0)	
	
    // Call getAmountsOut
	fmt.Printf("\nQuote on Uniswap V2 Router: %s\n", UniswapV2RouterAddress)
	pairAddr, err := getV2PairAddress(client, common.HexToAddress(UniswapV2RouterAddress), path[0], path[1])
	if (err == nil) {
		fmt.Printf("V2 PairAddress: %s\n", pairAddr.String())
		amountUniswV2, err = getAmountsOut(client, common.HexToAddress(UniswapV2RouterAddress), amountIn, path)
		if err != nil {
			fmt.Printf("Failed to get Uniswap getAmountsOut: %v\n", err)
		} else {
			fmt.Printf("V2 Uniswap amounts Out: %s (%.6f %s)\n", amountUniswV2.String(), toEther(amountUniswV2), token2Name)		

			price, priceImpact, err := getV2PriceImpact(client, common.HexToAddress(UniswapV2RouterAddress), amountIn, path[0], path[1])
			if err == nil {
				fmt.Printf("Price: %s - Impact: %s%%\n", price.Text('f', 6), priceImpact.Text('f', 6))		
			} else {
				fmt.Printf("V2 Uniswap unable to check price impact: %v\n", err)
			}
		}
	} else {
		fmt.Printf("V2 Uniswap no pair here\n")
	}
	
    // Call quoteExactInputSingle
	fmt.Printf("\nQuote on Uniswap V3 Quoter: %s\n", UniswapV3QuoterAddress)	
	pooladdr, fee, err := findV3PoolAddress(client, common.HexToAddress(UniswapV3QuoterAddress), path[0], path[1])
	if err == nil {
		amountUniswV3Fee = fee;
		fmt.Printf("V3 PoolAddress: %s - Fee: %s\n", pooladdr.String(), fee.String())		
		amountUniswV3, err = quoteExactInputSingle(client, common.HexToAddress(UniswapV3QuoterAddress), false, amountIn, path, fee, big. NewInt(0))
		if err != nil {
			fmt.Printf("Failed to get Uniswap quoteExactInputSingle: %v\n", err)
		} else {
			fmt.Printf("V3 Uniswap amounts Out: %s (%.6f %s)\n", amountUniswV3.String(), toEther(amountUniswV3), token2Name)
			price, priceImpact, err := getV3PriceImpact(client, common.HexToAddress(UniswapV3QuoterAddress), amountIn, path[0], path[1],  fee)
			if err == nil {
				fmt.Printf("Price: %s - Impact: %s%% - Fee: %s\n", price.Text('f', 6), priceImpact.Text('f', 6), fee.String())			
			} else {
				fmt.Printf("V3 Uniswap unable to check price impact: %v\n", err)
			}
		}
	} else {
		fmt.Printf("V3 Uniswap no pool here\n")
	}
	
	amountQuickV2 := big.NewInt(0)
	amountQuickV3 := big.NewInt(0)
	//amountQuickV3Fee := big.NewInt(0)	
	
    // Call getAmountsOut
	fmt.Printf("\nQuote on Quickswap V2 Router: %s\n", QuickswapV2RouterAddress)
	pairAddr, err = getV2PairAddress(client, common.HexToAddress(QuickswapV2RouterAddress), path[0], path[1])
	if (err == nil) {	
		fmt.Printf("V2 PairAddress: %s\n", pairAddr.String())
		amountQuickV2, err = getAmountsOut(client, common.HexToAddress(QuickswapV2RouterAddress), amountIn, path)
		if err != nil {
			fmt.Printf("Failed to get Quickswap getAmountsOut: %v\n", err)
		} else {
			fmt.Printf("V2 Quickswap amounts Out: %s (%.6f %s)\n", amountQuickV2.String(), toEther(amountQuickV2), token2Name)
			
			price, priceImpact, err  := getV2PriceImpact(client, common.HexToAddress(QuickswapV2RouterAddress), amountIn, path[0], path[1])
			if err == nil {
				fmt.Printf("Price: %s - Impact: %s%%\n", price.Text('f', 6), priceImpact.Text('f', 6))		
			} else {
				fmt.Printf("V2 Quickswap unable to check price impact: %v\n", err)
			}			
		}	
	} else {
		fmt.Printf("V2 Quickswap no pair here\n")
	}		
	
	// Call quoteExactInputSingle
	fmt.Printf("\nQuote on Quickswap V3 Quoter: %s\n", QuickswapV3QuoterAddress)
	pooladdr, fee, err = findV3PoolAddress(client, common.HexToAddress(QuickswapV3QuoterAddress), path[0], path[1])
	if err == nil {	
		//amountQuickV3Fee = fee;
		fmt.Printf("V3 PoolAddress: %s - Fee: %s\n", pooladdr.String(), fee.String())		
		amountQuickV3, err = quoteExactInputSingle(client, common.HexToAddress(QuickswapV3QuoterAddress), false, amountIn, path, fee, big. NewInt(0))
		if err != nil {
			fmt.Printf("Failed to get Quickswap quoteExactInputSingle: %v\n", err)
		} else {
			fmt.Printf("V3 Quickswap amounts Out: %s (%.6f %s)\n", amountQuickV3.String(), toEther(amountQuickV3), token2Name)	
			
			price, priceImpact, err := getV3PriceImpact(client, common.HexToAddress(QuickswapV3QuoterAddress), amountIn, path[0], path[1], fee)
			if err == nil {		
				fmt.Printf("Price: %s - Impact: %s%% - Fee: %s\n", price.Text('f', 6), priceImpact.Text('f', 6), fee.String())					
			} else {
				fmt.Printf("V3 Quickswap unable to check price impact: %v\n", err)
			}		
		}
	} else {
		fmt.Printf("V3 Quickswap no pool here\n")
	}		
	
	fmt.Printf("\n")
	quoteBack := big.NewInt(0)
	revertpath := []common.Address{path[1], path[0]}	
	if (amountUniswV3.Cmp(amountQuickV2) < 0) {
		quoteBack, err = quoteExactInputSingle(client, common.HexToAddress(UniswapV3QuoterAddress), false, amountQuickV2, revertpath, amountUniswV3Fee, big. NewInt(0))
		fmt.Printf("Swap buying from QuickV2 and selling to UniswV3 to get back: %.6f ETH (delta: %.6f ETH)\n", toEther(quoteBack), toEther(quoteBack.Sub(quoteBack, quoteAmount)))
		fmt.Printf("[[0,\"%s\",\"%s\",0,0],[1,\"%s\",\"%s\",%s,0]]\n%s\n%f\n", QuickswapV2RouterAddress, token1Address, UniswapV3RouterAddress, token2Address, amountUniswV3Fee.String(), amountIn.String(), toEther(amountIn))		
	} else {
		quoteBack, err = getAmountsOut(client, common.HexToAddress(QuickswapV2RouterAddress), amountUniswV3, revertpath)
		fmt.Printf("Swap buying from UniswV3 and selling to QuickV2 to get back: %.6f ETH (delta: %.6f ETH)\n", toEther(quoteBack), toEther(quoteBack.Sub(quoteBack, quoteAmount)))
		fmt.Printf("[[1,\"%s\",\"%s\",%s,0],[0,\"%s\",\"%s\",0,0]]\n%s\n%f\n", UniswapV3RouterAddress, token1Address, amountUniswV3Fee.String(), QuickswapV2RouterAddress, token2Address, amountIn.String(), toEther(amountIn))
	}	
	fmt.Printf("\n")
}
