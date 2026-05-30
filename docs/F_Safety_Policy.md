# F. Safety & Moderation Policy

## Pre-Market Moderation
*   **Template Restriction:** No free-text markets allowed for MVP.
*   **Category Enforcement:** Markets must fall into predefined safe categories: `TRAFFIC`, `WEATHER`, `GAMING_(SAFE)`, `SPORTS_PRACTICE`.
*   **Auto-Rejection Words:** Filter for any words associated with harm, self-harm, weapons, drugs, or illegal acts.

## Live Moderation
*   **AI Vision Scanning:** Sample 1 frame per second. Run against a "Safety Check" vision model (e.g., Azure Content Safety or AWS Rekognition).
*   **Triggers:** If `violence`, `weapons`, or `nudity` probability > 0.5:
    *   *Action:* Immediately pause stream, switch market state to `MANUAL_REVIEW`, lock betting.
*   **Report System:** If > 5% of viewers report the stream within a 30s window, auto-escalate to `MANUAL_REVIEW`.

## Resolution Rules
*   Never auto-resolve a market where a safety flag was raised during the event window.

## Functions

```typescript
function moderateMarket(marketDraft: MarketDraft): 'APPROVED' | 'REJECTED' | 'MODERATION_PENDING' {
    if (!allowedTemplates.includes(marketDraft.templateId)) return 'REJECTED';
    if (containsBannedWords(marketDraft.title)) return 'REJECTED';
    if (creatorTrustScore < 50) return 'MODERATION_PENDING';
    return 'APPROVED';
}

function moderateLiveFrame(frameData: FrameData): 'SAFE' | 'UNSAFE' | 'REQUIRES_REVIEW' {
    if (frameData.nudityScore > 0.8 || frameData.violenceScore > 0.8) return 'UNSAFE';
    if (frameData.nudityScore > 0.4 || frameData.violenceScore > 0.4) return 'REQUIRES_REVIEW';
    return 'SAFE';
}
```