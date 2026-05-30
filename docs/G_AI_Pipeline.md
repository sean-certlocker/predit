# G. Object Detection / AI Referee Pipeline

## Pipeline Architecture

1.  **Ingestion:** Stream media server (e.g., LiveKit) writes 2-second HLS chunks to a fast temporary storage (Redis or RAM disk) during the `EVENT_ACTIVE` window.
2.  **Processing Node (Python):**
    *   Pulls stream chunks.
    *   Extracts frames at 10 FPS.
3.  **Model Execution:**
    *   **Detection:** YOLOv8/v10 runs on each frame.
    *   **Tracking:** DeepSORT or ByteTrack assigns unique IDs to objects (e.g., `car_1`, `bus_2`).
4.  **Counting Logic (Virtual Line):**
    *   Market defines a vector `(A_x, A_y) -> (B_x, B_y)` as the "finish line".
    *   When tracking box centroid crosses the vector line, increment count.
5.  **Output Generation:**
    *   Aggregate total count over the `eventWindowSeconds`.
    *   Store "snapshot" frame of every line-crossing event for the audit log.
    *   POST result to Core API.

## AI Service Output Schema
```json
{
  "marketId": "uuid-here",
  "status": "resolved",
  "count": 23,
  "confidence": 0.91,
  "result": "YES",
  "events": [
    {
      "objectType": "car",
      "trackId": "car_12",
      "timestamp": "2026-05-30T12:00:04.123Z",
      "confidence": 0.94,
      "snapshotUrl": "s3://bucket/evidence/uuid-here/car_12.jpg"
    }
  ],
  "flags": []
}
```