# ATLAS Tokenomics Analysis: The Dual-Asset Economy

While the low-level Go implementation (`tokenomics.go`) appears simple at first glance—handling basic mints, burns, and fixed-fee calculations—the actual economic design detailed in the protocol planning is highly sophisticated.

ATLAS employs a **Dual-Asset Macroeconomic Model** designed to separate speculative financial value from social utility. This solves a major problem in Web3 social networks where high token prices price-out normal users from interacting.

---

## 1. TCOIN (The Hard Asset)
`TCOIN` is the Layer-1 native cryptocurrency. It acts as the immutable store of value, the security bond for validators, and the ultimate settlement layer.

*   **Hard Cap:** Fixed at **1,000,000,000 (1 Billion)** units. 
*   **Genesis Supply:** Pre-minted at block 0 and held by the Network Treasury. It cannot be generated out of thin air by users; new circulation only enters the market via validator block rewards.
*   **Base Fees:** Transaction fees are programmatically calculated via `CalculateFee` based on the payload size (e.g., `BaseFee(10) + len(Data)/100`).
*   **The Utility Anchor:** The primary utility of TCOIN is not just paying gas; it acts as a "Generator." Holding TCOIN in your wallet passively generates the network's second asset: **Data Units**.

## 2. Data Units / DU (The Fluid Asset)
Data Units (DUs) are not cryptocurrency tokens; they are environmental sub-layer metrics representing the "Energy" of an interaction. 

*   **Generation:** They are generated fluidly as a dividend of holding TCOIN.
*   **Social Physics:** Submitting a post *does not* burn DUs. Instead, the post acts as a sponge, absorbing the creator's contextual DUs. When another user interacts with the post (e.g., upvotes), their DUs undergo "Fusion" with the post's DUs. 
*   **The Burn & Metric Loop:** The resulting fused data packet is pushed into the ATLAS metric system and then *burned*. This creates highly valuable, verifiable macro-market data (e.g., "users with X trait interact heavily with Y product").
*   **Value Extraction:** Content creators whose objects facilitate massive DU burn can redeem that engagement for TCOIN rewards from the treasury emission pool. 

## 3. DeFi & Staking Mechanics (`defi_staking.go`)
The network natively implements a Proof-of-Stake reward pool to secure the chain and disincentivize market dumping.

*   **Lock-Up Period:** Staked TCOIN is subject to a hardcoded `7-day` lock-up period (`time.Hour * 24 * 7`). Unstaking before this period is cryptographically rejected by the node.
*   **Reward Rate:** The staking contract natively calculates rewards at an **10% APY** (`rewardRate: 0.1`). 
*   **Slashing:** (As noted in the consensus layer), malicious behavior burns the staker's reputation and ejects them, freezing their staked capital.

## 4. On-Chain Governance (`defi_staking.go`)
TCOIN acts as a governance token managed by the `GovernanceSystem`.
*   **Thresholds:** It requires a minimum of 1000 TCOIN to create a proposal, and 100 TCOIN to cast a vote.
*   **Quorum:** A proposal requires a 10% minimum network quorum and a 60% majority to pass.
*   **Execution Delay:** If passed, the system enforces a 100-block execution delay before the action (e.g., Treasury Transfer, Parameter Change) is processed natively by the chain.

---

### Conclusion
The tokenomics of ATLAS create a closed, self-sustaining loop:
1. **Hold TCOIN** -> 2. **Passively Generate DUs** -> 3. **Spend DUs on Social Interactions** -> 4. **Generate Valuable Macro-Metrics** -> 5. **Extract TCOIN Value from Metrics**. 

This completely bypasses the traditional Web2 model of selling user data to advertisers, instead redirecting the financial value of market data directly back into the `TCOIN` liquidity pool.
