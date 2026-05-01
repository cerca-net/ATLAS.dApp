import 'dart:async';
import 'package:flutter/material.dart';
import 'package:supabase_flutter/supabase_flutter.dart';
import '../auth_manager.dart';
import 'supabase_user_provider.dart';
import '/app_state.dart';
import '/services/blockchain/blockchain_service.dart';

export '../base_auth_user_provider.dart';

class SupabaseAuthManager extends AuthManager
    with EmailSignInManager {
      
  @override
  Future signOut() async {
    try {
      try {
        await BlockchainService().stopNode();
      } catch (e) {
        debugPrint('Warning: Failed to stop blockchain node gracefully on signout: $e');
      }
      
      // Clear wallet state from global app state
      FFAppState().isWalletConnected = false;
      FFAppState().walletAddress = '';
      FFAppState().walletBalance = 0.0;
      FFAppState().sessionToken = '';
    } catch (e) {
      debugPrint('Warning: Failed to clear session state on signout: $e');
    }
    
    return Supabase.instance.client.auth.signOut();
  }

  @override
  Future deleteUser(BuildContext context) async {
    // Supabase disables user deletion from client side by default for security.
    debugPrint('Supabase: deleteUser is typically a backend function.');
  }

  @override
  Future updateEmail({
    required String email,
    required BuildContext context,
  }) async {
    try {
      await Supabase.instance.client.auth.updateUser(
        UserAttributes(email: email),
      );
    } on AuthException catch (e) {
      ScaffoldMessenger.of(context).hideCurrentSnackBar();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: ${e.message}')),
      );
    }
  }

  @override
  Future resetPassword({
    required String email,
    required BuildContext context,
  }) async {
    try {
      await Supabase.instance.client.auth.resetPasswordForEmail(email);
    } on AuthException catch (e) {
      ScaffoldMessenger.of(context).hideCurrentSnackBar();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: ${e.message}')),
      );
      return null;
    }
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Password reset email sent')),
    );
  }

  @override
  Future<BaseAuthUser?> signInWithEmail(
    BuildContext context,
    String email,
    String password,
  ) async {
    try {
      final response = await Supabase.instance.client.auth.signInWithPassword(
        email: email,
        password: password,
      );
      return CercaendSupabaseUser(response.user);
    } on AuthException catch (e) {
      ScaffoldMessenger.of(context).hideCurrentSnackBar();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: ${e.message}')),
      );
      return null;
    }
  }

  @override
  Future<BaseAuthUser?> createAccountWithEmail(
    BuildContext context,
    String email,
    String password,
  ) async {
    try {
      final response = await Supabase.instance.client.auth.signUp(
        email: email,
        password: password,
      );
      if (response.user != null) {
        try {
          await Supabase.instance.client.from('users').upsert({
            'id': response.user!.id,
            'email': email,
            'created_time': DateTime.now().toIso8601String(),
          });
        } catch (dbError) {
          debugPrint('Error inserting user data into users table: $dbError');
        }
      }
      return CercaendSupabaseUser(response.user);
    } on AuthException catch (e) {
      ScaffoldMessenger.of(context).hideCurrentSnackBar();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: ${e.message}')),
      );
      return null;
    }
  }

  Future<BaseAuthUser?> signInWithGoogle(BuildContext context) async {
    try {
      final success = await Supabase.instance.client.auth.signInWithOAuth(
        OAuthProvider.google,
      );
      return null; // Will redirect via OAuth
    } on AuthException catch (e) {
      ScaffoldMessenger.of(context).hideCurrentSnackBar();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: ${e.message}')),
      );
      return null;
    }
  }
}
