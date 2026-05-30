# Live-Stream Prediction MVP - Architecture

## System Components

1. **Client Applications (Mobile-First Web/React)**
   - **Viewer App:** Streams WebRTC video, displays active markets, manages play-money balance, places stakes.
   - **Creator App:** Captures camera, goes live via WebRTC/RTMP, selects market templates, starts events.
   - **Admin Dashboard:** Moderates streams, reviews AI flags, manages users, resolves disputed markets.

2. **Edge / Ingest Layer**
   - **Media Server (e.g., LiveKit, Mediasoup, or Nginx-RTMP):** Receives WebRTC/RTMP, adds server-side timestamp watermarks, monitors stream health.
   - **CDN:** Distributes HLS/WebRTC to viewers.

3. **Core Backend (Go)**
   - **API Gateway:** REST/JSON using Go (Gin or Echo).
   - **Market Engine:** Go-based state machine for market lifecycle.
   - **Ledger Service:** High-performance append-only ledger in Go.
   - **Event Bus:** Redis/NATS for live updates.

4. **AI Referee Service (Python/Go)**
   - YOLO-based object detection.
   - Communicates via webhooks or gRPC to the Go Backend.

## Client Application (Native Android - Java)
- **Creator & Viewer Modes:** Unified native app.
- **Streaming:** RTMP/WebRTC using native Android libraries.
- **UI:** XML/Java based native UI for maximum performance.

5. **Storage Layer**
   - **Primary DB (PostgreSQL):** Users, markets, ledgers, streams.
   - **In-Memory Cache (Redis):** Active market pools, live stream health.
   - **Object Storage (S3):** Recorded stream chunks for audit, evidence snapshots.

## Diagram

```text
[Creator Mobile] ---WebRTC/RTMP---> [Media Server] ----> [Object Storage (S3)]
                                          |
                                          | (Stream Frames)
                                          v
                                    [AI Referee (Python, YOLO)]
                                          | (Resolution Events)
                                          v
[Viewer Mobile] <---WebSockets--->  [Core API (Node.js)] <---> [PostgreSQL] (Ledger, Markets)
[Admin Web]     <---REST/WS------>        ^
                                          |
                                    [Redis (Pub/Sub & Cache)]
```