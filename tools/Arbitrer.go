package main

import (
	"context"
	"fmt"
	"errors"
	"sort"
	//"log"
	"math/big"
	"strings"
	"os"
	"log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const (
	zero_address				= "0x0000000000000000000000000000000000000000"
	
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

type DexInterfaceType int

const (
    IUniswapV2Router DexInterfaceType = iota
	IUniswapV3RouterQuoter01
	IUniswapV3RouterQuoter02
	IUniswapV4PoolManager
	IQuickswapV3RouterQuoter
)

var dexInterfaceTypeMap = map[string]DexInterfaceType{
    "IUniswapV2Router": IUniswapV2Router,
    "IUniswapV3RouterQuoter01": IUniswapV3RouterQuoter01,
    "IUniswapV3RouterQuoter02": IUniswapV3RouterQuoter02,
	"IUniswapV4PoolManager": IUniswapV4PoolManager,
	"IQuickswapV3RouterQuoter": IQuickswapV3RouterQuoter,
}

type DexRouter struct {
    DexInterface		DexInterfaceType
    Name				string	
	RouterAddress		string	
	QuoterAddress		string	
}

type Contract struct {
    Name				string	
	Address				string	
}

func (c DexInterfaceType) String() string {
    switch c {
    case IUniswapV2Router:
        return "IUniswapV2Router"
    case IUniswapV3RouterQuoter01:
        return "IUniswapV3RouterQuoter01"
    case IUniswapV3RouterQuoter02:
        return "IUniswapV3RouterQuoter02"
    case IUniswapV4PoolManager:
        return "IUniswapV4PoolManager"		
    case IQuickswapV3RouterQuoter:
        return "IQuickswapV3RouterQuoter"			
    default:
        return "Unknown"
    }
}

func (c DexInterfaceType) Int() uint8 {
    switch c {
    case IUniswapV2Router:
        return 0
    case IUniswapV3RouterQuoter01:
        return 1
    case IUniswapV3RouterQuoter02:
        return 2
    case IUniswapV4PoolManager:
        return 3		
    case IQuickswapV3RouterQuoter:
        return 4			
    default:
        return 0
    }
}

func loadNetwork() (string, string, string, bool, error) {
    // Load the .env file
    err := godotenv.Load()
    if err != nil {
        return "", "", "", false, fmt.Errorf("error loading .env file: %v", err)
    }
	
	networkStr := os.Getenv("NETWORK")	
    if networkStr == "" {
        return "", "", "", false, fmt.Errorf("NETWORK not set in the .env file")
    }
	
	rpcStr := os.Getenv(networkStr+"_RPC")	
    if rpcStr == "" {
        return "", "", "", false, fmt.Errorf("%s_RPC not set in the .env file", networkStr)
    }
	
	quoteAmountStr := os.Getenv("QUOTE_AMOUNT")	
    if quoteAmountStr == "" {
        return "", "", "", false, fmt.Errorf("QUOTE_AMOUNT not set in the .env file")
    }	
	
	search_mixed_pools := false
	search_mixed_poolsStr := os.Getenv("SEARCH_MIXED_POOLS")
	if (search_mixed_poolsStr == "") {
		fmt.Printf("WARN: SEARCH_MIXED_POOLS not set in the .env file")
	}
	if (search_mixed_poolsStr == "YES") {
		search_mixed_pools = true
	}
	
	return networkStr, rpcStr, quoteAmountStr, search_mixed_pools, nil;
}

func loadTradeAddress(prefix string) (string, error) {
    // Load the .env file
    err := godotenv.Load()
    if err != nil {
        return "", fmt.Errorf("error loading .env file: %v", err)
    }
	
	tradeStr := os.Getenv(prefix+"_TRADE")	
    if tradeStr == "" {
        fmt.Printf("WARN: %_TRADE not set in the .env file", prefix)
    }	
	
	return tradeStr, nil;
}

func loadDexRouters(prefix string) ([]DexRouter, error) {
    // Load the .env file
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %v", err)
    }
	
    elementsStr := os.Getenv(prefix + "_DEXROUTERS")
    if elementsStr == "" {
        return nil, fmt.Errorf("%s_DEXROUTERS not set in the .env file", prefix)
    }

    // Split the elements by comma
    elementsArray := strings.Split(elementsStr, ",")

    // Create a slice of DexRouter structs
    elements := make([]DexRouter, len(elementsArray))
    for i, elem := range elementsArray {
        parts := strings.Split(strings.TrimSpace(elem), ":")
        if len(parts) != 4 {
            return nil, fmt.Errorf("invalid element format: %s", elem)
        }
        dexInterface, exists := dexInterfaceTypeMap[parts[0]]
        if !exists {
            return nil, fmt.Errorf("invalid dexInterface for element %u: %s", i, parts[0])
        }		
        name := parts[1]
		routerAddress := parts[2]
		quoterAddress := parts[3]
        elements[i] = DexRouter{DexInterface: dexInterface, Name: name, RouterAddress: routerAddress, QuoterAddress : quoterAddress }
    }

    return elements, nil
}

func loadNativeToken(prefix string) (Contract, error) {
    // Load the .env file
    err := godotenv.Load()
    if err != nil {
        return Contract{}, fmt.Errorf("error loading .env file: %v", err)
    }
	
    elementsStr := os.Getenv(prefix + "_NATIVE")
    if elementsStr == "" {
        return Contract{}, fmt.Errorf("%s_NATIVE not set in the .env file", prefix)
    }

    parts := strings.Split(strings.TrimSpace(elementsStr), ":")
    if len(parts) != 2 {
        return Contract{}, fmt.Errorf("invalid element format: %s", elementsStr)
    }
    name := parts[0]
	address := parts[1]
    element := Contract{Name: name, Address: address}

    return element, nil
}

func loadStableTokens(prefix string) ([]Contract, error) {
    // Load the .env file
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %v", err)
    }
	
    elementsStr := os.Getenv(prefix + "_STABLES")
    if elementsStr == "" {
        return nil, nil//fmt.Errorf("%s_STABLES not set in the .env file", prefix)
    }

    // Split the elements by comma
    elementsArray := strings.Split(elementsStr, ",")

    // Create a slice of DexRouter structs
    elements := make([]Contract, len(elementsArray))
    for i, elem := range elementsArray {
        parts := strings.Split(strings.TrimSpace(elem), ":")
        if len(parts) != 2 {
            return nil, fmt.Errorf("invalid element format: %s", elem)
        }
        name := parts[0]
		address := parts[1]
        elements[i] = Contract{Name: name, Address: address}
    }

    return elements, nil
}

func loadTestTokens(prefix string) ([]Contract, error) {
    // Load the .env file
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %v", err)
    }
	
    elementsStr := os.Getenv(prefix + "_TOKENS")
    if elementsStr == "" {
        return nil, fmt.Errorf("%s_TOKENS not set in the .env file", prefix)
    }

    // Split the elements by comma
    elementsArray := strings.Split(elementsStr, ",")

    // Create a slice of DexRouter structs
    elements := make([]Contract, len(elementsArray))
    for i, elem := range elementsArray {
        parts := strings.Split(strings.TrimSpace(elem), ":")
        if len(parts) != 2 {
            return nil, fmt.Errorf("invalid element format: %s", elem)
        }
        name := parts[0]
		address := parts[1]
        elements[i] = Contract{Name: name, Address: address}
    }

    return elements, nil
}

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

func createClient(rpc string) (*ethclient.Client, error) {
    client, err := ethclient.Dial(rpc)
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
	//fee_32 := uint32(fee.Uint64())
		
	var result []interface{}
	if (useV2) {
		// Set up the input parameters
		sqrtPriceLimitX96_160, _ := new(big.Int).SetString(sqrtPriceLimitX96.String(), 10)
		params := struct {
			TokenIn         common.Address
			TokenOut        common.Address
			Fee             *big.Int//uint32
			Amount			*big.Int
			SqrtPriceLimitX96 *big.Int
		}{
			TokenIn:         path[0],
			TokenOut:        path[1],
			Fee:             fee,//_32,
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
			fmt.Printf("Failed to call V2 contract function: %v\n", err)
			return quote, err
		}
	
	} else {
		err = UniswapV3Quoter.Call(&bind.CallOpts{
			Context: context.Background(),
		}, &result, "quoteExactInputSingle", path[0], path[1], fee, amountIn, sqrtPriceLimitX96)
		if err != nil {
			fmt.Printf("Failed to call V1 contract function: %v\n", err)
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
		return nil, nil, fmt.Errorf("Failed to get pair address: %v", err)
    }
	
	if (pairAddress.String() == zero_address) {
		return nil, nil, fmt.Errorf("Empty pair address")
	}
	
    reserve0, reserve1, err := getV2Reserves(client, pairAddress)
    if err != nil {
        fmt.Errorf("Failed to get reserves: %v\n", err)
		return nil, nil, err
    }
	fmt.Printf("   V2 Reserves: %.9f - %.9f\n", toEther(reserve0), toEther(reserve1))

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
	if (err == nil && addr.String() != zero_address) {
		return addr, big.NewInt(3000), nil		
	}
		
	addr, err = getV3PoolAddress(client, quoterAddress, tokenA, tokenB, big.NewInt(500))
	if (err == nil && addr.String() != zero_address) {
		return addr, big.NewInt(500), nil
	}

	addr, err = getV3PoolAddress(client, quoterAddress, tokenA, tokenB, big.NewInt(1000))
	if (err == nil && addr.String() != zero_address) {
		return addr, big.NewInt(1000), nil	
	}
		
	addr, err = getV3PoolAddress(client, quoterAddress, tokenA, tokenB, big.NewInt(10000))
	if (err == nil && addr.String() != zero_address) {
		return addr, big.NewInt(10000), nil			
	}
	
	addr, err = getV3PoolAddress(client, quoterAddress, tokenA, tokenB, nil)
	if (err == nil && addr.String() != zero_address) {
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

func getV3AlgebraPoolGlobalStateAndLiquidity(client *ethclient.Client, poolAddress common.Address) (*big.Int, *big.Int, *big.Int, *big.Int, bool, error) {
    poolABI, err := abi.JSON(strings.NewReader(QuickswapAlgebraPoolABI))
    if err != nil {
        return nil, nil, nil, nil, false, fmt.Errorf("failed to parse ABI: %v", err)
    }
	
	PoolContract := bind.NewBoundContract(poolAddress, poolABI, client, client, client)
    if err != nil {
        fmt.Printf("Failed to bind to Contract: %v\n", err)
		return nil, nil, nil, nil, false, err
    }
	
	var resultl []interface{}
	err = PoolContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &resultl, "liquidity")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return nil, nil, nil, nil, false, err
	}	
	
    liquidity := *abi.ConvertType(resultl[0], new(*big.Int)).(**big.Int)	
	
	var resultg []interface{}
	err = PoolContract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &resultg, "globalState")
	if err != nil {
		fmt.Printf("Failed to call contract function: %v\n", err)
		return nil, nil, nil, nil, false, err
	}

	/*
	var globalState struct {
        Price            *big.Int
        Tick             int32
        Fee              uint16
        TimepointIndex   uint16
        CommunityFeeToken0 uint8
        CommunityFeeToken1 uint8
        Unlocked         bool
    }
	*/	
	price := *abi.ConvertType(resultg[0], new(*big.Int)).(**big.Int)	
	tick := *abi.ConvertType(resultg[1], new(*big.Int)).(**big.Int)	
	fee := *abi.ConvertType(resultg[2], new(uint16)).(*uint16)	
	unlocked := *abi.ConvertType(resultg[6], new(bool)).(*bool)	

    return price, tick, liquidity, big.NewInt((int64)(fee)), unlocked, nil
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

/*
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
*/
func calculateV3AlgebraPoolPriceImpact(price *big.Int, tradeAmount *big.Int, reserve0 *big.Int, reserve1 *big.Int) *big.Float {
    initialPrice := new(big.Float).SetInt(price)
    initialPrice.Mul(initialPrice, initialPrice).Quo(initialPrice, big.NewFloat(1<<96))

    // Simulate price change after trade
	tradeAmountFloat := new(big.Float).SetInt(tradeAmount)
    newReserve1 := new(big.Float).Add(new(big.Float).SetInt(reserve1), tradeAmountFloat)
    newPrice := new(big.Float).Quo(new(big.Float).SetInt(reserve0), newReserve1)

    // Calculate price impact
    priceImpact := new(big.Float).Sub(initialPrice, newPrice)
    priceImpact.Quo(priceImpact, initialPrice)
    priceImpactPercentage := new(big.Float).Mul(priceImpact, big.NewFloat(100))
	return priceImpactPercentage;
}

func getV3PriceImpact(client *ethclient.Client, quoterAddress common.Address, amountIn *big.Int, tokenA common.Address, tokenB common.Address, fee *big.Int) (*big.Float, *big.Float, error) {

	//fmt.Printf("getV3PriceImpact\n")
	
	uniswapV3PoolAddress, err := getV3PoolAddress(client, quoterAddress, tokenA, tokenB, fee)
    if err != nil {        
		return nil, nil, fmt.Errorf("Failed to get pool address: %v\n", err)
    }
	
	if (uniswapV3PoolAddress.String() == zero_address) {        
		return nil, nil, fmt.Errorf("Empty pool address\n")	
	}
	 
	if (fee == nil) {
	
		gsprice, _/*tick*/, liquidity, fee, unlocked, err := getV3AlgebraPoolGlobalStateAndLiquidity(client, uniswapV3PoolAddress)
		if err != nil {
			fmt.Errorf("Failed to get AlgebraPoolGlobalStateAndLiquidity: %v\n", err)
			return nil, nil, err
		}
		fmt.Printf("   Liquidity: %.9f - Fee: %s - Unlocked: %t\n", toEther(liquidity), fee.String(), unlocked)
		if (unlocked == false) {			
			return nil, nil, fmt.Errorf("Pool is locked\n")
		}		

		reserve0, reserve1 := calculateV3AlgebraPoolReserves(gsprice, liquidity)
		if err != nil {
			fmt.Errorf("Failed to calculate V3AlgebraPoolReserves: %v\n", err)
		} else {
			fmt.Printf("   AlgebraPool Reserves: %.9f - %.9f\n", toEther(reserve0), toEther(reserve1))			
		}
		
		price := calculateV3AlgebraPoolPrice(gsprice)		
		priceImpact := calculateV3AlgebraPoolPriceImpact(gsprice, amountIn, reserve0, reserve1)		
		return price, priceImpact, nil
		
	} else {
	
		sqrtPriceX96, tick, unlocked, err := getV3PoolSlot0(client, uniswapV3PoolAddress)
		if err != nil {
			fmt.Errorf("Failed to get slot0: %v\n", err)
			return nil, nil, err
		}
		fmt.Printf("   SqrtPriceX96: %s - tick: %s - unlocked: %t\n", sqrtPriceX96.String(), tick.String(), unlocked)
		if (unlocked == false) {			
			return nil, nil, fmt.Errorf("Pool is locked\n")
		}			

		liquidity, err := getV3PoolLiquidity(client, uniswapV3PoolAddress)
		if err != nil {
			fmt.Errorf("Failed to get liquidity: %v\n", err)
			return nil, nil, err
		}
		fmt.Printf("   Liquidity: %.9f\n", toEther(liquidity))
		
		reserve0, reserve1 := calculateV3Reserves(sqrtPriceX96, liquidity)
		fmt.Printf("   V3 Reserves: %.9f - %.9f\n", toEther(reserve0), toEther(reserve1))	

		price, priceImpact := calculateV3PriceImpact(amountIn, sqrtPriceX96, liquidity)
		return price, priceImpact, nil
	}
}


func main() {

	network, rpc, quoteAmountStr, search_mixed_pools, err := loadNetwork();
	if err != nil {
		log.Fatalf("Failed to load Network: %v", err)
		return
	}	
	
	tradeAddress, err := loadTradeAddress(network)
	
	dexRouters, err := loadDexRouters(network)
	if err != nil {
		log.Fatalf("Failed to load DexRouters: %v", err)
		return
	}
    fmt.Println("Loaded Routers:")
    for _, elemr := range dexRouters {
        fmt.Printf("- Name: %s\n    Type: %s\n    Router: %s\n    Quoter: %s\n", elemr.DexInterface.String(), elemr.Name, elemr.RouterAddress, elemr.QuoterAddress)
    }	

	nativeToken, err := loadNativeToken(network)
	if err != nil {
		log.Fatalf("Failed to load NativeToken: %v", err)
		return
	}
	fmt.Println("Native Token:")
    fmt.Printf("- Name: %s, Address: %s\n", nativeToken.Name, nativeToken.Address)
	
	stableTokens, err := loadStableTokens(network)
	if err != nil {
		log.Fatalf("Failed to load StableTokens: %v", err)
		return
	}	
	fmt.Println("Loaded Stables:")
    for _, elems := range stableTokens {
        fmt.Printf("- Name: %s, Address: %s\n", elems.Name, elems.Address)
    }		

	testTokens, err := loadTestTokens(network)
	if err != nil {
		log.Fatalf("Failed to load TestTokens: %v", err)
		return
	}	
	fmt.Println("Loaded Tokens:")
    for _, elemt := range testTokens {
        fmt.Printf("- Name: %s, Address: %s\n", elemt.Name, elemt.Address)
    }		

	//quoteAmount := big.NewInt(quoteAmountStr)
	quoteAmount := big.NewInt(0)
	quoteAmount.SetString(quoteAmountStr, 10)
	//quoteAmount.Div(quoteAmount, big.NewInt(2))
	//quoteAmount.Mul(quoteAmount, big.NewInt(10))
	
	amountIn := quoteAmount	

    client, err := createClient(rpc)
    if err != nil {
        fmt.Printf("Error creating Ethereum client: %v\n", err)
		return
    }
	
	type quoteResult struct {
        Errored			bool
        Unlocked		bool		
		Price			*big.Float
        Fee				*big.Int
		Quote			*big.Int
		PriceImpact		*big.Float
		Reserve1		*big.Int
		Reserve2		*big.Int
    }
	
	for _, token1 := range testTokens {
		for _, token2 := range testTokens {
		
			if (token1.Address == token2.Address || (search_mixed_pools == false)) {
				token1 = nativeToken
			}
		
			fmt.Printf("\n\nTesting %s:%s\n     on %s:%s\nQuoting: %s (%.6f ETH)\n", token1.Name, token1.Address, token2.Name, token2.Address, amountIn.String(), toEther(amountIn))
			
			if (tradeAddress != "") {
				tknAmount, err := getTradeTokenBalance(client, common.HexToAddress(tradeAddress), common.HexToAddress(token1.Address))
				if err != nil {
					fmt.Printf("Trade Token1 %s: %s - Balance: %s (%.6f ETH)\n", token1.Name, token1.Address, tknAmount.String(), toEther(tknAmount))
				}
				tknAmount, err = getTradeTokenBalance(client, common.HexToAddress(tradeAddress), common.HexToAddress(token2.Address))
				if err != nil {
					fmt.Printf("Trade Token2 %s: %s - Balance: %s (%.6f ETH)\n", token2.Name, token2.Address, tknAmount.String(), toEther(tknAmount))
				}
			}			
	
			path := []common.Address{common.HexToAddress(token1.Address), common.HexToAddress(token2.Address)}	

			sortedAddress1, sortedAddress2 := SortAddresses(path[0], path[1])
			fmt.Printf("Sorted pair: %s - %s\n", sortedAddress1.String(), sortedAddress2.String())				
	
			results := make([]quoteResult, len(dexRouters))
			
			i := 0
			
			for _, dex := range dexRouters {
								
				fmt.Printf("Dex idx %d - Quote on: %s - Router: %s - Type: %s\n", i, dex.Name, dex.RouterAddress, dex.DexInterface.String())
				
				results[i].Errored = true
				results[i].Unlocked = false
				
				if (dex.DexInterface == IUniswapV2Router) {
					 // Call getAmountsOut
					pairAddr, err := getV2PairAddress(client, common.HexToAddress(dex.RouterAddress), path[0], path[1])
					if (err == nil && pairAddr.String() != zero_address) {
						fmt.Printf("   V2 PairAddress: %s\n", pairAddr.String())
						results[i].Quote, err = getAmountsOut(client, common.HexToAddress(dex.RouterAddress), amountIn, path)
						if err != nil {
							fmt.Printf("   Failed to get %s getAmountsOut: %v\n", dex.Name, err)
						} else {
							results[i].Errored = false						
							fmt.Printf("   V2 %s amounts Out: %s (%.6f %s)\n", dex.Name, results[i].Quote.String(), toEther(results[i].Quote), token2.Name)		
							results[i].Price, results[i].PriceImpact, err = getV2PriceImpact(client, common.HexToAddress(dex.RouterAddress), amountIn, path[0], path[1])
							if err == nil {
								results[i].Unlocked	= true
								fmt.Printf("   Price: %s - Impact: %s%%\n", results[i].Price.Text('f', 6), results[i].PriceImpact.Text('f', 6))								
							} else {
								fmt.Printf("   V2 %s unable to check price impact: %v\n", dex.Name, err)
							}
						}
					} else {
						fmt.Printf("   V2 %s no pair here\n", dex.Name)
					}					 
				//} else if (dex.DexInterface == IQuickswapV3RouterQuoter) {					
				} else if ((dex.DexInterface == IUniswapV3RouterQuoter01) || (dex.DexInterface == IUniswapV3RouterQuoter02) || (dex.DexInterface == IQuickswapV3RouterQuoter)) {
					// Call quoteExactInputSingle
					fmt.Printf("   Quote on %s V3 Quoter: %s\n", dex.Name, dex.QuoterAddress)	
					pooladdr, fee, err := findV3PoolAddress(client, common.HexToAddress(dex.QuoterAddress), path[0], path[1])
					if (err == nil && pooladdr.String() != zero_address) {
						results[i].Fee = fee;
						fmt.Printf("   V3 PoolAddress: %s - Fee: %s\n", pooladdr.String(), results[i].Fee.String())		
						results[i].Quote, err = quoteExactInputSingle(client, common.HexToAddress(dex.QuoterAddress), (dex.DexInterface == IUniswapV3RouterQuoter02), amountIn, path, results[i].Fee, big. NewInt(0))
						if err != nil {
							fmt.Printf("   Failed to get %s quoteExactInputSingle: %v\n", dex.Name, err)
						} else {
							results[i].Errored = false						
							fmt.Printf("   V3 %s amounts Out: %s (%.6f %s)\n", dex.Name, results[i].Quote.String(), toEther(results[i].Quote), token2.Name)
							results[i].Price, results[i].PriceImpact, err = getV3PriceImpact(client, common.HexToAddress(dex.QuoterAddress), amountIn, path[0], path[1], results[i].Fee)
							if err == nil {
								fmt.Printf("   Price: %s - Impact: %s%% - Fee: %s\n", results[i].Price.Text('f', 6), results[i].PriceImpact.Text('f', 6), results[i].Fee.String())
								results[i].Unlocked	= true
							} else {
								fmt.Printf("   V3 %s unable to check price impact: %v\n", dex.Name, err)
							}
						}
					} else {
						fmt.Printf("   V3 %s no pool here\n", dex.Name)
					}				
				} else if (dex.DexInterface == IUniswapV4PoolManager) {
					fmt.Printf("   WARN: IUniswapV4PoolManager unsupported yet\n")
				}
				
				i = i + 1
			}	
			
			minIdx := -1
			maxIdx := -1

			i = 0
			for _, res := range results {
				if ((res.Errored == false)&&(res.Unlocked == true)&&(res.PriceImpact.Cmp(big.NewFloat(0.02)) < 0)&&(res.PriceImpact.Cmp(big.NewFloat(-0.02)) > 0)) {
					if (minIdx == -1 || results[minIdx].Quote.Cmp(res.Quote) < 0) {
						minIdx = i
					}
					if (maxIdx == -1 || results[maxIdx].Quote.Cmp(res.Quote) > 0) {
						maxIdx = i
					}					
				}
				i = i + 1
			}

			fmt.Printf("\n")	
			if (minIdx != -1 && maxIdx != -1 && minIdx != maxIdx) {
					
				//dbg
				//fmt.Printf("minIdx: %d - maxIdx: %d\n", minIdx, maxIdx)
				
				quoteBack := big.NewInt(0)
				revertpath := []common.Address{path[1], path[0]}	
				
				if (dexRouters[maxIdx].DexInterface == IUniswapV2Router) {
					quoteBack, err = getAmountsOut(client, common.HexToAddress(dexRouters[maxIdx].RouterAddress), results[minIdx].Quote, revertpath)				
				} else if ((dexRouters[maxIdx].DexInterface == IUniswapV3RouterQuoter01) || (dexRouters[maxIdx].DexInterface == IUniswapV3RouterQuoter02) || (dexRouters[maxIdx].DexInterface == IQuickswapV3RouterQuoter)) {
					quoteBack, err = quoteExactInputSingle(client, common.HexToAddress(dexRouters[maxIdx].QuoterAddress), (dexRouters[maxIdx].DexInterface == IUniswapV3RouterQuoter02), results[minIdx].Quote, revertpath, results[maxIdx].Fee, big.NewInt(0))
				} else if (dexRouters[maxIdx].DexInterface == IUniswapV4PoolManager) {
					fmt.Printf("WARN: IUniswapV4PoolManager unsupported yet\n")
				}
				if results[minIdx].Fee == nil {
					results[minIdx].Fee = big.NewInt(0)
				}
				if results[maxIdx].Fee == nil {
					results[maxIdx].Fee = big.NewInt(0)
				}						
				fmt.Printf("Swap %s\\%s\n  buying %.6f ETH from %s and selling to %s to get back: %.6f ETH (delta: %.6f ETH)\n", token1.Name, token2.Name, toEther(quoteAmount), dexRouters[minIdx].Name, dexRouters[maxIdx].Name, toEther(quoteBack), toEther(quoteBack.Sub(quoteBack, quoteAmount)))
				fmt.Printf("[[%d,\"%s\",\"%s\",%s,0],[%d,\"%s\",\"%s\",%s,0]]\n%s\n%f\n", dexRouters[minIdx].DexInterface.Int(), dexRouters[minIdx].RouterAddress, token1.Address, results[minIdx].Fee.String(), dexRouters[maxIdx].DexInterface.Int(), dexRouters[maxIdx].RouterAddress, token2.Address, results[maxIdx].Fee.String(), amountIn.String(), toEther(amountIn))		
			
			} else {
				fmt.Printf("No Swap available\n")
			}			
			fmt.Printf("\n")
		}	
		if (search_mixed_pools == false)	{
			break;
		}		
	}
	
	fmt.Printf("\n")
}
