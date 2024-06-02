require('dotenv').config();
require("@nomicfoundation/hardhat-toolbox");

const MAINNET_PRIVATE_KEY = process.env.MAINNET_PRIVATE_KEY;
const TESTNET_PRIVATE_KEY = process.env.TESTNET_PRIVATE_KEY;


/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    version: '0.8.25',
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  networks: {
    mainnet: {
      url: `https://eth.meowrpc.com`,
      chainId: 1,
      accounts: [`0x${MAINNET_PRIVATE_KEY}`]
    },
    magma: {
      url: `https://turbo.magma-rpc.com/`,
      chainId: 6969696969,
      accounts: [`0x${TESTNET_PRIVATE_KEY}`]
    },
	localhost: {
      url: `http://127.0.0.1:8545`
    },
  },  
};
