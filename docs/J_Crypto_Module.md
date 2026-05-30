# J. Crypto Module & Feature Flags

## Policy Statement
This platform is NOT an unlicensed gambling product. All real-value capabilities are strictly gated behind compliance flags and must only be enabled in jurisdictions with appropriate licensure and KYC/AML partnerships.

## Feature Flags

```env
# .env config
PLAY_MONEY_ONLY=true
CRYPTO_ENABLED=false
REAL_MONEY_ENABLED=false
REQUIRE_KYC=false
```

## Future Architecture (When `CRYPTO_ENABLED=true`)

### 1. Deposits
*   **Provider:** Fireblocks or Coinbase Commerce.
*   **Flow:** User requests deposit address -> Sends USDT -> Webhook receives `CONFIRMED` status -> System adds `deposit` to ledger -> Balance updates.

### 2. Withdrawals
*   **Queue System:** Withdrawals are NOT instant. They enter a `PENDING_REVIEW` queue.
*   **Checks:**
    1.  `KYC_STATUS == 'VERIFIED'`
    2.  `USER_JURISDICTION` not in restricted list (e.g., US, UK).
    3.  `AML_RISK_SCORE` < threshold (Chainalysis/Elliptic API).
*   **Execution:** Admin approves -> API call to custodian to send funds -> `withdrawal` ledger entry created.

### 3. Smart Contracts (Optional Future)
*   For fully decentralized deployment, markets could be deployed as individual smart contracts, with the AI Referee acting as an Oracle (e.g., using Chainlink or a signed payload). The MVP will NOT use this to maintain speed and control.