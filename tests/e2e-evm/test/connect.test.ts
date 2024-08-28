import hre from "hardhat";
import { extendEnvironment } from "hardhat/config";
import { expect } from "chai";
import { Address, defineChain, Chain, PublicClientConfig, WalletClientConfig } from 'viem';
import { whaleAddress, userAddress } from './addresses';

describe("Viem Setup", function () {
  it("getChainID() matches the configured network", async function() {
    const publicClient = await hre.viem.getPublicClient();

    const expectedChainId = (() => {
      switch(hre.network.name) {
        case "hardhat":
          return 31337;
        case "kvtool":
          return 8888;
      }
    })();

    expect(hre.network.config.chainId).to.eq(expectedChainId);
    expect(await publicClient.getChainId()).to.eq(expectedChainId);
  });

  it("is configured with whale and user accounts", async function() {
    const walletClients = await hre.viem.getWalletClients();

    expect(walletClients.length).to.equal(2);
    expect(walletClients[0].account.address).to.equal(whaleAddress)
    expect(walletClients[1].account.address).to.equal(userAddress)
  });

  [
    {name: "whale", from: whaleAddress, to: userAddress, value: 1n},
    {name: "user", from: userAddress, to: whaleAddress, value: 1n},
  ].forEach((tc) => {
    it(`${tc.name} account can sign transactions`, async function() {
      const publicClient = await hre.viem.getPublicClient();
      const walletClient = await hre.viem.getWalletClient(tc.from);

      const txHash = await walletClient.sendTransaction({
        to: tc.to,
        value: tc.value,
      });
      const txReceipt = await publicClient.waitForTransactionReceipt({ hash: txHash });
      const tx = await publicClient.getTransaction({ hash: txHash });

      expect(txReceipt.status).to.eq('success');
      expect(txReceipt.gasUsed).to.eq(21000n);
      expect(tx.value).to.eq(tc.value)
    });
  });
});
