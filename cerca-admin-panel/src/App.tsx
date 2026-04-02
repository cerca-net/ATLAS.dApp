import { useState } from 'react';
import {
  LayoutDashboard,
  Cpu,
  Users,
  Database,
  Hash,
  Code,
  Coins,
  Scale,
  Network,
  Activity,
  History,
} from 'lucide-react';
import './App.css';

import { DashboardPage }       from './components/Dashboard';
import { NodeControlPage }     from './components/NodeControl';
import { PeersPage }           from './components/PeersPage';
import { BlocksPage }          from './components/BlocksPage';
import { TransactionsPage }    from './components/TransactionsPage';
import { ContractsPage }       from './components/ContractsPage';
import { FaucetPage }          from './components/Faucet';
import { TreasuryHistoryPage } from './components/TreasuryHistory';
import { ArbitrationPage }     from './components/Arbitration';

type Page =
  | 'dashboard'
  | 'blocks'
  | 'transactions'
  | 'contracts'
  | 'treasury'
  | 'treasury-history'
  | 'arbitration'
  | 'node-control'
  | 'peers';

interface NavItem {
  id: Page;
  label: string;
  icon: React.ReactNode;
  section?: string;
}

const NAV: NavItem[] = [
  // ── OVERVIEW ──────────────────────────────────
  { id: 'dashboard',         label: 'Dashboard',          icon: <LayoutDashboard size={18} />, section: 'OVERVIEW' },

  // ── BLOCKCHAIN ────────────────────────────────
  { id: 'blocks',            label: 'Block Explorer',     icon: <Database size={18} />,        section: 'BLOCKCHAIN' },
  { id: 'transactions',      label: 'Transactions',       icon: <Hash size={18} /> },
  { id: 'contracts',         label: 'System Contracts',   icon: <Code size={18} /> },

  // ── TREASURY ──────────────────────────────────
  { id: 'treasury',          label: 'Treasury',           icon: <Coins size={18} />,           section: 'TREASURY' },
  { id: 'treasury-history',  label: 'Tx History',         icon: <History size={18} /> },
  { id: 'arbitration',       label: 'Arbitration',        icon: <Scale size={18} /> },

  // ── NODE ──────────────────────────────────────
  { id: 'node-control',      label: 'Node Control',       icon: <Cpu size={18} />,             section: 'NODE' },
  { id: 'peers',             label: 'Peers & Validators', icon: <Users size={18} /> },
];

function App() {
  const [active, setActive] = useState<Page>('dashboard');

  return (
    <div className="app-container">
      <aside className="sidebar">
        <div className="sidebar-brand">
          <Network size={24} className="sidebar-icon" />
          <div>
            <div style={{ fontWeight: 700, fontSize: '0.95rem', lineHeight: 1.2 }}>ATLAS Admin</div>
            <div style={{ fontSize: '0.68rem', color: 'var(--text-muted)', marginTop: '1px' }}>CercaChain Testnet</div>
          </div>
        </div>

        <nav className="nav-links">
          {NAV.map((item) => (
            <div key={item.id}>
              {item.section && (
                <div style={{
                  fontSize: '0.63rem',
                  fontWeight: 700,
                  letterSpacing: '0.12em',
                  color: 'var(--text-muted)',
                  padding: '1rem 1rem 0.35rem',
                  textTransform: 'uppercase',
                  userSelect: 'none',
                }}>
                  {item.section}
                </div>
              )}
              <button
                className={`nav-item ${active === item.id ? 'active' : ''}`}
                onClick={() => setActive(item.id)}
              >
                <span className={active === item.id ? 'active-icon' : ''}>{item.icon}</span>
                {item.label}
              </button>
            </div>
          ))}
        </nav>

        <div className="sidebar-footer" style={{ marginTop: 'auto', fontSize: '0.75rem', color: 'var(--text-muted)', borderTop: '1px solid var(--border-subtle)', paddingTop: '1rem' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', marginBottom: '0.3rem' }}>
            <Activity size={11} color="var(--accent-success)" />
            <span style={{ color: 'var(--accent-success)' }}>API Connected</span>
          </div>
          <div>localhost:8080</div>
          <div>v0.0.1 · ATLAS.BC0.0.1</div>
        </div>
      </aside>

      <main className="main-content">
        <div className="ambient-glow-1" />
        <div className="ambient-glow-2" />

        <div className="view-container">
          {active === 'dashboard'        && <DashboardPage />}
          {active === 'blocks'           && <BlocksPage />}
          {active === 'transactions'     && <TransactionsPage />}
          {active === 'contracts'        && <ContractsPage />}
          {active === 'treasury'         && <FaucetPage />}
          {active === 'treasury-history' && <TreasuryHistoryPage />}
          {active === 'arbitration'      && <ArbitrationPage />}
          {active === 'node-control'     && <NodeControlPage />}
          {active === 'peers'            && <PeersPage />}
        </div>
      </main>
    </div>
  );
}

export default App;
