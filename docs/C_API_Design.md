# C. API Design

## Endpoints

### Markets
*   `POST /api/markets` (Creator): Create a new market from a template.
*   `GET /api/markets/:id` (Viewer): Get market details, rules, and current pool sizes.
*   `GET /api/streams/:streamId/markets` (Viewer): List active markets for a stream.

### Betting
*   `POST /api/markets/:id/stake` (Viewer): Place a stake.
    *   Body: `{ outcomeId: "uuid", amount: 10 }`
    *   *Returns: updated position and ledger balance.*

### Admin & Moderation
*   `POST /api/admin/markets/:id/moderate` (Admin): Approve/Reject market draft.
*   `POST /api/admin/markets/:id/resolve` (Admin): Manually resolve or void a market.

### Webhooks (Internal)
*   `POST /internal/ai/resolution` (AI Referee): Submit AI findings.
    *   Body: Match AI Service output JSON schema.
*   `POST /internal/stream/health` (Media Server): Report FPS drops, disconnects.