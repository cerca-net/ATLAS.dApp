import { useState, useEffect, useCallback } from 'react';
import { apiFetch, formatNumber } from '../api';
import { Coins, ShoppingBag, BarChart2, Vote } from 'lucide-react';

function ContractCard({ icon, title, address, color, children }: {
  icon: React.ReactNode; title: string; address?: string; color: string; children: React.ReactNode;
}) {
  return (
    <div className="glass-card" style={{ borderLeft: `4px solid ${color}` }}>
      <div className="glass-card-header">
        <div className="glass-card-title" style={{ color }}>
          {icon} {title}
        </div>
        {address && (
          <span style={{ fontFamily: 'monospace', fontSize: '0.75rem', color: 'var(--text-muted)' }}>{address}</span>
        )}
      </div>
      {children}
    </div>
  );
}

function StatRow({ label, value, accent }: { label: string; value: string | number; accent?: string }) {
  return (
    <div style={{ display: 'flex', justifyContent: 'space-between', padding: '0.6rem 0', borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
      <span style={{ color: 'var(--text-muted)', fontSize: '0.85rem' }}>{label}</span>
      <span style={{ fontWeight: 600, color: accent || '#e2e8f0' }}>{typeof value === 'number' ? formatNumber(value) : value}</span>
    </div>
  );
}

export function ContractsPage() {
  const [token, setToken] = useState<any>({});
  const [staking, setStaking] = useState<any>({});
  const [marketplace, setMarketplace] = useState<any>({});
  const [governance, setGovernance] = useState<any>({});

  const fetchContracts = useCallback(async () => {
    try {
      const [tRes, sRes, mRes, gRes] = await Promise.all([
        apiFetch('/token'),
        apiFetch('/staking'),
        apiFetch('/marketplace'),
        apiFetch('/governance-contract'),
      ]);
      if (tRes.ok) setToken(await tRes.json());
      if (sRes.ok) setStaking(await sRes.json());
      if (mRes.ok) setMarketplace(await mRes.json());
      if (gRes.ok) setGovernance(await gRes.json());
    } catch (e) { console.error(e); }
  }, []);

  useEffect(() => {
    fetchContracts();
    const id = setInterval(fetchContracts, 8000);
    return () => clearInterval(id);
  }, [fetchContracts]);

  return (
    <>
      <div className="content-header">
        <h1>System Contracts</h1>
        <p>Live state, balances, and configuration of all 4 CercaChain system smart contracts.</p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>

        {/* TOKEN CONTRACT */}
        <ContractCard icon={<Coins size={18} />} title="TCOIN — Token Contract" address={token.contract_address} color="var(--accent-primary)">
          <StatRow label="Token Name" value={token.name || '—'} />
          <StatRow label="Symbol" value={token.symbol || '—'} />
          <StatRow label="Total Supply" value={token.total_supply ?? 0} accent="var(--accent-primary)" />
          <StatRow label="Max Supply" value={token.max_supply ?? 0} />
          <StatRow label="Treasury Address" value={token.treasury_address ? token.treasury_address.substring(0, 16) + '…' : '—'} />
          <StatRow label="Treasury Balance" value={formatNumber(token.treasury_balance) + ' T'} accent="var(--accent-success)" />
          <div style={{ marginTop: '1rem', padding: '0.75rem', background: 'rgba(99,102,241,0.08)', borderRadius: '8px', fontSize: '0.78rem', color: 'var(--text-muted)' }}>
            Manages TCOIN minting, burning, and transfer authority. Treasury wallet is the sole minting authority.
          </div>
        </ContractCard>

        {/* STAKING CONTRACT */}
        <ContractCard icon={<BarChart2 size={18} />} title="Staking Contract" address={staking.staking_contract_address} color="#a78bfa">
          <StatRow label="Total Staked" value={staking.total_staked ?? 0} accent="#a78bfa" />
          <StatRow label="Min Stake" value={formatNumber(staking.minimum_stake) + ' T'} />
          <StatRow label="Max Validators" value={staking.max_validators ?? '∞'} />
          <StatRow label="Active Validators" value={staking.validator_count ?? (staking.validators || []).length} accent="var(--accent-success)" />
          <div style={{ marginTop: '1rem' }}>
            {(staking.validators || []).length > 0 && (
              <>
                <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Top Stakers</div>
                {(staking.validators || []).slice(0, 3).map((v: any, i: number) => (
                  <div key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '0.4rem 0', fontSize: '0.8rem' }}>
                    <span style={{ fontFamily: 'monospace', color: 'var(--text-muted)' }}>{(v.address || v.Address || '').substring(0, 14)}…</span>
                    <span style={{ color: '#a78bfa', fontWeight: 600 }}>{formatNumber(v.stake || v.Stake)} T</span>
                  </div>
                ))}
              </>
            )}
          </div>
          <div style={{ marginTop: '0.5rem', padding: '0.75rem', background: 'rgba(167,139,250,0.08)', borderRadius: '8px', fontSize: '0.78rem', color: 'var(--text-muted)' }}>
            Users stake TCOIN to register as validators. Validators produce blocks and earn rewards.
          </div>
        </ContractCard>

        {/* MARKETPLACE CONTRACT */}
        <ContractCard icon={<ShoppingBag size={18} />} title="Marketplace Contract" address={marketplace.contract_address} color="var(--accent-success)">
          <StatRow label="Total Orders" value={marketplace.total_orders ?? 0} accent="var(--accent-success)" />
          <StatRow label="Active Escrow" value={marketplace.active_escrow ?? 0} />
          <StatRow label="Total Volume" value={formatNumber(marketplace.total_volume ?? 0) + ' T'} accent="var(--accent-success)" />
          <StatRow label="Disputed Orders" value={marketplace.disputed_orders ?? 0} accent="var(--accent-danger)" />
          <StatRow label="Completed Orders" value={marketplace.completed_orders ?? 0} />
          <StatRow label="Refunded Orders" value={marketplace.refunded_orders ?? 0} />
          <div style={{ marginTop: '1rem', padding: '0.75rem', background: 'rgba(39,174,96,0.08)', borderRadius: '8px', fontSize: '0.78rem', color: 'var(--text-muted)' }}>
            Handles escrow locking, release on delivery confirmation, and dispute arbitration. Treasury mediates disputes.
          </div>
        </ContractCard>

        {/* GOVERNANCE CONTRACT */}
        <ContractCard icon={<Vote size={18} />} title="Governance Contract" address={governance.contract_address} color="#f472b6">
          <StatRow label="Total Proposals" value={governance.total_proposals ?? 0} />
          <StatRow label="Active Proposals" value={governance.active_proposals ?? 0} accent="#f472b6" />
          <StatRow label="Passed" value={governance.passed_proposals ?? 0} accent="var(--accent-success)" />
          <StatRow label="Rejected" value={governance.rejected_proposals ?? 0} accent="var(--accent-danger)" />
          <StatRow label="Min Quorum" value={governance.quorum_percentage ? governance.quorum_percentage + '%' : '—'} />
          <StatRow label="Voting Period" value={governance.voting_period_hours ? governance.voting_period_hours + 'h' : '—'} />
          <div style={{ marginTop: '1rem', padding: '0.75rem', background: 'rgba(244,114,182,0.08)', borderRadius: '8px', fontSize: '0.78rem', color: 'var(--text-muted)' }}>
            On-chain governance for network upgrades, parameter changes, and contract modifications by validators.
          </div>

          {(governance.proposals || []).length > 0 && (
            <div style={{ marginTop: '1rem' }}>
              <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Recent Proposals</div>
              {governance.proposals.slice(0, 3).map((p: any, i: number) => (
                <div key={i} style={{ padding: '0.5rem 0.75rem', background: 'rgba(255,255,255,0.03)', borderRadius: '6px', marginBottom: '0.4rem', fontSize: '0.82rem' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <span style={{ fontWeight: 600 }}>{p.id || `Proposal #${i + 1}`}</span>
                    <span style={{ color: p.status === 'active' ? '#f472b6' : p.status === 'passed' ? 'var(--accent-success)' : 'var(--accent-danger)', fontSize: '0.75rem' }}>{p.status}</span>
                  </div>
                  <div style={{ color: 'var(--text-muted)', fontSize: '0.78rem', marginTop: '0.2rem' }}>{p.description || p.title || '—'}</div>
                </div>
              ))}
            </div>
          )}
        </ContractCard>
      </div>
    </>
  );
}
