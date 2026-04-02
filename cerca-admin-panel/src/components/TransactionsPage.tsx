import { useState, useEffect, useCallback } from 'react';
import { apiFetch, shortHash, formatNumber } from '../api';
import { Hash, FileText, Clock, Filter } from 'lucide-react';

const TYPE_COLOR: Record<string, string> = {
  regular: '#38bdf8',
  call_contract: '#a78bfa',
  deploy_contract: '#f472b6',
  stake: 'var(--accent-success)',
  unstake: '#f59e0b',
  zk_proof: '#34d399',
};

export function TransactionsPage() {
  const [mempool, setMempool] = useState<any[]>([]);
  const [selected, setSelected] = useState<any>(null);
  const [filter, setFilter] = useState('');

  const fetchMempool = useCallback(async () => {
    try {
      const res = await apiFetch('/mempool');
      if (res.ok) {
        const data = await res.json();
        setMempool(Array.isArray(data) ? data : []);
      }
    } catch (e) { console.error(e); }
  }, []);

  useEffect(() => {
    fetchMempool();
    const id = setInterval(fetchMempool, 3000);
    return () => clearInterval(id);
  }, [fetchMempool]);

  const filtered = filter
    ? mempool.filter(tx =>
        (tx.sender || '').includes(filter) ||
        (tx.recipient || '').includes(filter) ||
        (tx.type || '').includes(filter) ||
        (tx.data || '').toLowerCase().includes(filter.toLowerCase())
      )
    : mempool;

  return (
    <>
      <div className="content-header">
        <h1>Transactions & Mempool</h1>
        <p>All pending transactions waiting to be included in the next block. Live updates every 3s.</p>
      </div>

      <div className="stat-grid" style={{ marginBottom: '1.5rem' }}>
        <div className="stat-card">
          <div className="stat-card-title"><Clock size={13} /> Pending</div>
          <div className="stat-card-value" style={{ color: '#f59e0b' }}>{mempool.length}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title">Regular</div>
          <div className="stat-card-value">{mempool.filter(t => t.type === 'regular').length}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title">Contract Calls</div>
          <div className="stat-card-value" style={{ color: '#a78bfa' }}>{mempool.filter(t => t.type === 'call_contract').length}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title">Staking</div>
          <div className="stat-card-value" style={{ color: 'var(--accent-success)' }}>{mempool.filter(t => t.type === 'stake' || t.type === 'unstake').length}</div>
        </div>
      </div>

      {selected && (
        <div className="glass-card" style={{ marginBottom: '1.5rem', animation: 'fadeIn 0.2s ease-out', borderLeft: '4px solid var(--accent-secondary)' }}>
          <div className="glass-card-header">
            <div className="glass-card-title"><FileText size={18} color="var(--accent-secondary)" /> Transaction Detail</div>
            <button onClick={() => setSelected(null)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-muted)', padding: '0.25rem' }}>✕</button>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '1rem', marginBottom: '1rem' }}>
            {[
              { label: 'Type', value: selected.type || 'regular', accent: TYPE_COLOR[selected.type] },
              { label: 'Amount', value: formatNumber(selected.amount) + ' T', accent: 'var(--accent-success)' },
              { label: 'Fee', value: formatNumber(selected.fee) + ' T', accent: 'var(--text-muted)' },
              { label: 'Nonce', value: selected.nonce ?? '—', accent: 'var(--text-muted)' },
              { label: 'Timestamp', value: selected.timestamp ? new Date(selected.timestamp * 1000).toLocaleString() : '—', accent: 'var(--text-muted)' },
              { label: 'Encrypted', value: selected.is_encrypted ? 'Yes' : 'No', accent: 'var(--text-muted)' },
            ].map(({ label, value, accent }) => (
              <div key={label} style={{ padding: '0.75rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px' }}>
                <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.2rem' }}>{label}</div>
                <div style={{ fontWeight: 600, color: accent || '#e2e8f0' }}>{value}</div>
              </div>
            ))}
          </div>

          {[
            { label: 'From (Sender)', value: selected.sender },
            { label: 'To (Recipient)', value: selected.recipient },
            { label: 'Sender Public Key', value: selected.senderPublicKey },
            { label: 'Signature', value: selected.signature },
          ].map(({ label, value }) => (
            <div key={label} style={{ padding: '0.75rem 1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', marginBottom: '0.5rem' }}>
              <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.2rem' }}>{label}</div>
              <div style={{ fontFamily: 'monospace', fontSize: '0.8rem', wordBreak: 'break-all' }}>{value || '—'}</div>
            </div>
          ))}

          {selected.data && (
            <div style={{ padding: '0.75rem 1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px' }}>
              <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.2rem' }}>Data Payload</div>
              <pre style={{ fontFamily: 'monospace', fontSize: '0.78rem', color: '#e2e8f0', wordBreak: 'break-all', whiteSpace: 'pre-wrap', margin: 0 }}>
                {(() => { try { return JSON.stringify(JSON.parse(selected.data), null, 2); } catch { return selected.data; } })()}
              </pre>
            </div>
          )}
        </div>
      )}

      <div className="glass-card">
        <div className="glass-card-header">
          <div className="glass-card-title"><Hash size={18} color="var(--accent-primary)" /> Mempool <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', fontWeight: 400 }}>(live · 3s)</span></div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <Filter size={15} color="var(--text-muted)" />
            <input
              className="glass-input"
              placeholder="Filter by address, type, data..."
              value={filter}
              onChange={e => setFilter(e.target.value)}
              style={{ height: '36px', width: '260px', fontSize: '0.82rem' }}
            />
          </div>
        </div>

        {filtered.length === 0 ? (
          <div style={{ color: 'var(--text-muted)', textAlign: 'center', padding: '3rem', fontSize: '0.9rem' }}>
            {mempool.length === 0 ? '✓ Mempool is empty — all transactions have been processed.' : 'No transactions match your filter.'}
          </div>
        ) : (
          <table className="table-glass">
            <thead>
              <tr><th>Type</th><th>From</th><th>To</th><th>Amount</th><th>Fee</th><th>Nonce</th><th>Data</th><th></th></tr>
            </thead>
            <tbody>
              {filtered.map((tx: any, i) => (
                <tr key={i} style={{ cursor: 'pointer' }} onClick={() => setSelected(tx)}>
                  <td>
                    <span style={{ padding: '0.2rem 0.5rem', borderRadius: '4px', fontSize: '0.72rem', fontWeight: 600, background: `${TYPE_COLOR[tx.type] || '#38bdf8'}22`, color: TYPE_COLOR[tx.type] || '#38bdf8' }}>
                      {tx.type || 'regular'}
                    </span>
                  </td>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{shortHash(tx.sender, 8)}</td>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{shortHash(tx.recipient, 8)}</td>
                  <td style={{ color: 'var(--accent-success)', fontWeight: 600 }}>{formatNumber(tx.amount)} T</td>
                  <td style={{ color: 'var(--text-muted)' }}>{tx.fee}</td>
                  <td style={{ color: 'var(--text-muted)' }}>{tx.nonce}</td>
                  <td style={{ maxWidth: '120px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', fontSize: '0.78rem', color: 'var(--text-muted)' }}>{tx.data || '—'}</td>
                  <td style={{ color: 'var(--text-muted)', fontSize: '0.75rem' }}>→</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </>
  );
}
