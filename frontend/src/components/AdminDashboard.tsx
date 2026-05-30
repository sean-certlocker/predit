import { useState, useEffect } from 'react'
import axios from 'axios'
import type { Market, User } from '../types'
import { 
  ShieldAlert, 
  Users, 
  BarChart3, 
  Activity, 
  AlertTriangle,
  CheckCircle2
} from 'lucide-react'

const AdminDashboard = ({ markets }: { markets: Market[] }) => {
  const [activeTab, setActiveTab] = useState<'streams' | 'queue' | 'users' | 'audit'>('streams')
  const [riskyUsers, setRiskyUsers] = useState<User[]>([])
  const [auditData, setAuditData] = useState<any>(null)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [usersRes, auditRes] = await Promise.all([
          axios.get('http://localhost:8080/api/users/risk'),
          axios.get('http://localhost:8080/api/ledger/audit')
        ])
        setRiskyUsers(usersRes.data)
        setAuditData(auditRes.data)
      } catch (err) {
        console.error('Failed to fetch admin data', err)
      }
    }
    fetchData()
  }, [])

  const handleTriggerAI = async (marketId: string) => {
    try {
      await axios.post(`http://localhost:8080/api/admin/trigger-resolution?market_id=${marketId}`)
      alert('Resolution triggered')
    } catch (err) {
      alert('Failed to trigger resolution')
    }
  }

  const handleModerate = async (marketId: string, action: 'approve' | 'reject') => {
    try {
      await axios.post(`http://localhost:8080/api/admin/market/moderate?market_id=${marketId}&action=${action}`)
      alert(`Market ${action}d`)
    } catch (err) {
      alert('Failed to moderate market')
    }
  }

  return (
    <div className="admin-dashboard">
      <nav className="admin-nav">
        <button 
          className={activeTab === 'streams' ? 'active' : ''} 
          onClick={() => setActiveTab('streams')}
        >
          <Activity size={18} /> Streams Matrix
        </button>
        <button 
          className={activeTab === 'queue' ? 'active' : ''} 
          onClick={() => setActiveTab('queue')}
        >
          <ShieldAlert size={18} /> Market Queue
        </button>
        <button 
          className={activeTab === 'users' ? 'active' : ''} 
          onClick={() => setActiveTab('users')}
        >
          <Users size={18} /> User Risk
        </button>
        <button 
          className={activeTab === 'audit' ? 'active' : ''} 
          onClick={() => setActiveTab('audit')}
        >
          <BarChart3 size={18} /> Ledger Audit
        </button>
      </nav>

      <div className="admin-content">
        {activeTab === 'streams' && (
          <div className="streams-matrix">
            <div className="grid">
              {markets.map(m => (
                <div key={m.id} className="stream-card">
                  <div className="stream-header">
                    <span className={`health-indicator ${m.health?.toLowerCase() || 'green'}`}></span>
                    <h4>{m.title}</h4>
                  </div>
                  <div className="stream-stats">
                    <div className="stat">
                      <span className="label">Status</span>
                      <span className="value">{m.status}</span>
                    </div>
                    <div className="stat">
                      <span className="label">Health</span>
                      <span className="value">{m.health || 'Stable'}</span>
                    </div>
                  </div>
                  {m.safety_flags && m.safety_flags.length > 0 && (
                    <div className="safety-warning">
                      <AlertTriangle size={14} />
                      Flags: {m.safety_flags.join(', ')}
                    </div>
                  )}
                  <button 
                    className="btn-resolve"
                    disabled={m.status !== 'OPEN'}
                    onClick={() => handleTriggerAI(m.id)}
                  >
                    Trigger Resolution
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'users' && (
          <div className="user-risk-panel">
            <h3>Suspicious Activity</h3>
            <table className="admin-table">
              <thead>
                <tr>
                  <th>User ID</th>
                  <th>Username</th>
                  <th>Trust Score</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {riskyUsers.map(u => (
                  <tr key={u.id}>
                    <td>{u.id}</td>
                    <td>{u.username}</td>
                    <td>
                      <span className={`trust-score ${u.trust_score > 50 ? 'high' : 'low'}`}>
                        {u.trust_score}
                      </span>
                    </td>
                    <td>
                      {u.suspicious ? (
                        <span className="status-flag red">Suspicious</span>
                      ) : (
                        <span className="status-flag yellow">Low Trust</span>
                      )}
                    </td>
                    <td>
                      <button className="btn-table-action">Shadowban</button>
                      <button className="btn-table-action">Review</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {activeTab === 'audit' && auditData && (
          <div className="audit-panel">
            <div className="audit-grid">
              <div className="audit-card">
                <span className="label">Total System Balance</span>
                <span className="value">{auditData.total_system_balance} Play</span>
              </div>
              <div className="audit-card">
                <span className="label">Platform Fees</span>
                <span className="value">{auditData.platform_fees} Play</span>
              </div>
              <div className="audit-card">
                <span className="label">Pending Stakes</span>
                <span className="value">{auditData.pending_stakes} Play</span>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'queue' && (
          <div className="market-queue-panel">
            <h3>Pending Market Approvals</h3>
            {markets.filter(m => m.status === 'DRAFT').length > 0 ? (
              <table className="admin-table">
                <thead>
                  <tr>
                    <th>Title</th>
                    <th>Resolution</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {markets.filter(m => m.status === 'DRAFT').map(m => (
                    <tr key={m.id}>
                      <td>{m.title}</td>
                      <td>{m.resolution_method}</td>
                      <td>
                        <button className="btn-approve" onClick={() => handleModerate(m.id, 'approve')}>Approve</button>
                        <button className="btn-reject" onClick={() => handleModerate(m.id, 'reject')}>Reject</button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <div className="empty-state">
                <CheckCircle2 size={48} />
                <p>Queue is clear. All markets approved.</p>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

export default AdminDashboard
