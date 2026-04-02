import { useState, useEffect, useCallback } from 'react';
import { apiFetch } from '../api';
import { Play, Square, PauseCircle, RefreshCw, Terminal, Wifi, WifiOff } from 'lucide-react';

const STATE_COLOR: Record<string, string> = {
  running: 'var(--accent-success)',
  paused: '#f59e0b',
  stopped: 'var(--accent-danger)',
  syncing: 'var(--accent-secondary)',
};

export function NodeControlPage() {
  const [status, setStatus] = useState<any>(null);
  const [logs, setLogs] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchStatus = useCallback(async () => {
    try {
      const [sRes, lRes] = await Promise.all([
        apiFetch('/node/status'),
        apiFetch('/node/logs?limit=100'),
      ]);
      if (sRes.ok) setStatus(await sRes.json());
      if (lRes.ok) {
        const d = await lRes.json();
        setLogs((d.logs || []).reverse()); // newest first
      }
    } catch (e) { console.error(e); }
  }, []);

  useEffect(() => {
    fetchStatus();
    const id = setInterval(fetchStatus, 3000);
    return () => clearInterval(id);
  }, [fetchStatus]);

  const nodeAction = async (action: 'start' | 'stop' | 'pause' | 'sync') => {
    setLoading(true);
    await apiFetch(`/node/${action}`, { method: 'POST' });
    await fetchStatus();
    setLoading(false);
  };

  const state = status?.state || 'unknown';
  const stateColor = STATE_COLOR[state] || 'var(--text-muted)';
  const isRunning = state === 'running';
  const isStopped = state === 'stopped';
  const isPaused = state === 'paused';

  return (
    <>
      <div className="content-header">
        <h1>Node Control</h1>
        <p>Manage the CercaChain main node lifecycle — start, stop, sync, and monitor logs.</p>
      </div>

      {/* Node State Card */}
      <div className="glass-card" style={{ marginBottom: '1.5rem', borderLeft: `4px solid ${stateColor}` }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem', flexWrap: 'wrap' }}>
          <div>
            <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)', marginBottom: '0.25rem' }}>Node State</div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <div style={{ width: 12, height: 12, borderRadius: '50%', background: stateColor, boxShadow: `0 0 10px ${stateColor}`, animation: isRunning ? 'pulse 2s infinite' : 'none' }} />
              <span style={{ fontWeight: 700, fontSize: '1.2rem', color: stateColor, textTransform: 'uppercase' }}>{state}</span>
            </div>
          </div>

          <div style={{ borderLeft: '1px solid rgba(255,255,255,0.1)', paddingLeft: '1.5rem' }}>
            <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Block Height</div>
            <div style={{ fontWeight: 700, fontSize: '1.1rem' }}>{status?.blockHeight ?? '—'}</div>
          </div>

          <div style={{ borderLeft: '1px solid rgba(255,255,255,0.1)', paddingLeft: '1.5rem' }}>
            <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Mempool</div>
            <div style={{ fontWeight: 700, fontSize: '1.1rem' }}>{status?.txPoolSize ?? '—'} txs</div>
          </div>

          <div style={{ borderLeft: '1px solid rgba(255,255,255,0.1)', paddingLeft: '1.5rem' }}>
            <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Validators</div>
            <div style={{ fontWeight: 700, fontSize: '1.1rem' }}>{status?.totalValidators ?? '—'}</div>
          </div>

          {status?.isValidator && (
            <div style={{ borderLeft: '1px solid rgba(255,255,255,0.1)', paddingLeft: '1.5rem' }}>
              <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>This Node</div>
              <div style={{ fontWeight: 600, color: 'var(--accent-success)' }}>✓ Validator</div>
            </div>
          )}

          {/* Controls */}
          <div style={{ marginLeft: 'auto', display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
            <button
              className="btn-success"
              onClick={() => nodeAction('start')}
              disabled={loading || isRunning}
              style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', opacity: (loading || isRunning) ? 0.5 : 1 }}
            >
              <Play size={16} /> Start
            </button>
            <button
              className="btn-secondary"
              onClick={() => nodeAction('pause')}
              disabled={loading || isPaused || isStopped}
              style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', opacity: (loading || isPaused || isStopped) ? 0.5 : 1 }}
            >
              <PauseCircle size={16} /> Pause
            </button>
            <button
              className="btn-secondary"
              onClick={() => nodeAction('sync')}
              disabled={loading || isStopped}
              style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', opacity: (loading || isStopped) ? 0.5 : 1 }}
            >
              <RefreshCw size={16} /> Sync
            </button>
            <button
              className="btn-danger"
              onClick={() => nodeAction('stop')}
              disabled={loading || isStopped}
              style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', opacity: (loading || isStopped) ? 0.5 : 1 }}
            >
              <Square size={16} /> Stop
            </button>
          </div>
        </div>
      </div>

      {/* Validator Info */}
      {status?.validatorAddress && (
        <div className="glass-card" style={{ marginBottom: '1.5rem', padding: '1rem 1.5rem', display: 'flex', gap: '1rem', alignItems: 'center' }}>
          {isRunning ? <Wifi size={20} color="var(--accent-success)" /> : <WifiOff size={20} color="var(--accent-danger)" />}
          <div>
            <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Validator Address</div>
            <div style={{ fontFamily: 'monospace', wordBreak: 'break-all' }}>{status.validatorAddress}</div>
          </div>
          <div style={{ marginLeft: 'auto', textAlign: 'right' }}>
            <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Stake</div>
            <div style={{ fontWeight: 700, color: 'var(--accent-success)' }}>{status.stakeAmount?.toLocaleString() ?? 0} TCOIN</div>
          </div>
        </div>
      )}

      {/* Log Stream */}
      <div className="glass-card">
        <div className="glass-card-header">
          <div className="glass-card-title"><Terminal size={18} color="var(--accent-primary)" /> Node Log Stream <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', fontWeight: 400 }}>(live · 3s)</span></div>
        </div>
        <div style={{ fontFamily: 'monospace', fontSize: '0.8rem', maxHeight: '450px', overflowY: 'auto', display: 'flex', flexDirection: 'column', gap: '0.3rem' }}>
          {logs.length === 0 ? (
            <div style={{ color: 'var(--text-muted)', padding: '1rem' }}>No logs yet. Start the node to begin logging.</div>
          ) : logs.map((log, i) => {
            const color = log.level === 'error' ? 'var(--accent-danger)'
              : log.level === 'success' ? 'var(--accent-success)'
              : log.level === 'warning' ? '#f59e0b'
              : 'var(--text-muted)';
            const prefix = log.level === 'error' ? '✗' : log.level === 'success' ? '✓' : log.level === 'warning' ? '⚠' : '·';
            return (
              <div key={i} style={{ display: 'flex', gap: '1rem', padding: '0.2rem 0.5rem', borderRadius: '4px', background: i % 2 === 0 ? 'rgba(255,255,255,0.02)' : 'transparent' }}>
                <span style={{ color: 'var(--text-muted)', flexShrink: 0 }}>{log.timestamp}</span>
                <span style={{ color, flexShrink: 0 }}>{prefix}</span>
                <span style={{ color: '#e2e8f0' }}>{log.message}</span>
              </div>
            );
          })}
        </div>
      </div>
    </>
  );
}
