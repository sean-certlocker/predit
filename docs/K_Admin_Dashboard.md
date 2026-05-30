# K. Admin Dashboard Design

## Views

1.  **Live Streams Matrix:**
    *   Grid of active streams.
    *   Health indicators (Green/Yellow/Red) based on FPS and disconnects.
    *   Active Safety Flags (e.g., "NSFW Model: 0.1").
2.  **Market Queue (Moderation_Pending):**
    *   List of requested markets from non-trusted creators.
    *   Buttons: [Approve] [Reject]
3.  **Disputed / Manual Review Markets:**
    *   Markets where AI confidence < threshold or viewers reported.
    *   Displays:
        *   Stream recording snippet (the `eventWindowSeconds`).
        *   AI Snapshot evidence.
        *   Current Pool sizes.
    *   Action Buttons: [Resolve YES] [Resolve NO] [VOID]
4.  **User Risk Dashboard:**
    *   List of users with suspicious betting patterns (e.g., high win rate on specific creators, late betting attempts).
5.  **Ledger Audit:**
    *   Real-time view of system liabilities and platform fee accumulation.