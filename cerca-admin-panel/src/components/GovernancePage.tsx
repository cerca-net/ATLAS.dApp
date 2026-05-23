import React, { useEffect, useState } from 'react';
import { Shield } from 'lucide-react';

interface Proposal {
  ID: string;
  Proposer: string;
  Description: string;
  State: string;
  VotesFor: number;
  VotesAgainst: number;
  StartBlock: number;
  EndBlock: number;
}

export const GovernancePage: React.FC = () => {
  const [proposals, setProposals] = useState<Proposal[]>([]);
  const [stats, setStats] = useState<any>(null);
  const [error, setError] = useState('');
  
  const [newProp, setNewProp] = useState({ proposer: '', description: '', actions: '', duration: 10 });
  const [vote, setVote] = useState({ proposalID: '', voter: '', choice: 'for', weight: 1 });

  const API_BASE = `http://localhost:${localStorage.getItem('selectedNodePort') || 8080}`;

  const fetchGovernance = async () => {
    try {
      const [propRes, statRes] = await Promise.all([
        fetch(`${API_BASE}/governance/proposals`),
        fetch(`${API_BASE}/governance/stats`)
      ]);
      if (propRes.ok) setProposals(await propRes.json());
      if (statRes.ok) setStats(await statRes.json());
    } catch (err: any) {
      setError(err.message || 'Failed to fetch governance data');
    }
  };

  useEffect(() => {
    fetchGovernance();
  }, []);

  const handleCreateProposal = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await fetch(`${API_BASE}/governance/submit-proposal`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newProp)
      });
      if (!res.ok) throw new Error('Submission failed');
      alert('Proposal submitted');
      fetchGovernance();
    } catch (err: any) {
      alert(err.message);
    }
  };

  const handleVote = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await fetch(`${API_BASE}/governance/vote`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(vote)
      });
      if (!res.ok) throw new Error('Vote failed');
      alert('Vote cast');
      fetchGovernance();
    } catch (err: any) {
      alert(err.message);
    }
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <h2 style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <Shield size={24} color="var(--accent-primary)" />
          Governance
        </h2>
        {error && <p style={{ color: 'red' }}>{error}</p>}
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', marginBottom: '1rem' }}>
        <div className="card">
          <h3>Stats</h3>
          {stats ? (
            <div>
              <p>Active Proposals: {stats.activeProposals}</p>
              <p>Total Voters: {stats.totalVoters}</p>
              <p>Voting Power: {stats.totalVotingPower}</p>
            </div>
          ) : <p>Loading stats...</p>}
        </div>
        <div className="card">
          <h3>Proposals</h3>
          {proposals.length === 0 ? <p>No active proposals</p> : (
            proposals.map(p => (
              <div key={p.ID} style={{ border: '1px solid var(--border-subtle)', padding: '0.5rem', marginBottom: '0.5rem', borderRadius: '4px' }}>
                <div style={{ fontWeight: 'bold' }}>{p.ID} ({p.State})</div>
                <div>{p.Description}</div>
                <div>For: {p.VotesFor} | Against: {p.VotesAgainst}</div>
              </div>
            ))
          )}
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
        <div className="card">
          <h3>Create Proposal</h3>
          <form onSubmit={handleCreateProposal} style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <input placeholder="Proposer Address" value={newProp.proposer} onChange={e => setNewProp({...newProp, proposer: e.target.value})} required />
            <input placeholder="Description" value={newProp.description} onChange={e => setNewProp({...newProp, description: e.target.value})} required />
            <input placeholder="Actions (JSON)" value={newProp.actions} onChange={e => setNewProp({...newProp, actions: e.target.value})} required />
            <input type="number" placeholder="Duration (blocks)" value={newProp.duration} onChange={e => setNewProp({...newProp, duration: parseInt(e.target.value)})} required />
            <button type="submit">Submit Proposal</button>
          </form>
        </div>
        <div className="card">
          <h3>Vote</h3>
          <form onSubmit={handleVote} style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <input placeholder="Proposal ID" value={vote.proposalID} onChange={e => setVote({...vote, proposalID: e.target.value})} required />
            <input placeholder="Voter Address" value={vote.voter} onChange={e => setVote({...vote, voter: e.target.value})} required />
            <select value={vote.choice} onChange={e => setVote({...vote, choice: e.target.value})}>
              <option value="for">For</option>
              <option value="against">Against</option>
            </select>
            <input type="number" placeholder="Weight" value={vote.weight} onChange={e => setVote({...vote, weight: parseInt(e.target.value)})} required />
            <button type="submit">Cast Vote</button>
          </form>
        </div>
      </div>
    </div>
  );
};
