import { useState, useEffect } from 'react'
import axios from 'axios'
import { Market } from './types'
import MarketViewer from './components/MarketViewer'
import { Activity, Wallet, Layout } from 'lucide-react'
import './App.css'

function App() {
  const [markets, setMarkets] = useState<Market[]>([])
  const [selectedMarket, setSelectedMarket] = useState<Market | null>(null)
  const [balance, setBalance] = useState<number>(0)
  const [loading, setLoading] = useState(true)

  const fetchMarkets = async () => {
    try {
      const res = await axios.get('http://localhost:8080/api/markets')
      setMarkets(res.data)
      if (res.data.length > 0 && !selectedMarket) {
        setSelectedMarket(res.data[0])
      }
    } catch (err) {
      console.error('Failed to fetch markets', err)
    } finally {
      setLoading(false)
    }
  }

  const fetchBalance = async () => {
    try {
      const res = await axios.get('http://localhost:8080/api/ledger/balance?user_id=user1')
      setBalance(res.data.balance)
    } catch (err) {
      console.error('Failed to fetch balance', err)
    }
  }

  useEffect(() => {
    fetchMarkets()
    fetchBalance()
    const interval = setInterval(fetchMarkets, 5000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="app-container">
      <header className="app-header">
        <div className="brand">
          <Activity className="logo-icon" />
          <h1>PREDIT</h1>
        </div>
        <div className="user-nav">
          <div className="balance-badge">
            <Wallet size={16} />
            <span>{balance.toFixed(2)} Play</span>
          </div>
          <button className="btn-primary">Connect</button>
        </div>
      </header>

      <main className="app-main">
        <aside className="sidebar">
          <div className="sidebar-header">
            <Layout size={18} />
            <h2>Markets</h2>
          </div>
          <div className="market-list">
            {loading ? (
              <p className="loading-text">Loading markets...</p>
            ) : (
              markets.map(m => (
                <div 
                  key={m.id} 
                  className={`market-item ${selectedMarket?.id === m.id ? 'active' : ''}`}
                  onClick={() => setSelectedMarket(m)}
                >
                  <p className="market-item-title">{m.title}</p>
                  <span className={`status-pill ${m.status.toLowerCase()}`}>
                    {m.status}
                  </span>
                </div>
              ))
            )}
          </div>
        </aside>

        <section className="viewer-section">
          {selectedMarket ? (
            <MarketViewer market={selectedMarket} />
          ) : (
            <div className="empty-state">
              <p>Select a market to start predicting</p>
            </div>
          )}
        </section>
      </main>
    </div>
  )
}

export default App
