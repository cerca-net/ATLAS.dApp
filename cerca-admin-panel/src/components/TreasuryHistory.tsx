import { useState, useEffect, useCallback } from 'react';
import { apiFetch, shortHash, formatNumber, timeAgo } from '../api';
import { ArrowUpRight, ArrowDownLeft, RefreshCw, Landmark } from 'lucide-react';

interface TxRecord {
  blockIndex: number;
  blockHash: string;
  type: string;
  sender: string;
  recipient: string;
  amount: number;
  fee: number;
  data: string;
  timestamp: number;
  direction: 'in' | 'out';
}

export function TreasuryHistoryPage() {
  const [history, setHistory] = useState<TxRecord[]>([]);
  const [treasuryAddr, setTreasuryAddr] = useState('');
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState<TxRecord | null>(null);

  const fetchHistory = useCallback(async () => {
    try {
      const res = await apiFetch('/admin/treasury-history');
      if (res.ok) {
        const data = await res.json();
        setTreasuryAddr(data.treasuryAddress || '');
        setHistory(data.transactions || []);
      }
    } catch (e) { console.error(e); }
    setLoading(false);
  }, []);

  useEffect(() => {
    fetchHistory();
    const id = setInterval(fetchHistory, 10000);
    return () => clearInterval(id);
  }, [fetchHistory]);

  const totalOut = history.filter(t => t.direction === 'out').reduce((s, t) => s + t.amount, 0);
  const totalIn  = history.filter(t => t.direction === 'in').reduce((s, t) => s + t.amount, 0);
  const totalFees = history.reduce((s, t) => s + t.fee, 0);

  return (
    <>
      <div className="content-header">
        <h1>Treasury Transaction History</h1>
        <p>All confirmed on-chain transactions from or to the treasury wallet.</p>
      </div>

      {/* Treasury Address */}
      {treasuryAddr && (
        <div style={{ padding: '0.875rem 1.25rem', background: 'rgba(255,255,255,0.03)', borderRadius: '12px', marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '1rem' }}>
          <Landmark size={20} color="var(--accent-success)" />
          <div>
            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.1rem' }}>Treasury Wallet Address</div>
            <div style={{ fontFamily: 'monospace', fontSize: '0.85rem', wordBreak: 'break-all' }}>{treasuryAddr}</div>
          </div>
          <button
            onClick={fetchHistory}
            style={{ marginLeft: 'auto', background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-muted)', display: 'flex', alignItems: 'center', gap: '0.4rem', fontSize: '0.8rem' }}
          >
            <RefreshCw size={14} /> Refresh
          </button>
        </div>
      )}

      {/* Stats */}
      <div className="stat-grid" style={{ marginBottom: '1.5rem' }}>
        <div className="stat-card">
          <div className="stat-card-title">Total Transactions</div>
          <div className="stat-card-value">{history.length}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title" style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
            <ArrowUpRight size={13} color="var(--accent-danger)" /> Total Sent
          </div>
          <div className="stat-card-value" style={{ color: 'var(--accent-danger)' }}>{formatNumber(totalOut)} T</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title" style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
            <ArrowDownLeft size={13} color="var(--accent-success)" /> Total Received
          </div>
          <div className="stat-card-value" style={{ color: 'var(--accent-success)' }}>{formatNumber(totalIn)} T</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title">Fees Paid</div>
          <div className="stat-card-value" style={{ color: 'var(--text-muted)' }}>{formatNumber(totalFees)} T</div>
        </div>
      </div>

      {/* Detail Panel */}
      {selected && (
        <div className="glass-card" style={{ marginBottom: '1.5rem', borderLeft: `4px solid ${selected.direction === 'out' ? 'var(--accent-danger)' : 'var(--accent-success)'}`, animation: 'fadeIn 0.2s ease-out' }}>
          <div className="glass-card-header">
            <div className="glass-card-title">
              {selected.direction === 'out'
                ? <ArrowUpRight size={18} color="var(--accent-danger)" />
                : <ArrowDownLeft size={18} color="var(--accent-success)" />}
              Transaction Detail — Block #{selected.blockIndex}
            </div>
            <button onClick={() => setSelected(null)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-muted)' }}>✕</button>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '0.75rem', marginBottom: '1rem' }}>
            {[
              { label: 'Direction', value: selected.direction === 'out' ? '⬆ SENT' : '⬇ RECEIVED', color: selected.direction === 'out' ? 'var(--accent-danger)' : 'var(--accent-success)' },
              { label: 'Amount', value: formatNumber(selected.amount) + ' TCOIN', color: 'var(--accent-success)' },
              { label: 'Fee', value: formatNumber(selected.fee) + ' T', color: 'var(--text-muted)' },
              { label: 'Type', value: selected.type || 'regular', color: '#a78bfa' },
              { label: 'Block', value: `#${selected.blockIndex}`, color: 'var(--accent-primary)' },
              { label: 'Age', value: selected.timestamp ? timeAgo(selected.timestamp) : '—', color: 'var(--text-muted)' },
            ].map(({ label, value, color }) => (
              <div key={label} style={{ padding: '0.75rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px' }}>
                <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', marginBottom: '0.2rem' }}>{label}</div>
                <div style={{ fontWeight: 600, color }}>{value}</div>
              </div>
            ))}
          </div>
          {[
            { label: 'From', value: selected.sender },
            { label: 'To', value: selected.recipient },
            { label: 'Block Hash', value: selected.blockHash },
          ].map(({ label, value }) => (
            <div key={label} style={{ padding: '0.6rem 0.875rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', marginBottom: '0.4rem' }}>
              <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', marginBottom: '0.1rem' }}>{label}</div>
              <div style={{ fontFamily: 'monospace', fontSize: '0.8rem', wordBreak: 'break-all' }}>{value || '—'}</div>
            </div>
          ))}
          {selected.data && (
            <div style={{ padding: '0.6rem 0.875rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', marginTop: '0.25rem' }}>
              <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', marginBottom: '0.1rem' }}>Data</div>
              <pre style={{ fontFamily: 'monospace', fontSize: '0.78rem', color: '#e2e8f0', whiteSpace: 'pre-wrap', wordBreak: 'break-all', margin: 0 }}>
                {(() => { try { return JSON.stringify(JSON.parse(selected.data), null, 2); } catch { return selected.data; } })()}
              </pre>
            </div>
          )}
        </div>
      )}

      {/* History Table */}
      <div className="glass-card">
        <div className="glass-card-header">
          <div className="glass-card-title">
            <Landmark size={18} color="var(--accent-success)" />
            Ledger <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', fontWeight: 400 }}>(10s refresh · click row for detail)</span>
          </div>
        </div>

        {loading ? (
          <div style={{ color: 'var(--text-muted)', textAlign: 'center', padding: '2rem' }}>Loading history...</div>
        ) : history.length === 0 ? (
          <div style={{ color: 'var(--text-muted)', textAlign: 'center', padding: '3rem' }}>
            No treasury transactions yet. Use the Treasury page to send TCOIN — it will appear here once included in a block.
          </div>
        ) : (
          <table className="table-glass">
            <thead>
              <tr><th>Dir</th><th>Block</th><th>Type</th><th>From</th><th>To</th><th>Amount</th><th>Fee</th><th>Age</th></tr>
            </thead>
            <tbody>
              {[...history].reverse().map((tx, i) => (
                <tr key={i} style={{ cursor: 'pointer' }} onClick={() => setSelected(tx)}>
                  <td>
                    {tx.direction === 'out'
                      ? <ArrowUpRight size={16} color="var(--accent-danger)" />
                      : <ArrowDownLeft size={16} color="var(--accent-success)" />}
                  </td>
                  <td style={{ color: 'var(--accent-primary)', fontWeight: 600 }}>#{tx.blockIndex}</td>
                  <td>
                    <span style={{ padding: '0.15rem 0.45rem', borderRadius: '4px', fontSize: '0.72rem', background: 'rgba(99,102,241,0.15)', color: 'var(--accent-secondary)' }}>
                      {tx.type || 'regular'}
                    </span>
                  </td>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{shortHash(tx.sender, 8)}</td>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{shortHash(tx.recipient, 8)}</td>
                  <td style={{ fontWeight: 700, color: tx.direction === 'out' ? 'var(--accent-danger)' : 'var(--accent-success)' }}>
                    {tx.direction === 'out' ? '−' : '+'}{formatNumber(tx.amount)} T
                  </td>
                  <td style={{ color: 'var(--text-muted)', fontSize: '0.82rem' }}>{tx.fee}</td>
                  <td style={{ color: 'var(--text-muted)', fontSize: '0.82rem' }}>{tx.timestamp ? timeAgo(tx.timestamp) : '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </>
  );
}
