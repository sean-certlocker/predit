# L. MVP Build Roadmap

## Phase 1: Core Scaffolding (Weeks 1-2)
*   Setup Node.js + Express API.
*   Setup PostgreSQL schema and Ledger engine.
*   Setup basic React Web App (Viewer & Creator views).
*   Implement `PLAY_MONEY_ONLY` authentication and wallets.

## Phase 2: Streaming & Templates (Weeks 3-4)
*   Integrate LiveKit for WebRTC streaming.
*   Implement hardcoded Market Templates (e.g., "Traffic > X").
*   Build the Market State Machine (`DRAFT` -> `SETTLED`).
*   Implement the Viewer UI to place play-money stakes.

## Phase 3: AI Referee Integration (Weeks 5-6)
*   Build Python service to pull HLS chunks from LiveKit.
*   Integrate YOLOv8 for vehicle counting.
*   Build webhook to send `ai_resolution_events` back to Core API.
*   Implement the `RESOLUTION_PENDING` -> `RESOLVED` logic based on confidence.

## Phase 4: Anti-Manipulation & Safety (Weeks 7-8)
*   Implement stream health monitoring (disconnect voids).
*   Integrate basic Vision API for NSFW/Violence detection.
*   Build basic Admin UI for `MANUAL_REVIEW`.
*   End-to-End testing with simulated streams.