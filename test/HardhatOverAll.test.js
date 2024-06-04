// npx hardhat test ./test/HardhatOverAll.test.js --network localhost; 
// run first (in another shell): npx hardhat node

const { expect } = require("chai");
const { ethers } = require("hardhat");

/*
const weiToEther = (n) => {
    return web3.utils.fromWei(n.toString(), 'ether');
}

async function getGasCosts(receipt) {
    const tx = await web3.eth.getTransaction(receipt.tx);  
    const gasPrice = new BN(tx.gasPrice);
    return gasPrice.mul(new BN(receipt.receipt.gasUsed));
}    
*/

describe("Overall Test", function () {	
		
	let Trade;
	let TradeAddr;
	let Trader;
	let TraderAddr;
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
	let initialSupply = 3000000;
	let initialDexReserve = 1000000;
	// Mock price (amount of tokenOut per 1 unit of tokenIn)
	let initialPrice = 1;
	//fee: The fee tier of the pool (e.g., 500, 3000, 10000 for 0.05%, 0.3%, 1% respectively).	
	let initialFee = 0;

	const NATIVE_TOKEN = "0x0000000000000000000000000000000000000000";
  
	beforeEach(async function () {
		[owner, addr1, addr2] = await ethers.getSigners();	
		
		const trade = await ethers.getContractFactory("Trade");
		Trade = await trade.deploy();
		await Trade.waitForDeployment();
		TradeAddr = await Trade.getAddress();
		
		const trader = await ethers.getContractFactory("Trader");
		Trader = await trader.deploy();
		await Trader.waitForDeployment();
		TraderAddr = await Trader.getAddress();
		
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
		
		//tokenA.transfer(TradeAddr, initialDexReserve);
		await tokenA.approve(TradeAddr, initialDexReserve);
		await Trade.depositToken(tokenA, initialDexReserve);
		//tokenB.transfer(TradeAddr, initialDexReserve);
		await tokenB.approve(TradeAddr, initialDexReserve);
		await Trade.depositToken(tokenB, initialDexReserve);		
		//tokenC.transfer(TradeAddr, initialDexReserve);	
		await tokenC.approve(TradeAddr, initialDexReserve);
		await Trade.depositToken(tokenC, initialDexReserve);
		//await owner.sendTransaction({
		//  to: TradeAddr,
		//  value: ethers.parseEther("3.0"),
		//});
		await Trade.depositEther({value: ethers.parseEther("3.0")});

		await dex1.setPairInfo(tokenAaddr, tokenBaddr, initialPrice, initialFee);
		await dex1.setPairInfo(tokenBaddr, tokenCaddr, initialPrice, initialFee);
		await dex1.setPairInfo(tokenCaddr, tokenAaddr, initialPrice, initialFee);		

		await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice, initialFee);
		await dex2.setPairInfo(tokenBaddr, tokenCaddr, initialPrice, initialFee);
		await dex2.setPairInfo(tokenCaddr, tokenAaddr, initialPrice, initialFee);	

		await dex1.setPairInfo(NATIVE_TOKEN, tokenAaddr, initialPrice, initialFee);
		await dex1.setPairInfo(NATIVE_TOKEN, tokenBaddr, initialPrice, initialFee);
		await dex1.setPairInfo(NATIVE_TOKEN, tokenCaddr, initialPrice, initialFee);
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
		 
		const availNative = await Trade.getEtherBalance();
		//console.log(availNative);
		const availTokenA = await Trade.getTokenBalance(tokenAaddr);
		//console.log(availTokenA);
		const availTokenB = await Trade.getTokenBalance(tokenBaddr);
		//console.log(availTokenB);
		const availTokenC = await Trade.getTokenBalance(tokenCaddr);
		//console.log(availTokenC);

		await Trader.AddDex([dex1addr, dex2addr], [/*DexInterfaceType.IUniswapV2Router*/0, /*DexInterfaceType.IUniswapV3Router*/1]);
		await Trader.AddTestTokens([tokenAaddr, tokenBaddr]);
		await Trader.AddTestStables([tokenCaddr]);
		await Trader.AddTestV3PoolFee(dex2addr, tokenAaddr, tokenBaddr, initialFee);
		
		const amtBack1 = await Trader.GetAmountOutMin(dex1addr, 0, tokenAaddr, tokenBaddr, initialDexReserve);
		expect(amtBack1.toString()).to.be.equal(initialDexReserve.toString());
        const amtBack2 = await Trader.GetAmountOutMin(dex2addr, initialFee, tokenBaddr, tokenAaddr, initialDexReserve);
		expect(amtBack2.toString()).to.be.equal(initialDexReserve.toString());
		const estimateDex1 = await Trader.EstimateDualDexTrade(tokenAaddr, tokenBaddr, dex1addr, 0, dex2addr, initialFee, initialDexReserve);
		expect(estimateDex1.toString()).to.be.equal(initialDexReserve.toString());
		
		await expect(
		  Trader.EstimateDualDexTrade(NATIVE_TOKEN, tokenCaddr, dex1addr, 0, dex2addr, initialFee, initialDexReserve)
		).to.be.revertedWith("Router does not support direct ETH swap");		
		
		const estimateDex1Native = await Trader.EstimateDualDexTrade(NATIVE_TOKEN, tokenCaddr, dex1addr, 0, dex1addr, 0, initialDexReserve);
		expect(estimateDex1Native.toString()).to.be.equal(initialDexReserve.toString());	
				
		const searchDex1 = await Trader.CrossStableSearch(dex1addr, tokenAaddr, ethers.parseEther("0.05"));
		expect(searchDex1[0].toString()).to.be.equal("0");
		expect(searchDex1[1].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex1[2].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex1[3].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
			
		await expect(
		  Trader.EstimateDualDexTrade(NATIVE_TOKEN, tokenCaddr, dex2addr, 0, dex1addr, 0, initialDexReserve)
		).to.be.revertedWith("Router does not support direct ETH swap");		
						
		const estimateDex2Native = await Trader.EstimateDualDexTrade(NATIVE_TOKEN, tokenCaddr, dex1addr, 0, dex1addr, 0, initialDexReserve);
		expect(estimateDex2Native.toString()).to.be.equal(initialDexReserve.toString());			
		
		const searchDex2 = await Trader.CrossStableSearch(dex2addr, tokenAaddr, ethers.parseEther("0.05"));
		expect(searchDex2[0].toString()).to.be.equal("0");
		expect(searchDex2[1].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex2[2].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex2[3].toString()).to.be.equal("0x0000000000000000000000000000000000000000");		
		
		//const instrade = await Trader.InstaTradeTokens(dex1addr, tokenAaddr, tokenCaddr, tokenBaddr, tokenAaddr, ethers.parseEther("0.1"), 0);
	  });
	});	
});

