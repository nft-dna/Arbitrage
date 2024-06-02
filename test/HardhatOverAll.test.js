// npx hardhat test ./test/HardhatOverAll.test.js --network localhost; 
// run first (in another shell): npx hardhat node
require('@nomiclabs/hardhat-truffle5');

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

const MockERC20 = artifacts.require('MockERC20');

describe("Overall Test", function () {	
		
	beforeEach(async function () {
		this.owner = await ethers.getSigners();		
		this.token = await MockERC20.new("MockToken", "MTK", 18, ethers.parseEther("1000"));
		await this.token.waitForDeployment();
		console.log(await this.token.getAddress());	
	});
		
	describe("MockERC20", async function () {
	  it("Should return the right name and symbol", async function () {
	  	console.log(await this.token.getAddress());
		expect(await this.token.name()).to.equal("MockToken");
		expect(await this.token.symbol()).to.equal("MTK");
		expect(await this.token.balanceOf(this.owner.address)).to.equal(ethers.parseEther("1000"));
	  });
	});
});

