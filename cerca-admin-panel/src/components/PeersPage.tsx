import { useState, useEffect, useCallback } from 'react';
import { apiFetch, shortHash, formatNumber } from '../api';
import { Wifi, Shield, PlusCircle } from 'lucide-react';

export function PeersPage() {
  const [peers, setPeers] = useState<any>({});
  const [validators, setValidators] = useState<any[]>([]);
  const [staking, setStaking] = useState<any>({});
  const [newPeerAddr, setNewPeerAddr] = useState('');
  const [connecting, setConnecting] = useState(false);
  const [connectMsg, setConnectMsg] = useState('');

  const fetchData = useCallback(async () => {
    try {
      const [pRes, vRes, sRes] = await Promise.all([
        apiFetch('/peers'),
        apiFetch('/validators'),
        apiFetch('/staking'),
      ]);
      if (pRes.ok) setPeers(await pRes.json());
      if (vRes.ok) setValidators(await vRes.json() || []);
      if (sRes.ok) setStaking(await sRes.json());
    } catch (e) { console.error(e); }
  }, []);

  useEffect(() => {
    fetchData();
    const id = setInterval(fetchData, 5000);
    return () => clearInterval(id);
  }, [fetchData]);

  const handleConnect = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newPeerAddr.trim()) return;
    setConnecting(true);
    setConnectMsg('');
    try {
      const res = await apiFetch('/connect-peer', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ peer_address: newPeerAddr }),
      });
      let data;
      try {
        data = await res.json();
      } catch (err) {
        data = { message: 'Invalid server response' };
      }
      setConnectMsg(data.message || (res.ok ? 'Connected!' : 'Failed to connect'));
      if (res.ok) { setNewPeerAddr(''); fetchData(); }
    } catch (e: any) {
      setConnectMsg('Error: ' + e.message);
    }
    setConnecting(false);
  };

  const peerList: any[] = peers.peers || [];
  const totalStaked = validators.reduce((sum, v: any) => sum + (Number(v.Stake || v.stake) || 0), 0);

  return (
    <>
      <div className="content-header">
        <h1>Peers & Validators</h1>
        <p>Monitor network topology, connected nodes, active validators and stake distribution.</p>
      </div>

      {/* Stats Row */}
      <div className="stat-grid" style={{ marginBottom: '1.5rem' }}>
        <div className="stat-card">
          <div className="stat-card-title"><Wifi size={13} /> Connected Peers</div>
          <div className="stat-card-value" style={{ color: 'var(--accent-success)' }}>{peers.count ?? 0}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title"><Shield size={13} /> Active Validators</div>
          <div className="stat-card-value" style={{ color: 'var(--accent-primary)' }}>{validators.length}</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title">Total Staked</div>
          <div className="stat-card-value">{formatNumber(totalStaked)} T</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-title">Min Stake</div>
          <div className="stat-card-value">{formatNumber(staking?.min_stake ?? 0)} T</div>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>
        {/* Peer List */}
        <div className="glass-card">
          <div className="glass-card-header">
            <div className="glass-card-title"><Wifi size={18} color="var(--accent-success)" /> Connected Peers</div>
          </div>

          {/* Connect New Peer */}
          <form onSubmit={handleConnect} style={{ display: 'flex', gap: '0.5rem', marginBottom: '1rem' }}>
            <input
              className="glass-input"
              placeholder="/ip4/x.x.x.x/tcp/9000/p2p/Q..."
              value={newPeerAddr}
              onChange={e => setNewPeerAddr(e.target.value)}
              style={{ flex: 1, height: '40px', fontSize: '0.82rem' }}
            />
            <button className="btn-primary" type="submit" disabled={connecting} style={{ padding: '0 1rem', display: 'flex', alignItems: 'center', gap: '0.4rem', whiteSpace: 'nowrap' }}>
              <PlusCircle size={15} /> {connecting ? '...' : 'Connect'}
            </button>
          </form>

          {connectMsg && (
            <div style={{ fontSize: '0.82rem', padding: '0.5rem 0.75rem', borderRadius: '6px', marginBottom: '1rem', background: connectMsg.startsWith('Error') ? 'rgba(239,68,68,0.1)' : 'rgba(39,174,96,0.1)', color: connectMsg.startsWith('Error') ? 'var(--accent-danger)' : 'var(--accent-success)' }}>
              {connectMsg}
            </div>
          )}

          {peerList.length === 0 ? (
            <div style={{ color: 'var(--text-muted)', textAlign: 'center', padding: '2rem', fontSize: '0.9rem' }}>
              No peers connected. Use the field above to add the first peer.
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.6rem' }}>
              {peerList.map((p: any, i) => (
                <div key={i} style={{ padding: '0.75rem 1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <div>
                    <div style={{ fontFamily: 'monospace', fontSize: '0.82rem', marginBottom: '0.2rem' }}>{shortHash(p.id || p.ID, 16)}</div>
                    <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{p.address || p.Address || 'unknown'}</div>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                    <div style={{ width: 8, height: 8, borderRadius: '50%', background: 'var(--accent-success)' }} />
                    <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>online</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Validator List */}
        <div className="glass-card">
          <div className="glass-card-header">
            <div className="glass-card-title"><Shield size={18} color="var(--accent-primary)" /> Registered Validators</div>
          </div>

          {validators.length === 0 ? (
            <div style={{ color: 'var(--text-muted)', textAlign: 'center', padding: '2rem', fontSize: '0.9rem' }}>
              No validators registered yet. Users register by staking TCOIN.
            </div>
          ) : (
            <table className="table-glass">
              <thead>
                <tr><th>Address</th><th>Stake</th><th>Status</th></tr>
              </thead>
              <tbody>
                {validators.map((v: any, i) => (
                  <tr key={i}>
                    <td style={{ fontFamily: 'monospace', fontSize: '0.82rem' }}>{shortHash(v.Address || v.address, 10)}</td>
                    <td style={{ color: 'var(--accent-success)', fontWeight: 600 }}>{formatNumber(v.Stake || v.stake)} T</td>
                    <td>
                      <span style={{ display: 'inline-flex', alignItems: 'center', gap: '0.3rem', fontSize: '0.75rem', padding: '0.2rem 0.6rem', borderRadius: '20px', background: 'rgba(39,174,96,0.15)', color: 'var(--accent-success)' }}>
                        <div style={{ width: 6, height: 6, borderRadius: '50%', background: 'var(--accent-success)' }} />
                        Active
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}

          {staking.staking_contract_address && (
            <div style={{ marginTop: '1rem', padding: '0.75rem 1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '8px', fontSize: '0.8rem' }}>
              <div style={{ color: 'var(--text-muted)', marginBottom: '0.25rem' }}>Staking Contract</div>
              <div style={{ fontFamily: 'monospace', wordBreak: 'break-all', color: '#a78bfa' }}>{staking.staking_contract_address}</div>
              <div style={{ marginTop: '0.5rem', display: 'flex', gap: '1rem', fontSize: '0.75rem', color: 'var(--text-muted)' }}>
                <span>Min stake: {formatNumber(staking.minimum_stake)} T</span>
                <span>Max validators: {staking.max_validators ?? '∞'}</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
