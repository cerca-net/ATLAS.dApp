import { createClient } from '@supabase/supabase-js';

// Load Supabase URL and Anon Key from environment variables (or fallbacks)
const SUPABASE_URL = import.meta.env.VITE_SUPABASE_URL || 'https://epawttrarbrpzmdbmxyn.supabase.co';
const SUPABASE_ANON_KEY = import.meta.env.VITE_SUPABASE_ANON_KEY || 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVwYXd0dHJhcmJycHptZGJteHluIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzY2MzA0NzMsImV4cCI6MjA5MjIwNjQ3M30.zhjusW5PGl4dRDsi38FtuFwCsfFw_wNSAG7oUuM_Dds';

export const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY);
