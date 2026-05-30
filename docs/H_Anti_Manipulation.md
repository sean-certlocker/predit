# H. Anti-Manipulation System

## Checks & Triggers

1.  **Stream Disconnects / Latency Abuse**
    *   *Rule:* If the stream drops for > 3 seconds cumulatively during `EVENT_ACTIVE`, VOID market.
    *   *Rule:* If client-server latency > 5 seconds during `OPEN` phase, pause betting.
2.  **Camera Movement**
    *   *Rule:* AI computes frame-to-frame optical flow. If significant global shift detected (camera knocked over or moved to change view), VOID market.
3.  **Camera Blocked**
    *   *Rule:* If > 80% of pixels become pure black or static for > 2 seconds during `EVENT_ACTIVE`, VOID market.
4.  **Replay / Uploaded Video Detection**
    *   *Rule:* Enforce WebRTC/RTMP native mobile ingest. Reject virtual webcams (OBS) for standard creator tier.
    *   *Rule:* Require creator to show a dynamic server-generated QR code or random word on screen before starting the first market of a stream to prove "liveness".
5.  **Late Betting**
    *   *Rule:* Enforce strict `LOCKED` phase. e.g., Event happens 12:00:10 to 12:01:10. Bets lock firmly at 12:00:05.
6.  **Creator Self-Dealing**
    *   *Rule:* Creator cannot place stakes on their own stream's markets.
7.  **Suspicious Spikes**
    *   *Rule:* If 90% of the pool is staked in the last 1 second before lock by a single entity, trigger `MANUAL_REVIEW` before resolution.