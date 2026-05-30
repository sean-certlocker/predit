import { useState } from 'react'
import { PlusCircle, Target, Clock as ClockIcon } from 'lucide-react'

const MarketCreator = ({ onCreated }: { onCreated: () => void }) => {
  const [title, setTitle] = useState('')
  const [window, setWindow] = useState(60)
  const [target, setTarget] = useState(20)
  const [isCreating, setIsCreating] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsCreating(true)
    try {
      // In a real app, this would be a POST to /api/markets
      // For MVP, we simulate it
      console.log('Creating market:', { title, window, target })
      alert('Market submitted for approval!')
      onCreated()
      setTitle('')
    } catch (err) {
      alert('Failed to create market')
    } finally {
      setIsCreating(false)
    }
  }

  return (
    <div className="market-creator">
      <h3>Create New Prediction</h3>
      <form onSubmit={handleSubmit} className="creator-form">
        <div className="form-group">
          <label>Question / Title</label>
          <input 
            type="text" 
            placeholder="e.g. Will 10 buses pass by?" 
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>
        <div className="form-row">
          <div className="form-group">
            <label><ClockIcon size={14} /> Window (sec)</label>
            <input 
              type="number" 
              value={window}
              onChange={(e) => setWindow(Number(e.target.value))}
            />
          </div>
          <div className="form-group">
            <label><Target size={14} /> Target Count</label>
            <input 
              type="number" 
              value={target}
              onChange={(e) => setTarget(Number(e.target.value))}
            />
          </div>
        </div>
        <button type="submit" className="btn-create" disabled={isCreating}>
          <PlusCircle size={18} />
          {isCreating ? 'Submitting...' : 'Submit to Queue'}
        </button>
      </form>
    </div>
  )
}

export default MarketCreator
