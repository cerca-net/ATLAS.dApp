// Central API config - change this to point to wherever the node is running
export const NODE_URL = 'http://localhost:8080';

export async function apiFetch(path: string, options?: RequestInit) {
  return fetch(`${NODE_URL}${path}`, options);
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
