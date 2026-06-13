import { useState, useEffect } from 'react';
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
  LogOut,
  Loader2,
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
import { GovernancePage }      from './components/GovernancePage';
import { LoginPage }           from './components/Login';
import { supabase }            from './supabaseClient';

type Page =
  | 'dashboard'
  | 'blocks'
  | 'transactions'
  | 'contracts'
  | 'treasury'
  | 'treasury-history'
  | 'arbitration'
  | 'governance'
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
  { id: 'governance',        label: 'Governance',         icon: <Scale size={18} /> },

  // ── NODE ──────────────────────────────────────
  { id: 'node-control',      label: 'Node Control',       icon: <Cpu size={18} />,             section: 'NODE' },
  { id: 'peers',             label: 'Peers & Validators', icon: <Users size={18} /> },
];

function App() {
  const [active, setActive] = useState<Page>('dashboard');
  const [session, setSession] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  // Monitor auth state changes
  useEffect(() => {
    // Check current session
    supabase.auth.getSession().then(({ data: { session: initialSession } }) => {
      if (initialSession) {
        const user = initialSession.user;
        const role = user.app_metadata?.role || user.user_metadata?.role;
        if (role === 'admin') {
          setSession(initialSession);
        } else {
          supabase.auth.signOut();
        }
      }
      setLoading(false);
    });

    // Handle authentication updates
    const { data: { subscription } } = supabase.auth.onAuthStateChange(async (_event, currentSession) => {
      if (currentSession) {
        const user = currentSession.user;
        const role = user.app_metadata?.role || user.user_metadata?.role;
        if (role === 'admin') {
          setSession(currentSession);
        } else {
          await supabase.auth.signOut();
          setSession(null);
        }
      } else {
        setSession(null);
      }
      setLoading(false);
    });

    return () => {
      subscription.unsubscribe();
    };
  }, []);

  const handleLogout = async () => {
    setLoading(true);
    await supabase.auth.signOut();
    setSession(null);
    setLoading(false);
  };

  if (loading) {
    return (
      <div className="login-overlay" style={{ flexDirection: 'column', gap: '1rem' }}>
        <Loader2 size={36} className="spinner" style={{ color: 'var(--accent-primary)' }} />
        <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>Loading Admin Workspace...</p>
      </div>
    );
  }

  if (!session) {
    return <LoginPage onAuthSuccess={(sess) => setSession(sess)} />;
  }

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

        <div className="sidebar-footer" style={{ marginTop: 'auto', fontSize: '0.75rem', color: 'var(--text-muted)', borderTop: '1px solid var(--border-subtle)', paddingTop: '1rem', display: 'flex', flexDirection: 'column', gap: '0.8rem' }}>
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', marginBottom: '0.3rem' }}>
              <Activity size={11} color="var(--accent-success)" />
              <span style={{ color: 'var(--accent-success)' }}>API Connected</span>
            </div>
            <div>localhost:8080</div>
            <div>v0.0.1 · ATLAS.BC0.0.1</div>
          </div>

          <button
            onClick={handleLogout}
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              gap: '0.5rem',
              width: '100%',
              background: 'rgba(239, 68, 68, 0.1)',
              border: '1px solid rgba(239, 68, 68, 0.2)',
              color: '#fca5a5',
              padding: '0.55rem 0.75rem',
              borderRadius: '8px',
              cursor: 'pointer',
              fontWeight: 600,
              fontSize: '0.75rem',
              transition: 'all 0.2s ease',
            }}
          >
            <LogOut size={13} />
            Sign Out
          </button>
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
          {active === 'governance'       && <GovernancePage />}
          {active === 'node-control'     && <NodeControlPage />}
          {active === 'peers'            && <PeersPage />}
        </div>
      </main>
    </div>
  );
}

export default App;
