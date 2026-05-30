import { useState } from 'react'
import { 
  LiveKitRoom, 
  VideoConference, 
  RoomAudioRenderer,
  ControlBar
} from '@livekit/components-react'
import '@livekit/components-styles'
import axios from 'axios'
import { Camera, VideoOff, Share2, Activity } from 'lucide-react'

const CreatorView = () => {
  const [token, setToken] = useState<string | null>(null)
  const [roomName, setRoomName] = useState('my-stream')
  const [isLive, setIsLive] = useState(false)

  const startStream = async () => {
    try {
      const res = await axios.get(`http://localhost:8080/api/streaming/token?room=${roomName}&identity=creator_${Math.floor(Math.random() * 1000)}`)
      setToken(res.data.token)
      setIsLive(true)
    } catch (err) {
      alert('Failed to start stream. Ensure LiveKit is running.')
    }
  }

  return (
    <div className="creator-view">
      <div className="creator-header">
        <h2><Camera size={24} /> Creator Studio</h2>
        <div className="stream-controls">
          <input 
            type="text" 
            value={roomName} 
            onChange={(e) => setRoomName(e.target.value)} 
            placeholder="Room Name"
            disabled={isLive}
          />
          {!isLive ? (
            <button className="btn-primary" onClick={startStream}>Go Live</button>
          ) : (
            <button className="btn-danger" onClick={() => setIsLive(false)}>End Stream</button>
          )}
        </div>
      </div>

      <div className="stream-preview">
        {isLive && token ? (
          <LiveKitRoom
            video={true}
            audio={true}
            token={token}
            serverUrl="ws://localhost:7880"
            onDisconnected={() => setIsLive(false)}
            style={{ height: '500px' }}
          >
            <VideoConference />
            <RoomAudioRenderer />
            <ControlBar />
          </LiveKitRoom>
        ) : (
          <div className="preview-placeholder">
            <VideoOff size={64} />
            <p>Camera is off. Click "Go Live" to start streaming.</p>
          </div>
        )}
      </div>

      <div className="creator-stats">
        <div className="stat-card">
          <Share2 size={18} />
          <span className="label">Viewers</span>
          <span className="value">0</span>
        </div>
        <div className="stat-card">
          <Activity size={18} />
          <span className="label">Bitrate</span>
          <span className="value">-- kbps</span>
        </div>
      </div>
    </div>
  )
}

export default CreatorView
