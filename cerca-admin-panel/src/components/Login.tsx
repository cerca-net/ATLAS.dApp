import React, { useState } from 'react';
import { Network, Lock, Mail, Loader2 } from 'lucide-react';
import { supabase } from '../supabaseClient';

interface LoginProps {
  onAuthSuccess: (session: any) => void;
}

export function LoginPage({ onAuthSuccess }: LoginProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const { data, error: authError } = await supabase.auth.signInWithPassword({
        email,
        password,
      });

      if (authError) {
        throw authError;
      }

      if (data.session) {
        const user = data.session.user;
        const role = user.app_metadata?.role || user.user_metadata?.role;
        
        if (role === 'admin') {
          onAuthSuccess(data.session);
        } else {
          // Sign out immediately if not admin to prevent unauthorised session caching
          await supabase.auth.signOut();
          setError('Access Denied: You do not have the required administrator role.');
        }
      } else {
        setError('No active session could be established.');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to authenticate.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-overlay">
      <div className="login-card">
        <div className="login-header">
          <div className="login-logo-container">
            <Network size={36} className="login-logo-icon" />
          </div>
          <h2>ATLAS Admin Control Plane</h2>
          <p>Provide team credentials to access the node environment</p>
        </div>

        <form onSubmit={handleLogin} className="login-form">
          {error && <div className="login-error">{error}</div>}

          <div className="login-input-group">
            <label htmlFor="email">Email Address</label>
            <div className="login-input-wrapper">
              <Mail size={18} className="login-input-icon" />
              <input
                id="email"
                type="email"
                placeholder="name@cercachain.net"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
          </div>

          <div className="login-input-group">
            <label htmlFor="password">Password</label>
            <div className="login-input-wrapper">
              <Lock size={18} className="login-input-icon" />
              <input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
          </div>

          <button type="submit" className="login-submit-btn" disabled={loading}>
            {loading ? (
              <>
                <Loader2 size={18} className="spinner" />
                Authenticating...
              </>
            ) : (
              'Enter Framework'
            )}
          </button>
        </form>

        <div className="login-footer">
          <span>v0.0.1 · Secured via Supabase Auth</span>
        </div>
      </div>
    </div>
  );
}
