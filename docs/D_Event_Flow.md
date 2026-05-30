# D. Market Lifecycle Event Flow

1.  **Creation (Creator):** Creator selects the "Traffic Count > 20" template. App sets `eventWindowSeconds=60`. Market created in `DRAFT` state.
2.  **Pre-Moderation (System):** System checks creator trust score and template ID. Moves to `APPROVED`.
3.  **Opening (Creator/System):** Creator clicks "Go". Market moves to `OPEN`. Broadcaster sends WS event to viewers.
4.  **Betting Phase (Viewers):** Viewers stake play-money. Ledger entries (`type: stake`) are created. Pool updates are broadcasted.
5.  **Locking (System):** `bettingClosesBeforeStartSeconds` (e.g. 5s) before event start, system moves market to `LOCKED`. Staking is rejected.
6.  **Event Active (System):** Event starts. State is `EVENT_ACTIVE`. AI Referee begins focused monitoring of stream segment.
7.  **Resolution Pending (System):** Event window ends. Market moves to `RESOLUTION_PENDING`. AI compiles final tracking data.
8.  **AI Decision (AI Referee):** AI posts result (e.g., `count: 23`, `confidence: 0.94`, `YES`).
9.  **Resolution (System):** If confidence > threshold & no stream health issues, state moves to `RESOLVED`.
10. **Settlement (System):** Ledger engine calculates pro-rata payouts. Creates `win` ledger entries for "YES" holders. Subtracts platform fee. Market state moves to `SETTLED`.