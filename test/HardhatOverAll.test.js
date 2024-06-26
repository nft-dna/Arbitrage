// npx hardhat test ./test/HardhatOverAll.test.js --network localhost; 
// run first (in another shell): npx hardhat node

const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Overall Test", function () {	
		
	let nativeToken;
	let NATIVE_TOKEN;
	let ZERO_ADDRESS = ethers.ZeroAddress;		
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
	let initialSupply = ethers.parseEther("10.0");
	let initialDexReserve = ethers.parseEther("1.0");
	// Mock price (amount of tokenOut per 1 unit of tokenIn in percentage points out of 100)
	// i.e 100 = 1 tokenIn costs as 1 tokenOut
	// i.e  50 = 1 tokenIn costs as half tokenOut
	let initialPrice = 100;
	//fee: The fee tier of the pool (e.g., 500, 3000, 10000 for 0.05%, 0.3%, 1% respectively).	
	let poolFee	= 0; // 3000;
	//tickSpacing: common tick spacing vlaues are 10, 60, 200
	let IU_V2_POOL = 0	
	let IU_V3_Q1_POOL = 1
	let IU_V3_Q2_POOL = 2	
	let IU_V4_POOL = 3
  
	beforeEach(async function () {
		[owner, addr1, addr2] = await ethers.getSigners();	
		
		//console.log(await ethers.provider.getBalance(owner));
		
		MockERC20 = await ethers.getContractFactory("MockERC20");
		nativeToken = await MockERC20.deploy("NATIVE_TOKEN", "WETH", 18, initialSupply);
		await nativeToken.waitForDeployment();
		NATIVE_TOKEN = await nativeToken.getAddress();				
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
		
		await tokenA.approve(TradeAddr, initialDexReserve);
		await Trade.depositToken(tokenAaddr, initialDexReserve);
		await tokenB.approve(TradeAddr, initialDexReserve);
		await Trade.depositToken(tokenBaddr, initialDexReserve);		
		await tokenC.approve(TradeAddr, initialDexReserve);
		await Trade.depositToken(tokenCaddr, initialDexReserve);

		//await Trade.depositEther({value: initialEthBalance});

		await dex1.setPairInfo(tokenAaddr, tokenBaddr, initialPrice, poolFee);
		await dex1.setPairInfo(tokenBaddr, tokenCaddr, initialPrice, poolFee);
		await dex1.setPairInfo(tokenCaddr, tokenAaddr, initialPrice, poolFee);		

		await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice, poolFee);
		await dex2.setPairInfo(tokenBaddr, tokenCaddr, initialPrice, poolFee);
		await dex2.setPairInfo(tokenCaddr, tokenAaddr, initialPrice, poolFee);	

		await dex1.setPairInfo(NATIVE_TOKEN, tokenAaddr, initialPrice, poolFee);
		await dex1.setPairInfo(NATIVE_TOKEN, tokenBaddr, initialPrice, poolFee);
		await dex1.setPairInfo(NATIVE_TOKEN, tokenCaddr, initialPrice, poolFee);
		
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
		 
		const availNative = await ethers.provider.getBalance(owner); // Trade.getEtherBalance();
		//console.log(availNative);
		const availTokenA = await Trade.getTokenBalance(tokenAaddr, owner);
		//console.log(availTokenA);
		const availTokenB = await Trade.getTokenBalance(tokenBaddr, owner);
		//console.log(availTokenB);
		const availTokenC = await Trade.getTokenBalance(tokenCaddr, owner);
		//console.log(availTokenC);

		await Trader.AddDex([dex1addr, dex2addr], [IU_V2_POOL, IU_V3_Q1_POOL],[dex1addr, dex2addr]);
		await Trader.AddTestTokens([tokenAaddr, tokenBaddr]);
		await Trader.AddTestStables([tokenCaddr]);
		await Trader.AddTestV3PoolFee(dex2addr, tokenAaddr, tokenBaddr, poolFee);
		
		const route1 = { Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
		const amtBack1 = BigInt(await Trader.GetAmountOutMin.staticCallResult(route1, tokenBaddr, initialDexReserve));	
		expect(amtBack1.toString()).to.be.equal(initialDexReserve.toString());
		const route2 = { Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }		
        const amtBack2 = BigInt(await Trader.GetAmountOutMin.staticCallResult(route2, tokenAaddr, initialDexReserve));
		expect(amtBack2.toString()).to.be.equal(initialDexReserve.toString());
		const routeData1 = [
			{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 },
			{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }
		]		
		const estimateDex1 = BigInt(await Trader.EstimateDualDexTradeGain.staticCallResult(routeData1, initialDexReserve));
		expect(estimateDex1.toString()).to.be.equal("0");

		const routeData2 = [
			{ Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
			{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenCaddr, poolFee: poolFee, tickSpacing: 0 }
		]	
		await expect(
		  Trader.EstimateDualDexTradeGain.staticCallResult(routeData2, initialDexReserve)
		).to.be.revertedWith("Router does not support direct ETH swap");		

		const routeData3 = [
			{ Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
			{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenCaddr, poolFee: poolFee, tickSpacing: 0 }
		]	
		const estimateDex1Native = BigInt(await Trader.EstimateDualDexTradeGain.staticCallResult(routeData3, initialDexReserve));
		expect(estimateDex1Native.toString()).to.be.equal("0");	
	
		const searchDex1 = await Trader.CrossStableSearch.staticCallResult(dex1addr, tokenAaddr, ethers.parseEther("0.05"));
		expect(searchDex1[0].toString()).to.be.equal("0");
		expect(searchDex1[1].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex1[2].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex1[3].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
	
		const routeData4 = [
			{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
			{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenCaddr, poolFee: poolFee, tickSpacing: 0 }
		]			
		await expect(
		  Trader.EstimateDualDexTradeGain.staticCallResult(routeData4, initialDexReserve)
		).to.be.revertedWith("Router does not support direct ETH swap");		
				
		const routeData5 = [
			{ Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
			{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenCaddr, poolFee: poolFee, tickSpacing: 0 }
		]							
		const estimateDex2Native = BigInt(await Trader.EstimateDualDexTradeGain.staticCallResult(routeData5, initialDexReserve));
		expect(estimateDex2Native.toString()).to.be.equal("0");			

		const searchDex2 = await Trader.CrossStableSearch.staticCallResult(dex2addr, tokenAaddr, ethers.parseEther("0.05"));
		expect(searchDex2[0].toString()).to.be.equal("0");
		expect(searchDex2[1].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex2[2].toString()).to.be.equal("0x0000000000000000000000000000000000000000");
		expect(searchDex2[3].toString()).to.be.equal("0x0000000000000000000000000000000000000000");		
		
		//const instrade = await Trader.InstaTradeTokens(dex1addr, tokenAaddr, tokenCaddr, tokenBaddr, tokenAaddr, ethers.parseEther("0.1"), 0);
	  });
	});	
	
    describe("Token Deposits and Withdrawals", function () {
        it("Should deposit tokens", async function () {
			// already done in beforeEach
            //await tokenA.approve(TradeAddr, initialDexReserve);
            //await Trade.depositToken(tokenAaddr, initialDexReserve);
            expect(await Trade.getTokenBalance(tokenAaddr, owner)).to.equal(initialDexReserve);
        });

        it("Should withdraw tokens", async function () {
			// already done in beforeEach
            //await tokenA.approve(TradeAddr, initialDexReserve);
            //await Trade.depositToken(tokenAaddr, initialDexReserve);
            await Trade.withdrawToken(tokenAaddr, initialDexReserve);
            expect(await Trade.getTokenBalance(tokenAaddr, owner)).to.equal(0);
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
        /*
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
		*/

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
			expect(await Trade.connect(owner).getTokenBalance(tokenAaddr, owner)).to.equal(initialDexReserve);
			expect(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1)).to.equal(initialDexReserve);
            await expect(Trade.connect(addr1).withdrawToken(tokenAaddr, initialDexReserve*BigInt(2))).to.be.revertedWith("Insufficient token balance");
            await expect(Trade.connect(addr1).withdrawToken(tokenAaddr, initialDexReserve))
                .to.emit(Trade, "WithdrawToken")
                .withArgs(tokenAaddr, addr1.address, initialDexReserve);			
        });				
    });	

	describe("Trade Safety Functions Security", function () {
        /*
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
		*/		
		
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
	
    describe("InstaTradeTokens UniswapV2 Function", function () {
        it("Should perform a token InstaTradeTokens", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			// add tokenB extra reserve for dex2
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]							
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialDexReserve, shouldGainAmount);			
				;
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a token InstaTradeTokens with a loss", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
								
			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			//console.log(initialBalance);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			//console.log(initialTotalBalance);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			
			
			// add tokenB extra reserve for dex2
			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }
			]					
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
			//await Trade.connect(addr1).InstaTradeTokens(tokenAaddr, tokenBaddr, dex1addr, V2_NO_POOL_FEE, dex2addr, V2_NO_POOL_FEE, initialDexReserve, 0);
			
			//const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			//console.log(tokenBalance);
			//const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			//console.log(tokenTotalBalance);
        });		
		
        it("Should revert a token InstaTradeTokens with 0 gain", async function () {
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });
		

		it("Should perform an Ether InstaTradeTokens (without payable call)", async function () {
			
			await nativeToken.connect(addr1).deposit({value: initialEthBalance});
			await nativeToken.connect(addr1).approve(TradeAddr, initialEthBalance);				
			await Trade.connect(addr1).depositToken(NATIVE_TOKEN, initialEthBalance);
            
			// should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			const initialBalance = await ethers.provider.getBalance(addr1); // Trade.connect(addr1).getEtherBalance();
			//const initialTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			//expect(initialTotalBalance).to.be.equal(initialEthBalance);
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			//const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			const initialBalanceN = await Trade.connect(addr1).getTokenBalance(NATIVE_TOKEN, addr1);

			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, 2*initialPrice, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			const shouldGainAmount = initialEthBalance;
			const route1 = { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 }
			const route2 = { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			const amtBack = BigInt(await Trader.GetAmountOutMin.staticCallResult(route1, route2.asset, initialEthBalance));
			const finalEthBalance = BigInt(await Trader.GetAmountOutMin.staticCallResult(route2, route1.asset, amtBack));
			expect(finalEthBalance).to.be.equal(initialEthBalance + shouldGainAmount);
					
			await tokenA.transfer(dex1addr, amtBack);	
			/*
			await owner.sendTransaction({
			  to: dex2addr,
			  value: initialEthBalance,
			});
			*/
			await nativeToken.deposit({value: BigInt(2)*initialEthBalance});
			await nativeToken.transfer(dex2addr, BigInt(2)*initialEthBalance);
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]							
			await /*expect(*/Trade.connect(addr1).InstaTradeTokens(routeData1, initialEthBalance, 0)//)
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialEthBalance, shouldGainAmount);			
				;
				
			await Trade.connect(addr1).withdrawToken(NATIVE_TOKEN, shouldGainAmount);
			await nativeToken.connect(addr1).withdraw(shouldGainAmount);
				
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			//expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr
			const finalBalanceN  = await Trade.connect(addr1).getTokenBalance(NATIVE_TOKEN, addr1);
			const finalBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();

            expect(finalBalance).to.be.above(initialBalance);
			//expect(finalBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			//const finalTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			//expect(finalTotalBalance).to.be.above(finalBalance);
			//expect(finalTotalBalance).to.be.above(initialTotalBalance);			
        });

		it("Should perform an Ether InstaTradeTokens (with payable call)", async function () {

			//const depositBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();	
			//await Trade.connect(addr1).withdrawEther(depositBalance);		
			const initialBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();			
			//expect(initialBalance).to.be.equal(0);
			//const initialTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			//expect(initialTotalBalance).to.be.equal(initialEthBalance);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, 2*initialPrice, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			const shouldGainAmount = initialEthBalance;
					
			const route1 = { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 }
			const route2 = { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }					
			const amtBack = BigInt(await Trader.GetAmountOutMin.staticCallResult(route1, route2.asset, initialEthBalance));
			const finalEthBalance = BigInt(await Trader.GetAmountOutMin.staticCallResult(route2, route1.asset, amtBack));
			expect(finalEthBalance).to.be.equal(initialEthBalance + shouldGainAmount);
			
			await tokenA.transfer(dex1addr, amtBack);	
			/*await owner.sendTransaction({
			  to: dex2addr,
			  value: initialEthBalance,
			});*/			
			await nativeToken.deposit({value: BigInt(2)*initialEthBalance});
			await nativeToken.transfer(dex2addr, BigInt(2)*initialEthBalance);
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]					
			await /*expect(*/Trade.connect(addr1).InstaTradeTokens(routeData1, initialEthBalance, 0, { value : initialEthBalance })//)
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialEthBalance, shouldGainAmount);	
				;				
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const finalBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();
            expect(finalBalance).to.be.above(initialBalance);
			//expect(finalBalance).to.be.equal(initialEthBalance + BigInt(shouldGainAmount));
			//const finalTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			//expect(finalTotalBalance).to.be.above(finalBalance);
			//expect(finalTotalBalance).to.be.above(initialTotalBalance);	
        });	
		
        it("Should revert an Ether InstaTradeTokens (without payable call) with a loss", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice/2, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);			
			/*
			const amtBack1 = await Trader.GetAmountOutMin.staticCallResult(dex1addr, 0, ZERO_ADDRESS, tokenAaddr, initialDexReserve);
			console.log(amtBack1);
			const balance1 = await ethers.provider.getBalance(dex1addr);
			console.log(balance1);			
			const amtBack2 = await Trader.GetAmountOutMin.staticCallResult(dex2addr, 0, tokenAaddr, ZERO_ADDRESS, amtBack1);
			console.log(amtBack2);
			const balance2 = await tokenA.balanceOf(dex2addr);
			console.log(balance2);
			*/
			await nativeToken.deposit({value: initialDexReserve});
			await nativeToken.transfer(dex2addr, initialDexReserve);
			
			await nativeToken.connect(addr1).deposit({value: initialDexReserve});
			await nativeToken.connect(addr1).approve(TradeAddr, initialDexReserve);				
			await Trade.connect(addr1).depositToken(NATIVE_TOKEN, initialDexReserve);
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]					
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert an Ether InstaTradeTokens (with payable call) with a loss", async function () {

			//const depositBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();	
			//await Trade.connect(addr1).withdrawEther(depositBalance);		
			const initialBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();			
			//expect(initialBalance).to.be.equal(0);

			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice/2, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			
			await nativeToken.deposit({value: initialDexReserve});
			await nativeToken.transfer(dex2addr, initialDexReserve);			
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0, { value : initialEthBalance })
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });				
		
        it("Should revert an Ether InstaTradeTokens (without payable call) with 0 gain", async function () {
			
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);

			await nativeToken.deposit({value: initialDexReserve});
			await nativeToken.transfer(dex2addr, initialDexReserve);
			
			await nativeToken.connect(addr1).deposit({value: initialDexReserve});
			await nativeToken.connect(addr1).approve(TradeAddr, initialDexReserve);				
			await Trade.connect(addr1).depositToken(NATIVE_TOKEN, initialDexReserve);		
            
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]			
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert an Ether InstaTradeTokens (with payable call) with 0 gain", async function () {
			
			//const depositBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();	
			//await Trade.connect(addr1).withdrawEther(depositBalance);		
			const initialBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();			
			//expect(initialBalance).to.be.equal(0);
			
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			
			await nativeToken.deposit({value: initialDexReserve});
			await nativeToken.transfer(dex2addr, initialDexReserve);			
            
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0, { value : initialEthBalance })
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });				
		
    });
	
    describe("InstaTradeTokens UniswapV3 Function", function () {
        it("Should perform a V3-V3 token InstaTradeTokens", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V3_Q1_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]						
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialDexReserve, shouldGainAmount);			
				;
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a V3-V3 token InstaTradeTokens with a loss", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			

			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V3_Q1_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }
			]					
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert a V3-V3 token InstaTradeTokens with 0 gain", async function () {
			
			const routeData1 = [
				{ Itype: IU_V3_Q1_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });
		
        it("Should perform a V3-V2 token InstaTradeTokens", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V3_Q1_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialDexReserve, shouldGainAmount);
				;				
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a V3-V2 token InstaTradeTokens with a loss", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			

			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V3_Q1_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }
			]		            
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert a V3-V2 token InstaTradeTokens with 0 gain", async function () {
			
			const routeData1 = [
				{ Itype: IU_V3_Q1_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });
			
        it("Should perform a V2-V3 token InstaTradeTokens", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialDexReserve, shouldGainAmount);			
				;
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a V2-V3 token InstaTradeTokens with a loss", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			

			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert a V2-V3 token InstaTradeTokens with 0 gain", async function () {
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });
		
    });	
	
    describe("InstaTradeTokens UniswapV4 Function", function () {
        it("Should perform a V4-V4 token InstaTradeTokens", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V4_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V4_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]			
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialDexReserve, shouldGainAmount);			
				;
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a V4-V4 token InstaTradeTokens with a loss", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			

			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V4_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V4_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }
			]					
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert a V4-V4 token InstaTradeTokens with 0 gain", async function () {
			
			const routeData1 = [
				{ Itype: IU_V4_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V4_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });
		
        it("Should perform a V4-V3 token InstaTradeTokens", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V4_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialDexReserve, shouldGainAmount);			
				;
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a V4-V3 token InstaTradeTokens with a loss", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			

			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V4_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }
			]			
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert a V4-V3 token InstaTradeTokens with 0 gain", async function () {
			
			const routeData1 = [
				{ Itype: IU_V4_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should perform a V3-V4 token InstaTradeTokens", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V3_Q1_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V4_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]					
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialDexReserve, shouldGainAmount);			
				;
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a V3-V4 token InstaTradeTokens with a loss", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			

			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V3_Q1_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V4_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert a V3-V4 token InstaTradeTokens with 0 gain", async function () {
			
			const routeData1 = [
				{ Itype: IU_V3_Q1_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V4_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]			
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should perform a V4-V2 token InstaTradeTokens", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V4_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialDexReserve, shouldGainAmount);			
				;
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a V4-V2 token InstaTradeTokens with a loss", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			

			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V4_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }
			]					
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert a V4-V2 token InstaTradeTokens with 0 gain", async function () {
			
			const routeData1 = [
				{ Itype: IU_V4_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]								
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });
			
        it("Should perform a V2-V4 token InstaTradeTokens", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V4_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]				
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, routeData1[0].asset, routeData1, initialDexReserve, shouldGainAmount);			
				;
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a V2-V4 token InstaTradeTokens with a loss", async function () {

			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			

			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V4_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 }
			]					
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert a V2-V4 token InstaTradeTokens with 0 gain", async function () {
			
			const routeData1 = [
				{ Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0 },
				{ Itype: IU_V4_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }
			]							
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData1, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
    });	

	describe("InstaTradeTokens Function", function () {
		
		it("Should perform an Ether 2dex InstaTradeTokens V2-V2 (without payable call)", async function () {

			await nativeToken.connect(addr1).deposit({value: initialEthBalance});
			await nativeToken.connect(addr1).approve(TradeAddr, initialEthBalance);				
			await Trade.connect(addr1).depositToken(NATIVE_TOKEN, initialEthBalance);
			
            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			const initialBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();
			//const initialTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			//expect(initialTotalBalance).to.be.equal(initialEthBalance);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, 2*initialPrice, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			const shouldGainAmount = initialEthBalance;
					
			const route1 = { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 }
			const route2 = { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }						
			const amtBack = BigInt(await Trader.GetAmountOutMin.staticCallResult(route1, route2.asset, initialEthBalance));
			const finalEthBalance = BigInt(await Trader.GetAmountOutMin.staticCallResult(route2, route1.asset, amtBack));
			expect(finalEthBalance).to.be.equal(initialEthBalance + shouldGainAmount);
			
			await tokenA.transfer(dex1addr, amtBack);	
			/*await owner.sendTransaction({
			  to: dex2addr,
			  value: initialEthBalance,
			});*/				
			await nativeToken.deposit({value: BigInt(2)*initialEthBalance});
			await nativeToken.transfer(dex2addr, BigInt(2)*initialEthBalance);
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await /*expect(*/Trade.connect(addr1).InstaTradeTokens(routeData, initialEthBalance, 0)//)
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, ZERO_ADDRESS, routeData, initialEthBalance, shouldGainAmount);			
				;
				
			await Trade.connect(addr1).withdrawToken(NATIVE_TOKEN, shouldGainAmount);
			await nativeToken.connect(addr1).withdraw(shouldGainAmount);				
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const finalBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();
            expect(finalBalance).to.be.above(initialBalance);
			//expect(finalBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			//const finalTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			//expect(finalTotalBalance).to.be.above(finalBalance);
			//expect(finalTotalBalance).to.be.above(initialTotalBalance);			
        });		

		it("Should perform an Ether 2dex InstaTradeTokens V2-V2 (with payable call)", async function () {

			//const depositBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();	
			//await Trade.connect(addr1).withdrawEther(depositBalance);		
			const initialBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();			
			//expect(initialBalance).to.be.equal(0);
			//const initialTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			//expect(initialTotalBalance).to.be.equal(initialEthBalance);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, 2*initialPrice, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			const shouldGainAmount = initialEthBalance;
					
			const route1 = { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0 }
			const route2 = { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0 }							
			const amtBack = BigInt(await Trader.GetAmountOutMin.staticCallResult(route1, route2.asset, initialEthBalance));
			const finalEthBalance = BigInt(await Trader.GetAmountOutMin.staticCallResult(route2, route1.asset, amtBack));
			expect(finalEthBalance).to.be.equal(initialEthBalance + shouldGainAmount);
			
			await tokenA.transfer(dex1addr, amtBack);	
			/*await owner.sendTransaction({
			  to: dex2addr,
			  value: initialEthBalance,
			});*/				
			await nativeToken.deposit({value: BigInt(2)*initialEthBalance});
			await nativeToken.transfer(dex2addr, BigInt(2)*initialEthBalance);
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialEthBalance, 0, { value : initialEthBalance }))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, ZERO_ADDRESS, routeData, initialEthBalance, shouldGainAmount);			
				;				
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const finalBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();
            expect(finalBalance).to.be.above(initialBalance);
			//expect(finalBalance).to.be.equal(initialEthBalance + BigInt(shouldGainAmount));
			//const finalTotalBalance = await Trade.connect(addr1).getTotalEtherBalance();
			//expect(finalTotalBalance).to.be.above(finalBalance);
			//expect(finalTotalBalance).to.be.above(initialTotalBalance);	
        });	
		
        it("Should revert an Ether 2dex InstaTradeTokens V2-V2 (without payable call) with a loss", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice/2, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);	

			await nativeToken.deposit({value: initialDexReserve});
			await nativeToken.transfer(dex2addr, initialDexReserve);
			
			await nativeToken.connect(addr1).deposit({value: initialDexReserve});
			await nativeToken.connect(addr1).approve(TradeAddr, initialDexReserve);			
			await Trade.connect(addr1).depositToken(NATIVE_TOKEN, initialDexReserve);			

            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert an Ether 2dex InstaTradeTokens V2-V2 (with payable call) with a loss", async function () {

			//const depositBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();	
			//await Trade.connect(addr1).withdrawEther(depositBalance);		
			const initialBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();			
			//expect(initialBalance).to.be.equal(0);

			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice/2, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			
			await nativeToken.deposit({value: initialDexReserve});
			await nativeToken.transfer(dex2addr, initialDexReserve);			
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  }
            ]			
            await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0, { value : initialEthBalance })
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });				
		
        it("Should revert an Ether 2dex InstaTradeTokens V2-V2 (without payable call) with 0 gain", async function () {
			
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			
			await nativeToken.deposit({value: initialDexReserve});
			await nativeToken.transfer(dex2addr, initialDexReserve);

			await nativeToken.connect(addr1).deposit({value: initialDexReserve});
			await nativeToken.connect(addr1).approve(TradeAddr, initialDexReserve);				
			await Trade.connect(addr1).depositToken(NATIVE_TOKEN, initialDexReserve);	
            
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  }
            ]			
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		
		
        it("Should revert an Ether 2dex InstaTradeTokens V2-V2 (with payable call) with 0 gain", async function () {
			
			//const depositBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();	
			//await Trade.connect(addr1).withdrawEther(depositBalance);		
			const initialBalance = await ethers.provider.getBalance(addr1); //Trade.connect(addr1).getEtherBalance();			
			//expect(initialBalance).to.be.equal(0);
			
			await dex1.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			await dex2.setPairInfo(tokenAaddr, NATIVE_TOKEN, initialPrice, poolFee);
			
			await nativeToken.deposit({value: initialDexReserve});
			await nativeToken.transfer(dex2addr, initialDexReserve);			
            
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: ZERO_ADDRESS, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  }
            ]			
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0, { value : initialEthBalance })
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });					

        it("Should perform a token 2dex InstaTradeTokens V2-V2 ", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			// add tokenB extra reserve for dex2
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, ZERO_ADDRESS, routeData, initialDexReserve, shouldGainAmount);			
				;           			
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a token 2dex InstaTradeTokens V2-V2 with a loss", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });								
			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			
			
			// add tokenB extra reserve for dex2
			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");

        });		
		
        it("Should revert a token 2dex InstaTradeTokens V2-V2 with 0 gain", async function () {
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });
		
        it("Should perform a token 3dex InstaTradeTokens V2-V3-V4 ", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });
			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);
			
			const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalanceA = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			
			// add tokenB extra reserve for dex2
			await tokenA.transfer(dex1addr, initialDexReserve);			
			await tokenB.transfer(dex2addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, 2*initialPrice, poolFee);
			const shouldGainAmount = initialDexReserve;
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenCaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V4_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0))
                //.to.emit(Trade, "InstaTraded")
                //.withArgs(addr1.address, ZERO_ADDRESS, routeData, initialDexReserve, shouldGainAmount);			
				;           			
            
			expect(initialBalanceA).to.be.equal(await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1));
			expect(initialTotalBalanceA).to.be.equal(await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr));

            const tokenBalance = await Trade.connect(addr1).getTokenBalance(tokenBaddr, addr1);
            expect(tokenBalance).to.be.above(initialBalance);
			expect(tokenBalance).to.be.equal(initialBalance + BigInt(shouldGainAmount));
			const tokenTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenBaddr);
			expect(tokenTotalBalance).to.be.above(tokenBalance);
			expect(tokenTotalBalance).to.be.above(initialTotalBalance);			
        });
		
        it("Should revert a token 3dex InstaTradeTokens V2-V3-V4 with a loss", async function () {

            // should already use the 'common' deposited amount
			//await Trade.depositEther({ value: initialPrice });								
			const initialBalance = await Trade.connect(addr1).getTokenBalance(tokenAaddr, addr1);
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(initialDexReserve);			
			
			// add tokenB extra reserve for dex2
			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V4_POOL, router: dex2addr, asset: tokenCaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");

        });		
		
        it("Should revert a token 3dex InstaTradeTokens V2-V3-V4 with 0 gain", async function () {
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V3_Q1_POOL, router: dex2addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V4_POOL, router: dex2addr, asset: tokenCaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await expect(Trade.connect(addr1).InstaTradeTokens(routeData, initialDexReserve, 0)
				).to.be.revertedWith("Trade Reverted, No Profit Made");
        });		

    });	
	
	describe("InstaTradeTokens Other Miscellaneus Functions", function () {
			
        /*
		it("Should NOT revert a token 2dex InstaTradeTokens V2-V2 with a loss", async function () {

			await tokenA.transfer(addr1, initialDexReserve);
			await tokenA.connect(addr1).approve(TradeAddr, initialDexReserve);
			await Trade.connect(addr1).depositToken(tokenAaddr, initialDexReserve);
            const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr,  addr1);
			expect(initialBalanceA).to.be.equal(initialDexReserve);	
			
			const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(BigInt(2)*initialDexReserve);			
			
			// add tokenB extra reserve for dex2
			await tokenA.transfer(dex2addr, initialDexReserve);			
			await tokenB.transfer(dex1addr, initialDexReserve);			
			await dex2.setPairInfo(tokenAaddr, tokenBaddr, initialPrice/2, poolFee);
			const shouldGainAmount = initialDexReserve;
			
            const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await Trade.connect(addr1).InstaTradeTokensChecked(routeData, initialDexReserve, 0, false)

			const finalTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.above(finalTotalBalance);	
            const finalBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr,  addr1);
			expect(initialBalanceA).to.be.above(finalBalanceA);				
        });		
		
        it("Should NOT revert a token 2dex InstaTradeTokens V2-V2 with 0 gain", async function () {
			
			await tokenA.transfer(addr1, initialDexReserve);
			await tokenA.connect(addr1).approve(TradeAddr, initialDexReserve);
			await Trade.connect(addr1).depositToken(tokenAaddr, initialDexReserve);
            const initialBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr,  addr1);
			expect(initialBalanceA).to.be.equal(initialDexReserve);				
			
            const initialTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(BigInt(2)*initialDexReserve);	
			
			const routeData = [
                { Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  },
                { Itype: IU_V2_POOL, router: dex2addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0  }
            ]
			await Trade.connect(addr1).InstaTradeTokensChecked(routeData, initialDexReserve, 0, false)
				
			const finalTotalBalance = await Trade.connect(addr1).getTotalTokenBalance(tokenAaddr);
			expect(initialTotalBalance).to.be.equal(finalTotalBalance);	
            const finalBalanceA = await Trade.connect(addr1).getTokenBalance(tokenAaddr,  addr1);
			expect(finalBalanceA).to.be.equal(initialBalanceA);				
        });
		*/
		
       it("InstaSwapTokens Token to Token", async function () {
		
			await tokenA.approve(TradeAddr, initialDexReserve);
			await Trade.depositToken(tokenAaddr, initialDexReserve);								
			const initialBalanceA = await Trade.getTokenBalance(tokenAaddr, owner);
			expect(initialBalanceA).to.be.equal(BigInt(2)*initialDexReserve);	
			const initialBalanceB = await Trade.getTokenBalance(tokenBaddr, owner);
			expect(initialBalanceB).to.be.equal(initialDexReserve);	
					
            const routeChain = { Itype: IU_V2_POOL, router: dex1addr, asset: tokenAaddr, poolFee: poolFee, tickSpacing: 0  }

			await Trade.InstaSwapTokens(routeChain, initialDexReserve, tokenBaddr, 0);

			const finalBalanceA = await Trade.getTokenBalance(tokenAaddr, owner);
			expect(initialBalanceA).to.be.above(finalBalanceA);	
			const finalBalanceB = await Trade.getTokenBalance(tokenBaddr, owner);
			expect(finalBalanceB).to.be.above(initialBalanceB);
        });		

       it("InstaSwapTokens Ether to Token", async function () {
							
			const initialBalanceN = await Trade.getTokenBalance(NATIVE_TOKEN, owner);
			await Trade.withdrawToken(NATIVE_TOKEN, initialBalanceN);
			
			await nativeToken.deposit({value: initialEthBalance});
			await nativeToken.transfer(dex1addr, initialEthBalance);			
			 
			const initialBalanceB = await Trade.getTokenBalance(tokenBaddr, owner);
			expect(initialBalanceB).to.be.equal(initialDexReserve);	
					
            const routeChain = { Itype: IU_V2_POOL, router: dex1addr, asset: NATIVE_TOKEN, poolFee: poolFee, tickSpacing: 0  }

			await Trade.InstaSwapTokens( routeChain, initialDexReserve, tokenBaddr, 0, { value : initialDexReserve });

			const finalBalanceB = await Trade.getTokenBalance(tokenBaddr, owner);
			expect(finalBalanceB).to.be.above(initialBalanceB);
        });			

       it("InstaSwapTokens Token to Ether", async function () {
							
			await nativeToken.connect(addr1).deposit({value: initialDexReserve});
			await nativeToken.connect(addr1).transfer(dex1addr, initialDexReserve);	
			
			const initialBalance = await ethers.provider.getBalance(owner);
			
			const initialBalanceB = await Trade.getTokenBalance(tokenBaddr, owner);
			expect(initialBalanceB).to.be.equal(initialDexReserve);	
					
            const routeChain = { Itype: IU_V2_POOL, router: dex1addr, asset: tokenBaddr, poolFee: poolFee, tickSpacing: 0  }

			await Trade.InstaSwapTokens( routeChain, initialDexReserve, ZERO_ADDRESS, 0);

			const finalBalance = await ethers.provider.getBalance(owner);
			expect(finalBalance).to.be.above(initialBalance);
        });	

    });		
	
});

