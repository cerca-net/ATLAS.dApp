import { useState, useEffect, useCallback } from 'react';
import { apiFetch, timeAgo, shortHash, formatNumber } from '../api';
import {
  Activity, Cpu, Database, Users, ArrowUp,
  CheckCircle, Clock, Zap, BarChart2
} from 'lucide-react';

interface DashMetrics {
  nodeState: string;
  blockHeight: number;
  txPoolSize: number;
  totalValidators: number;
  treasuryBalance: number;
  treasuryAddress: string;
  peerCount: number;
  totalSupply: number;
  marketplaceOrders: number;
  marketplaceVolume: number;
  latestBlocks: any[];
  latestTxs: any[];
}

const STATE_COLOR: Record<string, string> = {
  running: 'var(--accent-success)',
  paused: '#f59e0b',
  stopped: 'var(--accent-danger)',
  syncing: 'var(--accent-secondary)',
};

export function DashboardPage() {
  const [metrics, setMetrics] = useState<Partial<DashMetrics>>({});
  const [loading, setLoading] = useState(true);

  const fetchAll = useCallback(async () => {
    try {
      const [statusRes, treasuryRes, tokenRes, mktRes, blocksRes, mempoolRes, peersRes] = await Promise.all([
        apiFetch('/node/status'),
        apiFetch('/treasury'),
        apiFetch('/token'),
        apiFetch('/marketplace'),
        apiFetch('/blocks?limit=5'),
        apiFetch('/mempool'),
        apiFetch('/peers'),
      ]);

      const status = statusRes.ok ? await statusRes.json() : {};
      const treasury = treasuryRes.ok ? await treasuryRes.json() : {};
      const token = tokenRes.ok ? await tokenRes.json() : {};
      const mkt = mktRes.ok ? await mktRes.json() : {};
      const blocks = blocksRes.ok ? await blocksRes.json() : [];
      const mempool = mempoolRes.ok ? await mempoolRes.json() : [];
      const peers = peersRes.ok ? await peersRes.json() : {};

      setMetrics({
        nodeState: status.state || 'unknown',
        blockHeight: status.blockHeight ?? 0,
        txPoolSize: Array.isArray(mempool) ? mempool.length : 0,
        totalValidators: status.totalValidators ?? 0,
        treasuryBalance: treasury.balance ?? 0,
        treasuryAddress: treasury.address ?? '',
        peerCount: peers.count ?? 0,
        totalSupply: token.total_supply ?? 0,
        marketplaceOrders: mkt.total_orders ?? 0,
        marketplaceVolume: mkt.total_volume ?? 0,
        latestBlocks: Array.isArray(blocks) ? blocks.slice(0, 5) : [],
        latestTxs: Array.isArray(mempool) ? mempool.slice(0, 5) : [],
      });
    } catch (e) { console.error(e); }
    setLoading(false);
  }, []);

  useEffect(() => {
    fetchAll();
    const id = setInterval(fetchAll, 4000);
    return () => clearInterval(id);
  }, [fetchAll]);

  if (loading) return <div style={{ color: 'var(--text-muted)', padding: '3rem', textAlign: 'center' }}>Connecting to node...</div>;

  const stateColor = STATE_COLOR[metrics.nodeState || ''] || 'var(--text-muted)';

  return (
    <>
      <div className="content-header">
        <h1>Network Dashboard</h1>
        <p>Real-time overview of the CercaChain testnet health and activity.</p>
      </div>

      {/* Node State Banner */}
      <div style={{ padding: '0.75rem 1.5rem', background: 'rgba(255,255,255,0.04)', borderRadius: '12px', border: `1px solid ${stateColor}30`, display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '1.5rem' }}>
        <div style={{ width: 10, height: 10, borderRadius: '50%', background: stateColor, boxShadow: `0 0 8px ${stateColor}` }} />
        <span style={{ fontWeight: 600, color: stateColor, textTransform: 'uppercase', fontSize: '0.85rem', letterSpacing: '0.05em' }}>
          Node {metrics.nodeState}
        </span>
        <span style={{ color: 'var(--text-muted)', fontSize: '0.85rem' }}>· ATLAS.BC v0.0.1 · CercaChain Testnet · Live (4s refresh)</span>
      </div>

      {/* KPI Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '1rem', marginBottom: '1.5rem' }}>
        {[
          { icon: <Database size={18} />, label: 'Block Height', value: formatNumber(metrics.blockHeight), color: 'var(--accent-primary)' },
          { icon: <Activity size={18} />, label: 'Mempool', value: `${metrics.txPoolSize} txs`, color: 'var(--accent-secondary)' },
          { icon: <Users size={18} />, label: 'Validators', value: formatNumber(metrics.totalValidators), color: '#a78bfa' },
          { icon: <Zap size={18} />, label: 'Peers Online', value: formatNumber(metrics.peerCount), color: '#38bdf8' },
          { icon: <BarChart2 size={18} />, label: 'Treasury Balance', value: formatNumber(metrics.treasuryBalance) + ' T', color: 'var(--accent-success)' },
          { icon: <Cpu size={18} />, label: 'Total Supply', value: formatNumber(metrics.totalSupply) + ' T', color: '#f472b6' },
          { icon: <CheckCircle size={18} />, label: 'Marketplace Orders', value: formatNumber(metrics.marketplaceOrders), color: '#fb923c' },
          { icon: <ArrowUp size={18} />, label: 'Mkt. Volume', value: formatNumber(metrics.marketplaceVolume) + ' T', color: '#34d399' },
        ].map(({ icon, label, value, color }) => (
          <div key={label} className="stat-card" style={{ borderLeft: `3px solid ${color}` }}>
            <div className="stat-card-title" style={{ display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
              <span style={{ color }}>{icon}</span> {label}
            </div>
            <div className="stat-card-value" style={{ fontSize: '1.3rem', color }}>{value}</div>
          </div>
        ))}
      </div>

      {/* Latest Blocks + Mempool */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>
        <div className="glass-card">
          <div className="glass-card-header">
            <div className="glass-card-title"><Database size={18} color="var(--accent-primary)" /> Latest Blocks</div>
          </div>
          <table className="table-glass">
            <thead><tr><th>#</th><th>Hash</th><th>TXs</th><th>Age</th></tr></thead>
            <tbody>
              {(metrics.latestBlocks || []).length ? metrics.latestBlocks!.map((b: any) => (
                <tr key={b.Index ?? b.index}>
                  <td style={{ color: 'var(--accent-primary)', fontWeight: 600 }}>#{b.Index ?? b.index}</td>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{shortHash(b.Hash || b.hash)}</td>
                  <td>{(b.Transactions || b.transactions || []).length}</td>
                  <td style={{ color: 'var(--text-muted)' }}>{timeAgo(b.Timestamp || b.timestamp || 0)}</td>
                </tr>
              )) : <tr><td colSpan={4} style={{ color: 'var(--text-muted)', textAlign: 'center' }}>No blocks yet</td></tr>}
            </tbody>
          </table>
        </div>

        <div className="glass-card">
          <div className="glass-card-header">
            <div className="glass-card-title"><Clock size={18} color="var(--accent-secondary)" /> Pending Transactions</div>
          </div>
          {(metrics.latestTxs || []).length ? (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.6rem' }}>
              {metrics.latestTxs!.map((tx: any, i) => (
                <div key={i} style={{ padding: '0.6rem 0.8rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', fontSize: '0.82rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <div>
                    <span style={{ color: 'var(--text-muted)' }}>{tx.type || 'regular'}</span>
                    <span style={{ margin: '0 0.4rem', color: 'var(--text-muted)' }}>·</span>
                    <span style={{ fontFamily: 'monospace' }}>{shortHash(tx.sender, 6)} <ArrowRight /> {shortHash(tx.recipient, 6)}</span>
                  </div>
                  <span style={{ color: 'var(--accent-success)', fontWeight: 600 }}>{formatNumber(tx.amount)} T</span>
                </div>
              ))}
            </div>
          ) : (
            <div style={{ color: 'var(--text-muted)', textAlign: 'center', padding: '2rem', fontSize: '0.9rem' }}>Mempool is empty</div>
          )}
        </div>
      </div>

      {/* Treasury Info */}
      {metrics.treasuryAddress && (
        <div className="glass-card" style={{ marginTop: '1.5rem', display: 'flex', alignItems: 'center', gap: '1.5rem' }}>
          <div style={{ fontSize: '2rem' }}>🏦</div>
          <div>
            <div style={{ fontWeight: 600, marginBottom: '0.25rem' }}>Treasury Node Wallet</div>
            <div style={{ fontFamily: 'monospace', fontSize: '0.85rem', color: 'var(--text-muted)', wordBreak: 'break-all' }}>{metrics.treasuryAddress}</div>
          </div>
          <div style={{ marginLeft: 'auto', textAlign: 'right' }}>
            <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Available Balance</div>
            <div style={{ fontSize: '1.4rem', fontWeight: 700, color: 'var(--accent-success)' }}>{formatNumber(metrics.treasuryBalance)} TCOIN</div>
          </div>
        </div>
      )}
    </>
  );
}

function ArrowRight() {
  return <span style={{ margin: '0 0.25rem', color: 'var(--text-muted)' }}>→</span>;
}
