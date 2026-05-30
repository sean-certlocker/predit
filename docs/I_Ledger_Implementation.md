# I. Play-Money Ledger Implementation

## Core Principles
1. **Append-Only:** Balances are never updated directly. They are a computed view of the `ledger_entries` table.
2. **Double-Entry Simulation:** Every action creates matching positive and negative flows logically, though for MVP we track user balances.
3. **Pari-Mutuel Pool:** The odds are not fixed; they are determined by the ratio of stakes when the market locks.

## Example Flow: "Will >20 cars cross?"

**1. Staking Phase**
*   User A stakes 10 on YES.
    *   `INSERT INTO ledger_entries (user_id, amount, type) VALUES ('A', -10, 'stake')`
    *   `INSERT INTO positions (user_id, outcome, stake) VALUES ('A', 'YES', 10)`
*   User B stakes 30 on NO.
    *   `INSERT INTO ledger_entries (user_id, amount, type) VALUES ('B', -30, 'stake')`
    *   `INSERT INTO positions (user_id, outcome, stake) VALUES ('B', 'NO', 30)`

**2. Lock Phase**
*   Total Pool = 40.
*   YES Pool = 10, NO Pool = 30.

**3. Settlement Phase (Result = YES)**
*   System deducts 5% Platform Fee (2).
*   Remaining Pool = 38.
*   Winning Pool = YES (Total Stake: 10).
*   User A's Share = (10 / 10) = 100%.
*   User A Payout = 38.
*   System creates payout entry:
    *   `INSERT INTO ledger_entries (user_id, amount, type) VALUES ('A', 38, 'win')`
    *   System simulates Creator Fee (e.g., 20% of Platform Fee = 0.4):
    *   `INSERT INTO ledger_entries (user_id, amount, type) VALUES ('Creator', 0.4, 'creator_share')`

**4. Void Scenario**
*   If market is VOID.
*   Refund User A: `INSERT INTO ledger_entries (user_id, amount, type) VALUES ('A', 10, 'refund')`
*   Refund User B: `INSERT INTO ledger_entries (user_id, amount, type) VALUES ('B', 30, 'refund')`