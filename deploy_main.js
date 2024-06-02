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
	console.log("Trade address:", await trade.getAddress());	
   
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
  

