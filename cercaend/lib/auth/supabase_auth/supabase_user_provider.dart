import 'package:supabase_flutter/supabase_flutter.dart';
import 'package:rxdart/rxdart.dart';

import '../base_auth_user_provider.dart';

export '../base_auth_user_provider.dart';

class CercaendSupabaseUser extends BaseAuthUser {
  CercaendSupabaseUser(this.user);
  User? user;
  bool get loggedIn => user != null;

  @override
  AuthUserInfo get authUserInfo => AuthUserInfo(
        uid: user?.id,
        email: user?.email,
        phoneNumber: user?.phone,
      );

  @override
  Future? delete() => throw UnimplementedError();

  @override
  Future? updateEmail(String email) async {
    try {
      await Supabase.instance.client.auth.updateUser(
        UserAttributes(email: email),
      );
    } catch (_) {}
  }

  @override
  Future? sendEmailVerification() => throw UnimplementedError();

  @override
  bool get emailVerified {
    // Supabase returns an email_confirmed_at field
    return user?.emailConfirmedAt != null;
  }

  @override
  Future refreshUser() async {
    await Supabase.instance.client.auth.refreshSession();
  }

  @override
  Future? updatePassword(String newPassword) async {
    try {
      await Supabase.instance.client.auth.updateUser(
        UserAttributes(password: newPassword),
      );
    } catch (_) {}
  }
}

/// Helper stream for FlutterFlow to track auth state changes and pipe them into UI.
Stream<BaseAuthUser> cercaendSupabaseUserStream() => Supabase.instance.client.auth.onAuthStateChange
    .debounce((_) => TimerStream(true, const Duration(milliseconds: 50)))
    .map<BaseAuthUser>(
      (authState) => currentUser = CercaendSupabaseUser(authState.session?.user),
    );
