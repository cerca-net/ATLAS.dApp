// Central API config — uses environment variables for production builds.
// For Vite, prefix env vars with VITE_ in .env files.

export const NODE_URL = import.meta.env.VITE_NODE_URL || 'http://localhost:8080';
const ADMIN_API_KEY = import.meta.env.VITE_ADMIN_API_KEY || '';

/**
 * Fetch wrapper that adds the admin API key to all requests.
 * If VITE_ADMIN_API_KEY is not set, requests go without auth (dev mode).
 */
export async function apiFetch(path: string, options?: RequestInit) {
  const headers: Record<string, string> = {
    ...(options?.headers as Record<string, string>),
  };

  if (ADMIN_API_KEY) {
    headers['Authorization'] = `Bearer ${ADMIN_API_KEY}`;
  }

  return fetch(`${NODE_URL}${path}`, {
    ...options,
    headers,
  });
}

export function timeAgo(unixTimestamp: number): string {
  const diff = Math.floor(Date.now() / 1000) - unixTimestamp;
  if (diff < 2) return 'just now';
  if (diff < 60) return `${diff}s ago`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  return `${Math.floor(diff / 86400)}d ago`;
}

export function shortHash(hash: string = '', len = 10): string {
  if (!hash) return '—';
  if (hash.length <= len * 2) return hash;
  return hash.substring(0, len) + '…' + hash.substring(hash.length - 6);
}

export function formatNumber(n: number | undefined): string {
  if (n === undefined || n === null) return '—';
  return n.toLocaleString();
}

export type NodePage =
  | 'dashboard'
  | 'node-control'
  | 'peers'
  | 'blocks'
  | 'transactions'
  | 'contracts'
  | 'faucet'
  | 'arbitration';
