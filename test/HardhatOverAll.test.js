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
	let initialSupply = 3000000;
	let initialDexReserve = 1000000;
  
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
		await dex1.setPrice(tokenAaddr, tokenBaddr, ethers.parseEther("0.0001"));
		await dex1.setFee(tokenAaddr, tokenBaddr, 3000);
		await dex1.setPrice(tokenBaddr, tokenAaddr, ethers.parseEther("0.0001"));
		await dex1.setFee(tokenBaddr, tokenAaddr, 3000);
		await dex1.setPrice(tokenAaddr, tokenCaddr, ethers.parseEther("0.0001"));
		await dex1.setFee(tokenAaddr, tokenCaddr, 3000);
		await dex1.setPrice(tokenCaddr, tokenAaddr, ethers.parseEther("0.0001"));
		await dex1.setFee(tokenCaddr, tokenAaddr, 3000);
		await dex1.setPrice(tokenBaddr, tokenCaddr, ethers.parseEther("0.0001"));
		await dex1.setFee(tokenBaddr, tokenCaddr, 3000);
		await dex1.setPrice(tokenCaddr, tokenBaddr, ethers.parseEther("0.0001"));
		await dex1.setFee(tokenCaddr, tokenBaddr, 3000);

		await dex2.setPrice(tokenAaddr, tokenBaddr, ethers.parseEther("0.0001"));
		await dex2.setFee(tokenAaddr, tokenBaddr, 3000);
		await dex2.setPrice(tokenBaddr, tokenAaddr, ethers.parseEther("0.0001"));
		await dex2.setFee(tokenBaddr, tokenAaddr, 3000);
		await dex2.setPrice(tokenAaddr, tokenCaddr, ethers.parseEther("0.0001"));
		await dex2.setFee(tokenAaddr, tokenCaddr, 3000);
		await dex2.setPrice(tokenCaddr, tokenAaddr, ethers.parseEther("0.0001"));
		await dex2.setFee(tokenCaddr, tokenAaddr, 3000);
		await dex2.setPrice(tokenBaddr, tokenCaddr, ethers.parseEther("0.0001"));
		await dex2.setFee(tokenBaddr, tokenCaddr, 3000);
		await dex2.setPrice(tokenCaddr, tokenBaddr, ethers.parseEther("0.0001"));
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
		
		//const amtBack1 = await Trader.getAmountOutMin(dex1addr, 0, tokenAaddr, tokenBaddr, ethers.parseEther("0.1"));
		//console.log(amtBack1);
        //const amtBack2 = await Trader.getAmountOutMin(dex2addr, 3000, tokenBaddr, tokenAaddr, amtBack1);
		//console.log(amtBack1);
		const estimateDex1 = await Trader.EstimateDualDexTrade(tokenAaddr, tokenBaddr, dex1addr, 0, dex2addr, 3000, ethers.parseEther("0.1"));
		expect(estimateDex1.toString()).to.be.equal("0");
		
		const searchDex1 = await Trader.InstaSearch(dex1addr, tokenAaddr, ethers.parseEther("0.05"));
		expect(searchDex1[0].toString()).to.be.equal("0");
		expect(searchDex1[1].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex1[2].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex1[3].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
			
		const estimateDex2 = await Trader.EstimateDualDexTrade(tokenAaddr, tokenBaddr, dex2addr, 3000, dex1addr, 0, ethers.parseEther("0.1"));
		expect(estimateDex2.toString()).to.be.equal("0");		
		
		const searchDex2 = await Trader.InstaSearch(dex2addr, tokenAaddr, ethers.parseEther("0.05"));
		console.log(searchDex2);
		expect(searchDex2[0].toString()).to.be.equal("0");
		//expect(searchDex2[1].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		//expect(searchDex2[2].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		//expect(searchDex2[3].toString()).to.be.equal("0x0000000000000000000000000000000000000000");		
		
		//const instrade = await Trader.InstaTradeTokens(dex1addr, tokenAaddr, tokenCaddr, tokenBaddr, tokenAaddr, ethers.parseEther("0.1"), 0);
	  });
	});	
});

