// to deploy locally
// run: npx hardhat node on a terminal
// then run: npx hardhat run --network localhost scripts/deploy_main.js

async function main(network) {

	console.log('network: ', network.name);

	const [deployer] = await ethers.getSigners();
	console.log("Deploying contracts with the account:", deployer.address);
	console.log("Account balance:", await ethers.provider.getBalance(deployer));
	
	const Trade = await ethers.getContractFactory("Trade");
	const trade = await Trade.deploy();
	await trade.waitForDeployment();
	const TRADE_ADDRESS = await trade.getAddress();
	console.log("TRADE address:", TRADE_ADDRESS);	
	
	
	const MockDex = await ethers.getContractFactory("MockDEX");
	
	const v2_dex_1 = await MockDex.deploy();
	await v2_dex_1.waitForDeployment();
	const V2_DEX_1_ADDRESS = await v2_dex_1.getAddress();
	console.log("V2_DEX_1 address:", V2_DEX_1_ADDRESS);
	
	const v2_dex_2 = await MockDex.deploy();
	await v2_dex_2.waitForDeployment();
	const V2_DEX_2_ADDRESS = await v2_dex_2.getAddress();
	console.log("V2_DEX_2 address:", V2_DEX_2_ADDRESS);
	
	const v3_dex_1 = await MockDex.deploy();
	await v3_dex_1.waitForDeployment();
	const V3_DEX_1_ADDRESS = await v3_dex_1.getAddress();
	console.log("V3_DEX_1 address:", V3_DEX_1_ADDRESS);
	
	const v3_dex_2 = await MockDex.deploy();
	await v3_dex_2.waitForDeployment();
	const V3_DEX_2_ADDRESS = await v3_dex_2.getAddress();
	console.log("V3_DEX_2 address:", V3_DEX_2_ADDRESS);	


	const MockERC20 = await ethers.getContractFactory("MockERC20");
	
	const ERC20_1 = await MockERC20.deploy("ERC20_1", "ERC1", 18, 100000);
	await ERC20_1.waitForDeployment();
	const ERC20_1_ADDRESS =  await ERC20_1.getAddress();
	console.log("ERC20_1 address:", ERC20_1_ADDRESS);	
	
	const ERC20_2 = await MockERC20.deploy("ERC20_2", "ERC2", 18, 100000);
	await ERC20_2.waitForDeployment();
	const ERC20_2_ADDRESS =  await ERC20_2.getAddress();
	console.log("ERC20_2 address:", ERC20_2_ADDRESS);	

	const ERC20_3 = await MockERC20.deploy("ERC20_3", "ERC3", 18, 100000);
	await ERC20_3.waitForDeployment();
	const ERC20_3_ADDRESS =  await ERC20_3.getAddress();
	console.log("ERC20_3 address:", ERC20_3_ADDRESS);	
	
	const ERC20_4 = await MockERC20.deploy("ERC20_4", "ERC4", 18, 100000);
	await ERC20_4.waitForDeployment();
	const ERC20_4_ADDRESS =  await ERC20_4.getAddress();
	console.log("ERC20_4 address:", ERC20_4_ADDRESS);		
	
   
    const finalbalance = await ethers.provider.getBalance(deployer);
    console.log(`Deployer's balance: `, finalbalance);

  }
  
  // We recommend this pattern to be able to use async/await everywhere
  // and properly handle errors.
  main(network)
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });
  

