import { useState, useEffect, useCallback } from 'react';
import { Droplet, Send, AlertTriangle, CheckCircle } from 'lucide-react';
import { apiFetch } from '../api';

const NODE_URL = 'http://localhost:8080';

export function FaucetPage() {
  const [address, setAddress] = useState('');
  const [amount, setAmount] = useState('1000');
  const [lastTxHash, setLastTxHash] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [treasuryData, setTreasuryData] = useState({ balance: '...', epoch: '...', mempoolCount: '...' });

  const fetchTreasury = useCallback(async () => {
    try {
      const [treasRes, statusRes, mempoolRes] = await Promise.all([
        fetch(`${NODE_URL}/treasury`),
        fetch(`${NODE_URL}/monitoring/status`),
        fetch(`${NODE_URL}/mempool`),
      ]);
      if (treasRes.ok) {
        const d = await treasRes.json();
        const balance = typeof d.balance === 'number' ? d.balance.toLocaleString() + ' TCOIN' : (d.balance || '...');
        const epoch = statusRes.ok ? (await statusRes.json()).block_height ?? '...' : '...';
        const mCount = mempoolRes.ok ? (await mempoolRes.json() || []).length : '...';
        setTreasuryData({ balance, epoch: String(epoch), mempoolCount: String(mCount) });
      }
    } catch (e) {
      console.error(e);
    }
  }, []);

  useEffect(() => {
    fetchTreasury();
    const id = setInterval(fetchTreasury, 10000);
    return () => clearInterval(id);
  }, [fetchTreasury]);

  const handleDrain = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!address || !amount || submitting) return;
    setSubmitting(true);
    setLastTxHash('');
    try {
      const response = await apiFetch('/admin/faucet', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ address, amount: parseInt(amount, 10) }),
      });
      if (!response.ok) throw new Error(await response.text());
      const result = await response.json();
      setLastTxHash(result.txHash || '');
      fetchTreasury();
    } catch (err: any) {
      alert(`Faucet Error: ${err.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <div className="content-header">
        <h1>Node Interface & Faucet</h1>
        <p>Distribute TCOIN initial allocations directly from the network treasury node.</p>
      </div>

      <div className="stat-grid">
        <div className="stat-card">
          <div className="stat-card-title">Treasury Balance</div>
          <div className="stat-card-value">{treasuryData.balance}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title">Chain Height</div>
          <div className="stat-card-value">{treasuryData.epoch}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title">Pending Mempool</div>
          <div className="stat-card-value">{treasuryData.mempoolCount} txs</div>
        </div>
      </div>

      <div className="glass-card" style={{ maxWidth: '600px' }}>
        <div className="glass-card-header">
          <div className="glass-card-title">
            <Droplet size={20} color="var(--accent-primary)" />
            Direct Faucet Emission
          </div>
        </div>
        
        <form onSubmit={handleDrain}>
          <div className="input-group">
            <label className="input-label">Target Wallet Address</label>
            <input 
              className="glass-input"
              type="text" 
              placeholder="e.g. 0x04ca..." 
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              required
            />
          </div>

          <div className="input-group" style={{ marginBottom: '2rem' }}>
            <label className="input-label">Amount (TCOIN)</label>
            <input 
              className="glass-input"
              type="number" 
              placeholder="1000" 
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              required
            />
          </div>

          <button type="submit" className="btn-primary" style={{ width: '100%' }} disabled={submitting}>
            <Send size={18} />
            {submitting ? 'Broadcasting...' : 'Authorize Transfer'}
          </button>

          {lastTxHash && (
            <div style={{ marginTop: '1rem', padding: '0.75rem 1rem', background: 'rgba(39,174,96,0.1)', borderRadius: '8px', display: 'flex', alignItems: 'flex-start', gap: '0.75rem' }}>
              <CheckCircle size={18} color="var(--accent-success)" style={{ flexShrink: 0, marginTop: '2px' }} />
              <div style={{ fontSize: '0.8rem' }}>
                <div style={{ color: 'var(--accent-success)', fontWeight: 600 }}>Transaction broadcast successfully!</div>
                <div style={{ color: 'var(--text-muted)', fontFamily: 'monospace', wordBreak: 'break-all', marginTop: '0.25rem' }}>TxHash: {lastTxHash}</div>
              </div>
            </div>
          )}
        </form>

        <div style={{ marginTop: '1.5rem', padding: '1rem', background: 'rgba(239, 68, 68, 0.05)', borderRadius: '8px', display: 'flex', gap: '1rem', alignItems: 'flex-start' }}>
          <AlertTriangle color="var(--accent-danger)" size={20} />
          <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>
            Operations performed through the master node wallet bypass normal fee structures. Ensure exact wallet match before authorizing fund transfers.
          </p>
        </div>
      </div>
    </>
  );
}
