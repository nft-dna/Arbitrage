// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "./Interfaces.sol";
import "./Trade.sol";

/*

	UniswapV2			Factory Contract Address					V2Router02 Contract Address
	Mainnet				0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f	0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D
	Ethereum Sepolia	0xB7f907f7A9eBC822a80BD25E224be42Ce0A698A0	0x425141165d3DE9FEC831896C016617a52363b687
	Arbitrum			0xf1D7CC64Fb4452F05c498126312eBE29f30Fbcf9	0x4752ba5dbc23f44d87826276bf6fd6b1c372ad24
	Avalanche			0x9e5A52f57b3038F1B8EeE45F28b3C1967e22799C	0x4752ba5dbc23f44d87826276bf6fd6b1c372ad24
	BNB Chain			0x8909Dc15e40173Ff4699343b6eB8132c65e18eC6	0x4752ba5dbc23f44d87826276bf6fd6b1c372ad24
	Base				0x8909Dc15e40173Ff4699343b6eB8132c65e18eC6	0x4752ba5dbc23f44d87826276bf6fd6b1c372ad24
	Optimism			0x0c3c1c532F1e39EdF36BE9Fe0bE1410313E074Bf	0x4A7b5Da61326A6379179b40d00F57E5bbDC962c2
	Polygon				0x9e5A52f57b3038F1B8EeE45F28b3C1967e22799C	0xedf6066a2b290C185783862C7F4776A2C8077AD1
	Blast				0x5C346464d33F90bABaf70dB6388507CC889C1070	0xBB66Eb1c5e875933D44DAe661dbD80e5D9B03035

	UniswapV3	Mainnet, Polygon, Optimism, Arbitrum, Testnets Address
	UniswapV3Factory	0x1F98431c8aD98523631AE4a59f267346ea31F984
	Multicall2			0x5BA1e12693Dc8F9c48aAD8770482f4739bEeD696
	ProxyAdmin			0xB753548F6E010e7e680BA186F9Ca1BdAB2E90cf2
	TickLens			0xbfd8137f7d1516D3ea5cA83523914859ec47F573
	Quoter				0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6
	SwapRouter			0xE592427A0AEce92De3Edee1F18E0157C05861564
	NFTDescriptor		0x42B24A95702b9986e82d421cC3568932790A48Ec
	NonfungibleTokenPositionDescriptor	0x91ae842A5Ffd8d12023116943e72A606179294f3
	TransparentUpgradeableProxy		0xEe6A57eC80ea46401049E92587E52f5Ec1c24785
	NonfungiblePositionManager		0xC36442b4a4522E871399CD717aBDD847Ab11FE88
	V3Migrator			0xA5644E29708357803b5A882D272c41cC0dF92B34

	UniswapV3	Celo Address
	UniswapV3Factory	0xAfE208a311B21f13EF87E33A90049fC17A7acDEc
	Multicall2			0x633987602DE5C4F337e3DbF265303A1080324204
	ProxyAdmin			0xc1b262Dd7643D4B7cA9e51631bBd900a564BF49A
	TickLens			0x5f115D9113F88e0a0Db1b5033D90D4a9690AcD3D
	Quoter				0x82825d0554fA07f7FC52Ab63c961F330fdEFa8E8
	SwapRouter			0x5615CDAb10dc425a742d643d949a7F474C01abc4
	NFTDescriptor		0xa9Fd765d85938D278cb0b108DbE4BF7186831186
	NonfungibleTokenPositionDescriptor	0x644023b316bB65175C347DE903B60a756F6dd554
	TransparentUpgradeableProxy		0x505B43c452AA4443e0a6B84bb37771494633Fde9
	NonfungiblePositionManager		0x3d79EdAaBC0EaB6F08ED885C05Fc0B014290D95A
	V3Migrator			0x3cFd4d48EDfDCC53D3f173F596f621064614C582
	
	UniswapV4	Sepolia
	| Contract               | Address                                
	|------------------------|--------------------------------------------
	| PoolManager            | 0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9 
	| PoolSwapTest           | 0x0165878A594ca255338adfa4d48449f69242Eb8F 
	| PoolModifyPositionTest | 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707 
	| PoolDonateTest         | 0xa513E6E4b8f2a923D98304ec87F64353C4D5C853 

	WETH
	Mainnet
	Ethereum Mainnet 0xC02aaA39b223FE8D0A0e5C4F27ead9083C756Cc2
	Arbitrum Mainnet 0x82af49447d8a07e3bd95bd0d56f35241523fbab1
	Optimism Mainnet 0x4200000000000000000000000000000000000006
	Polygon (Matic) Mainnet 0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619
	BSC (Binance Smart Chain) Mainnet 0x2170ed0880ac9a755fd29b2688956bd959f933f8
	Avalanche Mainnet: 0x49D5c2BdFfac6CE2BFdB6640F4F80f226bc10bAB
	Fantom Mainnet 0x74b23882a30290451a17c44f4f05243b6b58c76d
	Testnets
	Ethereum Goerli Testnet 0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6
	Ethereum Sepolia Testnet 0xC778417E063141139Fce010982780140Aa0cD5Ab
	Arbitrum Goerli Testnet 0xe39Ab88f8A4777030A534146A9Ca3B52bd5D43A3
	Optimism Goerli Testnet 0x4200000000000000000000000000000000000006
	Polygon Mumbai Testnet 0xDfd5eC59A2F15b56fA139F5D44C11f4C8f869b60
	
	
	Polygon		Factory										Router										Quoter
	UniswapV3	0x1F98431c8aD98523631AE4a59f267346ea31F984	0xE592427A0AEce92De3Edee1F18E0157C05861564	0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6
	v2
	UniswapV2	0x9e5A52f57b3038F1B8EeE45F28b3C1967e22799C	0xedf6066a2b290C185783862C7F4776A2C8077AD1
	Quickswap	0x5757371414417b8c6caad45baef941abc7d3ab32	0xa5E0829CaCED8fFDD4De3c43696c57F7D7A678ff
	SushiSwap	0xc35DADB65012eC5796536bD9864eD8773aBc74C4	0x1b02da8cb0d097eb8d57a175b88c7d8b47997506
	Dfyn		0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f	0xA102072A4C07F06EC3B4900FDC4C7B80b6c57429
	KyberSwap	0xc9fAd10916B35b4d91Df3b19eC0f86e36dEdBA7b	0x546C79662E028B661Dfb4767664d0273184E4dD1
	other
	**Balancer	0x67D27634E44793FE63c467035E31ea8635117cd4	0xBA12222222228d8Ba445958a75a0704d566BF2C8	IBalancerVault { struct SwapReques
	
	WETH		0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619
	USDC		0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174
	USDT		0xC2132D05D31c914A87C6611C10748aEB04B58e8F
	DAI			0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063
	TUSD		0x2e1AD108fF1D8C782fcBbB89AAd783aC49586756
	BUSD		0xa8d394fe7380b8ce6145d5f85e6ac22d4e91acde		
	
	
	Ethereum	Factory										Router										Quoter
	UniswapV3	0x1F98431c8aD98523631AE4a59f267346ea31F984	0xE592427A0AEce92De3Edee1F18E0157C05861564	0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6
	v2
	UniswapV2	0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f	0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D
	SushiSwap	0xc35DADB65012eC5796536bD9864eD8773aBc74C4	0xd9e1CE17f2641F24aE83637ab66a2cca9C378B9F
	KyberSwap	*0x818E6FECD516Ecc3849DAf6845e3EC868087B755	0x7a250d5630b4cf539739df2c5dacaebc00c26ac0
	other		
	**Curve		0x7d86446ddb609ed0f5f8684acf30380a356b2b4c	0x8e764bE4288B842791989DB5b8ec067279829809	ICurveRouter { function exchange(
	**1inch		0x1111111254EEB25477B68fb85Ed929f73A960582	I1inchRouter { function swap
	**Balancer	0x9424B1412450D0f8Fc2255FAf6046b98213B76Bd	0xBA12222222228d8Ba445958a75a0704d566BF2C8	IBalancerVault { struct SwapReques
	**Bancor	0x8eE7D9235e01e6B42345120b5d270bDB763624C7	IBancorNetwork { function convert(
	
	WETH		0xC02aaA39b223FE8D0A0e5C4F27ead9083C756Cc2
	USDC		0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606EB48
	USDT		0xdAC17F958D2ee523a2206206994597C13D831ec7
	DAI			0x6B175474E89094C44Da98b954EedeAC495271d0F
	BUSD		0x4fabb145d64652a948d72533023f6e7a623c7c53
	TUSD		0x0000000000085d4780B73119b644AE5ecd22b376
	
*/

contract Trader {

    // Addresses
    address payable OWNER;
    address NATIVE_TOKEN;	
	
    constructor(address native_token) {
        OWNER = payable(msg.sender);
		NATIVE_TOKEN = native_token;
    }

    modifier onlyOwner() {
        require(msg.sender == OWNER, "caller is not the owner!");
        _;
    }
	
    address[] public dexAddresses; // Array to store dex addresses
    mapping(address => DexInterfaceType) public dexInterface;
	mapping(address => address) public v3quoters;    

    address [] public tokens;
    address [] public stables;
    // Mapping to store pools by token pairs
    mapping(address => mapping(address => mapping(address => uint24))) public tokenV3PoolsFee;
	
    function sortTokens(address tokenA, address tokenB) internal pure returns (address token0, address token1) {
        require(tokenA != tokenB, 'IDENTICAL_ADDRESSES');
        (token0, token1) = tokenA < tokenB ? (tokenA, tokenB) : (tokenB, tokenA);
        require(token0 != address(0), 'ZERO_ADDRESS');
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

    function AddDex(address[] calldata  _dex, DexInterfaceType[] calldata  _interface, address[] calldata  _dexQuoter) public onlyOwner {
        require ( _dex.length == _interface.length, "Invalid param");
        for (uint i=0; i<_dex.length; i++) {
            if (!IsDexAdded(_dex[i])) {
                dexAddresses.push(_dex[i]);
            }
            dexInterface[_dex[i]] = _interface[i];
			if (_dexQuoter[i] != address(0))
				v3quoters[_dex[i]] = _dexQuoter[i];
        }
    }

    function getDexCount() public view returns (uint256) {
        return dexAddresses.length;
    }

    function GetAmountOutMin(routeChain memory _routedata, address _tokenOut, uint256 _amountIn) public view returns (uint256 ) {
        address _tokenIn = _routedata.asset;
		if (_routedata.Itype == DexInterfaceType.IUniswapV4PoolManager) {
			require(_routedata.Itype != DexInterfaceType.IUniswapV4PoolManager, "not implemented here yet");
			return 0;
		} else if (_routedata.Itype == DexInterfaceType.IUniswapV3Router) {
			require(address(0x0) != _tokenIn, "Router does not support direct ETH swap");
			require(address(0x0) != _tokenOut, "Router does not support direct ETH swap");
            uint256 result = IQuoter(v3quoters[_routedata.router]).quoteExactInputSingle(_tokenIn, _tokenOut ,_routedata.poolFee, _amountIn, 0);
            return result;
        } else { // DexInterfaceType.IUniswapV2Router
            //uint256 result = 0;            
            address[] memory path;
            path = new address[](2);
            path[0] = _tokenIn == address(0x0) ? NATIVE_TOKEN : _tokenIn;
            path[1] = _tokenOut == address(0x0) ? NATIVE_TOKEN : _tokenOut;            
            //try IUniswapV2Router(_router).getAmountsOut(_amountIn, path) returns (uint256[] memory amountOutMins) {
            //    result = amountOutMins[path.length-1];
            //} catch {
            //}
            //return result;      			
			uint256[] memory amountOutMins = IUniswapV2Router(_routedata.router).getAmountsOut(_amountIn, path);
			return amountOutMins[path.length-1];      
        }
    }

    function EstimateDualDexTradeGain(routeChain[] calldata _routedata, uint256 _fromAmount) external view returns (uint256) {
		require ( _routedata.length == 2, "Invalid param");
        uint256 amtBack1 = GetAmountOutMin(_routedata[0], _routedata[1].asset, _fromAmount);
        uint256 amtBack2 = GetAmountOutMin(_routedata[1], _routedata[0].asset, amtBack1);
		if (amtBack2 < _fromAmount)
			return 0;
        return amtBack2 - _fromAmount;
    }
  
	
    function AmountBack(
        routeChain[] memory _routedata,
        uint256 amountIn
    ) internal view returns (uint256) {
		require ( _routedata.length == 4, "Invalid param");
		_routedata[0].poolFee = getTestV3PoolFee(_routedata[0].router, _routedata[0].asset, _routedata[1].asset);
        uint256 amtBack = GetAmountOutMin(_routedata[0], _routedata[1].asset, amountIn);
		_routedata[1].poolFee = getTestV3PoolFee(_routedata[1].router, _routedata[1].asset, _routedata[2].asset);
        amtBack = GetAmountOutMin(_routedata[1], _routedata[2].asset, amtBack);
		_routedata[2].poolFee = getTestV3PoolFee(_routedata[2].router, _routedata[2].asset, _routedata[3].asset);
        amtBack = GetAmountOutMin(_routedata[2], _routedata[3].asset, amtBack);
		_routedata[3].poolFee = getTestV3PoolFee(_routedata[3].router, _routedata[3].asset, _routedata[0].asset);
        amtBack = GetAmountOutMin(_routedata[3], _routedata[0].asset, amtBack);
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
								routeChain[] memory _routedata;
								_routedata = new routeChain[](4);
								_routedata[0].router = _router;	
								_routedata[0].asset = _baseAsset;
								_routedata[1].router = _router;	
								_routedata[1].asset = tokens[i1];	
								_routedata[2].router = _router;	
								_routedata[2].asset = stables[i2];	
								_routedata[3].router = _router;	
								_routedata[3].asset = tokens[i3];	
								uint256 amtBack = AmountBack(_routedata, _amount);
								if (amtBack > _amount && amtBack > maxAmtBack) {
									maxAmtBack = amtBack;
									token1 = tokens[i1];
									token2 = stables[i2];
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
