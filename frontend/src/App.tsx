import { useState, useEffect } from 'react'
import axios from 'axios'
import type { Market } from './types'
import MarketViewer from './components/MarketViewer'
import AdminDashboard from './components/AdminDashboard'
import MarketCreator from './components/MarketCreator'
import CreatorView from './components/CreatorView'
import { Activity, Wallet, Layout, Settings, Plus, Camera } from 'lucide-react'
import './App.css'

function App() {
  const [markets, setMarkets] = useState<Market[]>([])
  const [selectedMarket, setSelectedMarket] = useState<Market | null>(null)
  const [balance, setBalance] = useState<number>(0)
  const [loading, setLoading] = useState(true)
  const [view, setView] = useState<'viewer' | 'admin' | 'creator'>('viewer')
  const [showCreator, setShowCreator] = useState(false)

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
          <button 
            className={`btn-secondary ${view === 'creator' ? 'active' : ''}`}
            onClick={() => setView('creator')}
          >
            <Camera size={16} /> Studio
          </button>
          <button 
            className={`btn-secondary ${view === 'admin' ? 'active' : ''}`}
            onClick={() => setView('admin')}
          >
            <Settings size={16} /> Admin
          </button>
          <div className="balance-badge" onClick={() => setView('viewer')}>
            <Wallet size={16} />
            <span>{balance.toFixed(2)} Play</span>
          </div>
          <button className="btn-primary">Connect</button>
        </div>
      </header>

      <main className="app-main">
        {view === 'viewer' && (
          <>
            <aside className="sidebar">
              <div className="sidebar-header">
                <Layout size={18} />
                <h2>Markets</h2>
                <button 
                  className={`btn-icon ${showCreator ? 'active' : ''}`} 
                  onClick={() => setShowCreator(!showCreator)}
                >
                  <Plus size={18} />
                </button>
              </div>
              <div className="market-list">
                {showCreator && <MarketCreator onCreated={() => {
                  setShowCreator(false)
                  fetchMarkets()
                }} />}
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
          </>
        )}

        {view === 'admin' && <AdminDashboard markets={markets} />}
        
        {view === 'creator' && <CreatorView />}
      </main>
    </div>
  )
}

export default App
