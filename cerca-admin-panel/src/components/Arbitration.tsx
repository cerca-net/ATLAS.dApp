import { useState, useEffect, useCallback } from 'react';
import { Scale, CheckCircle, Send, AlertCircle, RefreshCw, ShoppingBag } from 'lucide-react';
import { apiFetch, shortHash, formatNumber, timeAgo } from '../api';

// OrderInfo JSON tags from backend: order_id, buyer, seller, amount, fee, status, created_at
interface Dispute {
  order_id: string;
  buyer: string;
  seller: string;
  amount: number;
  fee: number;
  status: number;
  created_at: number;
}

export function ArbitrationPage() {
  const [disputes, setDisputes] = useState<Dispute[]>([]);
  const [selected, setSelected] = useState<Dispute | null>(null);
  const [resolving, setResolving] = useState(false);
  const [lastResult, setLastResult] = useState<{ ok: boolean; msg: string } | null>(null);
  const [apiError, setApiError] = useState('');
  const [mktInfo, setMktInfo] = useState<any>({});

  const fetchDisputes = useCallback(async () => {
    try {
      const [dRes, mRes] = await Promise.all([
        apiFetch('/admin/disputes'),
        apiFetch('/marketplace'),
      ]);
      if (dRes.ok) {
        const data = await dRes.json();
        setDisputes(data.disputes || []);
        setApiError('');
      } else {
        setApiError(`API error ${dRes.status}: ${await dRes.text()}`);
      }
      if (mRes.ok) setMktInfo(await mRes.json());
    } catch (e: any) {
      setApiError('Cannot reach node API: ' + e.message);
    }
  }, []);

  useEffect(() => {
    fetchDisputes();
    const id = setInterval(fetchDisputes, 6000);
    return () => clearInterval(id);
  }, [fetchDisputes]);

  const handleResolve = async (payBuyer: boolean) => {
    if (!selected || resolving) return;
    setResolving(true);
    setLastResult(null);
    try {
      const res = await apiFetch('/admin/resolve-dispute', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          order_id: selected.order_id,
          pay_buyer: payBuyer,
        }),
      });
      if (res.ok) {
        setLastResult({ ok: true, msg: `Dispute ${selected.order_id} resolved — ${payBuyer ? 'Buyer refunded' : 'Seller paid'}. TX queued.` });
        setDisputes(prev => prev.filter(d => d.order_id !== selected.order_id));
        setSelected(null);
      } else {
        const txt = await res.text();
        setLastResult({ ok: false, msg: `Resolution failed: ${txt}` });
      }
    } catch (e: any) {
      setLastResult({ ok: false, msg: 'Error: ' + e.message });
    }
    setResolving(false);
  };

  return (
    <>
      <div className="content-header">
        <h1>Arbitration Panel</h1>
        <p>Review on-chain disputed escrow orders and authorize resolution as referee.</p>
      </div>

      {/* Stats */}
      <div className="stat-grid" style={{ marginBottom: '1.5rem' }}>
        <div className="stat-card" style={{ borderLeft: '3px solid var(--accent-danger)' }}>
          <div className="stat-card-title"><AlertCircle size={13} /> Active Disputes</div>
          <div className="stat-card-value" style={{ color: 'var(--accent-danger)' }}>{disputes.length}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title"><ShoppingBag size={13} /> Total Orders</div>
          <div className="stat-card-value">{formatNumber(mktInfo.total_orders ?? 0)}</div>
        </div>
        <div className="stat-card" style={{ borderLeft: '3px solid var(--accent-success)' }}>
          <div className="stat-card-title">Completed Orders</div>
          <div className="stat-card-value" style={{ color: 'var(--accent-success)' }}>{formatNumber(mktInfo.completed_orders ?? 0)}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title">Locked Escrow</div>
          <div className="stat-card-value">
            {formatNumber(disputes.reduce((s, d) => s + (d.amount || 0), 0))} T
          </div>
        </div>
      </div>

      {/* API error banner */}
      {apiError && (
        <div style={{ padding: '0.75rem 1rem', marginBottom: '1rem', background: 'rgba(239,68,68,0.1)', borderRadius: '8px', color: 'var(--accent-danger)', fontSize: '0.85rem', display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
          <AlertCircle size={16} /> {apiError}
        </div>
      )}

      {/* Result message */}
      {lastResult && (
        <div style={{ padding: '0.75rem 1rem', marginBottom: '1rem', background: lastResult.ok ? 'rgba(16,185,129,0.1)' : 'rgba(239,68,68,0.1)', borderRadius: '8px', color: lastResult.ok ? 'var(--accent-success)' : 'var(--accent-danger)', fontSize: '0.85rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span>{lastResult.msg}</span>
          <button onClick={() => setLastResult(null)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'inherit', fontSize: '1rem' }}>✕</button>
        </div>
      )}

      {/* Resolution panel */}
      {selected && (
        <div className="glass-card" style={{ marginBottom: '1.5rem', borderLeft: '4px solid var(--accent-secondary)', animation: 'fadeIn 0.2s ease-out' }}>
          <div className="glass-card-header">
            <div className="glass-card-title">
              <Scale size={18} color="var(--accent-secondary)" />
              Resolving: {selected.order_id}
            </div>
            <button onClick={() => setSelected(null)} className="btn-secondary" style={{ padding: '0.3rem 0.75rem' }}>Close</button>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '0.75rem', marginBottom: '1.5rem' }}>
            {[
              { label: 'Order ID', value: selected.order_id, accent: 'var(--accent-primary)' },
              { label: 'Locked Amount', value: formatNumber(selected.amount) + ' TCOIN', accent: 'var(--accent-success)' },
              { label: 'Platform Fee', value: formatNumber(selected.fee) + ' TCOIN', accent: 'var(--text-muted)' },
              { label: 'Created', value: selected.created_at ? timeAgo(selected.created_at) : '—', accent: 'var(--text-muted)' },
              { label: 'Buyer', value: shortHash(selected.buyer, 12), accent: '#38bdf8' },
              { label: 'Seller', value: shortHash(selected.seller, 12), accent: '#a78bfa' },
            ].map(({ label, value, accent }) => (
              <div key={label} style={{ padding: '0.75rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px' }}>
                <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', marginBottom: '0.2rem' }}>{label}</div>
                <div style={{ fontWeight: 600, color: accent, fontFamily: label === 'Buyer' || label === 'Seller' ? 'monospace' : 'inherit' }}>{value}</div>
              </div>
            ))}
          </div>

          <div style={{ padding: '0.75rem 1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', marginBottom: '1.5rem', fontSize: '0.82rem', color: 'var(--text-muted)' }}>
            You are executing <code style={{ color: '#e2e8f0', background: 'rgba(255,255,255,0.08)', padding: '0.1rem 0.4rem', borderRadius: '4px' }}>resolveDispute</code> on the Marketplace contract via the Treasury wallet. This action is irreversible and will be confirmed on-chain.
          </div>

          <div style={{ display: 'flex', gap: '1rem' }}>
            <button
              className="btn-success"
              style={{ flex: 1, padding: '0.875rem', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '0.5rem', fontSize: '0.95rem', opacity: resolving ? 0.5 : 1 }}
              onClick={() => handleResolve(true)}
              disabled={resolving}
            >
              <CheckCircle size={18} /> {resolving ? 'Submitting...' : 'Refund Buyer'}
            </button>
            <button
              className="btn-primary"
              style={{ flex: 1, padding: '0.875rem', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '0.5rem', fontSize: '0.95rem', opacity: resolving ? 0.5 : 1 }}
              onClick={() => handleResolve(false)}
              disabled={resolving}
            >
              <Send size={18} /> {resolving ? 'Submitting...' : 'Pay Seller'}
            </button>
          </div>
        </div>
      )}

      {/* Disputes Table */}
      <div className="glass-card">
        <div className="glass-card-header">
          <div className="glass-card-title">
            <AlertCircle size={18} color="var(--accent-danger)" />
            Open Disputes <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', fontWeight: 400 }}>(live · 6s)</span>
          </div>
          <button onClick={fetchDisputes} className="btn-secondary" style={{ padding: '0.35rem 0.75rem', display: 'flex', alignItems: 'center', gap: '0.4rem', fontSize: '0.82rem' }}>
            <RefreshCw size={13} /> Refresh
          </button>
        </div>

        {disputes.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '3rem', color: 'var(--text-muted)' }}>
            <Scale size={40} style={{ opacity: 0.3, marginBottom: '1rem' }} />
            <div style={{ fontWeight: 600, marginBottom: '0.5rem' }}>No active disputes</div>
            <div style={{ fontSize: '0.85rem' }}>
              Disputes appear here when a marketplace order is flagged via the <code style={{ color: '#e2e8f0' }}>disputeOrder</code> contract call.
              {mktInfo.total_orders === 0 && ' No marketplace orders exist yet.'}
            </div>
          </div>
        ) : (
          <table className="table-glass">
            <thead>
              <tr>
                <th>Order ID</th>
                <th>Buyer</th>
                <th>Seller</th>
                <th>Locked</th>
                <th>Fee</th>
                <th>Age</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {disputes.map((d) => (
                <tr key={d.order_id}>
                  <td style={{ fontWeight: 700, color: 'var(--accent-primary)' }}>{d.order_id}</td>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem', color: '#38bdf8' }}>{shortHash(d.buyer, 10)}</td>
                  <td style={{ fontFamily: 'monospace', fontSize: '0.8rem', color: '#a78bfa' }}>{shortHash(d.seller, 10)}</td>
                  <td style={{ color: 'var(--accent-success)', fontWeight: 600 }}>{formatNumber(d.amount)} T</td>
                  <td style={{ color: 'var(--text-muted)' }}>{formatNumber(d.fee)} T</td>
                  <td style={{ color: 'var(--text-muted)', fontSize: '0.82rem' }}>{d.created_at ? timeAgo(d.created_at) : '—'}</td>
                  <td>
                    <button
                      className="btn-primary"
                      style={{ padding: '0.35rem 0.8rem', fontSize: '0.82rem' }}
                      onClick={() => { setSelected(d); setLastResult(null); }}
                    >
                      Review
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </>
  );
}
