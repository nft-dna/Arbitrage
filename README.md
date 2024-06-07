# ARBITRAGE TOOLS

## Project Setup
```
npm install

Create .env file
TESTNET_PRIVATE_KEY='0000000000000000000000000000000000000000000000000000000000000000'
MAINNET_PRIVATE_KEY='0000000000000000000000000000000000000000000000000000000000000000'
```

## Compile
```
npx hardhat compile
```

## Test
```
npx hardhat test ./test/HardhatOverAll.test.js

SEPOLIA: 0x2dbFb325fDa043085DC6005072AEBCa6c8A18dA9
	0x425141165d3DE9FEC831896C016617a52363b687 V2Router02 Contract Address
	0x3bFA4769FB09eefC5a80d6E87c3B9C650f7Ae48E V3 SwapRouter02
	0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14 WETH
	0x0000000000000000000000000000000000000000
	0xb4f1737af37711e9a5890d9510c9bb60e170cb0d DAI
	
	InstaTrade (routeChain[] calldata _routedata, uint256 _startAmount, uint deadlineDeltaSec) payable
	call example on Remix
	_routedata[] routeChain is
		DexInterfaceType Itype;
		address router;
		address asset;
		uint24 poolFee;
		int24 tickSpacing;
	[[0, "0x425141165d3DE9FEC831896C016617a52363b687", "0x0000000000000000000000000000000000000000", 0, 0],[1, "0x3bFA4769FB09eefC5a80d6E87c3B9C650f7Ae48E", "0xb4f1737af37711e9a5890d9510c9bb60e170cb0d", 3000, 0]]
	_startAmount (same as paybale value) i.e 3000000000000000 (0,003 ETH)
	deadlineDeltaSec 0
```
