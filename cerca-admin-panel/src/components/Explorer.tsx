import { useState, useEffect, useCallback } from 'react';
import { Search, Database, FileText, Cpu, Activity } from 'lucide-react';

const NODE_URL = 'http://localhost:8080';

function timeAgo(unixTimestamp: number): string {
  const diff = Math.floor(Date.now() / 1000) - unixTimestamp;
  if (diff < 60) return `${diff}s ago`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  return `${Math.floor(diff / 3600)}h ago`;
}

export function ExplorerPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [searchResult, setSearchResult] = useState<any>(null);
  const [blocks, setBlocks] = useState<any[]>([]);
  const [nodeStatus, setNodeStatus] = useState<any>(null);
  const [mempool, setMempool] = useState<any[]>([]);
  const [searching, setSearching] = useState(false);

  const fetchData = useCallback(async () => {
    try {
      const [blocksRes, statusRes, mempoolRes] = await Promise.all([
        fetch(`${NODE_URL}/blocks?limit=8`),
        fetch(`${NODE_URL}/monitoring/status`),
        fetch(`${NODE_URL}/mempool`),
      ]);

      if (blocksRes.ok) setBlocks(await blocksRes.json() || []);
      if (statusRes.ok) setNodeStatus(await statusRes.json());
      if (mempoolRes.ok) {
        const mTxs = await mempoolRes.json();
        setMempool(Array.isArray(mTxs) ? mTxs : []);
      }
    } catch (e) {
      console.error('Explorer polling error:', e);
    }
  }, []);

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000); // poll every 5s
    return () => clearInterval(interval);
  }, [fetchData]);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchTerm.trim()) return;
    setSearching(true);
    setSearchResult(null);
    try {
      // Try block search
      const blockRes = await fetch(`${NODE_URL}/block?hash=${searchTerm}`);
      if (blockRes.ok) {
        const data = await blockRes.json();
        setSearchResult({ type: 'Block', data });
        setSearching(false); return;
      }
      // Try transaction
      const txRes = await fetch(`${NODE_URL}/transaction?hash=${searchTerm}`);
      if (txRes.ok) {
        const data = await txRes.json();
        setSearchResult({ type: 'Transaction', data });
        setSearching(false); return;
      }
      // Try order (marketplace)
      const orderRes = await fetch(`${NODE_URL}/marketplace?order_id=${searchTerm}`);
      if (orderRes.ok) {
        const data = await orderRes.json();
        if (data.order_info && data.order_info !== 'Order not found') {
          setSearchResult({ type: 'Order', data: data.order_info });
          setSearching(false); return;
        }
      }
      // Try balance/address
      const balRes = await fetch(`${NODE_URL}/balance?address=${searchTerm}`);
      if (balRes.ok) {
        const data = await balRes.json();
        setSearchResult({ type: 'Address', data });
        setSearching(false); return;
      }
      setSearchResult({ type: 'Not Found', data: null });
    } catch (e) {
      setSearchResult({ type: 'Error', data: null });
    }
    setSearching(false);
  };

  return (
    <>
      <div className="content-header">
        <h1>Entity State & Block Explorer</h1>
        <p>Monitor transaction execution, state tree roots, and node memory pool directly from the core.</p>
      </div>

      {/* Node Status Strip */}
      {nodeStatus && (
        <div className="stat-grid" style={{ marginBottom: '1.5rem' }}>
          <div className="stat-card">
            <div className="stat-card-title" style={{ display: 'flex', alignItems: 'center', gap: '0.4rem' }}><Cpu size={14} /> Node Status</div>
            <div className="stat-card-value" style={{ color: 'var(--accent-success)', fontSize: '1.1rem' }}>
              {nodeStatus.status || 'Online'}
            </div>
          </div>
          <div className="stat-card">
            <div className="stat-card-title" style={{ display: 'flex', alignItems: 'center', gap: '0.4rem' }}><Database size={14} /> Chain Height</div>
            <div className="stat-card-value" style={{ fontSize: '1.1rem' }}>{nodeStatus.block_height ?? blocks.length}</div>
          </div>
          <div className="stat-card">
            <div className="stat-card-title" style={{ display: 'flex', alignItems: 'center', gap: '0.4rem' }}><Activity size={14} /> Mempool</div>
            <div className="stat-card-value" style={{ fontSize: '1.1rem' }}>{mempool.length} txs</div>
          </div>
        </div>
      )}

      {/* Search */}
      <div className="glass-card" style={{ marginBottom: '2rem' }}>
        <form style={{ display: 'flex', gap: '1rem' }} onSubmit={handleSearch}>
          <div style={{ flex: 1, position: 'relative' }}>
            <Search size={18} color="var(--text-muted)" style={{ position: 'absolute', top: '50%', transform: 'translateY(-50%)', left: '16px' }} />
            <input
              className="glass-input"
              type="text"
              placeholder="Search by block hash, TX hash, address, or order ID..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              style={{ width: '100%', paddingLeft: '48px', height: '48px' }}
            />
          </div>
          <button className="btn-primary" type="submit" style={{ padding: '0 2rem' }} disabled={searching}>
            {searching ? 'Searching...' : 'Inspect Entity'}
          </button>
        </form>

        {searchResult && (
          <div style={{ marginTop: '1.5rem', padding: '1rem', background: 'rgba(255,255,255,0.04)', borderRadius: '8px', borderLeft: `4px solid ${searchResult.type === 'Not Found' || searchResult.type === 'Error' ? 'var(--accent-danger)' : 'var(--accent-primary)'}` }}>
            <div style={{ fontWeight: '600', marginBottom: '0.5rem', color: searchResult.type === 'Not Found' ? 'var(--accent-danger)' : 'var(--accent-primary)' }}>
              {searchResult.type === 'Not Found' ? '⚠ No entity found for that query.' : `✓ ${searchResult.type} Found`}
            </div>
            {searchResult.data && (
              <pre style={{ fontSize: '0.8rem', color: 'var(--text-muted)', overflowX: 'auto', maxHeight: '200px', margin: 0 }}>
                {JSON.stringify(searchResult.data, null, 2)}
              </pre>
            )}
          </div>
        )}
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 2fr) minmax(0, 1fr)', gap: '1.5rem' }}>
        {/* Blocks Table */}
        <div className="glass-card">
          <div className="glass-card-header">
            <div className="glass-card-title">
              <Database size={20} color="var(--accent-primary)" />
              Recent Blocks <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', fontWeight: 400 }}>(live · 5s)</span>
            </div>
          </div>

          <table className="table-glass">
            <thead>
              <tr>
                <th>Height</th>
                <th>Hash</th>
                <th>Transactions</th>
                <th>Validator</th>
                <th>Age</th>
              </tr>
            </thead>
            <tbody>
              {blocks.length > 0 ? blocks.map((b: any) => (
                <tr key={b.Index ?? b.index}>
                  <td style={{ color: 'var(--accent-primary)', fontWeight: '600' }}>#{b.Index ?? b.index}</td>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>
                    {((b.Hash || b.hash) || '').substring(0, 12)}...
                  </td>
                  <td>{(b.Transactions || b.transactions || []).length} txs</td>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--text-muted)' }}>
                    {((b.Validator || b.validator) || 'GENESIS').substring(0, 10)}...
                  </td>
                  <td style={{ color: 'var(--text-muted)' }}>{timeAgo(b.Timestamp || b.timestamp || 0)}</td>
                </tr>
              )) : (
                <tr>
                  <td colSpan={5} style={{ textAlign: 'center', color: 'var(--text-muted)' }}>
                    {blocks === null ? 'Connecting to node...' : 'No blocks yet.'}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* System Contracts + Mempool */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          <div className="glass-card">
            <div className="glass-card-header">
              <div className="glass-card-title">
                <FileText size={20} color="var(--accent-secondary)" />
                System Contracts
              </div>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
              {[
                { name: 'CONTRACT_TOKEN_SYSTEM', desc: 'Mints & burns TCOIN allocation.', color: 'var(--accent-primary)' },
                { name: 'CONTRACT_STAKING_SYSTEM', desc: 'Node staking infrastructure.', color: 'var(--accent-secondary)' },
                { name: 'CONTRACT_MARKETPLACE_SYSTEM', desc: 'Handles escrow & order dispute.', color: 'var(--accent-success)' },
                { name: 'CONTRACT_GOVERNANCE_SYSTEM', desc: 'On-chain proposals & voting.', color: '#a78bfa' },
              ].map(c => (
                <div key={c.name} style={{ padding: '0.75rem 1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', borderLeft: `3px solid ${c.color}` }}>
                  <div style={{ fontWeight: '600', marginBottom: '0.25rem', fontSize: '0.85rem' }}>{c.name}</div>
                  <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{c.desc}</div>
                </div>
              ))}
            </div>
          </div>

          <div className="glass-card">
            <div className="glass-card-header">
              <div className="glass-card-title">
                <Activity size={20} color="var(--accent-primary)" />
                Mempool ({mempool.length})
              </div>
            </div>
            {mempool.length === 0 ? (
              <div style={{ color: 'var(--text-muted)', textAlign: 'center', padding: '1rem', fontSize: '0.9rem' }}>No pending transactions.</div>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem', maxHeight: '200px', overflowY: 'auto' }}>
                {mempool.slice(0, 10).map((tx: any, i) => (
                  <div key={i} style={{ padding: '0.5rem 0.75rem', background: 'rgba(255,255,255,0.03)', borderRadius: '6px', fontSize: '0.8rem' }}>
                    <span style={{ color: 'var(--text-muted)' }}>{tx.type || 'regular'} · </span>
                    <span style={{ fontFamily: 'monospace' }}>{(tx.sender || '').substring(0, 8)}… → {(tx.recipient || '').substring(0, 8)}…</span>
                    <span style={{ float: 'right', color: 'var(--accent-secondary)' }}>{tx.amount} TCOIN</span>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
