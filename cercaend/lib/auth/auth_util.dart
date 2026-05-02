/// Canonical auth utilities for CercaChain.
///
/// This file is the correct import path for auth utilities.
/// It re-exports the Supabase-backed auth implementation that currently
/// lives at the legacy `firebase_auth/auth_util.dart` path.
///
/// All NEW code should import from here:
///   import '/auth/auth_util.dart';
///
/// Existing imports of '/auth/firebase_auth/auth_util.dart' continue to work
/// but should be migrated to this path over time.
export 'firebase_auth/auth_util.dart';
