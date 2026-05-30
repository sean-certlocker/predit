import { useState, useEffect } from 'react'
import axios from 'axios'
import type { Market } from '../types'
import { 
  LiveKitRoom, 
  VideoConference, 
  RoomAudioRenderer,
  ControlBar
} from '@livekit/components-react'
import '@livekit/components-styles'
import { TrendingUp, TrendingDown, Clock, ShieldCheck } from 'lucide-react'

interface MarketViewerProps {
  market: Market
}

const MarketViewer = ({ market }: MarketViewerProps) => {
  const [token, setToken] = useState<string | null>(null)
  const [stake, setStake] = useState<number>(10)
  const [isPlacing, setIsPlacing] = useState(false)

  useEffect(() => {
    const fetchToken = async () => {
      try {
        const res = await axios.get(`http://localhost:8080/api/streaming/token?room=${market.id}&identity=viewer_${Math.floor(Math.random() * 1000)}`)
        setToken(res.data.token)
      } catch (err) {
        console.error('Failed to fetch stream token', err)
      }
    }
    fetchToken()
  }, [market.id])

  const handlePlaceStake = (outcome: 'YES' | 'NO') => {
    setIsPlacing(true)
    // Simulate API call
    setTimeout(() => {
      alert(`Placed ${stake} Play stake on ${outcome}`)
      setIsPlacing(false)
    }, 1000)
  }

  return (
    <div className="market-viewer">
      <div className="stream-container">
        {token ? (
          <LiveKitRoom
            video={true}
            audio={true}
            token={token}
            serverUrl="ws://localhost:7880" // Default local LiveKit URL
            data-lk-theme="default"
            style={{ height: '400px' }}
          >
            <VideoConference />
            <RoomAudioRenderer />
            <ControlBar />
          </LiveKitRoom>
        ) : (
          <div className="stream-placeholder">
            <div className="pulse-icon">
              <Clock size={48} />
            </div>
            <p>Waiting for Stream...</p>
            <small>Ensure backend is running with LiveKit configured</small>
          </div>
        )}
      </div>

      <div className="market-details">
        <div className="details-header">
          <h2>{market.title}</h2>
          <div className="meta-info">
            <span className="info-item">
              <Clock size={14} />
              {market.event_window_seconds}s Window
            </span>
            <span className="info-item">
              <ShieldCheck size={14} />
              AI Resolved
            </span>
          </div>
        </div>

        <div className="staking-card">
          <h3>Place Your Stake</h3>
          <div className="stake-input-group">
            <label>Amount (Play Money)</label>
            <input 
              type="number" 
              value={stake} 
              onChange={(e) => setStake(Number(e.target.value))}
              min={1}
            />
          </div>
          
          <div className="action-buttons">
            <button 
              className="btn-stake yes"
              disabled={isPlacing || market.status !== 'OPEN'}
              onClick={() => handlePlaceStake('YES')}
            >
              <TrendingUp size={18} />
              YES
            </button>
            <button 
              className="btn-stake no"
              disabled={isPlacing || market.status !== 'OPEN'}
              onClick={() => handlePlaceStake('NO')}
            >
              <TrendingDown size={18} />
              NO
            </button>
          </div>
          {market.status !== 'OPEN' && (
            <p className="status-warning">Betting is currently closed ({market.status})</p>
          )}
        </div>
      </div>
    </div>
  )
}

export default MarketViewer
