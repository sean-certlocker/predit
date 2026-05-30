import asyncio
import requests
from fastapi import FastAPI, BackgroundTasks
from pydantic import BaseModel
import random
import time

app = FastAPI()

CORE_API_URL = "http://localhost:8080/internal/ai/resolution"

class ResolutionTask(BaseModel):
    market_id: str
    stream_url: str
    window_seconds: int

@app.post("/tasks")
async def create_task(task: ResolutionTask, background_tasks: BackgroundTasks):
    background_tasks.add_task(process_resolution, task)
    return {"status": "accepted", "market_id": task.market_id}

async def process_resolution(task: ResolutionTask):
    print(f"Starting resolution for market {task.market_id}")
    
    # PHASE 4/5: Advanced Simulation
    # 1. Simulate Object Tracking over the window
    start_time = time.time()
    total_detections = 0
    objects_tracked = set()
    
    # 2. Simulate Motion Tracking (Optical Flow)
    camera_instability_score = 0
    
    while time.time() - start_time < task.window_seconds:
        # Simulate frame processing
        await asyncio.sleep(1) 
        
        # Simulate object detection & tracking (Phase 3/4)
        new_objects = random.randint(0, 3)
        total_detections += new_objects
        for _ in range(new_objects):
            objects_tracked.add(f"obj_{random.randint(1000, 9999)}")
        
        # Simulate motion tracking (Phase 4: Anti-Manipulation)
        # 5% chance of a "camera knock" event per second
        if random.random() < 0.05:
            camera_instability_score += random.uniform(0.5, 1.0)

    # 3. Decision Logic
    count = len(objects_tracked)
    confidence = random.uniform(0.85, 0.99)
    flags = []
    
    # Phase 4 Rule: Camera stability
    if camera_instability_score > 1.5:
        flags.append("CAMERA_MOVEMENT_DETECTED")
        confidence = 0.4 # Force manual review or void
        
    # Phase 4 Rule: Safety Check (NSFW/Violence)
    if random.random() < 0.01: # 1% simulation of bad content
        flags.append("NSFW")
    
    payload = {
        "market_id": task.market_id,
        "count": count,
        "confidence": confidence,
        "flags": flags
    }
    
    try:
        res = requests.post(CORE_API_URL, json=payload)
        print(f"Sent resolution to core: {res.status_code}")
    except Exception as e:
        print(f"Failed to send resolution: {e}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
