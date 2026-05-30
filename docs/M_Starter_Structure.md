# M. Starter Code Structure

```text
predi/
├── docs/                      # Architectural documents
├── server/                    # Node.js / Express Core Backend
│   ├── src/
│   │   ├── index.ts           # Entry point
│   │   ├── routes/            # API endpoints (markets, stakes, admin)
│   │   ├── services/          # LedgerService, MarketMachine, MediaService
│   │   ├── models/            # DB schemas/types
│   │   └── webhooks/          # Internal AI endpoints
│   ├── package.json
│   └── tsconfig.json
├── ai-referee/                # Python / YOLO Service
│   ├── main.py                # Fast API listener
│   ├── detector.py            # YOLO inference logic
│   ├── tracker.py             # Object tracking logic
│   └── requirements.txt
├── client/                    # React / Next.js Web App
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   └── hooks/
│   └── package.json
├── docker-compose.yml         # Local dev: Postgres, Redis, LiveKit
└── .env.example               # Feature flags (PLAY_MONEY_ONLY=true)
```