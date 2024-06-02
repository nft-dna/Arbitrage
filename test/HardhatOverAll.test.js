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
		
	let MockERC20;
	let token;
	let owner;
	let addr1;
	let addr2;
	let initialSupply = ethers.parseEther("1000");
  
	beforeEach(async function () {
		[owner, addr1, addr2] = await ethers.getSigners();	
		MockERC20 = await ethers.getContractFactory("MockERC20");
		token = await MockERC20.deploy("MockToken", "MTK", 18, initialSupply);
		await token.waitForDeployment();		
		//console.log(await token.getAddress());	
	});
		
	describe("MockERC20", async function () {
	  it("Should return the right name and symbol", async function () {
	  	//console.log(await token.getAddress());
		expect(await token.name()).to.equal("MockToken");
		expect(await token.symbol()).to.equal("MTK");
		expect(await token.balanceOf(owner.address)).to.equal(ethers.parseEther("1000"));
	  });
	});
});

