import { useState, useEffect, useCallback } from 'react';
import { apiFetch, timeAgo, shortHash, formatNumber } from '../api';
import { Database, ChevronRight, X, ArrowLeft } from 'lucide-react';

export function BlocksPage() {
  const [blocks, setBlocks] = useState<any[]>([]);
  const [selected, setSelected] = useState<any>(null);
  const [page, setPage] = useState(0);
  const pageSize = 15;

  const fetchBlocks = useCallback(async () => {
    try {
      const res = await apiFetch(`/blocks?limit=100&offset=0`);
      if (res.ok) {
        const data = await res.json();
        setBlocks(Array.isArray(data) ? data : []);
      }
    } catch (e) { console.error(e); }
  }, []);

  useEffect(() => {
    fetchBlocks();
    const id = setInterval(fetchBlocks, 5000);
    return () => clearInterval(id);
  }, [fetchBlocks]);

  const paged = blocks.slice(page * pageSize, page * pageSize + pageSize);
  const totalPages = Math.ceil(blocks.length / pageSize);

  return (
    <>
      <div className="content-header">
        <h1>Block Explorer</h1>
        <p>Inspect every block, its transactions, validator, and cryptographic hash.</p>
      </div>

      {selected ? (
        /* Block Detail View */
        <div style={{ animation: 'fadeIn 0.2s ease-out' }}>
          <button
            onClick={() => setSelected(null)}
            className="btn-secondary"
            style={{ marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}
          >
            <ArrowLeft size={16} /> Back to blocks
          </button>

          <div className="glass-card" style={{ marginBottom: '1.5rem' }}>
            <div className="glass-card-header">
              <div className="glass-card-title">
                <Database size={18} color="var(--accent-primary)" />
                Block #{selected.Index ?? selected.index}
              </div>
              <button onClick={() => setSelected(null)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-muted)' }}>
                <X size={18} />
              </button>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '1rem', marginBottom: '1.5rem' }}>
              {[
                { label: 'Block Height', value: `#${selected.Index ?? selected.index}` },
                { label: 'Timestamp', value: new Date((selected.Timestamp || selected.timestamp || 0) * 1000).toLocaleString() },
                { label: 'Transaction Count', value: (selected.Transactions || selected.transactions || []).length },
                { label: 'Confirmed', value: timeAgo(selected.Timestamp || selected.timestamp || 0) },
              ].map(({ label, value }) => (
                <div key={label} style={{ padding: '1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px' }}>
                  <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.3rem' }}>{label}</div>
                  <div style={{ fontWeight: 600 }}>{value}</div>
                </div>
              ))}
            </div>

            <div style={{ padding: '1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', marginBottom: '1rem' }}>
              <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.3rem' }}>Block Hash</div>
              <div style={{ fontFamily: 'monospace', fontSize: '0.82rem', wordBreak: 'break-all', color: 'var(--accent-primary)' }}>{selected.Hash || selected.hash || '—'}</div>
            </div>

            <div style={{ padding: '1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', marginBottom: '1rem' }}>
              <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.3rem' }}>Previous Hash</div>
              <div style={{ fontFamily: 'monospace', fontSize: '0.82rem', wordBreak: 'break-all', color: 'var(--text-muted)' }}>{selected.PrevHash || selected.prevHash || '—'}</div>
            </div>

            <div style={{ padding: '1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px' }}>
              <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.3rem' }}>Validator (Public Key)</div>
              <div style={{ fontFamily: 'monospace', fontSize: '0.82rem', wordBreak: 'break-all', color: '#a78bfa' }}>{selected.Validator || selected.validator || '—'}</div>
            </div>
          </div>

          {/* Transactions in Block */}
          <div className="glass-card">
            <div className="glass-card-header">
              <div className="glass-card-title">Transactions in Block ({(selected.Transactions || selected.transactions || []).length})</div>
            </div>
            {(selected.Transactions || selected.transactions || []).length === 0 ? (
              <div style={{ color: 'var(--text-muted)', textAlign: 'center', padding: '1.5rem' }}>No transactions in this block.</div>
            ) : (
              <table className="table-glass">
                <thead>
                  <tr><th>Type</th><th>From</th><th>To</th><th>Amount</th><th>Fee</th><th>Data</th></tr>
                </thead>
                <tbody>
                  {(selected.Transactions || selected.transactions).map((tx: any, i: number) => (
                    <tr key={i}>
                      <td><span style={{ padding: '0.2rem 0.5rem', background: 'rgba(99,102,241,0.2)', borderRadius: '4px', fontSize: '0.75rem', color: 'var(--accent-secondary)' }}>{tx.type || 'regular'}</span></td>
                      <td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{shortHash(tx.sender, 8)}</td>
                      <td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{shortHash(tx.recipient, 8)}</td>
                      <td style={{ color: 'var(--accent-success)', fontWeight: 600 }}>{formatNumber(tx.amount)} T</td>
                      <td style={{ color: 'var(--text-muted)' }}>{tx.fee}</td>
                      <td style={{ maxWidth: '150px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', fontSize: '0.78rem', color: 'var(--text-muted)' }}>{tx.data || '—'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </div>
      ) : (
        /* Block List */
        <>
          <div className="glass-card">
            <div className="glass-card-header">
              <div className="glass-card-title">
                <Database size={18} color="var(--accent-primary)" />
                All Blocks ({blocks.length}) <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', fontWeight: 400 }}>(live · 5s)</span>
              </div>
              <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
                <button className="btn-secondary" onClick={() => setPage(p => Math.max(0, p - 1))} disabled={page === 0} style={{ padding: '0.4rem 0.8rem' }}>←</button>
                <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>Page {page + 1}/{totalPages || 1}</span>
                <button className="btn-secondary" onClick={() => setPage(p => Math.min(totalPages - 1, p + 1))} disabled={page >= totalPages - 1} style={{ padding: '0.4rem 0.8rem' }}>→</button>
              </div>
            </div>

            <table className="table-glass">
              <thead>
                <tr>
                  <th>Height</th>
                  <th>Hash</th>
                  <th>Prev Hash</th>
                  <th>Txs</th>
                  <th>Validator</th>
                  <th>Age</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {paged.length ? paged.map((b: any) => (
                  <tr key={b.Index ?? b.index} style={{ cursor: 'pointer' }} onClick={() => setSelected(b)}>
                    <td style={{ color: 'var(--accent-primary)', fontWeight: 700 }}>#{b.Index ?? b.index}</td>
                    <td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{shortHash(b.Hash || b.hash)}</td>
                    <td style={{ fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--text-muted)' }}>{shortHash(b.PrevHash || b.prevHash)}</td>
                    <td>{(b.Transactions || b.transactions || []).length}</td>
                    <td style={{ fontFamily: 'monospace', fontSize: '0.8rem', color: '#a78bfa' }}>{shortHash(b.Validator || b.validator, 8)}</td>
                    <td style={{ color: 'var(--text-muted)' }}>{timeAgo(b.Timestamp || b.timestamp || 0)}</td>
                    <td><ChevronRight size={16} color="var(--text-muted)" /></td>
                  </tr>
                )) : (
                  <tr><td colSpan={7} style={{ textAlign: 'center', color: 'var(--text-muted)' }}>No blocks found. Start the node to begin mining.</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </>
      )}
    </>
  );
}
