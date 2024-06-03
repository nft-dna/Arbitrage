// npx hardhat test ./test/HardhatOverAll.test.js --network localhost; 
// run first (in another shell): npx hardhat node

const { expect } = require("chai");
const { ethers } = require("hardhat");

const weiToEther = (n) => {
    return web3.utils.fromWei(n.toString(), 'ether');
}

async function getGasCosts(receipt) {
    const tx = await web3.eth.getTransaction(receipt.tx);  
    const gasPrice = new BN(tx.gasPrice);
    return gasPrice.mul(new BN(receipt.receipt.gasUsed));
}    

describe("Overall Test", function () {	
		
	//let Trade;
	let Trader;
	let Traderaddr;
	let MockDex;
	let dex1;
	let dex1addr;
	let dex2;
	let dex2addr;
	let MockERC20;
	let tokenA;
	let tokenAaddr;
	let tokenB;
	let tokenBaddr;
	let tokenC;
	let tokenCaddr;
	let owner;
	let addr1;
	let addr2;
	let initialSupply = ethers.parseEther("1000");
	let initialDexReserve = ethers.parseEther("300");
  
	beforeEach(async function () {
		[owner, addr1, addr2] = await ethers.getSigners();	
		
		const Trade = await ethers.getContractFactory("Trade");
		Trader = await Trade.deploy();
		await Trader.waitForDeployment();
		Traderaddr = await Trader.getAddress();
		
		MockDex = await ethers.getContractFactory("MockDEX");	
		dex1 = await MockDex.deploy();
		await dex1.waitForDeployment();	
		dex2 = await MockDex.deploy();
		await dex2.waitForDeployment();
		
		MockERC20 = await ethers.getContractFactory("MockERC20");
		tokenA = await MockERC20.deploy("MockTokenA", "MTKA", 18, initialSupply);
		await tokenA.waitForDeployment();
		tokenAaddr = await tokenA.getAddress();
		tokenB = await MockERC20.deploy("MockTokenB", "MTKB", 18, initialSupply);
		await tokenB.waitForDeployment();
		tokenBaddr = await tokenB.getAddress();
		tokenC = await MockERC20.deploy("MockTokenC", "MTKC", 18, initialSupply);
		await tokenC.waitForDeployment();	
		tokenCaddr = await tokenC.getAddress();
		
		dex1addr = await dex1.getAddress();
		dex2addr = await dex2.getAddress();
		tokenA.transfer(dex1addr, initialDexReserve);
		tokenA.transfer(dex2addr, initialDexReserve);
		tokenB.transfer(dex1addr, initialDexReserve);
		tokenB.transfer(dex2addr, initialDexReserve);
		tokenC.transfer(dex1addr, initialDexReserve);
		tokenC.transfer(dex2addr, initialDexReserve);	
		await owner.sendTransaction({
		  to: dex1addr,
		  value: ethers.parseEther("3.0"),
		});		
		await owner.sendTransaction({
		  to: dex2addr,
		  value: ethers.parseEther("3.0"),
		});
		
		//tokenA.transfer(Traderaddr, initialDexReserve);
		await tokenA.approve(Traderaddr, initialDexReserve);
		await Trader.depositToken(tokenA, initialDexReserve);
		//tokenB.transfer(Traderaddr, initialDexReserve);
		await tokenB.approve(Traderaddr, initialDexReserve);
		await Trader.depositToken(tokenB, initialDexReserve);		
		//tokenC.transfer(Traderaddr, initialDexReserve);	
		await tokenC.approve(Traderaddr, initialDexReserve);
		await Trader.depositToken(tokenC, initialDexReserve);
		//await owner.sendTransaction({
		//  to: Traderaddr,
		//  value: ethers.parseEther("3.0"),
		//});
		await Trader.depositEther({value: ethers.parseEther("3.0")});

		
		//fee: The fee tier of the pool (e.g., 500, 3000, 10000 for 0.05%, 0.3%, 1% respectively).
		await dex1.setPrice(tokenAaddr, tokenBaddr, ethers.parseEther("0.1"));
		await dex1.setFee(tokenAaddr, tokenBaddr, 3000);
		await dex1.setPrice(tokenBaddr, tokenAaddr, ethers.parseEther("0.1"));
		await dex1.setFee(tokenBaddr, tokenAaddr, 3000);
		await dex1.setPrice(tokenAaddr, tokenCaddr, ethers.parseEther("0.1"));
		await dex1.setFee(tokenAaddr, tokenCaddr, 3000);
		await dex1.setPrice(tokenCaddr, tokenAaddr, ethers.parseEther("0.1"));
		await dex1.setFee(tokenCaddr, tokenAaddr, 3000);
		await dex1.setPrice(tokenBaddr, tokenCaddr, ethers.parseEther("0.1"));
		await dex1.setFee(tokenBaddr, tokenCaddr, 3000);
		await dex1.setPrice(tokenCaddr, tokenBaddr, ethers.parseEther("0.1"));
		await dex1.setFee(tokenCaddr, tokenBaddr, 3000);

		await dex2.setPrice(tokenAaddr, tokenBaddr, ethers.parseEther("0.1"));
		await dex2.setFee(tokenAaddr, tokenBaddr, 3000);
		await dex2.setPrice(tokenBaddr, tokenAaddr, ethers.parseEther("0.1"));
		await dex2.setFee(tokenBaddr, tokenAaddr, 3000);
		await dex2.setPrice(tokenAaddr, tokenCaddr, ethers.parseEther("0.1"));
		await dex2.setFee(tokenAaddr, tokenCaddr, 3000);
		await dex2.setPrice(tokenCaddr, tokenAaddr, ethers.parseEther("0.1"));
		await dex2.setFee(tokenCaddr, tokenAaddr, 3000);
		await dex2.setPrice(tokenBaddr, tokenCaddr, ethers.parseEther("0.1"));
		await dex2.setFee(tokenBaddr, tokenCaddr, 3000);
		await dex2.setPrice(tokenCaddr, tokenBaddr, ethers.parseEther("0.1"));
		await dex2.setFee(tokenCaddr, tokenBaddr, 3000);		
	});
		
	describe("MockERC20", async function () {
	  it("Should return the right name and symbol", async function () {
	  	//console.log(await token.getAddress());
		expect(await tokenA.name()).to.equal("MockTokenA");
		expect(await tokenA.symbol()).to.equal("MTKA");
		//expect(await tokenA.balanceOf(owner.address)).to.equal(ethers.parseEther("1000"));
	  });
	});
	
	describe("DexSetup", async function () {
	  it("Set Dex and Token references", async function () {
		 
		const availNative = await Trader.getEtherBalance();
		//console.log(availNative);
		const availTokenA = await Trader.getTokenBalance(tokenAaddr);
		//console.log(availTokenA);
		const availTokenB = await Trader.getTokenBalance(tokenBaddr);
		//console.log(availTokenB);
		const availTokenC = await Trader.getTokenBalance(tokenCaddr);
		//console.log(availTokenC);

		await Trader.AddDex([dex1addr, dex2addr], [/*DexInterfaceType.IUniswapV2Router*/0, /*DexInterfaceType.IUniswapV3Router*/1]);
		await Trader.AddTestTokens([tokenAaddr, tokenBaddr]);
		await Trader.AddTestStables([tokenCaddr]);
		await Trader.AddTestV3PoolFee(tokenAaddr, tokenBaddr, 3000);
		
		const amountOutMin = await Trader.getAmountOutMin(dex1addr, 0, tokenAaddr, tokenBaddr, ethers.parseEther("0.05"));
		//console.log(amountOutMin);

		//const estimate = await Trader.EstimateDualDexTrade(tokenAaddr, tokenBaddr, dex1addr, 0, dex2addr,3000, ethers.parseEther("0.1"));
	  	//console.log(estimate);
		
		//await Trader.InstaSearch(dex1addr, tokenAaddr, ethers.parseEther("0.05"));
		//InstaTradeTokens(address _router1, address _baseAsset, address _token2, address _token3, address _token4, uint256 _amount, uint deadlineDeltaSec)	
	  });
	});	
});

