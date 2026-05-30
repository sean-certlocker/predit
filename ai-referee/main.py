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
    
    # Simulate pulling chunks and running YOLOv8
    # In a real implementation, we would use cv2.VideoCapture(task.stream_url)
    # and run model(frame) for each frame in the window.
    
    await asyncio.sleep(task.window_seconds)
    
    # Mocked results
    count = random.randint(15, 25)
    confidence = random.uniform(0.8, 0.99)
    
    payload = {
        "market_id": task.market_id,
        "count": count,
        "confidence": confidence,
        "flags": []
    }
    
    try:
        res = requests.post(CORE_API_URL, json=payload)
        print(f"Sent resolution to core: {res.status_code}")
    except Exception as e:
        print(f"Failed to send resolution: {e}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
