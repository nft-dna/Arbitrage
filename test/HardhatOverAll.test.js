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
	let initialEthBalance = ethers.parseEther("3.0");
	let initialSupply = ethers.parseEther("10.0"); //= 5000000;
	let initialDexReserve = ethers.parseEther("1.0"); //= 1000000;
	// Mock price (amount of tokenOut per 1 unit of tokenIn in percentage points out of 100)
	// i.e 100 = 1 tokenIn costs as 1 tokenOut
	// i.e  50 = 1 tokenIn costs as half tokenOut
	let initialPrice = 100;
	//fee: The fee tier of the pool (e.g., 500, 3000, 10000 for 0.05%, 0.3%, 1% respectively).	
	let initialFee = 0;

	let NATIVE_TOKEN;
	let ZERO_ADDRESS = ethers.ZeroAddress;
  
	beforeEach(async function () {
		[owner, addr1, addr2] = await ethers.getSigners();	
		
	
		MockERC20 = await ethers.getContractFactory("MockERC20");
		const weth = await MockERC20.deploy("NATIVE_TOKEN", "WETH", 18, initialSupply);
		await weth.waitForDeployment();
		NATIVE_TOKEN = await weth.getAddress();				
		tokenA = await MockERC20.deploy("MockTokenA", "MTKA", 18, initialSupply);
		await tokenA.waitForDeployment();
		tokenAaddr = await tokenA.getAddress();
		tokenB = await MockERC20.deploy("MockTokenB", "MTKB", 18, initialSupply);
		await tokenB.waitForDeployment();
		tokenBaddr = await tokenB.getAddress();
		tokenC = await MockERC20.deploy("MockTokenC", "MTKC", 18, initialSupply);
		await tokenC.waitForDeployment();	
		tokenCaddr = await tokenC.getAddress();
		
		const trade = await ethers.getContractFactory("Trade");
		Trade = await trade.deploy(NATIVE_TOKEN);
		await Trade.waitForDeployment();
		TradeAddr = await Trade.getAddress();
		
		const trader = await ethers.getContractFactory("Trader");
		Trader = await trader.deploy(NATIVE_TOKEN);
		await Trader.waitForDeployment();
		TraderAddr = await Trader.getAddress();		
		
		MockDex = await ethers.getContractFactory("MockDEX");	
		dex1 = await MockDex.deploy(NATIVE_TOKEN);
		await dex1.waitForDeployment();	
		dex2 = await MockDex.deploy(NATIVE_TOKEN);
		await dex2.waitForDeployment();		
		
		dex1addr = await dex1.getAddress();
		dex2addr = await dex2.getAddress();		
		await tokenA.transfer(dex1addr, initialDexReserve);
		await tokenA.transfer(dex2addr, initialDexReserve);
		await tokenB.transfer(dex1addr, initialDexReserve);
		await tokenB.transfer(dex2addr, initialDexReserve);
		await tokenC.transfer(dex1addr, initialDexReserve);
		await tokenC.transfer(dex2addr, initialDexReserve);	
		await owner.sendTransaction({
		  to: dex1addr,
		  value: initialEthBalance,
		});		
		await owner.sendTransaction({
		  to: dex2addr,
		  value: initialEthBalance,
		});
		
		//tokenA.transfer(TradeAddr, initialDexReserve);
		await tokenA.approve(TradeAddr, initialDexReserve);
		await Trade.depositToken(tokenAaddr, initialDexReserve);
		//tokenB.transfer(TradeAddr, initialDexReserve);
		await tokenB.approve(TradeAddr, initialDexReserve);
		await Trade.depositToken(tokenBaddr, initialDexReserve);		
		//tokenC.transfer(TradeAddr, initialDexReserve);	
		await tokenC.approve(TradeAddr, initialDexReserve);
		await Trade.depositToken(tokenCaddr, initialDexReserve);
		//await owner.sendTransaction({
		//  to: TradeAddr,
		//  value: initialEthBalance,
		//});
		await Trade.depositEther({value: initialEthBalance});

		await dex1.setPairInfo(tokenAaddr, tokenBaddr, initialPrice, initialFee);
		await dex1.setPairInfo(tokenBaddr, tokenCaddr, initialPrice, initialFee);
		await dex1.setPairInfo(tokenCaddr, tokenAaddr, initialPrice, initialFee);		

		await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice, initialFee);
		await dex2.setPairInfo(tokenBaddr, tokenCaddr, initialPrice, initialFee);
		await dex2.setPairInfo(tokenCaddr, tokenAaddr, initialPrice, initialFee);	

		await dex1.setPairInfo(NATIVE_TOKEN, tokenAaddr, initialPrice, initialFee);
		await dex1.setPairInfo(NATIVE_TOKEN, tokenBaddr, initialPrice, initialFee);
		await dex1.setPairInfo(NATIVE_TOKEN, tokenCaddr, initialPrice, initialFee);
		
		// dex2 is IUniswapV3, so it does not support direct ETH swap
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
		const estimateDex1 = await Trader.EstimateDualDexTradeGain(tokenAaddr, tokenBaddr, dex1addr, 0, dex2addr, initialFee, initialDexReserve);
		expect(estimateDex1.toString()).to.be.equal("0");

		await expect(
		  Trader.EstimateDualDexTradeGain(ZERO_ADDRESS, tokenCaddr, dex1addr, 0, dex2addr, initialFee, initialDexReserve)
		).to.be.revertedWith("Router does not support direct ETH swap");		

		const estimateDex1Native = await Trader.EstimateDualDexTradeGain(ZERO_ADDRESS, tokenCaddr, dex1addr, 0, dex1addr, 0, initialDexReserve);
		expect(estimateDex1Native.toString()).to.be.equal("0");	
	
		const searchDex1 = await Trader.CrossStableSearch(dex1addr, tokenAaddr, ethers.parseEther("0.05"));
		expect(searchDex1[0].toString()).to.be.equal("0");
		expect(searchDex1[1].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex1[2].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex1[3].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
	
		await expect(
		  Trader.EstimateDualDexTradeGain(ZERO_ADDRESS, tokenCaddr, dex2addr, 0, dex1addr, 0, initialDexReserve)
		).to.be.revertedWith("Router does not support direct ETH swap");		
				
		const estimateDex2Native = await Trader.EstimateDualDexTradeGain(ZERO_ADDRESS, tokenCaddr, dex1addr, 0, dex1addr, 0, initialDexReserve);
		expect(estimateDex2Native.toString()).to.be.equal("0");			

		const searchDex2 = await Trader.CrossStableSearch(dex2addr, tokenAaddr, ethers.parseEther("0.05"));
		expect(searchDex2[0].toString()).to.be.equal("0");
		expect(searchDex2[1].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex2[2].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex2[3].toString()).to.be.equal("0x0000000000000000000000000000000000000000");		
		
		//const instrade = await Trader.InstaTradeTokens(dex1addr, tokenAaddr, tokenCaddr, tokenBaddr, tokenAaddr, ethers.parseEther("0.1"), 0);
	  });
	});	
	
    describe("Ether Deposits and Withdrawals", function () {
        it("Should deposit Ether", async function () {
			// already done in beforeEach
            //await Trade.depositEther({ value: initialEthBalance });
            expect(await Trade.getEtherBalance()).to.equal(initialEthBalance);
        });

        it("Should withdraw Ether", async function () {
			// already done in beforeEach
            //await Trade.depositEther({ value: initialEthBalance });
            await Trade.withdrawEther(initialEthBalance);
            expect(await Trade.getEtherBalance()).to.equal(0);
        });

        it("Should emit DepositEther event", async function () {
            await expect(Trade.depositEther({ value: initialEthBalance }))
                .to.emit(Trade, "DepositEther")
                .withArgs(owner.address, initialEthBalance);
        });

        it("Should emit WithdrawEther event", async function () {
            await Trade.depositEther({ value: initialEthBalance });
            await expect(Trade.withdrawEther(initialEthBalance))
                .to.emit(Trade, "WithdrawEther")
                .withArgs(owner.address, initialEthBalance);
        });

        it("Should revert if Ether balance is insufficient", async function () {
            await expect(Trade.withdrawEther(ethers.parseEther("4.0"))).to.be.revertedWith("Insufficient Ether balance");
        });		
    });

    describe("Token Deposits and Withdrawals", function () {
        it("Should deposit tokens", async function () {
			// already done in beforeEach
            //await tokenA.approve(TradeAddr, initialDexReserve);
            //await Trade.depositToken(tokenAaddr, initialDexReserve);
            expect(await Trade.getTokenBalance(tokenAaddr)).to.equal(initialDexReserve);
        });

        it("Should withdraw tokens", async function () {
			// already done in beforeEach
            //await tokenA.approve(TradeAddr, initialDexReserve);
            //await Trade.depositToken(tokenAaddr, initialDexReserve);
            await Trade.withdrawToken(tokenAaddr, initialDexReserve);
            expect(await Trade.getTokenBalance(tokenAaddr)).to.equal(0);
        });

        it("Should emit DepositToken event", async function () {
            await tokenA.approve(TradeAddr, initialDexReserve);
            await expect(Trade.depositToken(tokenAaddr, initialDexReserve))
                .to.emit(Trade, "DepositToken")
                .withArgs(tokenAaddr, owner.address, initialDexReserve);
        });

        it("Should emit WithdrawToken event", async function () {
            await tokenA.approve(TradeAddr, initialDexReserve);
            await Trade.depositToken(tokenAaddr, initialDexReserve);
            await expect(Trade.withdrawToken(tokenAaddr, initialDexReserve))
                .to.emit(Trade, "WithdrawToken")
                .withArgs(tokenAaddr, owner.address, initialDexReserve);
        });

        it("Should revert if token balance is insufficient", async function () {
            await expect(Trade.withdrawToken(tokenAaddr, BigInt(2)*initialDexReserve)).to.be.revertedWith("Insufficient token balance");
        });
    });
	
    describe("Deposits Security", function () {
        it("Should not allow another user to withdraw Ether balance", async function () {
            await Trade.depositEther({ value: ethers.parseEther("1") });
            await expect(Trade.connect(addr1).withdrawEther(ethers.parseEther("1"))).to.be.revertedWith("Insufficient Ether balance");
        });
		
        it("Should not allow another user to withdraw more than his Ether balance", async function () {
			await Trade.connect(addr1).depositEther({ value: ethers.parseEther("1") });
			expect(await Trade.connect(owner).getTotalEtherBalance()).to.equal(ethers.parseEther("4"));
			expect(await Trade.connect(owner).getEtherBalance()).to.equal(ethers.parseEther("3"));
			expect(await Trade.connect(addr1).getEtherBalance()).to.equal(ethers.parseEther("1"));
            await expect(Trade.connect(addr1).withdrawEther(ethers.parseEther("2"))).to.be.revertedWith("Insufficient Ether balance");
            await expect(Trade.connect(addr1).withdrawEther(ethers.parseEther("1")))
                .to.emit(Trade, "WithdrawEther")
                .withArgs(addr1.address, ethers.parseEther("1"));			
        });		

        it("Should not allow another user to withdraw token balance", async function () {
			await tokenA.transfer(dex1addr, initialDexReserve);
            await tokenA.approve(TradeAddr, initialDexReserve);
            await Trade.depositToken(tokenAaddr, initialDexReserve);
            await expect(Trade.connect(addr1).withdrawToken(tokenAaddr, initialDexReserve)).to.be.revertedWith("Insufficient token balance");
        });
		
        it("Should not allow another user to withdraw more than his token balance", async function () {
			await tokenA.transfer(addr1, initialDexReserve);
			await tokenA.connect(addr1).approve(TradeAddr, initialDexReserve);
			await Trade.connect(addr1).depositToken(tokenAaddr, initialDexReserve);
			expect(await Trade.connect(owner).getTotalTokenBalance(tokenAaddr)).to.equal(initialDexReserve*BigInt(2));
			expect(await Trade.connect(owner).getTokenBalance(tokenAaddr)).to.equal(initialDexReserve);
			expect(await Trade.connect(addr1).getTokenBalance(tokenAaddr)).to.equal(initialDexReserve);
            await expect(Trade.connect(addr1).withdrawToken(tokenAaddr, initialDexReserve*BigInt(2))).to.be.revertedWith("Insufficient token balance");
            await expect(Trade.connect(addr1).withdrawToken(tokenAaddr, initialDexReserve))
                .to.emit(Trade, "WithdrawToken")
                .withArgs(tokenAaddr, addr1.address, initialDexReserve);			
        });				
    });	

	describe("Trade Safety Functions Security", function () {
        it("Recover Ether balance", async function () {
            const balance = await Trade.getTotalEtherBalance();
            await Trade.connect(owner).safeWithdrawEther(balance);
			expect(await Trade.getTotalEtherBalance()).to.equal(0);
        });
		
        it("No other could recover Ether balance", async function () {
            const balance = await Trade.getTotalEtherBalance();
            await expect(Trade.connect(addr1).safeWithdrawEther(balance)).to.be.revertedWith("caller is not the owner!");
			expect(await Trade.getTotalEtherBalance()).to.equal(balance);
        });		
		
        it("Recover token balance", async function () {
            const balance = await Trade.getTotalTokenBalance(tokenAaddr);
            await Trade.connect(owner).safeWithdrawToken(tokenAaddr, balance);
			expect(await Trade.getTotalTokenBalance(tokenAaddr)).to.equal(0);
        });
		
        it("No other could recover token balance", async function () {
            const balance = await Trade.getTotalTokenBalance(tokenAaddr);
            await expect(Trade.connect(addr1).safeWithdrawToken(tokenAaddr, balance)).to.be.revertedWith("caller is not the owner!");
			expect(await Trade.getTotalTokenBalance(tokenAaddr)).to.equal(balance);
        });			
		
    });	
	
    describe("DualDexTrade Function", function () {
        it("Should perform a token DualDexTrade", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			// add tokenB extra reserve for dex2
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, initialFee);
			const shouldGainAmount = initialDexReserve;
			
            await expect(Trade.connect(addr1).DualDexTrade(tokenBaddr, tokenAaddr, dex1addr, initialFee, dex2addr, initialFee, initialDexReserve, 0))
                .to.emit(Trade, "DualDexTraded")
                .withArgs(addr1.address, tokenBaddr, tokenAaddr, dex1addr, dex2addr, initialDexReserve, shouldGainAmount);			
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a token DualDexTrade with a loss", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
								
			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr);
			//console.log(initialBalance);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			//console.log(initialTotalBalance);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			
			
			// add tokenB extra reserve for dex2
			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, initialFee);
			const shouldGainAmount = initialDexReserve;
			
            await expect(Trade.connect(addr1).DualDexTrade(tokenAaddr, tokenBaddr, dex1addr, initialFee, dex2addr, initialFee, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
			//await Trade.connect(addr1).DualDexTrade(tokenAaddr, tokenBaddr, dex1addr, initialFee, dex2addr, initialFee, initialDexReserve, 0);
			
			//const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr);
			//console.log(tokenBalance);
			//const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			//console.log(tokenTotalBalance);
        });		
		
        it("Should revert a token DualDexTrade with 0 gain", async function () {
			
            await expect(Trade.connect(addr1).DualDexTrade(tokenBaddr, tokenAaddr, dex1addr, initialFee, dex2addr, initialFee, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });
		

		it("Should perform an Ether DualDexTrade (without payable call)", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			const initialBalance = await Trade.connect(addr1).getEtherBalance();
			const initialTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			expect(initialTotalBalance).to.be.equal(initialEthBalance);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, 2*initialPrice, initialFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, initialFee);
			const shouldGainAmount = initialEthBalance;
					
			const amtBack = await Trader.GetAmountOutMin(dex1addr, 0, ZERO_ADDRESS, tokenAaddr, initialEthBalance);
			const finalEthBalance = await Trader.GetAmountOutMin(dex2addr, 0, tokenAaddr, ZERO_ADDRESS, amtBack);
			expect(finalEthBalance).to.be.equal(initialEthBalance + shouldGainAmount);
			
			await tokenA.transfer(dex1addr, amtBack);	
			
			await expect(Trade.connect(addr1).DualDexTrade(ZERO_ADDRESS, tokenAaddr, dex1addr, initialFee, dex2addr, initialFee, initialEthBalance, 0))
                .to.emit(Trade, "DualDexTraded")
                .withArgs(addr1.address, ZERO_ADDRESS, tokenAaddr, dex1addr, dex2addr, initialEthBalance, shouldGainAmount);			
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const finalBalance = await Trade.connect(addr1).getEtherBalance();
            expect(finalBalance).to.be.above(initialBalance);
			expect(finalBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const finalTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			expect(finalTotalBalance).to.be.above(finalBalance);
			expect(finalTotalBalance).to.be.above(initialTotalBalance);			
        });

        /*
		it("Should perform an Ether DualDexTrade (with payable call)", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			const initialBalance = await Trade.connect(addr1).getEtherBalance();
			const initialTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			expect(initialTotalBalance).to.be.equal(0);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await dex2.setPairInfo(NATIVE_TOKEN, tokenAaddr, 2*initialPrice, initialEthBalance);
			const shouldGainAmount = initialEthBalance;
			
            await expect(Trade.connect(addr1).DualDexTrade({value: initialEthBalance} NATIVE_TOKEN, tokenAaddr, dex1addr, initialFee, dex2addr, initialFee, initialEthBalance, 0))
                .to.emit(Trade, "DualDexTraded")
                .withArgs(addr1.address, NATIVE_TOKEN, tokenAaddr, dex1addr, dex2addr, initialEthBalance, shouldGainAmount);			
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const finalBalance = await Trade.connect(addr1).getEtherBalance();
            expect(finalBalance).to.be.above(initialBalance);
			expect(finalBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const finalTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			expect(finalTotalBalance).to.be.above(finalBalance);
			expect(finalTotalBalance).to.be.above(initialTotalBalance);			
        });	
		*/

		
    });

    /*
	describe("InstaTradeTokens Function", function () {
        it("Should perform an InstaTradeTokens", async function () {
            // This test assumes mock routers and tokens are set up to simulate multi-hop trades
            await Trade.connect(addr1).depositEther({ value: ethers.parseEther("1") });

            // Mock the route data
            // This is a placeholder for actual swap logic; adapt as needed
            const routeData = [
                { router: dex1addr, asset: NATIVE_TOKEN, poolFee: initialFee },
                { router: dex1addr, asset: tokenAaddr, poolFee: initialFee },
                { router: dex2addr, asset: NATIVE_TOKEN, poolFee: initialFee },
            ];

            await Trade.connect(addr1).InstaTradeTokens(routeData, ethers.parseEther("1"), 0);

            // Check the results (mock expected behavior)
            const etherBalance = await Trade.connect(addr1).getEtherBalance();
            expect(etherBalance).to.be.above(ethers.parseEther("1"));
        });
    });	
	*/
	
});

